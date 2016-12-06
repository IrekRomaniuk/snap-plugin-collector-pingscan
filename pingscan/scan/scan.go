package scan
//Based on https://gist.github.com/kotakanbe/d3059af990252ba89a82
import (
	"os/exec"
	"github.com/IrekRomaniuk/snap-plugin-collector-pingscan/pingscan/targets"
	"fmt"
)

func Ping(hosts []string) int {
	concurrentMax := 200
	pingChan := make(chan string, concurrentMax)
	pongChan := make(chan string, len(hosts))
	doneChan := make(chan []string)
	//fmt.Printf("concurrentMax=%d hosts=%d -> %s...%s\n", concurrentMax, len(hosts), hosts[0], hosts[len(hosts) - 1])
	for i := 0; i < concurrentMax; i++ {
		go sendingPing(pingChan, pongChan)
	}

	go receivePong(len(hosts), pongChan, doneChan)

	for _, ip := range hosts {
		pingChan <- ip
	//fmt.Println("sent: ", ip)
	}
	alives := <-doneChan
	result := targets.DeleteEmpty(alives)

	fmt.Printf("\n%d/%d %d\n", len(result),len(hosts),concurrentMax)
	return len(result)
}

func sendingPing(pingChan <-chan string, pongChan chan <- string) {
	for ip := range pingChan {
		_, err := exec.Command("ping", "-c", "1", "-w", "1", ip).Output()
		if err == nil {
			pongChan <- ip
			//fmt.Printf("%s is alive\n", ip)
		} else {
			pongChan <- ""
			//fmt.Printf("%s is dead\n", ip)
		}
	}
}

func receivePong(pongNum int, pongChan <-chan string, doneChan chan <- []string) {
	var alives []string
	for i := 0; i < pongNum; i++ {
		ip := <-pongChan
		//fmt.Println("received: ", ip)
		alives = append(alives, ip)
	}
	doneChan <- alives
}