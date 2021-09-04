/*
 Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"github.com/cavaliercoder/grab"
	"os"
	"os/exec"
	"strings"
)

const QXSL = "qxsl.exe"

//go:embed qxsl.url
var url string

//go:embed qxsl.fmt
var fileExtFilter string

func init() {
	OnLaunchEvent = onLaunchEvent
	OnImportEvent = onImportEvent
	OnExportEvent = onExportEvent
	FileExtFilter = fileExtFilter
}

func onLaunchEvent() {
	if _, err := os.Stat(QXSL); err != nil {
		grab.Get(".", strings.TrimSpace(url))
	}
}

func onImportEvent(source, target string) error {
	return exec.Command(QXSL, "format", source, target, "zbin").Run()
}

func onExportEvent(source, format string) error {
	return exec.Command(QXSL, "format", source, source, format).Run()
}
