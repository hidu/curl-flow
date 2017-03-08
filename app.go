package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/hidu/curl-flow/internal"
	"io"
	"log"
	"os"
)

var conc = flag.Int("c", 1, "concurrency Number of multiple requests to make")
var url = flag.String("url", "", "test url")
var n = flag.Int("n", 0, "Number of requests to perform")
var t = flag.Uint("t", 10, "Timeout of request")
var useUi = flag.Bool("ui", false, "use termui")

func main() {
	flag.Parse()
	fmt.Println("start")

	client := internal.NewClient(*conc)
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
		r, _ := req.AsHttpRequest()
		for i := 0; i < *n; i++ {
			client.AddRequest(r)
		}
	} else {
		buf := bufio.NewReaderSize(os.Stdin, 8192)
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
			r, _ := req.AsHttpRequest()
			client.AddRequest(r)
		}
	}

	client.Wait()
}
