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
var CALL string

//go:embed rttest.dat
var cityMultiList string

var ws = recws.RecConn{
	KeepAliveTimeout: 30 * time.Second,
}

var (
	LOCALE = time.Local
	BINARY = websocket.BinaryMessage
	stopCh chan bool
)

type Station struct {
	Call  string `json:"call"`
	Score int    `json:"score"`
	Total int    `json:"total"`
}

type Stations map[string]([]Station)

func init() {
	stopCh = make(chan bool)
	CityMultiList = cityMultiList
	OnAssignEvent = onAssignEvent
	OnDetachEvent = onDetachEvent
	OnInsertEvent = onInsertEvent
	OnDeleteEvent = onDeleteEvent
	h.HandleFunc("/", wait)
	go h.ListenAndServe(":8873", nil)
}

func wait(w h.ResponseWriter, r *h.Request) {
	r.ParseForm()
	UID = r.FormValue("id")
	if UID != "" {
		SetINI(SEC, KEY, UID)
		onAssignEvent("", "")
	}
}

func onAssignEvent(contest, configs string) {
	UID = GetINI(SEC, KEY)
	CALL = Query("{C}")
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

func onDetachEvent(contest, configs string) {
	if ws.IsConnected() {
		close(stopCh)
		ws.Close()
	}
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
		submit(INSERT, qso.Dump(LOCALE))
		DisplayToast("insert %s", qso.GetCall())
	}
}

func onDeleteEvent(qso *QSO) {
	if ws.IsConnected() {
		submit(DELETE, qso.Dump(LOCALE))
		DisplayToast("delete %s", qso.GetCall())
	}
}

func RealTimeStreamHandlerInfiniteLoop() {
	for ok := true; ok; _, ok = <-stopCh {
		_, data, err := ws.ReadMessage()
		if err == nil {
			var stations Stations
			json.Unmarshal(data, &stations)
			OnRealTimeStreamEvent(stations)
		}
	}
}

func OnRealTimeStreamEvent(stations Stations) {
	for section, stations := range stations {
		sort.Slice(stations, func(i, j int) bool {
			total_i := stations[i].Total
			total_j := stations[j].Total
			return total_i > total_j
		})
		for _, station := range stations {
			if station.Call == CALL {
				Display(section, station)
			}
		}
	}
}

func Display(section string, station Station) {
	DisplayToast("%s: %d", section, station.Total)
}
