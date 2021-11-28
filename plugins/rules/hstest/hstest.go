/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import _ "embed"

var hsmults int

//go:embed hstest.dat
var cityMultiList string

func init() {
	CityMultiList = cityMultiList
	OnAssignEvent = onAssignEvent
	OnInsertEvent = onInsertEvent
	OnDeleteEvent = onDeleteEvent
	OnAcceptEvent = onAcceptEvent
	OnPointsEvent = onPointsEvent
	AllowBand(M7, M21, M50, M144, M430)
	AllowMode(CW, SSB, FM, AM)
	AllowRcvd(`^(\d{2,})(HS|C)$`)
}

func onAssignEvent(contest, configs string) {
	hsmults = 0
}

func onInsertEvent(qso *QSO) {
	if qso.GetMul2() == "HS" {
		hsmults += 1
	}
}

func onDeleteEvent(qso *QSO) {
	if qso.GetMul2() == "HS" {
		hsmults -= 1
	}
}

func onAcceptEvent(qso *QSO) {
	rcvd := qso.GetRcvdGroups()
	qso.SetMul1(rcvd[1])
	qso.SetMul2(rcvd[2])
	if qso.Mode == CW {
		qso.Score = 3
	} else {
		qso.Score = 1
	}
}

func onPointsEvent(score, mults int) int {
	return score * (mults + hsmults)
}
