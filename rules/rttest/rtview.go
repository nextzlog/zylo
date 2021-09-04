/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"bufio"
	_ "embed"
	"github.com/tadvi/winc"
	"math"
	"strconv"
	"strings"
)

const RT = "rttest"

var (
	form *winc.Form
	tabs *winc.TabView
)

//go:embed rtview.tab
var sects string
var views = make(map[string]ScoreView)

type ScoreItem struct {
	ranking int
	station Station
}

type ScoreView struct {
	list *winc.ListView
}

func (item ScoreItem) Text() (text []string) {
	text = append(text, strconv.Itoa(item.ranking))
	text = append(text, item.station.Call)
	text = append(text, strconv.Itoa(item.station.Score))
	text = append(text, strconv.Itoa(item.station.Total))
	return
}

func (item ScoreItem) ImageIndex() int {
	return 0
}

func ShowLeaderWindow() {
	form = winc.NewForm(nil)
	form.SetText("ATS-4 Stream Service")
	x, _ := strconv.Atoi(GetINI(RT, "x"))
	y, _ := strconv.Atoi(GetINI(RT, "y"))
	w, _ := strconv.Atoi(GetINI(RT, "w"))
	h, _ := strconv.Atoi(GetINI(RT, "h"))
	if w <= 0 || h <= 0 {
		w = 300
		h = 300
	}
	form.SetSize(w, h)
	if x <= 0 || y <= 0 {
		form.Center()
	} else {
		form.SetPos(x, y)
	}
	tabs = winc.NewTabView(form)
	dock := winc.NewSimpleDock(form)
	dock.Dock(tabs, winc.Top)
	dock.Dock(tabs.Panels(), winc.Fill)
	reader := strings.NewReader(sects)
	buffer := bufio.NewScanner(reader)
	for buffer.Scan() {
		AddSection(buffer.Text())
	}
	form.Show()
	return
}

func CloseLeaderWindow() {
	x, y := form.Pos()
	w, h := form.Size()
	SetINI(RT, "x", strconv.Itoa(x))
	SetINI(RT, "y", strconv.Itoa(y))
	SetINI(RT, "w", strconv.Itoa(w))
	SetINI(RT, "h", strconv.Itoa(h))
	form.Close()
}

func AddSection(section string) (view ScoreView) {
	panel := tabs.AddPanel(section)
	view.list = winc.NewListView(panel)
	view.list.EnableEditLabels(false)
	view.list.AddColumn("rank", 120)
	view.list.AddColumn("call", 120)
	view.list.AddColumn("score", 120)
	view.list.AddColumn("total", 120)
	dock := winc.NewSimpleDock(panel)
	dock.Dock(view.list, winc.Fill)
	views[section] = view
	tabs.SetCurrent(0)
	return
}

func (view ScoreView) Update(stations []Station) {
	view.list.DeleteAllItems()
	count := 0
	worst := math.MaxInt64
	for num, station := range stations {
		if worst > station.Total {
			worst = station.Total
			count = num + 1
		}
		view.list.AddItem(ScoreItem{
			ranking: count,
			station: station,
		})
	}
}

func Display(section string, stations []Station) {
	views[section].Update(stations)
}
