/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"zylo/reiwa"
)

var hsmults int

//go:embed hstest.dat
var cityMultiList string

func init() {
	reiwa.CityMultiList = cityMultiList
	reiwa.OnAssignEvent = onAssignEvent
	reiwa.OnInsertEvent = onInsertEvent
	reiwa.OnDeleteEvent = onDeleteEvent
	reiwa.OnAcceptEvent = onAcceptEvent
	reiwa.OnPointsEvent = onPointsEvent
	reiwa.AllowBand(reiwa.M7, reiwa.M21, reiwa.M50, reiwa.M144, reiwa.M430)
	reiwa.AllowModeRange(reiwa.CW, reiwa.AM)
	reiwa.AllowRcvd(`^(\d{2,})(HS|C)$`)
}

func onAssignEvent(contest, configs string) {
	hsmults = 0
}

func onInsertEvent(qso *reiwa.QSO) {
	if qso.GetMul2() == "HS" {
		hsmults += 1
	}
}

func onDeleteEvent(qso *reiwa.QSO) {
	if qso.GetMul2() == "HS" {
		hsmults -= 1
	}
}

func onAcceptEvent(qso *reiwa.QSO) {
	rcvd := qso.GetRcvdGroups()
	qso.SetMul1(rcvd[1])
	qso.SetMul2(rcvd[2])
	if qso.Mode == reiwa.CW {
		qso.Score = 3
	} else {
		qso.Score = 1
	}
}

func onPointsEvent(score, mults int) int {
	return score * (mults + hsmults)
}
