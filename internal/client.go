package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/hidu/goutils/time_util"
)

type Client struct {
	reqsChan    chan *Request
	statistics  *Statistics
	timeout     time.Duration
	concurrency int
	ui          *UI

	logErrTime time.Time
	logErrRw   sync.RWMutex
	workerWg   sync.WaitGroup

	needDetail bool
}

func NewClient(concurrency int, detail bool) *Client {
	cs := &Client{
		reqsChan:    make(chan *Request, concurrency*10),
		concurrency: concurrency,
		statistics:  NewStatistics(concurrency),
		logErrTime:  time.Now(),
		needDetail:  detail,
	}
	return cs
}

func (c *Client) AddRequest(req *Request) {
	c.reqsChan <- req
}

func (c *Client) SetTimeout(sec int) {
	c.timeout = time.Duration(sec) * time.Second
}

func (c *Client) NextRequest() *Request {
	return <-c.reqsChan
}

func (c *Client) Start() {
	for i := 0; i < c.concurrency; i++ {
		c.workerWg.Add(1)
		go c.worker(c.reqsChan, i)
	}
	time_util.SetInterval(func() {
		c.PrintStatistics()
	}, 5)

}
func (c *Client) UI() *UI {
	if c.ui == nil {
		c.ui, _ = NewUI(c.statistics)
	}
	return c.ui
}

func (c *Client) Wait() {
	close(c.reqsChan)
	c.workerWg.Wait()
	c.statistics.Stop()
	c.PrintStatistics()
}

func (c *Client) worker(jobs <-chan *Request, worker_id int) {
	for req := range jobs {
		client := &http.Client{
			Timeout: c.timeout,
		}
		s := NewRequestStatus()
		hReq, err := req.AsHttpRequest()

		if err != nil {
			s.Status(0)
			c.statistics.AddStatus(s)
			c.logError(0, req.Raw(), "")
			continue
		}

		resp, err := client.Do(hReq)

		if err != nil {
			s.Status(0)
			c.statistics.AddStatus(s)
			c.logError(0, req.Raw(), "")
			continue
		}

		bd, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		s.Status(resp.StatusCode)
		c.statistics.AddStatus(s)

		c.statistics.AddResponseSize(len(bd))

		if resp.StatusCode != http.StatusOK {
			c.logError(resp.StatusCode, req.Raw(), string(bd))
		}

		if c.needDetail {
			log.Println("request_detail", req.Raw(), string(bd))
		}
	}
	c.workerWg.Done()
}

func (c *Client) logError(statusCode int, req string, resp string) {
	now := time.Now()
	if now.Sub(c.logErrTime).Seconds() < 5 {
		return
	}
	log.Println("faild_request_sample,status=", statusCode, "request=", req, "response=", resp)
	c.logErrRw.Lock()
	defer c.logErrRw.Unlock()
	c.logErrTime = now
}

func (c *Client) PrintStatistics() {
	msg := fmt.Sprintf("conc=%d,%s", c.concurrency, c.statistics.StatusTxt())
	log.Println(msg)
}
