package cast

import (
	"context"
	"fmt"
	"goplay2/audio"
	"goplay2/config"
	"goplay2/globals"
	"goplay2/handlers"
	"goplay2/homekit"
	"goplay2/ptp"
	"goplay2/rtsp"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/brutella/hc/event"
	"github.com/grandcat/zeroconf"
	"github.com/vishen/go-chromecast/dns"
)

/*	CastDevice extends dns.CastEntry
 *	by gathering the MAC Address as well the associated
 *	net.Interface for the Chromecast Device
 */
type CastDevice struct {
	dns.CastEntry
	netInterface *net.Interface
	HwAddr       net.HardwareAddr
}

/*	CastDeviceEntry extends dns.CastDNSEntry
 *	by providing the MAC Address as well the associated
 *	net.Interface for the Chromecast Device
 */
type DeviceEntry interface {
	dns.CastDNSEntry
	NetInterface() *net.Interface
	HardwareAddr() net.HardwareAddr
	Entry() *CastDevice
}

func (d CastDevice) NetInterface() {

}

func (d CastDevice) Entry() *CastDevice {
	return &d
}

func (d CastDevice) HardwareAddr() net.HardwareAddr {
	return d.HwAddr
}

type RTSPHandle = handlers.Rstp

type AirplayDevice struct {
	s        *zeroconf.Server
	clock    *ptp.VirtualClock
	ptp      *ptp.Server
	audioBuf *audio.Ring
	player   *audio.Player
	rtsp     *rtsp.Server
	rHandle  *RTSPHandle
	// TODO: do we need an event server?
}

type AirplayReciever interface {
	Create(device *DeviceEntry) *AirplayDevice
	Registration(name *string, net *net.Interface) *zeroconf.Server
	Clock() *ptp.VirtualClock
	PtpServer() *ptp.Server
	RingBuffer() *audio.Ring
	Player() *audio.Player
	RTSPHandler() *handlers.Rstp
	RunRTSP() error
}

func (a *AirplayDevice) Create(device *DeviceEntry) *AirplayDevice {
	a := &AirplayDevice{}
	if device != nil {
		// TODO: Fix
		a.Registration(&(*device).Entry().Name, nil)
	}
	return a
}

func (a *AirplayDevice) Registration(name *string, net *net.Interface) *zeroconf.Server {
	// Create 0-cfg
	if a.s != nil {
		return a.s
	}
	a.s, err := zeroconf.Register(*name, "_airplay._tcp", "local.",
		7000, homekit.Device.ToRecords(), []net.Interface{*net})
	if err != nil {
		panic(fmt.Errorf("Registration: %s", err))
	}
	return a.s
}

func (a AirplayDevice) Clock(delay *int64) *ptp.VirtualClock {
	if a.clock != nil {
		return a.clock
	}
	if delay != nil {
		a.clock = ptp.NewVirtualClock(*delay)
	} else {
		a.clock = ptp.NewVirtualClock(0)
	}
	return a.clock
}

func (a AirplayDevice) PtpServer() *ptp.Server {
	if a.ptp != nil {
		return a.ptp
	}
	a.ptp = ptp.NewServer(a.Clock(nil))
	return a.ptp
}
func (a AirplayDevice) RingBuffer() *audio.Ring {
	if a.audioBuf != nil {
		return a.audioBuf
	}
	// Divided by 100 -> average size of a RTP packet
	a.audioBuf = audio.NewRing(globals.BufferSize / 100)
	return a.audioBuf
}

func (a AirplayDevice) Player(clockDelay *int64) *audio.Player {
	if a.player != nil {
		return a.player
	}
	a.player = audio.NewPlayer(a.Clock(clockDelay), a.audioBuf)
	return a.player
}

func (a AirplayDevice) RTSPHandler() *handlers.Rstp {
	if a.rHandle != nil {
		return a.rHandle
	}
	a.rHandle, err = handlers.NewRstpHandler(a.DeviceName(), a.Player())
	if err != nil {
		fmt.Errorf("RTSPHandler: %s", err)
	}
	return a.rHandle
}
func (a AirplayDevice) RunRTSP() error {
	err := rtsp.RunRtspServer(*a.rHandle)
	if err != nil {
		panic(fmt.Errorf("RunRTSP: %s", err))
	}
	return nil
}

// discoverLocalInterfaces disovers interfaces used
// by the device executing this function
func discoverLocalInterfaces() {
	netFaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, face := range netFaces {
		addrs := face.Addrs()
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP)
			}
		}
	}
}

func GetCastDevices(ifaceName string, timeoutSec uint) {
	var iface *net.Interface
	var err error
	if ifaceName != "" {
		if iface, err = net.InterfaceByName(ifaceName); err != nil {
			exit("unable to find interface %q: %v", ifaceName, err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSec))
	defer cancel()
	// Returns a <- chan CastEntry
	//? Refrence: https://thomasnguyen.site/function-returning-channel-pattern-in-go
	castEntryChan, err := dns.DiscoverCastDNSEntries(ctx, iface)
	if err != nil {
		exit("unable to discover chromecast devices: %v", err)
	}
	// TODO: get MAC address
	//? See go-chromecast.cmd.ls for use
	for i, d := range castEntryChan {
		fmt.Printf("%d) device=%q device_name=%q address=\"%s:%d\" uuid=%q",
			i+1, d.Device, d.DeviceName, d.AddrV4, d.Port, d.UUID)
	}
}

func CreateVirtualAirplayReciever(device *CastDevice) {
	// Interface of chromecast
	iFace, err := net.InterfaceByName(ifName)
	if err != nil {
		panic(err)
	}

	// uses net.InterfaceByName
	macAddress := strings.ToUpper(device.HwAddr.String())
	ipAddresses, err := iFace.Addrs()
	if err != nil {
		panic(err)
	}
	var a AirplayDevice = &AirplayDevice{}
	// Register ZeroConf Service
	server := a.Registration()
	delay := 100
	// clock := ptp.NewVirtualClock(delay)
	a.Clock(delay)

	// ptpServer := ptp.NewServer(clock)
	a.PtpServer()

	// Divided by 100 -> average size of a RTP packet
	// audioBuffer := audio.NewRing(globals.BufferSize / 100)
	a.RingBuffer()

	// player := audio.NewPlayer(clock, audioBuffer)
	a.Player()

	wg := new(sync.WaitGroup)
	wg.Add(4)

	go func() {
		a.Player(nil).Run()
		wg.Done()
	}()

	go func() {
		event.RunEventServer()
		wg.Done()
	}()

	go func() {
		a.PtpServer().Serve()
		wg.Done()
	}()

	go func() {
		handler, e := handlers.NewRstpHandler(config.Config.DeviceName, player)
		if e != nil {
			panic(e)
		}
		e = rtsp.RunRtspServer(handler)
		if e != nil {
			panic(e)
		}
		wg.Done()
	}()

	wg.Wait()
}

// DiscoverCastDNSEntries will return a channel with any cast dns entries
// found.
func DiscoverCastDNSEntries(ctx context.Context, iface *net.Interface) (<-chan CastDevice, error) {
	var opts = []zeroconf.ClientOption{}
	// Act as a client using a Network Interface
	if iface != nil {
		opts = append(opts, zeroconf.SelectIfaces([]net.Interface{*iface}))
	}
	resolver, err := zeroconf.NewResolver(opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create new zeroconf resolver: %w", err)
	}
	castDNSEntriesChan := make(chan CastEntry, 5)
	entriesChan := make(chan *zeroconf.ServiceEntry, 5)
	go func() {
		// look for client's on the Network Interface that support Chromecast
		if err := resolver.Browse(ctx, "_googlecast._tcp", "local", entriesChan); err != nil {
			log.WithError(err).Error("unable to browser for mdns entries")
			return
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(castDNSEntriesChan)
				return
			case entry := <-entriesChan:
				if entry == nil {
					continue
				}
				castEntry := CastDevice{
					Port: entry.Port,
					Host: entry.HostName,
					//? Added by Jason
					HwAddr: entry.HardwareAddr,
				}
				if len(entry.AddrIPv4) > 0 {
					castEntry.AddrV4 = entry.AddrIPv4[0]
				}
				if len(entry.AddrIPv6) > 0 {
					castEntry.AddrV6 = entry.AddrIPv6[0]
				}
				infoFields := make(map[string]string, len(entry.Text))
				for _, value := range entry.Text {
					if kv := strings.SplitN(value, "=", 2); len(kv) == 2 {
						key := kv[0]
						val := kv[1]

						infoFields[key] = val

						switch key {
						case "fn":
							castEntry.DeviceName = decode(val)
						case "md":
							castEntry.Device = decode(val)
						case "id":
							castEntry.UUID = val
						}
					}
				}
				castEntry.InfoFields = infoFields
				castDNSEntriesChan <- castEntry
			}
		}
	}()
	return castDNSEntriesChan, nil
}

// decode attempts to decode the passed in string using escaped utf8 bytes.
// some DNS entries for other languages seem to include utf8 escape sequences as
// part of the name.
func decode(val string) string {
	if strings.Index(val, "\\") == -1 {
		return val
	}

	var (
		r        []rune
		toDecode []byte
	)

	decodeRunes := func() {
		if len(toDecode) > 0 {
			for len(toDecode) > 0 {
				rr, size := utf8.DecodeRune(toDecode)
				r = append(r, rr)
				toDecode = toDecode[size:]
			}
			toDecode = []byte{}
		}
	}

	for i := 0; i < len(val); {
		if val[i] == '\\' {
			if i+3 < len(val) {
				v, err := strconv.Atoi(val[i+1 : i+4])
				if err == nil {
					toDecode = append(toDecode, byte(v))
					i += 4
					continue
				}
			}
		}
		decodeRunes()
		r = append(r, rune(val[i]))
		i++
	}
	decodeRunes()
	return string(r)
}
