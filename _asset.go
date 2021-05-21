/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/
package main

/*
#include <stdlib.h>

typedef void (*Insert)(void*);
typedef void (*Delete)(void*);
typedef void (*Update)(void*);
typedef void (*Filter)(char*);
typedef void (*Cities)(char*);
typedef void (*Editor)(char*);
typedef void (*Button)(char*);

inline void doInsert(void *p, Insert f) {
	f(p);
}

inline void doDelete(void *p, Delete f) {
	f(p);
}

inline void doUpdate(void *p, Update f) {
	f(p);
}

inline void doFilter(char *t, Filter f) {
	f(t);
}

inline void doCities(char *t, Cities f) {
	f(t);
}

inline void doEditor(char *t, Editor f) {
	f(t);
}

inline void doButton(char *t, Button f) {
	f(t);
}
*/
import "C"
import (
	_ "embed"
	"unsafe"
	"github.com/nextzlog/zylo"
)

//go:embed qxsl.exe
var qbin []byte

/*
 an embedded QxSL library.
 */
var qxsl *zylo.QxSL

func free(value *C.char) {
	C.free(unsafe.Pointer(value))
}

//export zylo_handle_launch
func zylo_handle_launch() {
	defer zylo.CapturePanic()
	zlaunch()
}

//export zylo_handle_finish
func zylo_handle_finish() {
	defer zylo.CapturePanic()
	if qxsl != nil {
		defer qxsl.Close()
	}
	zfinish()
}

//export zylo_permit_insert
func zylo_permit_insert(fun C.Insert) {
	defer zylo.CapturePanic()
	zylo.InsertQSO = func(qso *zylo.QSO) {
		C.doInsert(unsafe.Pointer(qso), fun)
	}
}

//export zylo_permit_delete
func zylo_permit_delete(fun C.Delete) {
	defer zylo.CapturePanic()
	zylo.DeleteQSO = func(qso *zylo.QSO) {
		C.doDelete(unsafe.Pointer(qso), fun)
	}
}

//export zylo_permit_update
func zylo_permit_update(fun C.Update) {
	defer zylo.CapturePanic()
	zylo.UpdateQSO = func(qso *zylo.QSO) {
		C.doUpdate(unsafe.Pointer(qso), fun)
	}
}

//export zylo_permit_filter
func zylo_permit_filter(fun C.Filter) {
	defer zylo.CapturePanic()
	qxsl, _ = zylo.NewQxSL(qbin)
	if qxsl != nil {
		ex, _ := qxsl.Filter()
		value := C.CString(ex)
		defer free(value)
		C.doFilter(value, fun)
	}
}

//export zylo_permit_cities
func zylo_permit_cities(fun C.Cities) {
	defer zylo.CapturePanic()
	value := C.CString(zcities())
	defer free(value)
	C.doCities(value, fun)
}

//export zylo_permit_editor
func zylo_permit_editor(fun C.Editor) {
	defer zylo.CapturePanic()
	zylo.HookEditor = func(name string) {
		value := C.CString(name)
		defer free(value)
		C.doEditor(value, fun)
	}
}

//export zylo_permit_button
func zylo_permit_button(fun C.Button) {
	defer zylo.CapturePanic()
	zylo.HookButton = func(name string) {
		value := C.CString(name)
		defer free(value)
		C.doButton(value, fun)
	}
}

//export zylo_handle_import
func zylo_handle_import(source, target *C.char) bool {
	defer zylo.CapturePanic()
	f := "zbin"
	s := C.GoString(source)
	d := C.GoString(target)
	return qxsl != nil && qxsl.Format(s, d, f) != nil
}

//export zylo_handle_export
func zylo_handle_export(source, format *C.char) bool {
	defer zylo.CapturePanic()
	s := C.GoString(source)
	f := C.GoString(format)
	return qxsl != nil && qxsl.Format(s, s, f) != nil
}

//export zylo_handle_attach
func zylo_handle_attach(test, path *C.char) {
	defer zylo.CapturePanic()
	t := C.GoString(test)
	c := C.GoString(path)
	zattach(t, c)
}

//export zylo_handle_detach
func zylo_handle_detach() {
	defer zylo.CapturePanic()
	zdetach()
}

//export zylo_handle_insert
func zylo_handle_insert(ptr uintptr) {
	defer zylo.CapturePanic()
	zinsert(zylo.ToQSO(ptr))
}

//export zylo_handle_delete
func zylo_handle_delete(ptr uintptr) {
	defer zylo.CapturePanic()
	zdelete(zylo.ToQSO(ptr))
}

//export zylo_handle_verify
func zylo_handle_verify(ptr uintptr) {
	defer zylo.CapturePanic()
	zverify(zylo.ToQSO(ptr))
}

//export zylo_handle_points
func zylo_handle_points(score, mults int) int {
	defer zylo.CapturePanic()
	return zpoints(score, mults)
}

//export zylo_handle_editor
func zylo_handle_editor(key int, name *C.char) bool {
	defer zylo.CapturePanic()
	return zeditor(key, C.GoString(name))
}

//export zylo_handle_button
func zylo_handle_button(btn int, name *C.char) bool {
	defer zylo.CapturePanic()
	return zbutton(btn, C.GoString(name))
}
