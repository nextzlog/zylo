/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"zylo/reiwa"
)

var days = make(map[int]int)

//go:embed tmtest.dat
var cityMultiList string

func init() {
	reiwa.CityMultiList = cityMultiList
	reiwa.OnInsertEvent = onInsertEvent
	reiwa.OnDeleteEvent = onDeleteEvent
	reiwa.OnAcceptEvent = onAcceptEvent
	reiwa.OnPointsEvent = onPointsEvent
	reiwa.AllowBandRange(reiwa.M50, reiwa.G10UP)
	reiwa.AllowModeRange(reiwa.CW, reiwa.AM)
	reiwa.AllowCall(`^\w{3,}`)
	reiwa.AllowRcvd(`^\d{3,}$`)
}

func onInsertEvent(qso *reiwa.QSO) {
	days[qso.GetTime().YearDay()] += 1
}

func onDeleteEvent(qso *reiwa.QSO) {
	d := qso.GetTime().YearDay()
	if days[d] > 1 {
		days[d] -= 1
	} else {
		delete(days, d)
	}
}

func mult(qso *reiwa.QSO) string {
	call := qso.GetCallSign()
	return call[len(call)-1:]
}

func score(qso *reiwa.QSO) byte {
	switch qso.Band {
	case reiwa.M50:
		return 1
	case reiwa.M144:
		return 1
	case reiwa.M430:
		return 1
	case reiwa.M1200:
		return 2
	case reiwa.M2400:
		return 5
	case reiwa.M5600:
		return 10
	case reiwa.G10UP:
		return 20
	default:
		return 0
	}
}

func onAcceptEvent(qso *reiwa.QSO) {
	qso.Score = score(qso)
	qso.SetMul1(mult(qso))
}

func onPointsEvent(score, mults int) int {
	return score * mults * len(days)
}
