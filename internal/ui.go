package internal

import (
	"fmt"
	"github.com/gizak/termui"
)

type UI struct {
	statistics *Statistics
	totalQps   []float64
	minuteQps  []float64
}

func NewUI(s *Statistics) (*UI, error) {
	err := termui.Init()
	if err != nil {
		return nil, err
	}
	return &UI{
		statistics: s,
		totalQps:   []float64{},
		minuteQps:  []float64{},
	}, nil
}

func (ui *UI) Init() {

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		ui.Close()
	})

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		ui.Close()
	})

	lc1 := termui.NewLineChart()
	lc1.BorderLabel = "total qps"
	lc1.Mode = "dot"
	lc1.Data = ui.totalQps
	lc1.Width = 80
	lc1.Height = 12
	lc1.X = 0
	lc1.DotStyle = '+'
	lc1.AxesColor = termui.ColorWhite
	lc1.LineColor = termui.ColorYellow | termui.AttrBold

	lc2 := termui.NewLineChart()
	lc2.BorderLabel = "minute qps"
	lc2.Mode = "dot"
	lc2.Data = ui.minuteQps
	lc2.Width = 80
	lc2.Height = 12
	lc2.X = 0
	lc2.Y = 13
	lc2.DotStyle = '*'
	lc2.AxesColor = termui.ColorWhite
	lc2.LineColor = termui.ColorYellow | termui.AttrBold

	termui.Render(lc2, lc1)

	termui.Handle("/timer/1s", func(e termui.Event) {
		lc1.Data = ui.getTotalQps()
		lc2.Data = ui.getMinuteQps()
		termui.Render(lc2, lc1)
	})

	go termui.Loop()
}

func (ui *UI) getTotalQps() []float64 {
	q := ui.statistics.TotalQps()
	ui.totalQps = append(ui.totalQps, q)
	if len(ui.totalQps) > 100 {
		ui.totalQps = ui.totalQps[1:]
	}
	return ui.totalQps
}
func (ui *UI) getMinuteQps() []float64 {
	q := ui.statistics.MinuteQps()
	ui.minuteQps = append(ui.minuteQps, q)
	if len(ui.minuteQps) > 100 {
		ui.minuteQps = ui.minuteQps[1:]
	}
	return ui.minuteQps
}

func (ui *UI) Close() {
	termui.StopLoop()
	termui.Close()
	fmt.Println("ui close")
}

func (ui *UI) Refresh() {
}
