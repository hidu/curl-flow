package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hidu/curl-flow/internal"
)

const version = "20210121"

var conc = flag.Int("c", 1, "concurrency Number of multiple requests to make")
var url = flag.String("url", "", "test url,When no flow is used")
var n = flag.Int("n", 1, "Number of each requests to perform")
var t = flag.Uint("t", 10, "Timeout of request (second)")
var useUi = flag.Bool("ui", false, "use termui")
var detail = flag.Bool("detail", false, "print request detail to log")

func init() {
	df := flag.Usage
	flag.Usage = func() {
		df()
		fmt.Println("\n    site : https://github.com/hidu/curl-flow")
		fmt.Println(" version :", version)
	}
}

func main() {
	flag.Parse()
	log.Println("start")

	client := internal.NewClient(*conc, *detail)
	client.SetTimeout(int(*t))
	client.Start()

	if *useUi {
		ui := client.UI()
		if ui != nil {
			ui.Init()
			defer ui.Close()
		}
	}

	if *url != "" {
		req := internal.NewRequest(*url, "GET")
		for i := 0; i < *n; i++ {
			client.AddRequest(req)
		}
	} else {
		buf := bufio.NewReaderSize(os.Stdin, 8192)
		reqId := 0
		for {
			line, err := buf.ReadBytes('\n')
			if err == io.EOF {
				break
			}

			req, jerr := internal.NewRequestJson(line)
			if jerr != nil {
				log.Println("parse request failed:", jerr, ",input:", string(line))
				continue
			}

			for i := 0; i < *n; i++ {
				client.AddRequest(req)
				reqId++
			}

			// 			if(*n > 1 && reqId >= *n){
			// 				break
			// 			}
		}
	}

	client.Wait()
	log.Println("stopped")
}
