package scan

//Based on https://gist.github.com/kotakanbe/d3059af990252ba89a82
import (
	"fmt"
	"net"
	"os"
	"sync/atomic"

	"time"

	fastping "github.com/tatsushid/go-fastping"
)

// Ping takes a slice of IP addresses and return an int count of those that
// respond to ping.  The MaxRTT is set to 4 seconds.
func Ping(hosts []string) int {
	p := fastping.NewPinger()
	p.MaxRTT = 4 * time.Second
	var successCount, failCount uint64
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		atomic.AddUint64(&successCount, 1)
		// fmt.Printf("IP Addr: %s receive, RTT: %v  successCount: %v \n", addr.String(), rtt, successCount)
	}
	p.OnIdle = func() {
		atomic.AddUint64(&failCount, 1)
		// fmt.Println("timed out - finish")
	}

	for _, ip := range hosts {
		// fmt.Printf("adding ip: %v \n", ip)
		err := p.AddIP(ip)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error adding IP (%v): %v", ip, err)
		}
	}

	err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during Ping.Run(): %v", err)
	}

	return int(successCount)
}
