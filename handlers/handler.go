package handlers

import (
	"goplay2/audio"
	"goplay2/homekit"
	"goplay2/pairing"
	"goplay2/rtsp"
	"log"

	"github.com/brutella/hc/hap"
)

type Rstp struct {
	streams map[string]*audio.Server
	pairing *pairing.Controller
	player  *audio.Player
}

func NewRstpHandler(deviceName string, player *audio.Player) (*Rstp, error) {

	ctrl, err := pairing.NewController(deviceName)
	if err != nil {
		return nil, err
	}
	return &Rstp{
		streams: make(map[string]*audio.Server),
		pairing: ctrl,
		player:  player,
	}, nil
}

func (r *Rstp) OnConnOpen(conn *rtsp.Conn) {
	conn.SetNetConn(hap.NewConnection(conn.NetConn(), homekit.Server.Context))
}

func (r *Rstp) OnRequest(conn *rtsp.Conn, request *rtsp.Request) {
	log.Printf("request received : %s %s body %d", request.Method, request.URL, len(request.Body))
}

func (r *Rstp) Handle(conn *rtsp.Conn, req *rtsp.Request) (*rtsp.Response, error) {
	switch req.Method {
	case "GET":
		return r.OnGetWeb(req)
	case "POST":
		return r.OnPostWeb(conn, req)
	case "OPTIONS": // TODO
		return r.OnPostWeb(conn, req)
	case "SETUP":
		return r.OnSetupWeb(req)
	case "GET_PARAMETER":
		return r.OnGetParameterWeb(req)
	case "SET_PARAMETER":
		return r.OnSetParameterWeb(req)
	case "RECORD":
		return r.OnRecordWeb(req)
	case "SETPEERS":
		return r.OnSetPeerWeb(req)
	case "SETRATEANCHORTIME":
		return r.OnSetRateAnchorTime(req)
	case "FLUSHBUFFERED":
		return r.OnFlushBuffered(req)
	case "TEARDOWN":
		return r.OnTeardownWeb(req)
	}
	return &rtsp.Response{StatusCode: rtsp.StatusNotImplemented}, nil
}

func (r *Rstp) OnResponse(conn *rtsp.Conn, resp *rtsp.Response) {
	log.Printf("response sent : body %d", len(resp.Body))
}

// The ANNOUNCE request tells the RTSP server about stream properties using [SDP]. Codec informations and encryption keys are of particular interest.

// The OPTIONS request asks the RTSP server for its supported methods. Apple TV supports the following methods:
// ANNOUNCE, SETUP, RECORD, PAUSE, FLUSH, TEARDOWN, OPTIONS, GET_PARAMETER, SET_PARAMETER, POST and GET.
// OPTIONS * RTSP/1.0
// CSeq: 3
// User-Agent:
// RTSP/1.0 200 OK
// Public: ANNOUNCE, SETUP, RECORD, PAUSE, FLUSH, TEARDOWN, OPTIONS,
// GET_PARAMETER, SET_PARAMETER, POST, GET
// Server: AirTunes/130.14
// CSeq: 3
