/*
Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	"github.com/tcnksm/go-latest"
	"regexp"
	"zylo/reiwa"
)

const (
	USER = "jr8ppg"
	REPO = "zLog"
)

func init() {
	reiwa.OnLaunchEvent = onLaunchEvent
}

var tag = latest.GithubTag{
	Owner:             USER,
	Repository:        REPO,
	FixVersionStrFunc: FixVersionString,
}

func FixVersionString(number string) string {
	rep := regexp.MustCompile(`ZLOG(\d)(\d)(\d)(\d+)`)
	return rep.ReplaceAllString(number, "$1.$2.$3.$4")
}

func onLaunchEvent() {
	go fetchReleases()
}

func fetchReleases() {
	res, _ := latest.Check(&tag, reiwa.Query("{V}"))
	if res != nil && res.Outdated {
		reiwa.DisplayToast("%s is now available", res.Current)
	}
}
