/*
Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"zylo/reiwa"
)

var enabled bool
var hsmults int
var hstable map[string]int

//go:embed hstest.dat
var cityMultiList string

func init() {
	hstable = make(map[string]int)
	reiwa.CityMultiList = cityMultiList
	reiwa.OnAssignEvent = onAssignEvent
	reiwa.OnDetachEvent = onDetachEvent
	reiwa.OnInsertEvent = onInsertEvent
	reiwa.OnDeleteEvent = onDeleteEvent
	reiwa.OnAcceptEvent = onAcceptEvent
	reiwa.OnPointsEvent = onPointsEvent
	reiwa.AllowBand(reiwa.M7, reiwa.M21, reiwa.M50, reiwa.M144, reiwa.M430)
	reiwa.AllowModeRange(reiwa.CW, reiwa.AM)
	reiwa.AllowRcvd(`^(\d{2,})(HS|C)$`)
}

func onAssignEvent(contest, configs string) {
	enabled = true
}

func onDetachEvent(contest, configs string) {
	clear(hstable)
	hsmults = 0
	enabled = false
}

func onInsertEvent(qso *reiwa.QSO) {
	if enabled && qso.GetMul2() == "HS" {
		call := qso.GetCall()
		if n, ok := hstable[call]; ok {
			hstable[call] = n + 1
		} else {
			hstable[call] = 1
			hsmults += 1
		}
	}
}

func onDeleteEvent(qso *reiwa.QSO) {
	if enabled && qso.GetMul2() == "HS" {
		call := qso.GetCall()
		if n := hstable[call]; n <= 1 {
			delete(hstable, call)
			hsmults -= 1
		} else {
			hstable[call] = n - 1
		}
	}
}

func onAcceptEvent(qso *reiwa.QSO) {
	rcvd := qso.GetRcvdGroups()
	qso.SetMul1(rcvd[1])
	if rcvd[2] == "HS" {
		qso.SetMul2("HS")
	} else {
		qso.SetMul2("")
	}
	if qso.Mode == reiwa.CW {
		qso.Score = 3
	} else {
		qso.Score = 1
	}
}

func onPointsEvent(score, mults int) int {
	return score * (mults + hsmults)
}
