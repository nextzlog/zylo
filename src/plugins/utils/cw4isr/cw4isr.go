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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gen2brain/malgo"
	"github.com/thoas/go-funk"
	"github.com/wcharczuk/go-chart/v2"
	"image/color"
	"net/url"
	"unsafe"
	"zylo/morse"
)

const (
	INTERVAL_MS = 200
	MAX_HISTORY = 100
	VOL_MAX_VAL = 100
)

const (
	NAME = "CW4ISR Morse Decoder"
	HREF = "https://use.zlog.org"
	LINK = "Download Latest zLog"
)

const GRAPH_WIDTH = 500

var (
	ctx *malgo.AllocatedContext
	dev *Device
)

var (
	items []morse.Message
	table map[string]*Device
	names []string
	graph [][]float64
	level float64
)

var (
	his *widget.List
	lab *widget.Label
	osc *canvas.Image
	spa *canvas.Raster
)

var dcb = malgo.DeviceCallbacks{
	Data: onSignalEvent,
}

type Device struct {
	pointer unsafe.Pointer
	capture *malgo.Device
	decoder morse.Decoder
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
	dev.decoder = morse.DefaultDecoder(int(dev.capture.SampleRate()))
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
	messages := dev.decoder.Read(readSignedInt(in))
	for _, m := range messages {
		miss := true
		for n, p := range items {
			freq := m.Freq == p.Freq
			time := m.Time == p.Time
			if freq && time {
				items[n] = m
				miss = false
			}
		}
		if miss {
			items = append(items, m)
		}
	}
	if len(items) > MAX_HISTORY {
		items = items[len(items)-MAX_HISTORY:]
	}
	graph = append(graph, dev.decoder.Spec...)
	if len(graph) > GRAPH_WIDTH {
		graph = graph[len(graph)-GRAPH_WIDTH:]
	}
	level = 0.0
	for _, row := range graph {
		for _, v := range row {
			if v > level {
				level = v
			}
		}
	}
	his.Refresh()
	spa.Refresh()
}

func readSignedInt(signal []byte) (result []float64) {
	for _, b := range funk.Chunk(signal, 4).([][]byte) {
		v := binary.LittleEndian.Uint32(b)
		result = append(result, float64(int32(v)))
	}
	return
}

func updateSpec(x, y, w, h int) (pixel color.Color) {
	x = x * GRAPH_WIDTH / w
	y = (h - y) * 100 / h
	value := 0.0
	width := GRAPH_WIDTH - len(graph)
	if x > width {
		value = graph[x-width][y] / level
	}
	return color.RGBA{
		R: uint8(255 * value),
		G: uint8(255 * value),
		B: 0,
		A: 255,
	}
}

func length() int {
	return len(items)
}

func create() fyne.CanvasObject {
	return widget.NewLabel("")
}

func update(id widget.ListItemID, obj fyne.CanvasObject) {
	obj.(*widget.Label).SetText(morse.CodeToText(items[id].Code))
}

func choice(id widget.ListItemID) {
	x := make([]float64, len(items[id].Data))
	for n, _ := range x {
		x[n] = float64(n)
	}
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: x,
				YValues: items[id].Data,
				Style: chart.Style{
					FillColor: chart.GetDefaultColor(0),
				},
			},
		},
		XAxis: chart.XAxis{
			Style: chart.Style{
				Hidden: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Hidden: true,
			},
		},
		Width:  int(osc.Size().Width),
		Height: int(osc.Size().Height),
	}
	buffer := &chart.ImageWriter{}
	graph.Render(chart.PNG, buffer)
	image, _ := buffer.Image()
	osc.Image = image
	osc.Refresh()
	lab.SetText(morse.CodeToText(items[id].Code))
}

func restart(name string) {
	if dev != nil {
		dev.Stop()
	}
	dev = table[name]
	dev.Listen()
}

func clear() {
	items = nil
	his.Refresh()
}

func volume(vol float64) {
	dev.decoder.Mute = vol / VOL_MAX_VAL
}

func main() {
	ctx, _ = malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	app := app.New()
	app.Settings().SetTheme(theme.DarkTheme())
	win := app.NewWindow(NAME)
	ref, _ := url.Parse(HREF)
	btm := widget.NewHyperlink(LINK, ref)
	osc = &canvas.Image{}
	lab = widget.NewLabel("")
	spa = canvas.NewRasterWithPixels(updateSpec)
	his = widget.NewList(length, create, update)
	his.OnSelected = choice
	table, names = DeviceList()
	sel := widget.NewSelect(names, restart)
	vol := widget.NewSlider(0, VOL_MAX_VAL)
	btn := widget.NewButton("clear", clear)
	vol.OnChanged = volume
	sel.SetSelectedIndex(0)
	vol.SetValue(0.3 * VOL_MAX_VAL)
	lhs := container.NewBorder(nil, btn, nil, nil, his)
	rhs := container.NewBorder(lab, nil, nil, nil, osc)
	hsp := container.NewHSplit(lhs, rhs)
	vsp := container.NewVSplit(hsp, spa)
	bar := container.NewBorder(nil, nil, sel, nil, vol)
	out := container.NewBorder(bar, btm, nil, nil, vsp)
	hsp.SetOffset(0.2)
	win.SetContent(out)
	win.Resize(fyne.NewSize(640, 480))
	win.ShowAndRun()
}
