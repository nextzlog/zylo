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
	"zylo/reiwa"
)

const (
	WINDOW_SIZE = 600
	MARKER_SIZE = 16
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
	reiwa.PluginName = "maplot"
	reiwa.OnLaunchEvent = onLaunchEvent
	reiwa.OnAttachEvent = onAttachEvent
	reiwa.OnInsertEvent = onInsertEvent
}

func onLaunchEvent() {
	createWindow()
	ctx = sm.NewContext()
	ctx.SetSize(form.ClientWidth(), form.ClientHeight())
	yaml.UnmarshalStrict([]byte(cityYaml), &cityMap)
	yaml.UnmarshalStrict([]byte(codeYaml), &codeMap)
	reiwa.RunDelphi(runDelphi)
	reiwa.HandleButton(MAPLOT_MENU, func(num int) {
		form.Show()
		update()
	})
}

func onAttachEvent(contest, config string) {
	file = reiwa.Query("{F}.png")
}

func onInsertEvent(qso *reiwa.QSO) {
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
	form = newForm(nil)
	form.SetSize(WINDOW_SIZE, WINDOW_SIZE)
	form.EnableSizable(false)
	form.EnableMaxButton(false)
	pane = newPanel(form)
	view = winc.NewImageView(pane)
	dock := winc.NewSimpleDock(form)
	dock.Dock(pane, winc.Fill)
	return
}