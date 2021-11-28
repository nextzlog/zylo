/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pkg/browser"
	"github.com/recws-org/recws"
	h "net/http"
	"os"
	"sort"
	"time"
)

const (
	ATS = "https://realtime.allja1.org"
	WSS = "wss://realtime.allja1.org/agent/%s"
)

const (
	SEC = "ATS"
	KEY = "UID"
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

type Station struct {
	Call  string `json:"call"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}

type Sections map[string]([]Station)

func init() {
	stopCh = make(chan bool)
	CityMultiList = cityMultiList
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

func onAssignEvent(contest, configs string) {
	ShowLeaderWindow()
	connectWebSocketAPI()
	go server.ListenAndServe()
}

func onDetachEvent(contest, configs string) {
	CloseLeaderWindow()
	if ws.IsConnected() {
		close(stopCh)
		ws.Close()
	}
	server.Close()
}

const (
	INSERT = 0
	DELETE = 1
)

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
		Display(section, stations)
	}
}
