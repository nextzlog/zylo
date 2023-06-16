/*
Copyright (C) 2022 JA1ZLO.
*/
package main

import (
	_ "embed"
	"github.com/tadvi/winc"
	"gopkg.in/yaml.v2"
	"regexp"
	"zylo/reiwa"
	"zylo/win32"
)

const MAPLOT_MENU = "MainForm.MainMenu.MaplotMenu"

const SIZE = 800

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
}

type City struct {
	YX map[int]map[int]int
}

var (
	codeMap map[string]Code
	cityMap map[string]City
	enabled []City
)

var (
	rcvd = regexp.MustCompile("\\d*")
	face = winc.RGB(0xff, 0x00, 0x00)
	edge = winc.RGB(0x00, 0x00, 0x00)
)

var (
	form *winc.Form
	view *winc.Canvas
)

func init() {
	reiwa.PluginName = "maplot"
	reiwa.OnLaunchEvent = onLaunchEvent
	reiwa.OnInsertEvent = onInsertEvent
}

func onLaunchEvent() {
	createWindow()
	yaml.UnmarshalStrict([]byte(cityYaml), &cityMap)
	yaml.UnmarshalStrict([]byte(codeYaml), &codeMap)
	reiwa.RunDelphi(runDelphi)
	reiwa.HandleButton(MAPLOT_MENU, func(num int) {
		form.Show()
		onUpdateEvent(nil)
	})
}

func onInsertEvent(qso *reiwa.QSO) {
	n := rcvd.FindString(qso.GetRcvd())
	if code, ok := codeMap[n]; ok {
		mark(code, face)
	}
}

func onUpdateEvent(ev *winc.Event) {
	fill(cityMap["JA"], edge)
	for _, code := range enabled {
		fill(code, face)
	}
}

func mark(code Code, color winc.Color) {
	for _, city := range code.Cities {
		if pt, ok := cityMap[city]; ok {
			enabled = append(enabled, pt)
			fill(pt, face)
		}
	}
}

func fill(city City, color winc.Color) {
	br := winc.NewSolidColorBrush(color)
	for y, runs := range city.YX {
		for x, width := range runs {
			line(y, x, width, br)
		}
	}
	br.Dispose()
}

func line(y, x, w int, brush *winc.Brush) {
	view.FillRect(winc.NewRect(x, y, x+w, y+1), brush)
}

func createWindow() {
	win32.DefaultWindowW = SIZE
	win32.DefaultWindowH = SIZE
	form = win32.NewForm(nil)
	view = winc.NewCanvasFromHwnd(form.Handle())
	form.OnSize().Bind(onUpdateEvent)
	return
}
