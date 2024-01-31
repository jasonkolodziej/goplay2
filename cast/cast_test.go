package cast

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"
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

// ! Good func
func Test_BrowseVsLookup(t *testing.T) {
	ifaceName := "eth0"
	iface, _ := net.InterfaceByName(ifaceName)
	var opts = []zeroconf.ClientOption{}
	// // Act as a client using a Network Interface
	if iface != nil {
		opts = append(opts, zeroconf.SelectIfaces([]net.Interface{*iface}))
	}
	resolver, err := zeroconf.NewResolver(opts...)
	if err != nil {
		t.Errorf("unable to create new zeroconf resolver: %s", err)
	}

	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			mapB, err := json.Marshal(*entry)
			if err != nil {
				t.Errorf("unable to serialize: %s", err)
			}
			t.Logf("Content: %s", string(mapB))
		}
		t.Log("No more entries.")
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	err = resolver.Browse(ctx, "_googlecast._tcp", "local", entries)
	if err != nil {
		t.Fatalf("Failed to browse: %s", err.Error())
	}

	<-ctx.Done()
	// Wait some additional time to see debug messages on go routine shutdown.
	// time.Sleep(1 * time.Second)
}
