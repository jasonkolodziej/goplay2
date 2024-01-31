package cast

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
	log "github.com/sirupsen/logrus"
	"github.com/vishen/go-chromecast/dns"
)

/*	CastDevice extends dns.CastEntry
 *	by gathering the MAC Address as well the associated
 *	net.Interface for the Chromecast Device
 */
type CastDevice struct {
	svcEntry *zeroconf.ServiceEntry
	HwAddr   net.HardwareAddr
	dns.CastEntry
}

type Port uint

const (
	CHROMECAST       Port = 8009
	CHROMECAST_GROUP Port = 32187
)

/*	CastDeviceEntry extends dns.CastDNSEntry
 *	by providing the MAC Address as well the associated
 *	net.Interface for the Chromecast Device
 */
type DeviceEntry interface {
	dns.CastDNSEntry
	// NetInterface() *net.Interface
	HardwareAddr() net.HardwareAddr
	ServiceEntry() *zeroconf.ServiceEntry
	FromEntry(svc *zeroconf.ServiceEntry) *CastDevice
	svcEntryTemplate() *zeroconf.ServiceEntry
	// Entry() *CastDevice
}

func (d CastDevice) ServiceEntry() *zeroconf.ServiceEntry {
	return nil
}

func (d CastDevice) svcEntryTemplate() *zeroconf.ServiceEntry {
	d.svcEntry = zeroconf.NewServiceEntry(
		"",                 //! Instance will need to be changed later
		"_googlecast._tcp", // Service type
		"local")            // Domain name
	d.svcEntry.Port = 8009
	return d.svcEntry
}

func (d CastDevice) FromEntry(svcEntry *zeroconf.ServiceEntry) *CastDevice {
	return &CastDevice{
		CastEntry: dns.CastEntry{
			Port: svcEntry.Port,
			Name: svcEntry.HostName,
			// TODO: should we check?
			AddrV4:     svcEntry.AddrIPv4[0],
			AddrV6:     svcEntry.AddrIPv6[0],
			InfoFields: TxtRecordHelper(svcEntry.Text),
			DeviceName: decode(TxtRecordHelper(svcEntry.Text)["fn"]),
			Device:     decode(TxtRecordHelper(svcEntry.Text)["md"]),
			UUID:       TxtRecordHelper(svcEntry.Text)["id"],
		},
		svcEntry: svcEntry,
	}
}

func TxtRecordHelper(r []string) (info map[string]string) {
	info = map[string]string{}
	for _, value := range r {
		if kv := strings.SplitN(value, "=", 2); len(kv) == 2 {
			key := kv[0]
			val := kv[1]
			info[key] = val
			// switch key {
			// case "fn":
			// 	castEntry.DeviceName = decode(val)
			// case "md":
			// 	castEntry.Device = decode(val)
			// case "id":
			// 	castEntry.UUID = val
			// }
		}
	}
	return
}

// func (d CastDevice) Entry() *dns.CastEntry {
// 	return &d.CastEntry
// }

func (d CastDevice) HardwareAddr() net.HardwareAddr {
	return d.HwAddr
}

func GetCastDevices(ifaceName string, timeoutSec uint) {
	iface, _ := net.InterfaceByName(ifaceName)
	var opts = []zeroconf.ClientOption{}
	// // Act as a client using a Network Interface
	if iface != nil {
		opts = append(opts, zeroconf.SelectIfaces([]net.Interface{*iface}))
	}
	resolver, err := zeroconf.NewResolver(opts...)
	if err != nil {
		log.Errorf("unable to create new zeroconf resolver: %s", err)
	}
	// //? Refrence: https://thomasnguyen.site/function-returning-channel-pattern-in-go
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			fmt.Printf("%v\n", entry)
		}
		//log.Log("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSec))
	defer cancel()
	err = resolver.Browse(ctx, "_googlecast._tcp", "local", entries)
	if err != nil {
		panic(fmt.Errorf("Failed to browse: %s", err.Error()))
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	// time.Sleep(1 * time.Second)
	// var iface *net.Interface
	// var err error
	// if ifaceName != "" {
	// 	if iface, err = net.InterfaceByName(ifaceName); err != nil {
	// 		log.Debugf("unable to find interface %q: %v", ifaceName, err)
	// 	}
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSec))
	// defer cancel()
	// // Returns a <- chan CastEntry
	// //? Refrence: https://thomasnguyen.site/function-returning-channel-pattern-in-go
	// castEntryChan, err := dns.DiscoverCastDNSEntries(ctx, iface)
	// if err != nil {
	// 	log.Debugf("unable to discover chromecast devices: %v", err)
	// }
	// // TODO: get MAC address
	// ii := 1
	// //? See go-chromecast.cmd.ls for use
	// for d := range castEntryChan {
	// 	fmt.Printf("%d) device=%q device_name=%q address=\"%s:%d\" uuid=%q",
	// 		ii, d.Device, d.DeviceName, d.AddrV4, d.Port, d.UUID)
	// 	ii++
	// }
	// if ii == 1 {
	// 	log.Error("no cast devices found on network")
	// }
}
