package benchmark

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const requests int = 8 << 30

func send(wg *sync.WaitGroup, url string) {
	http.Get(url)
	wg.Done()
}

func main() {
	var addr *string = flag.String("-a", "localhost:9000", "load balancer address")
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(requests)
	fmt.Printf("Starting benchmark on %s\n", *addr)
	var start time.Time = time.Now()

	for i := 0; i < requests; i++ {
		go send(&wg, *addr)
	}

	wg.Wait()
	var end time.Time = time.Now()
	var elapsed time.Duration = end.Sub(start)

	fmt.Printf("\n\n======= BENCHMARK RESPONSE =======")
	fmt.Printf("Elapsed time: %s\n", elapsed.String())
	fmt.Printf("Number of requests: %d\n", requests)
	fmt.Printf("Requests per millisecond: %d\n", elapsed.Milliseconds() / int64(requests))
}