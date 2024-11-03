/*
Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"os/exec"
	"zylo/reiwa"
)

const QXSL = "./qxsl.exe"

//go:embed qxsl.fmt
var fileExtFilter string

func init() {
	reiwa.OnImportEvent = onImportEvent
	reiwa.OnExportEvent = onExportEvent
	reiwa.FileExtFilter = fileExtFilter
}

func onImportEvent(source, target string) error {
	return exec.Command(QXSL, source, target, "zbin").Run()
}

func onExportEvent(source, format string) error {
	return exec.Command(QXSL, source, source, format).Run()
}
