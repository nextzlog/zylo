/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
*/
package main

/*
typedef void (*Insert)(void*);
typedef void (*Delete)(void*);
typedef void (*Filter)(char*);

inline void doInsert(void *p, Insert f) {
	f(p);
}

inline void doDelete(void *p, Delete f) {
	f(p);
}

inline void doFilter(char *t, Filter f) {
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

//export zylo_to_zlog_launch
func zylo_to_zlog_launch() {
	defer zylo.CapturePanic()
	qxsl, _ = zylo.NewQxSL(qbin)
	zlaunch()
}

//export zylo_to_zlog_finish
func zylo_to_zlog_finish() {
	defer zylo.CapturePanic()
	if qxsl != nil {
		defer qxsl.Close()
	}
	zfinish()
}

//export zlog_to_zylo_insert
func zlog_to_zylo_insert(fun C.Insert) {
	defer zylo.CapturePanic()
	zylo.InsertQSO = func(qso *zylo.QSO) {
		C.doInsert(unsafe.Pointer(qso), fun)
	}
}

//export zlog_to_zylo_delete
func zlog_to_zylo_delete(fun C.Delete) {
	defer zylo.CapturePanic()
	zylo.DeleteQSO = func(qso *zylo.QSO) {
		C.doDelete(unsafe.Pointer(qso), fun)
	}
}

//import zlog_to_zylo_filter
func zlog_to_zylo_filter(fun C.Filter) {
	defer zylo.CapturePanic()
	if qxsl != nil {
		ex, _ := qxsl.Filter()
		value := C.CString(ex)
		C.doFilter(value, fun)
	}
}

//export zylo_to_zlog_import
func zylo_to_zlog_import(source, target *C.char) (err bool) {
	defer zylo.CapturePanic()
	f := "zbin"
	s := C.GoString(source)
	d := C.GoString(target)
	err = qxsl != nil && qxsl.Format(s, d, f) != nil
	return
}

//export zylo_to_zlog_export
func zylo_to_zlog_export(source, format *C.char) (err bool) {
	defer zylo.CapturePanic()
	s := C.GoString(source)
	f := C.GoString(format)
	err = qxsl != nil && qxsl.Format(s, s, f) != nil
	return
}

//export zylo_to_zlog_attach
func zylo_to_zlog_attach(test, path *C.char) {
	defer zylo.CapturePanic()
	t := C.GoString(test)
	c := C.GoString(path)
	zattach(t, c)
}

//export zylo_to_zlog_detach
func zylo_to_zlog_detach() {
	defer zylo.CapturePanic()
	zdetach()
}

//export zylo_to_zlog_verify
func zylo_to_zlog_verify(ptr uintptr, size int) (score int) {
	defer zylo.CapturePanic()
	score = zverify(zylo.ToLog(ptr, size))
	return
}

//export zylo_to_zlog_update
func zylo_to_zlog_update(ptr uintptr, size int) (total int) {
	defer zylo.CapturePanic()
	total = zupdate(zylo.ToLog(ptr, size))
	return
}

//export zylo_to_zlog_insert
func zylo_to_zlog_insert(ptr uintptr, size int) {
	defer zylo.CapturePanic()
	zinsert(zylo.ToLog(ptr, size))
}

//export zylo_to_zlog_delete
func zylo_to_zlog_delete(ptr uintptr, size int) {
	defer zylo.CapturePanic()
	zdelete(zylo.ToLog(ptr, size))
}

//export zylo_to_zlog_kpress
func zylo_to_zlog_kpress(key int, source *C.char) (block bool) {
	defer zylo.CapturePanic()
	block = zkpress(key, C.GoString(source))
	return
}

//export zylo_to_zlog_fclick
func zylo_to_zlog_fclick(btn int, source *C.char) (block bool) {
	defer zylo.CapturePanic()
	block = zfclick(btn, C.GoString(source))
	return
}
