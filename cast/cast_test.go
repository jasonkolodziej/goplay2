package cast

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/mdns"
	"github.com/vishen/go-chromecast/dns"
)

func Test_discoverLocalNetInterfaces(t *testing.T) {
	interfaces := discoverLocalInterfaces()
	if len(interfaces) > 0 {
		t.Log(interfaces)
	} else {
		t.Fail()
	}
}

func Test_GetCastDevices(t *testing.T) {
	ifaceName := "eth0"
	// GetCastDevices(inter, 3)
	var iface *net.Interface
	var err error
	if ifaceName != "" {
		if iface, err = net.InterfaceByName(ifaceName); err != nil {
			t.Fatalf("unable to find interface %q: %v", ifaceName, err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(10))
	defer cancel()
	// Returns a <- chan CastEntry
	//? Refrence: https://thomasnguyen.site/function-returning-channel-pattern-in-go
	castEntryChan, err := dns.DiscoverCastDNSEntries(ctx, iface)
	if err != nil {
		t.Fatalf("unable to discover chromecast devices: %v", err)
	}
	// TODO: get MAC address
	ii := 1
	//? See go-chromecast.cmd.ls for use
	for d := range castEntryChan {
		t.Logf("%d) device=%q device_name=%q address=\"%s:%d\" uuid=%q",
			ii, d.Device, d.DeviceName, d.AddrV4, d.Port, d.UUID)
		ii++
	}
	if ii == 1 {
		t.Error("no cast devices found on network")
	}
}

func TestGeneralDiscovery(t *testing.T) {
	service := "_googlecast._tcp"
	// domain := "local"
	// waitTime := 10
	// Discover all services on the network (e.g. _workstation._tcp)
	// resolver, err := zeroconf.NewResolver(nil)
	// mdns.Query(mdns.DefaultParams(service))
	// q := mdns.DefaultParams(service)
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			t.Logf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	if e := mdns.Lookup(service, entriesCh); e != nil {
		t.Fatalf("Lookup failed with: %s", e)
	}
	close(entriesCh)
}
