/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"os"
	"strings"
)

var hsmults int

//go:embed hstest.dat
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
	hsmults = 0
	bin, _ := os.ReadFile(Query("{F}"))
	for _, qso := range LoadZLO(bin) {
		onInsertEvent(&qso)
	}
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

func splitAndSetMuls(qso *QSO, code string) {
	rcvd := qso.GetRcvd()
	head := strings.TrimSuffix(rcvd, code)
	if strings.HasSuffix(rcvd, code) {
		qso.SetMul1(head)
		qso.SetMul2(code)
	}
}

func onVerifyEvent(qso *QSO) {
	splitAndSetMuls(qso, "C")
	splitAndSetMuls(qso, "HS")
	if qso.Dupe {
		qso.Score = 0
	} else if qso.Mode == CW {
		qso.Score = 3
	} else {
		qso.Score = 1
	}
}

func onPointsEvent(score, mults int) int {
	return score * (mults + hsmults)
}
