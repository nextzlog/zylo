/*
Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"strings"
	"zylo/reiwa"
)

//go:embed aktest.dat
var cityMultiList string

func init() {
	reiwa.CityMultiList = cityMultiList
	reiwa.OnAcceptEvent = onAcceptEvent
	reiwa.AllowBandRange(reiwa.K3500, reiwa.M430)
	reiwa.AllowModeRange(reiwa.CW, reiwa.AM)
	reiwa.AllowRcvd(`^(\d{4,})(M?)$`)
}

func isMember(mul, mem string) bool {
	return mem == "M"
}

func isInCity(mul, mem string) bool {
	return mul == "0102"
}

func isInPref(mul, mem string) bool {
	return strings.HasPrefix(mul, "01")
}

func score(mul, mem string) byte {
	switch {
	case isMember(mul, mem):
		return 9
	case isInCity(mul, mem):
		return 9
	case isInPref(mul, mem):
		return 6
	default:
		return 3
	}
}

func onAcceptEvent(qso *reiwa.QSO) {
	rcvd := qso.GetRcvdGroups()
	mul, mem := rcvd[1], rcvd[2]
	qso.Score = score(mul, mem)
	qso.SetMul1(mul)
}
