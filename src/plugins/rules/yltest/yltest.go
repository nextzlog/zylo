/*
Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"strconv"
	"zylo/reiwa"
)

//go:embed yltest.dat
var cityMultiList string

func init() {
	reiwa.CityMultiList = cityMultiList
	reiwa.OnAcceptEvent = onAcceptEvent
	reiwa.AllowBandRange(reiwa.K1900, reiwa.M1200)
	reiwa.AllowModeRange(reiwa.CW, reiwa.AM)
	reiwa.AllowCall(`^\w{3,}`)
	reiwa.AllowRcvd(`^\d{4,}$`)
	reiwa.AllowSent(`^\d{4,}$`)
}

const (
	OM = 0
	YL = 1
	JL = 2
)

func codeToClass(code int) int {
	switch {
	case code < 2000:
		return OM
	case code > 5000:
		return JL
	default:
		return YL
	}
}

func scoreWithOM(rcvd int) int {
	switch codeToClass(rcvd) {
	case JL:
		return 5
	case YL:
		return 1
	default:
		return 0
	}
}

func scoreWithYL(rcvd int) int {
	switch codeToClass(rcvd) {
	case JL:
		return 5
	case YL:
		return 5
	default:
		return 1
	}
}

func scoreWithJL(rcvd int) int {
	return 1
}

func score(rcvd, sent int) int {
	switch sent {
	case JL:
		return scoreWithJL(rcvd)
	case YL:
		return scoreWithYL(rcvd)
	default:
		return scoreWithOM(rcvd)
	}
}

func onAcceptEvent(qso *reiwa.QSO) {
	qso.SetMul1(qso.GetCall()[:3])
	qso.SetMul2(qso.GetRcvd())
	rcvd, _ := strconv.Atoi(qso.GetRcvd())
	sent, _ := strconv.Atoi(qso.GetSent())
	qso.Score = byte(score(rcvd, sent))
	if qso.Score == 0 {
		qso.SetNote("OM-to-OM QSO")
	}
}
