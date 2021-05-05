/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"github.com/nextzlog/zylo"
	mapset "github.com/deckarep/golang-set"
)

func zlaunch() {
	zylo.Notify("CQ!")
}

func zfinish() {
	zylo.Notify("Bye")
}

func zattach(test string, path string) {
	zylo.Notify("%s opened", test)
}

func zdetach() {
	zylo.Notify("contest closed")
}

func zverify(list zylo.Log) (score int) {
	for _, qso := range list {
		call := qso.GetCall()
		rcvd := qso.GetRcvd()
		qso.SetMul1(rcvd)
		if call != "" && rcvd != "" {
			score = 1
		}
	}
	return
}

func zupdate(list zylo.Log) (total int) {
	calls := mapset.NewSet()
	mults := mapset.NewSet()
	for _, qso := range list {
		calls.Add(qso.GetCall())
		mults.Add(qso.GetMul1())
	}
	score := calls.Cardinality()
	multi := mults.Cardinality()
	total = score * multi
	return
}

func zinsert(list zylo.Log) {
	for _, qso := range list {
		zylo.Notify("insert %s", qso.GetCall())
	}
}

func zdelete(list zylo.Log) {
	for _, qso := range list {
		zylo.Notify("delete %s", qso.GetCall())
	}
}

func zkpress(key int, source string) (block bool) {
	return
}

func zfclick(btn int, source string) (block bool) {
	zylo.Notify("button %s (%d) clicked", source, btn)
	return
}

func main() {}
