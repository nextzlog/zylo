/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"os"
	"regexp"
)

var days = make(map[int]int)
var code = regexp.MustCompile(`^\d{3,}$`)
var call = regexp.MustCompile(`^\w{3,}`)

//go:embed tmtest.dat
var cityMultiList string

func init() {
	CityMultiList = cityMultiList
	OnAssignEvent = onAssignEvent
	OnInsertEvent = onInsertEvent
	OnDeleteEvent = onDeleteEvent
	OnVerifyEvent = onVerifyEvent
	OnPointsEvent = onPointsEvent
}

func onAssignEvent(contest, configs string) {
	bin, _ := os.ReadFile(Query("{F}"))
	for _, qso := range LoadZLO(bin) {
		onInsertEvent(&qso)
	}
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

func valid(qso *QSO) bool {
	b1 := code.MatchString(qso.GetRcvd())
	b2 := call.MatchString(qso.GetCall())
	return b1 && b2 && score(qso) > 0
}

func onVerifyEvent(qso *QSO) {
	if !qso.Dupe && valid(qso) {
		qso.Score = score(qso)
		qso.SetMul1(mult(qso))
	} else {
		qso.Score = 0
		qso.SetMul1("")
	}
}

func onPointsEvent(score, mults int) int {
	return score * mults * len(days)
}
