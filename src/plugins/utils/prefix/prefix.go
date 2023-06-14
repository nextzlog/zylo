/*
 Copyright (C) 2022 JA1ZLO.
*/
package main

import (
	_ "embed"
	"gopkg.in/yaml.v2"
	"zylo/reiwa"
)

//go:embed prefix.yaml
var preYaml string

var preMap map[string]string

func init() {
	reiwa.PluginName = "prefix"
	reiwa.OnLaunchEvent = onLaunchEvent
	reiwa.OnInsertEvent = onInsertEvent
}

func onLaunchEvent() {
	yaml.UnmarshalStrict([]byte(preYaml), &preMap)
	for key, value := range preMap {
		if value == "Japan" {
			delete(preMap, key)
		}
	}
}

func onInsertEvent(qso *reiwa.QSO) {
	call := qso.GetCall()
	for num := len(call); num > 0; num-- {
		if n, ok := preMap[call[:num]]; ok {
			reiwa.DisplayToast("%s: %s", call, n)
			break
		}
	}
}
