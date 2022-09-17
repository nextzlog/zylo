/*
 Copyright (C) 2022 JA1ZLO.
*/
package main

import (
	_ "embed"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/tadvi/winc"
	"gopkg.in/yaml.v2"
	"image/color"
	"os"
	"regexp"
	"strconv"
)

const (
	WINDOW_SIZE = 600
	MARKER_SIZE = 16
	MAPLOT_NAME = "maplot"
	MAPLOT_MENU = "MainForm.MainMenu.MaplotMenu"
)

var (
	//go:embed code.yaml
	codeYaml string
	//go:embed city.yaml
	cityYaml string
	//go:embed maplot.pas
	runDelphi string
)

type Code struct {
	Name   string
	Cities []string
	marked bool
}

type City struct {
	Lat float64
	Lon float64
}

var (
	codeMap map[string]Code
	cityMap map[string]City
)

var (
	number = regexp.MustCompile("[0-9]*")
	marker = color.RGBA{0xff, 0, 0, 0xff}
)

var ctx *sm.Context

var (
	file string
	temp *os.File
	form *winc.Form
	view *winc.ImageView
	pane *winc.Panel
)

func init() {
	OnLaunchEvent = onLaunchEvent
	OnFinishEvent = onFinishEvent
	OnAttachEvent = onAttachEvent
	OnInsertEvent = onInsertEvent
}

func onLaunchEvent() {
	createWindow()
	ctx = sm.NewContext()
	ctx.SetSize(form.ClientWidth(), form.ClientHeight())
	yaml.UnmarshalStrict([]byte(cityYaml), &cityMap)
	yaml.UnmarshalStrict([]byte(codeYaml), &codeMap)
	RunDelphi(runDelphi)
	HandleButton(MAPLOT_MENU, func(num int) {
		form.Show()
		update()
	})
}

func onFinishEvent() {
	x, y := form.Pos()
	SetINI(MAPLOT_NAME, "x", strconv.Itoa(x))
	SetINI(MAPLOT_NAME, "y", strconv.Itoa(y))
}

func onAttachEvent(contest, config string) {
	file = Query("{F}.png")
}

func onInsertEvent(qso *QSO) {
	rcvd := number.FindString(qso.GetRcvd())
	if code, ok := codeMap[rcvd]; ok {
		enable(code, marker)
	}
	if form.Visible() {
		update()
	}
}

func enable(code Code, marker color.RGBA) {
	for _, city := range code.Cities {
		if pt, ok := cityMap[city]; ok {
			p := s2.LatLngFromDegrees(pt.Lat, pt.Lon)
			m := sm.NewMarker(p, marker, MARKER_SIZE)
			if !code.marked {
				ctx.AddMarker(m)
				code.marked = true
			}
		}
	}
}

func update() {
	if img, e := ctx.Render(); e == nil {
		if gg.SavePNG(file, img) == nil {
			view.DrawImageFile(file)
			pane.Invalidate(true)
		}
	}
}

func createWindow() {
	form = winc.NewForm(nil)
	form.SetText("Maplot")
	icon, err := winc.ExtractIcon("zlog.exe", 0)
	if err == nil {
		form.SetIcon(0, icon)
	}
	x, _ := strconv.Atoi(GetINI(MAPLOT_NAME, "x"))
	y, _ := strconv.Atoi(GetINI(MAPLOT_NAME, "y"))
	form.SetSize(WINDOW_SIZE, WINDOW_SIZE)
	form.EnableSizable(false)
	form.EnableMaxButton(false)
	if x <= 0 || y <= 0 {
		form.Center()
	} else {
		form.SetPos(x, y)
	}
	pane = winc.NewPanel(form)
	view = winc.NewImageView(pane)
	form.OnClose().Bind(closeWindow)
	dock := winc.NewSimpleDock(form)
	dock.Dock(pane, winc.Fill)
	return
}

func closeWindow(arg *winc.Event) {
	form.Hide()
}
