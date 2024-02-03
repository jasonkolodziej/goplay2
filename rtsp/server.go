package rtsp

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"goplay2/globals"
	"log"
	"net"
	"reflect"
)

const (
	serverConnReadBufferSize  = 4096
	serverConnWriteBufferSize = 4096
)

type Handler interface {
	Handle(conn *Conn, req *Request) (*Response, error)
	OnRequest(conn *Conn, request *Request)
	OnResponse(conn *Conn, resp *Response)
	OnConnOpen(conn *Conn)
}

type Server struct {
	h  Handler
	bw *bufio.Writer
	br *bufio.Reader
}

type Conn struct {
	c net.Conn
}

func (c *Conn) NetConn() net.Conn {
	return c.c
}

func (c *Conn) Close() error {
	return c.c.Close()
}

func (c *Conn) SetNetConn(conn net.Conn) {
	c.c = conn
}

func RunRtspServer(handlers Handler) (err error) {

	s := &Server{
		h: handlers,
	}
	if l, err := net.Listen("tcp4", ":7000"); err == nil {
		defer l.Close()
		for {

			conn, err := l.Accept()
			if err != nil {
				globals.ErrLog.Println("Error accepting: ", err.Error())
				return err
			}
			rConn := &Conn{
				c: conn,
			}
			s.h.OnConnOpen(rConn)
			go s.handleRstpConnection(rConn)
		}
	}
	return err
}

func (s *Server) handleRstpConnection(conn *Conn) {
	defer conn.Close()

	s.br = bufio.NewReaderSize(conn.NetConn(), serverConnReadBufferSize)
	s.bw = bufio.NewWriterSize(conn.NetConn(), serverConnWriteBufferSize)

	for {
		request, err := parseRequest(s.br)
		if err != nil {
			globals.ErrLog.Printf("Error parsing RSTP request %v \n", err)
			return
		}
		s.h.OnRequest(conn, request)
		request.Log()
		response, err := s.h.Handle(conn, request)
		if err != nil {
			globals.ErrLog.Printf("Error handling RSTP request %v \n", err)
			return
		}
		response.Log()
		err = s.flushResponse(conn, request, response)
		if err != nil {
			globals.ErrLog.Printf("Error flusing RSTP response %v \n", err)
			return
		}
	}

}

func (s *Server) flushResponse(conn *Conn, req *Request, resp *Response) error {
	if resp.Header == nil {
		resp.Header = make(Header)
	}
	resp.Header["CSeq"] = req.Header["CSeq"]
	resp.Header["Server"] = HeaderValue{"AirTunes/366.0"}
	s.h.OnResponse(conn, resp)
	return resp.Write(s.bw)
}

func parseRequest(br *bufio.Reader) (*Request, error) {
	var req Request
	var err error
	if err = req.Read(br); err != nil {
		return nil, err
	}
	return &req, nil
}

func BodyHelper(obj any) {
	ref := reflect.ValueOf(obj)
	typ := reflect.TypeOf(obj)
	switch ref.Kind() {
	// if its a pointer, resolve its value
	case reflect.Ptr:
		typ = reflect.PtrTo(typ)
		ref = reflect.Indirect(ref)
	case reflect.Interface:
		typ = typ.Elem()
		ref = ref.Elem()
	case reflect.Struct:
		break
	default:
		// should double check we now
		// have a struct (could still be anything)
		log.Fatal("unexpected type")

	}
	log.Printf("Type %s: %v\n",
		typ.Name(), ref)
}

func (req *Request) Log() {
	// switch req.Method {
	// case "GET":
	// 	log.Printf("%s %s %s\n", req.Method, req.URL, req.URL.Scheme)
	// case "POST":
	// 	log.Printf("%s %s %s\n", req.Method, req.Path, req.URL.Scheme)
	// }
	var headers string
	for k, v := range req.Header {
		headers = string(fmt.Appendf([]byte(headers), "%s: %v\n", k, v))
	}
	// if req.Method == "POST" && req.Query != "" {
	// 	log.Printf("%s\n", req.Query)
	// }
	log.Printf("%s", fmt.Sprintf("\n%s %s %s\n%s%s\nBody length: %d\n", req.Method, req.URL,
		req.URL.Scheme, headers, hex.EncodeToString(req.Body), len(req.Body)))
}
func (req *Response) Log() {
	var headers string
	for k, v := range req.Header {
		headers = string(fmt.Appendf([]byte(headers), "%s: %v\n", k, v))
	}
	log.Printf("%s", fmt.Sprintf("\n%s %d %s\n%s%s\nBody length: %d\n", rtspProtocol10, req.StatusCode,
		StatusMessages[req.StatusCode], headers, hex.EncodeToString(req.Body), len(req.Body)))
}
