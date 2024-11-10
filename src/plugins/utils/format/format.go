/*
Copyright (C) 2020 JA1ZLO.
*/
package main

import (
	_ "embed"
	"os/exec"
	"path/filepath"
	"syscall"
	"zylo/reiwa"
)

const (
	ZYLO = "zylo"
	PATH = "path"
	QXSL = "qxsl.exe"
)

//go:embed qxsl.fmt
var fileExtFilter string

func init() {
	reiwa.OnImportEvent = onImportEvent
	reiwa.OnExportEvent = onExportEvent
	reiwa.FileExtFilter = fileExtFilter
}

func onImportEvent(source, target string) error {
	return invoke(source, target, "zbin")
}

func onExportEvent(source, format string) error {
	return invoke(source, source, format)
}

func invoke(source, target, format string) (err error) {
	exe := filepath.Join(reiwa.GetINI(ZYLO, PATH), QXSL)
	exe, _ = filepath.Abs(exe)
	cmd := exec.Command(exe, source, target, format)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	if err = cmd.Run(); err != nil {
		reiwa.DisplayModal(err.Error())
	}
	return
}
