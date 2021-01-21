package internal

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Statistics struct {
	StartTime        time.Time
	RsChan           chan *RequestStatus
	Concurrency      int
	RequestDone      uint64
	ResponseSize     uint64
	StatusAll        map[int]uint64
	LastMinuteStatus []*RequestStatus
	rw               sync.RWMutex
	workerWg         sync.WaitGroup
}

func NewStatistics(concurrency int) *Statistics {
	s := &Statistics{
		StartTime:    time.Now(),
		StatusAll:    make(map[int]uint64),
		RequestDone:  0,
		ResponseSize: 0,
		Concurrency:  concurrency,
		RsChan:       make(chan *RequestStatus, 2*concurrency),
	}
	s.workerWg.Add(1)
	go s.addStatusWorker()

	return s
}

func (s *Statistics) Stop() {
	close(s.RsChan)
	s.workerWg.Wait()
}

func (s *Statistics) addStatusWorker() {
	deal := func(rs *RequestStatus) {
		s.rw.Lock()
		defer s.rw.Unlock()
		s.RequestDone++
		statusCode := rs.StatusCode

		if _, has := s.StatusAll[statusCode]; !has {
			s.StatusAll[statusCode] = 0
		}
		s.StatusAll[statusCode]++

		s.LastMinuteStatus = append(s.LastMinuteStatus, rs)
		var expirePos int
		for i, v := range s.LastMinuteStatus {
			expirePos = i
			if !v.IsExpire() {
				break
			}
		}
		if expirePos > 0 {
			s.LastMinuteStatus = s.LastMinuteStatus[expirePos:]
		}
	}

	for rs := range s.RsChan {
		deal(rs)
	}
	s.workerWg.Done()
}
func (s *Statistics) AddStatus(rs *RequestStatus) {
	s.RsChan <- rs
}

func (s *Statistics) AddResponseSize(size int) {
	atomic.AddUint64(&s.ResponseSize, uint64(size))
}

func (s *Statistics) NowUsed() time.Duration {
	return time.Now().Sub(s.StartTime)
}

func (s *Statistics) TotalQps() float64 {
	return float64(s.RequestDone) / s.NowUsed().Seconds()
}

func (s *Statistics) MinuteQps() float64 {
	last := s.MinuteStatus()
	l := len(last)
	if l < 1 {
		return 0
	}
	first := last[0]
	// 	lastOne := last[l-1]
	// 	_used := lastOne.Time.Sub(first.Time).Seconds()
	_used := time.Now().Sub(first.StartTime).Seconds()
	if _used == 0 {
		return 1
	}
	return float64(l) / _used
}

func (s *Statistics) MinuteStatus() []*RequestStatus {
	s.rw.RLock()
	defer s.rw.RUnlock()
	var last []*RequestStatus
	for _, rs := range s.LastMinuteStatus {
		if !rs.IsExpire() {
			last = append(last, rs)
		}
	}
	return last
}

func (s *Statistics) MinuteStatusCode() map[int]uint64 {
	last := s.MinuteStatus()
	m := make(map[int]uint64)
	for _, rs := range last {
		if _, has := m[rs.StatusCode]; !has {
			m[rs.StatusCode] = 0
		}
		m[rs.StatusCode]++
	}
	return m
}

func (s *Statistics) StatusTxt() string {
	var msg []string
	s.rw.RLock()
	msg = append(msg, fmt.Sprintf("requests=%d", s.RequestDone))
	bs_all, _ := json.Marshal(s.StatusAll)
	s.rw.RUnlock()

	msg = append(msg, fmt.Sprintf("qps_avg=%.1f", s.TotalQps()))
	msg = append(msg, fmt.Sprintf("qps_minute=%.1f", s.MinuteQps()))

	msg = append(msg, fmt.Sprintf("status_all=%s", strings.Replace(string(bs_all), `"`, "", -1)))

	bs_minute, _ := json.Marshal(s.MinuteStatusCode())
	msg = append(msg, fmt.Sprintf("status_minute=%s", strings.Replace(string(bs_minute), `"`, "", -1)))
	// 	msg=append(msg,fmt.Sprintf("resp_size=%d", s.ResponseSize))

	return strings.Join(msg, ", ")
}

type RequestStatus struct {
	StartTime  time.Time
	EndTime    time.Time
	StatusCode int
}

func NewRequestStatus() *RequestStatus {
	return &RequestStatus{
		StartTime: time.Now(),
	}
}
func (rs *RequestStatus) Status(statusCode int) {
	rs.StatusCode = statusCode
	rs.EndTime = time.Now()
}

func (rs *RequestStatus) IsExpire() bool {
	return time.Now().Sub(rs.StartTime).Seconds() > 60.0
}
