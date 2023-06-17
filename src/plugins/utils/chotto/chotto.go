/*
Copyright (C) 2022 JA1ZLO.
*/
package main

import (
	"encoding/binary"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/gen2brain/malgo"
	"github.com/thoas/go-funk"
	"github.com/wcharczuk/go-chart/v2"
	"math"
	"net/url"
	"unsafe"
	"zylo/morse"
)

const (
	INTERVAL_MS = 200
	MAX_HISTORY = 100
)

const (
	NAME = "Chotto Wakaru CW"
	HREF = "https://use.zlog.org"
	LINK = "Download Latest zLog"
)

var (
	ctx *malgo.AllocatedContext
	dev *Device
)

var (
	table map[string]*Device
	names []string
)

var (
	data []morse.Message
	his  *widget.List
)

var (
	lab *widget.Label
	osc *canvas.Image
)

var dcb = malgo.DeviceCallbacks{
	Data: onSignalEvent,
}

type Device struct {
	pointer unsafe.Pointer
	capture *malgo.Device
	monitor morse.Monitor
}

func (dev *Device) Config() (cfg malgo.DeviceConfig) {
	cfg = malgo.DefaultDeviceConfig(malgo.Capture)
	cfg.PeriodSizeInMilliseconds = INTERVAL_MS
	cfg.Capture.Format = malgo.FormatS32
	cfg.Capture.Channels = 1
	cfg.Capture.DeviceID = dev.pointer
	return
}

func (dev *Device) Listen() {
	dev.capture, _ = malgo.InitDevice(ctx.Context, dev.Config(), dcb)
	dev.monitor = morse.DefaultMonitor(int(dev.capture.SampleRate()))
	dev.capture.Start()
}

func (dev *Device) Stop() {
	dev.capture.Uninit()
}

func DeviceList() (table map[string]*Device, names []string) {
	if ctx != nil {
		table = make(map[string]*Device)
		devs, _ := ctx.Devices(malgo.Capture)
		for _, dev := range devs {
			names = append(names, dev.Name())
			table[dev.Name()] = &Device{
				pointer: dev.ID.Pointer(),
			}
		}
	}
	return
}

func onSignalEvent(out, in []byte, frames uint32) {
	messages := dev.monitor.Read(readSignedInt(in))
	for _, m := range messages {
		miss := true
		for n, p := range data {
			freq := m.Freq == p.Freq
			life := m.Life >= p.Life
			if freq && life {
				data[n] = m
				miss = false
			} else if freq {
				data[n].Life = math.MaxInt32
			}
		}
		if miss {
			data = append(data, m)
		}
	}
	if len(data) > MAX_HISTORY {
		data = data[len(data)-MAX_HISTORY:]
	}
	his.Refresh()
}

func readSignedInt(signal []byte) (result []float64) {
	for _, b := range funk.Chunk(signal, 4).([][]byte) {
		v := binary.LittleEndian.Uint32(b)
		result = append(result, float64(int32(v)))
	}
	return
}

func length() int {
	return len(data)
}

func create() fyne.CanvasObject {
	return widget.NewLabel("")
}

func update(id widget.ListItemID, obj fyne.CanvasObject) {
	obj.(*widget.Label).SetText(morse.CodeToText(data[id].Code))
}

func choice(id widget.ListItemID) {
	x := make([]float64, len(data[id].Data))
	for n, _ := range x {
		x[n] = float64(n)
	}
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: data[id].Data,
			},
		},
	}
	buffer := &chart.ImageWriter{}
	graph.Render(chart.PNG, buffer)
	image, _ := buffer.Image()
	osc.Image = image
	osc.Refresh()
	lab.SetText(morse.CodeToText(data[id].Code))
}

func restart(name string) {
	if dev != nil {
		dev.Stop()
	}
	dev = table[name]
	dev.Listen()
}

func clear() {
	data = nil
	his.Refresh()
}

func main() {
	ctx, _ = malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	app := app.New()
	win := app.NewWindow(NAME)
	ref, _ := url.Parse(HREF)
	btm := widget.NewHyperlink(LINK, ref)
	osc = &canvas.Image{}
	lab = widget.NewLabel("")
	his = widget.NewList(length, create, update)
	his.OnSelected = choice
	table, names = DeviceList()
	sel := widget.NewSelect(names, restart)
	btn := widget.NewButton("clear", clear)
	sel.SetSelectedIndex(0)
	lhs := container.NewBorder(nil, btn, nil, nil, his)
	rhs := container.NewBorder(lab, nil, nil, nil, osc)
	hsp := container.NewHSplit(lhs, rhs)
	out := container.NewBorder(sel, btm, nil, nil, hsp)
	win.SetContent(out)
	win.Resize(fyne.NewSize(640, 480))
	win.ShowAndRun()
}
