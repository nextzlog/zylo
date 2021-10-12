/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"os/exec"
)

const QXSL = "qxsl.exe"

//go:embed qxsl.fmt
var fileExtFilter string

func init() {
	OnImportEvent = onImportEvent
	OnExportEvent = onExportEvent
	FileExtFilter = fileExtFilter
}

func onImportEvent(source, target string) error {
	return exec.Command(QXSL, "format", source, target, "zbin").Run()
}

func onExportEvent(source, format string) error {
	return exec.Command(QXSL, "format", source, source, format).Run()
}
