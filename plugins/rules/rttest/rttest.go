/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tadvi/winc"
	"github.com/pkg/browser"
	"github.com/recws-org/recws"
	"math"
	h "net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	ATS = "https://realtime.allja1.org"
	WSS = "wss://realtime.allja1.org/agent/%s"
)

const (
	INSERT = 0
	DELETE = 1
)

const (
	SEC = "ATS"
	KEY = "UID"
)

const (
	RTTEST_NAME = "rttest"
	RTTEST_MENU = "MainForm.MainMenu.RttestMenu"
)

var UID string

//go:embed rttest.dat
var cityMultiList string

var ws = recws.RecConn{
	KeepAliveTimeout: 30 * time.Second,
}

var (
	BINARY = websocket.BinaryMessage
	stopCh chan bool
	server *h.Server
)

var (
	form *winc.Form
	tabs *winc.TabView
)

//go:embed rttest.pas
var runDelphi string

//go:embed rttest.tab
var sects string
var views = make(map[string]ScoreView)

type ScoreItem struct {
	ranking int
	station Station
}

type ScoreView struct {
	list *winc.ListView
}

type Station struct {
	Call  string `json:"call"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}

type Sections map[string]([]Station)

func init() {
	stopCh = make(chan bool)
	CityMultiList = cityMultiList
	OnLaunchEvent = onLaunchEvent
	OnFinishEvent = onFinishEvent
	OnAssignEvent = onAssignEvent
	OnDetachEvent = onDetachEvent
	OnInsertEvent = onInsertEvent
	OnDeleteEvent = onDeleteEvent
	h.HandleFunc("/", wait)
	server = &h.Server{Addr: ":8873"}
	AllowBandRange(K3500, M50)
	AllowMode(CW, SSB, FM, AM)
	AllowRcvd(`^\d{2,}^`)
}

func wait(w h.ResponseWriter, r *h.Request) {
	r.ParseForm()
	UID = r.FormValue("id")
	if UID != "" {
		SetINI(SEC, KEY, UID)
		connectWebSocketAPI()
	}
}

func connectWebSocketAPI() {
	UID = GetINI(SEC, KEY)
	ws.Dial(fmt.Sprintf(WSS, UID), nil)
	if ws.GetDialError() != nil {
		DisplayModal("authenticate via ATS-4")
		browser.OpenURL(ATS)
	} else {
		go RealTimeStreamHandlerInfiniteLoop()
		DisplayModal("successfully connected")
		binary, _ := os.ReadFile(Query("{F}"))
		submit(INSERT, binary)
	}
}

func onLaunchEvent() {
	createWindow()
	RunDelphi(runDelphi)
	HandleButton(RTTEST_MENU, func(num int) {
		form.Show()
	})
}

func onFinishEvent() {
	x, y := form.Pos()
	w, h := form.Size()
	SetINI(RTTEST_NAME, "x", strconv.Itoa(x))
	SetINI(RTTEST_NAME, "y", strconv.Itoa(y))
	SetINI(RTTEST_NAME, "w", strconv.Itoa(w))
	SetINI(RTTEST_NAME, "h", strconv.Itoa(h))
}

func onAssignEvent(contest, configs string) {
	form.Show()
	connectWebSocketAPI()
	go server.ListenAndServe()
}

func onDetachEvent(contest, configs string) {
	if ws.IsConnected() {
		close(stopCh)
		ws.Close()
	}
	server.Close()
}

func submit(cmd byte, data []byte) {
	msg := append([]byte{cmd}, data...)
	err := ws.WriteMessage(BINARY, msg)
	if err != nil {
		DisplayModal(err.Error())
	}
}

func onInsertEvent(qso *QSO) {
	if ws.IsConnected() {
		submit(INSERT, DumpZLO(*qso))
		DisplayToast("insert %s", qso.GetCall())
	}
}

func onDeleteEvent(qso *QSO) {
	if ws.IsConnected() {
		submit(DELETE, DumpZLO(*qso))
		DisplayToast("delete %s", qso.GetCall())
	}
}

func RealTimeStreamHandlerInfiniteLoop() {
	for {
		select {
		case <-stopCh:
			return
		default:
			_, data, err := ws.ReadMessage()
			if err == nil {
				var sections Sections
				json.Unmarshal(data, &sections)
				OnRealTimeStreamEvent(sections)
			}
		}
	}
}

func OnRealTimeStreamEvent(sections Sections) {
	for section, stations := range sections {
		sort.Slice(stations, func(i, j int) bool {
			total_i := stations[i].Total
			total_j := stations[j].Total
			return total_i > total_j
		})
		display(section, stations)
	}
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

func createWindow() {
	form = winc.NewForm(nil)
	tabs = winc.NewTabView(form)
	form.SetText("Real-Time Contest")
	x, _ := strconv.Atoi(GetINI(RTTEST_NAME, "x"))
	y, _ := strconv.Atoi(GetINI(RTTEST_NAME, "y"))
	w, _ := strconv.Atoi(GetINI(RTTEST_NAME, "w"))
	h, _ := strconv.Atoi(GetINI(RTTEST_NAME, "h"))
	exec, _ := os.Executable()
	icon, _ := winc.ExtractIcon(exec, 0)
	form.SetIcon(0, icon)
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
	form.OnClose().Bind(closeWindow)
	dock := winc.NewSimpleDock(form)
	dock.Dock(tabs, winc.Top)
	dock.Dock(tabs.Panels(), winc.Fill)
	reader := strings.NewReader(sects)
	buffer := bufio.NewScanner(reader)
	for buffer.Scan() {
		addSection(buffer.Text())
	}
	return
}

func closeWindow(arg *winc.Event) {
	form.Hide()
}

func addSection(section string) (view ScoreView) {
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

func (view ScoreView) update(stations []Station) {
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

func display(section string, stations []Station) {
	views[section].update(stations)
}
