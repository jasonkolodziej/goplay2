package other

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/hashicorp/mdns"
)

var ifaceName = "Wi-Fi"

func Test_discoverLocalNetInterfaces(t *testing.T) {
	interfaces := discoverLocalInterfaces()
	if len(interfaces) > 0 {
		for _, in := range interfaces {
			m, err := in.MulticastAddrs()
			if err != nil {
				t.Fatal(err)
			}
			addr, err := in.Addrs()
			if err != nil {
				t.Fatal(err)
			}
			t.Logf("Interface: Name: %s, MulticastAddr: %s, HwAddr: %s, IPAddr: %v\n",
				in.Name, m, in.HardwareAddr.String(), addr)
		}
	} else {
		t.Fail()
	}
}

func Test_GetCastDevices(t *testing.T) {
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
	castEntryChan, err := DiscoverCastDNSEntries(ctx, iface)
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

func Test_BrowseVsLookup(t *testing.T) {
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
			t.Logf("%v", entry)
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

// other_test.go:80: &{{Google-Home-Mini-4829b1998821bef6d3d757c7daa3b6d0 _googlecast._tcp local _googlecast._tcp.local. Google-Home-Mini-4829b1998821bef6d3d757c7daa3b6d0._googlecast._tcp.local. _services._dns-sd._udp.local.} 4829b199-8821-bef6-d3d7-57c7daa3b6d0.local. 8009 [] 120 [] [fe80::3a8b:59ff:fe57:229]}
//     other_test.go:80: &{{Chromecast-3d1743a93d5b1fc2ea4952aa152823ba _googlecast._tcp local _googlecast._tcp.local. Chromecast-3d1743a93d5b1fc2ea4952aa152823ba._googlecast._tcp.local. _services._dns-sd._udp.local.} 3d1743a9-3d5b-1fc2-ea49-52aa152823ba.local. 8009 [] 120 [192.168.2.48] []}
//     other_test.go:80: &{{Chromecast-342e4af2383f515fb337d64d79f33747 _googlecast._tcp local _googlecast._tcp.local. Chromecast-342e4af2383f515fb337d64d79f33747._googlecast._tcp.local. _services._dns-sd._udp.local.} 342e4af2-383f-515f-b337-d64d79f33747.local. 8009 [id=342e4af2383f515fb337d64d79f33747 cd=638580A1E73D8B64483932F7CBDBF010 rm= ve=05 md=Chromecast ic=/setup/icon.png fn=Bedroom TV ca=465413 st=0 bs=FA8F4D7B4098 nf=1 rs=] 120 [192.168.2.112] []}
//     other_test.go:80: &{{Google-Home-a548ff5ad1fac1941101acb5a1204788 _googlecast._tcp local _googlecast._tcp.local. Google-Home-a548ff5ad1fac1941101acb5a1204788._googlecast._tcp.local. _services._dns-sd._udp.local.} a548ff5a-d1fa-c194-1101-acb5a1204788.local. 8009 [id=a548ff5ad1fac1941101acb5a1204788 cd=C7153E0902929806731E9FD35BB81C1E rm= ve=05 md=Google Home ic=/setup/icon.png fn=Kitchen speaker ca=199172 st=0 bs=FA8FCA8766F6 nf=1 rs=] 120 [192.168.2.152] []}
//     other_test.go:80: &{{Google-Cast-Group-AC7A2ADB0FCD4009BDAE07F386F136A4 _googlecast._tcp local _googlecast._tcp.local. Google-Cast-Group-AC7A2ADB0FCD4009BDAE07F386F136A4._googlecast._tcp.local. _services._dns-sd._udp.local.} a548ff5a-d1fa-c194-1101-acb5a1204788.local. 32187 [id=AC7A2ADB-0FCD-4009-BDAE-07F386F136A4 cd=AC7A2ADB-0FCD-4009-BDAE-07F386F136A4 rm= ve=05 md=Google Cast Group ic=/setup/icon.png fn=inside ca=199204 st=0 bs=FA8FCA8766F6 nf=1 rs=] 120 [192.168.2.152] []}
