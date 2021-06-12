/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"strings"
	"github.com/nextzlog/zylo"
)

var score = 0

func zlaunch() {
	zylo.Notify("CQ!")
}

func zfinish() {
	zylo.Notify("Bye")
}

func zattach(test string, path string) {
	zylo.HookButton("CWPlayButton")
	zylo.HookButton("CWStopButton")
	zylo.HookEditor("CallsignEdit")
	zylo.Notify(test)
}

func zdetach() {
	score = 0
}

func zinsert(qso *zylo.QSO) {
	score += int(qso.Score)
}

func zdelete(qso *zylo.QSO) {
	score -= int(qso.Score)
}

func zverify(qso *zylo.QSO) {
	rcvd := qso.GetRcvd()
	qso.SetMul1(rcvd)
	if qso.Dupe {
		qso.Score = 0
	} else {
		qso.Score = 1
	}
}

func zcities() string {
	var list []string;
	list = append(list, "100105 Bunkyo")
	list = append(list, "100110 Meguro")
	return strings.Join(list, "\n")
}

func zpoints(score, mults int) int {
	return score * mults
}

func zeditor(key int, name string) bool {
	return rune(key) == ' '
}

func zbutton(btn int, name string) bool {
	zylo.Notify("button %s clicked", name)
	return false
}

func main() {}
