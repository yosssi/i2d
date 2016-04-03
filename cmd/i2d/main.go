package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	inFilePath  string
	outFilePath string
	threadNum   int
	maxRetryNum int

	ips       []string
	ipDomains = map[string]string{}
)

type ipDomain struct {
	ip     string
	domain string
}

func lookupAddr(addr string, retryNum int) string {
	names, err := net.LookupAddr(addr)
	if (err != nil || len(names) == 0 || names[0] == addr) && retryNum < maxRetryNum {
		return lookupAddr(addr, retryNum+1)
	}
	if len(names) == 0 {
		return addr
	}
	return strings.TrimSuffix(names[0], ".")
}

func parse() {
	ifp := flag.String("in", "", "input file path")
	ofp := flag.String("out", "", "output file path")
	tn := flag.Int("n", 1, "number of threads")
	mrn := flag.Int("r", 0, "max number of retries")

	flag.Parse()

	inFilePath = *ifp
	outFilePath = *ofp
	threadNum = *tn
	maxRetryNum = *mrn
}

func read() {
	data, err := ioutil.ReadFile(inFilePath)
	if err != nil {
		panic(err)
	}
	rows := strings.Split(strings.TrimSpace(string(data)), "\n")
	ips = make([]string, len(rows), len(rows))
	for i, ip := range rows {
		ips[i] = ip
	}
}

func lookup(in <-chan string, out chan<- ipDomain, wg *sync.WaitGroup) {
	defer wg.Done()
	for ip := range in {
		out <- ipDomain{ip, lookupAddr(ip, 0)}
	}
}

func write() {
	file, err := os.Create(outFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for _, ip := range ips {
		if _, err := writer.WriteString(ipDomains[ip] + "\n"); err != nil {
			panic(err)
		}
	}
	writer.Flush()
}

func main() {
	parse()
	read()

	in := make(chan string, 1000000)
	out := make(chan ipDomain, 1000000)
	wg := new(sync.WaitGroup)
	for i := 0; i < threadNum; i++ {
		wg.Add(1)
		go lookup(in, out, wg)
	}

	end := make(chan struct{})

	go func() {
		defer func() {
			end <- struct{}{}
		}()
		for id := range out {
			ipDomains[id.ip] = id.domain
		}
	}()
	for _, ip := range ips {
		in <- ip
	}

	close(in)

	wg.Wait()

	close(out)

	<-end

	write()
}
