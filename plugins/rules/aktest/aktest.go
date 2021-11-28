/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"strings"
)

//go:embed aktest.dat
var cityMultiList string

func init() {
	CityMultiList = cityMultiList
	OnAcceptEvent = onAcceptEvent
	AllowBandRange(K3500, M430)
	AllowMode(CW, SSB, FM, AM)
	AllowRcvd(`^(\d{4,})(M?)$`)
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

func onAcceptEvent(qso *QSO) {
	rcvd := qso.GetRcvdGroups()
	mul, mem := rcvd[1], rcvd[2]
	qso.Score = score(mul, mem)
	qso.SetMul1(mul)
}
