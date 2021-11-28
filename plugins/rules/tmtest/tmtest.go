/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import _ "embed"

var days = make(map[int]int)

//go:embed tmtest.dat
var cityMultiList string

func init() {
	CityMultiList = cityMultiList
	OnInsertEvent = onInsertEvent
	OnDeleteEvent = onDeleteEvent
	OnAcceptEvent = onAcceptEvent
	OnPointsEvent = onPointsEvent
	AllowBandRange(M50, G10UP)
	AllowMode(CW, SSB, FM, AM)
	AllowCall(`^\w{3,}`)
	AllowRcvd(`^\d{3,}$`)
}

func onInsertEvent(qso *QSO) {
	days[qso.GetTime().YearDay()] += 1
}

func onDeleteEvent(qso *QSO) {
	d := qso.GetTime().YearDay()
	if days[d] > 1 {
		days[d] -= 1
	} else {
		delete(days, d)
	}
}

func mult(qso *QSO) string {
	call := qso.GetCallSign()
	return call[len(call)-1:]
}

func score(qso *QSO) byte {
	switch qso.Band {
	case M50:
		return 1
	case M144:
		return 1
	case M430:
		return 1
	case M1200:
		return 2
	case M2400:
		return 5
	case M5600:
		return 10
	case G10UP:
		return 20
	default:
		return 0
	}
}

func onAcceptEvent(qso *QSO) {
	qso.Score = score(qso)
	qso.SetMul1(mult(qso))
}

func onPointsEvent(score, mults int) int {
	return score * mults * len(days)
}
