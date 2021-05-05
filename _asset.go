/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/
package main

/*
typedef char* text;

typedef text (*ShowInputDialog)(text, text);
typedef void (*SetFilterString)(text, text);

inline text doShowInputDialog(text l, text v, ShowInputDialog f) {
	return f(l, v);
}

inline void doSetFilterString(text i, text e, SetFilterString f) {
	f(i, e);
}
*/
import "C"
import (
	_ "embed"
	"github.com/nextzlog/zylo"
)

//go:embed qxsl.exe
var qbin []byte

/*
 an embedded QxSL library.
 */
var qxsl *zylo.QxSL

/*
 a bridge function of Delphi InputBox.
 */
var ibox C.ShowInputDialog

//export _zylo_export_launch
func _zylo_export_launch() {
	defer zylo.CapturePanic()
	qxsl, _ = zylo.NewQxSL(qbin)
	zlaunch()
}

//export _zylo_export_finish
func _zylo_export_finish() {
	defer zylo.CapturePanic()
	if qxsl != nil {
		defer qxsl.Close()
	}
	zfinish()
}

//export _zylo_export_dialog
func _zylo_export_dialog(fun C.ShowInputDialog) {
	defer zylo.CapturePanic()
	zylo.Ibox = func(l, v string) (string, bool){
		lab := C.CString(l)
		val := C.CString(v)
		str := C.doShowInputDialog(lab, val, fun)
		if str != nil {
			return C.GoString(str), true
		} else {
			return "", false
		}
	}
}

//export _zylo_export_filter
func _zylo_export_filter(fun C.SetFilterString) {
	defer zylo.CapturePanic()
	if qxsl != nil {
		ex, _ := qxsl.Filter()
		value := C.CString(ex)
		C.doSetFilterString(value, value, fun)
	}
}

//export _zylo_export_import
func _zylo_export_import(source, target *C.char) (err bool) {
	defer zylo.CapturePanic()
	f := "zbin"
	s := C.GoString(source)
	d := C.GoString(target)
	err = qxsl != nil && qxsl.Format(s, d, f) != nil
	return
}

//export _zylo_export_export
func _zylo_export_export(source, format *C.char) (err bool) {
	defer zylo.CapturePanic()
	s := C.GoString(source)
	f := C.GoString(format)
	err = qxsl != nil && qxsl.Format(s, s, f) != nil
	return
}

//export _zylo_export_attach
func _zylo_export_attach(test, path *C.char) {
	defer zylo.CapturePanic()
	t := C.GoString(test)
	c := C.GoString(path)
	zattach(t, c)
}

//export _zylo_export_detach
func _zylo_export_detach() {
	defer zylo.CapturePanic()
	zdetach()
}

//export _zylo_export_verify
func _zylo_export_verify(ptr uintptr, size int) (score int) {
	defer zylo.CapturePanic()
	score = zverify(zylo.ToLog(ptr, size))
	return
}

//export _zylo_export_update
func _zylo_export_update(ptr uintptr, size int) (total int) {
	defer zylo.CapturePanic()
	total = zupdate(zylo.ToLog(ptr, size))
	return
}

//export _zylo_export_insert
func _zylo_export_insert(ptr uintptr, size int) {
	defer zylo.CapturePanic()
	zinsert(zylo.ToLog(ptr, size))
}

//export _zylo_export_delete
func _zylo_export_delete(ptr uintptr, size int) {
	defer zylo.CapturePanic()
	zdelete(zylo.ToLog(ptr, size))
}

//export _zylo_export_kpress
func _zylo_export_kpress(key int, source *C.char) (block bool) {
	defer zylo.CapturePanic()
	block = zkpress(key, C.GoString(source))
	return
}

//export _zylo_export_fclick
func _zylo_export_fclick(btn int, source *C.char) (block bool) {
	defer zylo.CapturePanic()
	block = zfclick(btn, C.GoString(source))
	return
}
