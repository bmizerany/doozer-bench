package main

import (
	"flag"
	"fmt"
	"github.com/ha/doozer"
	"os"
	"time"
	"math"
)

var (
	c = flag.Int("c", 5, "The number of clients to run")
	n = flag.Int("n", 500, "The number of total operations to run")
	a = flag.String("a", "", "The doozer cluster to bind to")
)

type cb func(id, iter int, start, end int64, err os.Error)

func main() {
	flag.Parse()

	if *a == "" {
		*a = os.Getenv("DOOZER_URI")
	}

	opc := *n / *c
	done := make(chan bool)

	f := func(id, iter int, start, end int64, err os.Error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %s\n", err.String())
			fmt.Println(id, -1)
		} else {
			fmt.Println(time.Nanoseconds(), id, iter, end - start)
		}

		if iter == (opc - 1) {
			done <- true
		}
	}

	for i := 0; i < *c; i++ {
		cl, err := doozer.DialUri(*a)
		if err != nil {
			panic(err)
		}
		go set(cl, i, opc, "/bench/set", []byte("foo"), f)
	}

	for i := 0; i < *c; i++ {
		<-done
	}
}

func set(cl *doozer.Conn, id, iter int, path string, value []byte, f cb) {
	for i := 0; i < iter; i++ {
		s := time.Nanoseconds()
		_, err := cl.Set(path, math.MaxInt64, value)
		e := time.Nanoseconds()
		f(id, i, s, e, err)
	}
}
