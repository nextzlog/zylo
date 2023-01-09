/*
 Copyright (C) 2022 JA1ZLO.
*/
package main

import (
	_ "embed"
	"encoding/binary"
	"github.com/gen2brain/malgo"
	"github.com/tadvi/winc"
	"github.com/thoas/go-funk"
	"strings"
	"zylo/morse"
	"zylo/reiwa"
	"zylo/win32"
)

const (
	CHOTTO_MENU = "MainForm.MainMenu.ChottoMenu"
	INTERVAL_MS = 500
)

const FONT = "MS UI Gothic"

//go:embed chotto.pas
var runDelphi string

var (
	dev *malgo.Device
	ctx *malgo.AllocatedContext
)

var monitor morse.Monitor

var (
	form *winc.Form
	view *winc.Label
)

func init() {
	reiwa.PluginName = "chotto"
	reiwa.OnLaunchEvent = onLaunchEvent
	reiwa.OnFinishEvent = onFinishEvent
}

func onLaunchEvent() {
	createWindow()
	reiwa.RunDelphi(runDelphi)
	reiwa.HandleButton(CHOTTO_MENU, onButtonEvent)
}

func onFinishEvent() {
	closeWindow(nil)
}

func onButtonEvent(num int) {
	form.Show()
	if err := createMonitor(); err != nil {
		reiwa.DisplayModal(err.Error())
	} else {
		dev.Start()
	}
}

func onDecodeEvent(signal []float64) {
	text := []string{}
	for _, m := range monitor.Read(signal) {
		text = append(text, morse.CodeToText(m.Code))
	}
	view.SetText(strings.Join(text, "\n"))
}

func onSignalEvent(out, in []byte, frames uint32) {
	go onDecodeEvent(readSignedInt(in))
}

func readSignedInt(signal []byte) (result []float64) {
	for _, b := range funk.Chunk(signal, 4).([][]byte) {
		v := binary.LittleEndian.Uint32(b)
		result = append(result, float64(int32(v)))
	}
	return
}

func DeviceConfig() (cfg malgo.DeviceConfig) {
	cfg = malgo.DefaultDeviceConfig(malgo.Capture)
	cfg.PeriodSizeInMilliseconds = INTERVAL_MS
	cfg.Capture.Format = malgo.FormatS32
	cfg.Capture.Channels = 1
	return
}

func createMonitor() (err error) {
	ctx, err = malgo.InitContext(nil, malgo.ContextConfig{}, nil)
	if err != nil {
		return
	}
	dcb := malgo.DeviceCallbacks{
		Data: onSignalEvent,
	}
	dev, err = malgo.InitDevice(ctx.Context, DeviceConfig(), dcb)
	if err != nil {
		return
	}
	monitor = morse.DefaultMonitor(int(dev.SampleRate()))
	return
}

func createWindow() {
	form = win32.NewForm(nil)
	view = winc.NewLabel(form)
	form.SetText("Chotto Wakaru CW")
	form.OnClose().Bind(closeWindow)
	dock := winc.NewSimpleDock(form)
	dock.Dock(view, winc.Fill)
	view.SetText("")
	view.SetFont(winc.NewFont(FONT, 24, 0))
	return
}

func closeWindow(event *winc.Event) {
	form.Hide()
	if dev != nil {
		dev.Uninit()
		dev = nil
	}
	if ctx != nil {
		ctx.Uninit()
		ctx.Free()
		ctx = nil
	}
}
