package main

import (
	"encoding/line"
	"flag"
	"fmt"
	"os"
	"math"
	"bytes"
	"strconv"
	"sort"
)


var Delim = []byte{' '}

type Int64Array []int64

func (a Int64Array) Len() int {
	return len(a)
}

func (a Int64Array) Less(i, j int) bool {
	return a[i] < a[j]
}

func (a Int64Array) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}


func main() {
	if flag.NArg() < 1 {
		fmt.Println("usage: doozer-report FILE")
		os.Exit(1)
	}

	fn := flag.Arg(0)
	f, err := os.OpenFile(fn, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ln := 0
	ts := Int64Array(make([]int64, 0))
	r := line.NewReader(f, math.MaxInt8)
	for {
		ln += 1
		line, _, err := r.ReadLine()
		if err == os.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		fds := bytes.Split(line, Delim, -1)

		if len(fds) < 5 {
			fmt.Printf("invalid line (%d): %s\n", ln, line)
			continue
		}

		ns, err := strconv.Atoi64(string(fds[4]))
		if err != nil {
			// TODO: be more gracefull
			panic(err)
		}

		ts = append(ts, ns)
	}

	sort.Sort(ts)

	n := float64(len(ts))
	ftp := int((n * 50 / 100)+0.5)
	nnp := int((n * 99 / 100)+0.5)

	fmt.Printf("50%% %d.%03d sec\n", ts[ftp]/1e9, ts[ftp] % 1e9 / 1e6)
	fmt.Printf("99%% %d.%03d sec\n", ts[nnp]/1e9, ts[nnp] % 1e9 / 1e6)
}
