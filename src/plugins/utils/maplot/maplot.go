/*
Copyright (C) 2022 JA1ZLO.
*/
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/pkg/browser"
	"net/http"
	"regexp"
	"zylo/reiwa"
)

type Marker struct {
	Code string
	Drop bool
}

const MAPLOT_MENU = "MainForm.MainMenu.MaplotMenu"

//go:embed maplot.pas
var runDelphi string

var (
	markers []Marker
	pattern = regexp.MustCompile("[0-9]+")
)

var server = &http.Server{Addr: ":49599"}

func init() {
	reiwa.PluginName = "maplot"
	reiwa.OnLaunchEvent = onLaunchEvent
	reiwa.OnFinishEvent = onFinishEvent
	reiwa.OnDetachEvent = onDetachEvent
	reiwa.OnInsertEvent = onInsertEvent
	reiwa.OnDeleteEvent = onDeleteEvent
}

func onLaunchEvent() {
	http.HandleFunc("/", serve)
	go server.ListenAndServe()
	reiwa.RunDelphi(runDelphi)
	reiwa.HandleButton(MAPLOT_MENU, onBrowseEvent)
}

func onFinishEvent() {
	server.Close()
}

func onDetachEvent(contest, configs string) {
	markers = nil
}

func onInsertEvent(qso *reiwa.QSO) {
	marker := Marker{pattern.FindString(qso.GetRcvd()), false}
	markers = append(markers, marker)
}

func onDeleteEvent(qso *reiwa.QSO) {
	marker := Marker{pattern.FindString(qso.GetRcvd()), true}
	markers = append(markers, marker)
}

func onBrowseEvent(num int) {
	browser.OpenURL("https://jg1vpp.github.io/qth.zlog.org")
}

func serve(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	if text, err := json.Marshal(markers); err != nil {
		reiwa.DisplayModal(err.Error())
	} else {
		fmt.Fprintf(writer, string(text))
	}
}
