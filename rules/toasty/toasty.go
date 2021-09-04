/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import _ "embed"

//go:embed toasty.dat
var cityMultiList string

func init() {
	CityMultiList = cityMultiList
	OnLaunchEvent = onLaunchEvent
	OnFinishEvent = onFinishEvent
	OnAttachEvent = onAttachEvent
	OnInsertEvent = onInsertEvent
	OnDeleteEvent = onDeleteEvent
}

func onLaunchEvent() {
	DisplayToast("CQ!")
	HandleButton("CWPlayButton", onButtonEvent)
	HandleButton("CWStopButton", onButtonEvent)
	HandleEditor("CallsignEdit", onEditorEvent)
}

func onFinishEvent() {
	DisplayToast("Bye")
}

func onAttachEvent(contest, configs string) {
	DisplayToast(contest)
}

func onInsertEvent(qso *QSO) {
	DisplayToast("insert %s", qso.GetCall())
}

func onDeleteEvent(qso *QSO) {
	DisplayToast("delete %s", qso.GetCall())
}

func onButtonEvent(btn int) {
	DisplayToast("button (%d) clicked", btn)
}

func onEditorEvent(key int) {
	if rune(key) == ' ' {
		DisplayToast(Query("$B"))
	}
}
