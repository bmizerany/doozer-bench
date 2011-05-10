package main

import (
	"fmt"
	"flag"
	"github.com/ha/doozer"
	"math"
	"os"
	"testing"
)

var uri = os.Getenv("DOOZER_URI")

// Flags
var (
	n = flag.Int("n", 5, "the number of clients to attach")
)

func main() {
	flag.StringVar(&uri, "a", uri, "the doozer cluster to bind to")
	flag.Parse()

	r := testing.Benchmark(Benchmark5DoozerConClientSet)
	fmt.Printf("%d\top/sec\n", 1e9/r.NsPerOp())
	fmt.Printf("%d\tns/op\n", r.NsPerOp())
}

func Benchmark5DoozerConClientSet(b *testing.B) {
	b.StopTimer()

	var cls = make([]*doozer.Conn, *n)
	for i := 0; i < *n; i++ {
		cls[i] = dial()
	}

	c := make(chan bool, b.N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		i := i
		go func() {
			cls[i%len(cls)].Set("/test", math.MaxInt64, nil)
			c <- true
		}()
	}
	for i := 0; i < b.N; i++ {
		<-c
	}
}

var usage1 = `
Environment:

   DOOZER_URI - the doozer cluster to bind to

`

func usage() {
	fmt.Fprintln(os.Stderr, "Use: doozer-bench [options]")
	fmt.Fprintln(os.Stderr, "")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, usage1)
	os.Exit(1)
}

func dial() *doozer.Conn {
	c, err := doozer.DialUri(uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "! %s\n", err.String())
		fmt.Fprintln(os.Stderr, "")
		usage()
	}
	return c
}
