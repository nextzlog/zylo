/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/
package win32

import (
	"fmt"
	"github.com/tadvi/winc"
	"github.com/tadvi/winc/w32"
	"os"
	"strconv"
	"zylo/reiwa"
)

/*
 この拡張機能の専用のクラス名を持つフォームを構築します。
*/
func NewForm(parent winc.Controller) *winc.Form {
	name := fmt.Sprintf("zylo_%s_form", reiwa.PluginName)
	winc.RegClassOnlyOnce(name)
	form := new(winc.Form)
	form.SetParent(parent)
	ex := uint(w32.WS_EX_CONTROLPARENT)
	st := uint(w32.WS_OVERLAPPEDWINDOW)
	exe, _ := os.Executable()
	ico, _ := winc.ExtractIcon(exe, 0)
	form.SetHandle(winc.CreateWindow(name, parent, ex, st))
	winc.RegMsgHandler(form)
	form.SetIcon(0, ico)
	form.SetIsForm(true)
	form.SetText(reiwa.PluginName)
	form.SetFont(winc.DefaultFont)
	x, _ := strconv.Atoi(reiwa.GetINI(reiwa.PluginName, "x"))
	y, _ := strconv.Atoi(reiwa.GetINI(reiwa.PluginName, "y"))
	w, _ := strconv.Atoi(reiwa.GetINI(reiwa.PluginName, "w"))
	h, _ := strconv.Atoi(reiwa.GetINI(reiwa.PluginName, "h"))
	form.OnClose().Bind(func(arg *winc.Event) {
		reiwa.SetINI(reiwa.PluginName, "x", strconv.Itoa(x))
		reiwa.SetINI(reiwa.PluginName, "y", strconv.Itoa(y))
		reiwa.SetINI(reiwa.PluginName, "w", strconv.Itoa(w))
		reiwa.SetINI(reiwa.PluginName, "h", strconv.Itoa(h))
		form.Hide()
	})
	if w <= 0 || h <= 0 {
		w = 300
		h = 300
	}
	form.SetSize(w, h)
	if x <= 0 || y <= 0 {
		form.Center()
	} else {
		form.SetPos(x, y)
	}
	return form
}

/*
 この拡張機能の専用のクラス名を持つパネルを構築します。
*/
func NewPanel(parent winc.Controller) *winc.Panel {
	name := fmt.Sprintf("zylo_%s_panel", reiwa.PluginName)
	winc.RegClassOnlyOnce(name)
	pane := new(winc.Panel)
	pane.SetParent(parent)
	ex := uint(w32.WS_EX_CONTROLPARENT)
	st := uint(w32.WS_CHILD | w32.WS_VISIBLE)
	pane.SetHandle(winc.CreateWindow(name, parent, ex, st))
	winc.RegMsgHandler(pane)
	pane.SetFont(winc.DefaultFont)
	return pane
}
