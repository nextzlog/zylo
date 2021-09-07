/*
 zLogの拡張機能を開発するためのフレームワークです。
*/
package zylo

/*
#include <stdlib.h>
#include "zutils.h"
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"gopkg.in/go-toast/toast.v1"
	"gopkg.in/ini.v1"
	"io"
	"math"
	"runtime/debug"
	"strings"
	"time"
	"unsafe"
)

/*
 問合わせの返り値の長さの限度です。
*/
const ResponseCapacity = 256

/*
 設定を保管するファイルの名前です。
*/
const SettingsFileName = "zlog.ini"

var zone = time.Local

// event handlers
var buttons = make(map[int]func(int))
var editors = make(map[int]func(int))

func main() {}

//export zylo_allow_insert
func zylo_allow_insert(callback C.InsertCallBack) {
	C.insertCallBack = callback
}

//export zylo_allow_delete
func zylo_allow_delete(callback C.DeleteCallBack) {
	C.deleteCallBack = callback
}

//export zylo_allow_update
func zylo_allow_update(callback C.UpdateCallBack) {
	C.updateCallBack = callback
}

//export zylo_allow_dialog
func zylo_allow_dialog(callback C.DialogCallBack) {
	C.dialogCallBack = callback
}

//export zylo_allow_access
func zylo_allow_access(callback C.AccessCallBack) {
	C.accessCallBack = callback
}

//export zylo_allow_button
func zylo_allow_button(callback C.ButtonCallBack) {
	C.buttonCallBack = callback
}

//export zylo_allow_editor
func zylo_allow_editor(callback C.EditorCallBack) {
	C.editorCallBack = callback
}

//export zylo_query_format
func zylo_query_format(callback C.FormatCallBack) {
	defer zylo_recover_capture_panic()
	f := C.CString(FileExtFilter)
	defer C.free(unsafe.Pointer(f))
	C.callFormat(f, callback)
}

//export zylo_query_cities
func zylo_query_cities(callback C.CitiesCallBack) {
	defer zylo_recover_capture_panic()
	c := C.CString(CityMultiList)
	defer C.free(unsafe.Pointer(c))
	C.callCities(c, callback)
}

//export zylo_launch_event
func zylo_launch_event() {
	defer zylo_recover_capture_panic()
	OnLaunchEvent()
}

//export zylo_finish_event
func zylo_finish_event() {
	defer zylo_recover_capture_panic()
	OnFinishEvent()
}

//export zylo_import_event
func zylo_import_event(source, target *C.char) {
	defer zylo_recover_capture_panic()
	src := C.GoString(source)
	tgt := C.GoString(target)
	if OnImportEvent(src, tgt) != nil {
		DisplayModal("failed to load %s", src)
	}
}

//export zylo_export_event
func zylo_export_event(target, format *C.char) {
	defer zylo_recover_capture_panic()
	tgt := C.GoString(target)
	fmt := C.GoString(format)
	if OnExportEvent(tgt, fmt) != nil {
		DisplayModal("failed to save %s", tgt)
	}
}

//export zylo_offset_event
func zylo_offset_event(offset int) {
	zone = time.FixedZone("", -offset*60)
}

//export zylo_attach_event
func zylo_attach_event(test, path *C.char) {
	defer zylo_recover_capture_panic()
	t := C.GoString(test)
	c := C.GoString(path)
	OnAttachEvent(t, c)
}

//export zylo_assign_event
func zylo_assign_event(test, path *C.char) {
	defer zylo_recover_capture_panic()
	t := C.GoString(test)
	c := C.GoString(path)
	OnAssignEvent(t, c)
}

//export zylo_detach_event
func zylo_detach_event(test, path *C.char) {
	defer zylo_recover_capture_panic()
	t := C.GoString(test)
	c := C.GoString(path)
	OnDetachEvent(t, c)
}

//export zylo_insert_event
func zylo_insert_event(ptr uintptr) {
	defer zylo_recover_capture_panic()
	OnInsertEvent((*QSO)(unsafe.Pointer(ptr)))
}

//export zylo_delete_event
func zylo_delete_event(ptr uintptr) {
	defer zylo_recover_capture_panic()
	OnDeleteEvent((*QSO)(unsafe.Pointer(ptr)))
}

//export zylo_verify_event
func zylo_verify_event(ptr uintptr) {
	defer zylo_recover_capture_panic()
	OnVerifyEvent((*QSO)(unsafe.Pointer(ptr)))
}

//export zylo_points_event
func zylo_points_event(pts, muls int) int {
	defer zylo_recover_capture_panic()
	return OnPointsEvent(pts, muls)
}

//export zylo_button_event
func zylo_button_event(comp, btn int) {
	defer zylo_recover_capture_panic()
	if h, ok := buttons[comp]; ok {
		h(btn)
	}
}

//export zylo_editor_event
func zylo_editor_event(comp, key int) {
	defer zylo_recover_capture_panic()
	if h, ok := editors[comp]; ok {
		h(key)
	}
}

func zylo_recover_capture_panic() {
	if err := recover(); err != nil {
		DisplayModal(string(debug.Stack()))
	}
}

func zylo_add_button_handler(name string) (evID int) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	return int(C.callButton(n))
}

func zylo_add_editor_handler(name string) (evID int) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	return int(C.callEditor(n))
}

/*
 QSO構造体です。
*/
type QSO struct {
	time  float64
	call  [13]byte
	sent  [31]byte
	rcvd  [31]byte
	void  byte
	sRST  int16
	rRST  int16
	ID    int32
	Mode  byte
	Band  byte
	Pow1  byte
	mul1  [31]byte
	mul2  [31]byte
	new1  bool
	new2  bool
	Score byte
	name  [15]byte
	note  [65]byte
	isCQ  bool
	Dupe  bool
	rsv1  byte
	TxID  byte
	Pow2  int32
	rsv2  int32
	rsv3  int32
}

/*
 QSO構造体の通信方式の列挙子です。
*/
const (
	CW    = 0
	SSB   = 1
	FM    = 2
	AM    = 3
	RTTY  = 4
	OTHER = 5
)

/*
 QSO構造体の周波数帯の列挙子です。
*/
const (
	M1_9  = 0
	M3_5  = 1
	M7    = 2
	M10   = 3
	M14   = 4
	M18   = 5
	M21   = 6
	M24   = 7
	M28   = 8
	M50   = 9
	M144  = 10
	M430  = 11
	M1200 = 12
	M2400 = 13
	M5600 = 14
	G10UP = 15
)

/*
 交信時刻を返します。
*/
func (qso *QSO) GetTime() time.Time {
	t := math.Abs(qso.time)
	h := time.Duration((t - float64(int(t))) * 24)
	d := time.Date(1899, 12, 30, 0, 0, 0, 0, zone)
	return d.Add(h*time.Hour).AddDate(0, 0, int(t))
}

/*
 呼出符号を返します。
*/
func (qso *QSO) GetCall() string {
	return stringFromDtoC(qso.call[:])
}

/*
 送信した番号を返します。
*/
func (qso *QSO) GetSent() string {
	return stringFromDtoC(qso.sent[:])
}

/*
 受信した番号を返します。
*/
func (qso *QSO) GetRcvd() string {
	return stringFromDtoC(qso.rcvd[:])
}

/*
 運用者名を返します。
*/
func (qso *QSO) GetName() string {
	return stringFromDtoC(qso.name[:])
}

/*
 備考を返します。
*/
func (qso *QSO) GetNote() string {
	return stringFromDtoC(qso.note[:])
}

/*
 第1マルチプライヤを返します。
*/
func (qso *QSO) GetMul1() string {
	return stringFromDtoC(qso.mul1[:])
}

/*
 第2マルチプライヤを返します。
*/
func (qso *QSO) GetMul2() string {
	return stringFromDtoC(qso.mul2[:])
}

/*
 第2マルチプライヤを返します。
*/
func (qso *QSO) SetCall(value string) {
	copy(qso.call[:], stringFromCtoD(value))
}

/*
 送信した番号を設定します。
*/
func (qso *QSO) SetSent(value string) {
	copy(qso.sent[:], stringFromCtoD(value))
}

/*
 受信した番号を設定します。
*/
func (qso *QSO) SetRcvd(value string) {
	copy(qso.rcvd[:], stringFromCtoD(value))
}

/*
 運用者名を設定します。
*/
func (qso *QSO) SetName(value string) {
	copy(qso.name[:], stringFromCtoD(value))
}

/*
 備考を設定します。
*/
func (qso *QSO) SetNote(value string) {
	copy(qso.note[:], stringFromCtoD(value))
}

/*
 第1マルチプライヤを設定します。
*/
func (qso *QSO) SetMul1(value string) {
	copy(qso.mul1[:], stringFromCtoD(value))
}

/*
 第2マルチプライヤを設定します。
*/
func (qso *QSO) SetMul2(value string) {
	copy(qso.mul2[:], stringFromCtoD(value))
}

func stringFromDtoC(f []byte) string {
	v := string(f[1 : int(f[0])+1])
	return strings.TrimSpace(v)
}

func stringFromCtoD(v string) []byte {
	return append([]byte{byte(len(v))}, v...)
}

/*
 QSO構造体をヘッダ情報なしで書き込みます。
*/
func (qso *QSO) DumpWithoutHead(w io.Writer) {
	binary.Write(w, binary.LittleEndian, qso)
}

/*
 QSO構造体をヘッダ情報なしで読み取ります。
*/
func (qso *QSO) LoadWithoutHead(r io.Reader) {
	raw := make([]byte, 256)
	binary.Read(r, binary.LittleEndian, raw)
	*qso = *(*QSO)(unsafe.Pointer(&raw[0]))
}

/*
 QSO列をヘッダ情報付きのバイト列に変換します。
*/
func DumpZLO(qso ...QSO) (bin []byte) {
	_, diff := time.Now().In(zone).Zone()
	log := append(make([]QSO, 1), qso...)
	log[0].sRST = int16(-diff / 60)
	buf := new(bytes.Buffer)
	for _, qso := range log {
		qso.DumpWithoutHead(buf)
	}
	return buf.Bytes()
}

/*
 ヘッダ情報付きのバイト列をQSO列に変換します。
*/
func LoadZLO(bin []byte) (logs []QSO) {
	buf := bytes.NewReader(bin)
	new(QSO).LoadWithoutHead(buf)
	for buf.Len() > 0 {
		qso := QSO{}
		qso.LoadWithoutHead(buf)
		logs = append(logs, qso)
	}
	return logs
}

/*
 指定された交信記録を追加します。
*/
func (qso *QSO) Insert() {
	C.callInsert(unsafe.Pointer(qso))
}

/*
 指定された交信記録を削除します。
*/
func (qso *QSO) Delete() {
	C.callDelete(unsafe.Pointer(qso))
}

/*
 指定された交信記録を更新します。
*/
func (qso *QSO) Update() {
	C.callUpdate(unsafe.Pointer(qso))
}

/*
 指定された設定の内容を取得します。
*/
func GetINI(section, key string) string {
	init, _ := ini.LooseLoad(SettingsFileName)
	return init.Section(section).Key(key).String()
}

/*
 指定された設定の内容を変更します。
*/
func SetINI(section, key, value string) {
	init, _ := ini.LooseLoad(SettingsFileName)
	init.Section(section).Key(key).SetValue(value)
	init.SaveTo(SettingsFileName)
}

/*
 指定された文字列をSJISに変換します。
*/
func UnicodeToShiftJIS(utf string) (string, error) {
	return japanese.ShiftJIS.NewEncoder().String(utf)
}

/*
 指定された文字列を通知欄に表示します。
*/
func DisplayToast(msg string, args ...interface{}) {
	msg = fmt.Sprintf(msg, args...)
	msg, _ = UnicodeToShiftJIS(msg)
	toast := toast.Notification{
		AppID:   "zLog",
		Title:   "ZyLO",
		Message: msg,
	}
	toast.Push()
}

/*
 指定された文字列を対話的に表示します。
*/
func DisplayModal(msg string, args ...interface{}) {
	text := C.CString(fmt.Sprintf(msg, args...))
	defer C.free(unsafe.Pointer(text))
	C.callDialog(text)
}

/*
 指定されたクエリで問合わせを行います。
*/
func Query(text string) string {
	buf := make([]byte, ResponseCapacity+1)
	copy(buf[:ResponseCapacity], text[:])
	C.callAccess(unsafe.Pointer(&buf[0]))
	return string(buf[:bytes.IndexByte(buf, 0)])
}

/*
 対応済みの書式の名称と拡張子のリストを設定し、
 インポート及びエクスポート機能を有効化します。
*/
var FileExtFilter string

/*
 市区町村や国や地域の番号のリストを指定します。
*/
var CityMultiList string

/*
 zLogの起動時に呼び出されます。
*/
var OnLaunchEvent = func() {}

/*
 zLogの終了時に呼び出されます。
*/
var OnFinishEvent = func() {}

/*
 交信記録をzLogでインポート可能な書式に変換します。
*/
var OnImportEvent = func(source, target string) error {
	return nil
}

/*
 zLogがエクスポートした交信記録の書式を変換します。
*/
var OnExportEvent = func(source, format string) error {
	return nil
}

/*
 コンテストを開いた直後に呼び出されます。
*/
var OnAttachEvent = func(contest, configs string) {}

/*
 得点計算の権限が移譲された場合に呼び出されます。
*/
var OnAssignEvent = func(contest, configs string) {}

/*
 コンテストを閉じた直後に呼び出されます。
*/
var OnDetachEvent = func(contest, configs string) {}

/*
 交信記録が追加された際に呼び出されます。
 修正時はまず削除、次に追加が行われます。
*/
var OnInsertEvent = func(qso *QSO) {}

/*
 交信記録が削除された際に呼び出されます。
 修正時はまず削除、次に追加が行われます。
*/
var OnDeleteEvent = func(qso *QSO) {}

/*
 交信の得点やマルチプライヤを検査する時に呼び出されます。
 編集中の交信記録に対し、必要なら何度でも呼び出されます。
*/
var OnVerifyEvent = func(qso *QSO) {
	rcvd := qso.GetRcvd()
	qso.SetMul1(rcvd)
	if qso.Dupe {
		qso.Score = 0
	} else {
		qso.Score = 1
	}
}

/*
 総得点を計算します。
 引数は交信の合計得点と第1マルチプライヤの異なり数です。
*/
var OnPointsEvent = func(score, mults int) int {
	return score * mults
}

/*
 指定された名前のボタンにイベントハンドラを登録します。
 起動時のみ登録できます。それ以後の登録は無視されます。
*/
func HandleButton(name string, handler func(int)) {
	buttons[zylo_add_button_handler(name)] = handler
}

/*
 指定された名前の記入欄にイベントハンドラを登録します。
 起動時のみ登録できます。それ以後の登録は無視されます。
*/
func HandleEditor(name string, handler func(int)) {
	editors[zylo_add_editor_handler(name)] = handler
}
