package other

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/grandcat/zeroconf"

	log "github.com/sirupsen/logrus"
	"github.com/vishen/go-chromecast/dns"
)

/*	CastDevice extends dns.CastEntry
 *	by gathering the MAC Address as well the associated
 *	net.Interface for the Chromecast Device
 */
type CastDevice struct {
	netInterface *net.Interface
	HwAddr       net.HardwareAddr
	dns.CastEntry
}

/*	CastDeviceEntry extends dns.CastDNSEntry
 *	by providing the MAC Address as well the associated
 *	net.Interface for the Chromecast Device
 */
type DeviceEntry interface {
	dns.CastDNSEntry
	NetInterface() *net.Interface
	HardwareAddr() net.HardwareAddr
	// Entry() *CastDevice
}

func (d CastDevice) NetInterface() *net.Interface {
	return nil
}

func (d CastDevice) Entry() *dns.CastEntry {
	return nil
}

func (d CastDevice) HardwareAddr() net.HardwareAddr {
	return d.HwAddr
}

// discoverLocalInterfaces disovers interfaces used
// by the device executing this function
func discoverLocalInterfaces() []net.Interface {
	var ret []net.Interface
	netFaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, face := range netFaces {
		addrs, err := face.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
				fmt.Println(ipNet.IP)
				ret = append(ret, face)
			}
		}
	}
	return ret
}

func GetCastDevices(ifaceName string, timeoutSec uint) {
	var iface *net.Interface
	var err error
	if ifaceName != "" {
		if iface, err = net.InterfaceByName(ifaceName); err != nil {
			log.Debugf("unable to find interface %q: %v", ifaceName, err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSec))
	defer cancel()
	// Returns a <- chan CastEntry
	//? Refrence: https://thomasnguyen.site/function-returning-channel-pattern-in-go
	castEntryChan, err := dns.DiscoverCastDNSEntries(ctx, iface)
	if err != nil {
		log.Debugf("unable to discover chromecast devices: %v", err)
	}
	// TODO: get MAC address
	ii := 1
	//? See go-chromecast.cmd.ls for use
	for d := range castEntryChan {
		fmt.Printf("%d) device=%q device_name=%q address=\"%s:%d\" uuid=%q",
			ii, d.Device, d.DeviceName, d.AddrV4, d.Port, d.UUID)
		ii++
	}
	if ii == 1 {
		log.Error("no cast devices found on network")
	}
}

// DiscoverCastDNSEntries will return a channel with any cast dns entries
// found.
func DiscoverCastDNSEntries(ctx context.Context, iface *net.Interface) (<-chan CastDevice, error) {
	// var opts = []zeroconf.ClientOption{zeroconf.SelectIPTraffic(zeroconf.IPv4)}
	// Act as a client using a Network Interface
	// if iface != nil {
	// 	opts = append(opts, zeroconf.SelectIfaces([]net.Interface{*iface}))
	// }
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create new zeroconf resolver: %w", err)
	}
	castDNSEntriesChan := make(chan CastDevice, 5)
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
					// //? Added by Jason
					// HwAddr: net.HardwareAddr(),
					CastEntry: dns.CastEntry{
						Port: entry.Port,
						Host: entry.HostName,
					},
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
