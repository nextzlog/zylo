/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22
 * License: The MIT License since 2021 October 28 (see LICENSE)
 * Author: Journal of Hamradio Informatics (http://pafelog.net)
*******************************************************************************/
package main

/*
#include <stdlib.h>
typedef void (*InsertCB)(void*);
typedef void (*DeleteCB)(void*);
typedef void (*UpdateCB)(void*);
typedef void (*DialogCB)(char*);
typedef void (*NotifyCB)(char*);
typedef void (*AccessCB)(void*);
typedef long (*HandleCB)(char*);
typedef long (*ButtonCB)(char*);
typedef long (*EditorCB)(char*);
typedef long (*ScriptCB)(char*);
typedef void (*FormatCB)(char*);
typedef void (*CitiesCB)(char*);

inline void insert(void *qso, InsertCB cb) {
	if(cb) cb(qso);
}

inline void delete(void *qso, DeleteCB cb) {
	if(cb) cb(qso);
}

inline void update(void *qso, UpdateCB cb) {
	if(cb) cb(qso);
}

inline void dialog(char *str, DialogCB cb) {
	if(cb) cb(str);
}

inline void notify(char *str, NotifyCB cb) {
	if(cb) cb(str);
}

inline void access(void *str, AccessCB cb) {
	if(cb) cb(str);
}

inline long handle(char *str, HandleCB cb) {
	if(cb) return cb(str);
}

inline long button(char *str, ButtonCB cb) {
	if(cb) return cb(str);
}

inline long editor(char *str, EditorCB cb) {
	if(cb) return cb(str);
}

inline long script(char *str, ScriptCB cb) {
	if(cb) return cb(str);
}

inline void format(char *str, FormatCB cb) {
	if(cb) cb(str);
}

inline void cities(char *str, CitiesCB cb) {
	if(cb) cb(str);
}
*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gopkg.in/ini.v1"
	"io"
	"math"
	"os"
	"regexp"
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

var insertCB C.InsertCB
var deleteCB C.DeleteCB
var updateCB C.UpdateCB
var dialogCB C.DialogCB
var notifyCB C.NotifyCB
var accessCB C.AccessCB
var handleCB C.HandleCB
var buttonCB C.ButtonCB
var editorCB C.EditorCB
var scriptCB C.ScriptCB

func zylo_count_lost_cb() (cnt int) {
	cnt += zylo_btoi(insertCB == nil)
	cnt += zylo_btoi(deleteCB == nil)
	cnt += zylo_btoi(updateCB == nil)
	cnt += zylo_btoi(dialogCB == nil)
	cnt += zylo_btoi(notifyCB == nil)
	cnt += zylo_btoi(accessCB == nil)
	cnt += zylo_btoi(handleCB == nil)
	cnt += zylo_btoi(buttonCB == nil)
	cnt += zylo_btoi(editorCB == nil)
	cnt += zylo_btoi(scriptCB == nil)
	return
}

func zylo_btoi(v bool) int {
	if v {
		return 1
	} else {
		return 0
	}
}

func main() {}

//export zylo_allow_insert
func zylo_allow_insert(callback C.InsertCB) {
	insertCB = callback
}

//export zylo_allow_delete
func zylo_allow_delete(callback C.DeleteCB) {
	deleteCB = callback
}

//export zylo_allow_update
func zylo_allow_update(callback C.UpdateCB) {
	updateCB = callback
}

//export zylo_allow_dialog
func zylo_allow_dialog(callback C.DialogCB) {
	dialogCB = callback
}

//export zylo_allow_notify
func zylo_allow_notify(callback C.NotifyCB) {
	notifyCB = callback
}

//export zylo_allow_access
func zylo_allow_access(callback C.AccessCB) {
	accessCB = callback
}

//export zylo_allow_handle
func zylo_allow_handle(callback C.HandleCB) {
	handleCB = callback
}

//export zylo_allow_button
func zylo_allow_button(callback C.ButtonCB) {
	buttonCB = callback
}

//export zylo_allow_editor
func zylo_allow_editor(callback C.EditorCB) {
	editorCB = callback
}

//export zylo_allow_script
func zylo_allow_script(callback C.ScriptCB) {
	scriptCB = callback
}

//export zylo_query_format
func zylo_query_format(callback C.FormatCB) {
	defer zylo_recover_capture_panic()
	f := C.CString(FileExtFilter)
	defer C.free(unsafe.Pointer(f))
	C.format(f, callback)
}

//export zylo_query_cities
func zylo_query_cities(callback C.CitiesCB) {
	defer zylo_recover_capture_panic()
	c := C.CString(CityMultiList)
	defer C.free(unsafe.Pointer(c))
	C.cities(c, callback)
}

//export zylo_launch_event
func zylo_launch_event() bool {
	defer zylo_recover_capture_panic()
	if zylo_count_lost_cb() == 0 {
		OnLaunchEvent()
		return true
	} else {
		return false
	}
}

//export zylo_finish_event
func zylo_finish_event() bool {
	defer zylo_recover_capture_panic()
	OnFinishEvent()
	return false
}

//export zylo_window_event
func zylo_window_event(msg uintptr) {
	defer zylo_recover_capture_panic()
	OnWindowEvent(msg)
}

//export zylo_import_event
func zylo_import_event(source, target *C.char) bool {
	defer zylo_recover_capture_panic()
	src := C.GoString(source)
	tgt := C.GoString(target)
	return OnImportEvent(src, tgt) == nil
}

//export zylo_export_event
func zylo_export_event(target, format *C.char) bool {
	defer zylo_recover_capture_panic()
	tgt := C.GoString(target)
	fmt := C.GoString(format)
	return OnExportEvent(tgt, fmt) == nil
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
	bin, _ := os.ReadFile(Query("{F}"))
	for _, qso := range LoadZLO(bin) {
		OnInsertEvent(&qso)
	}
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
	return int(C.button(n, buttonCB))
}

func zylo_add_editor_handler(name string) (evID int) {
	n := C.CString(name)
	defer C.free(unsafe.Pointer(n))
	return int(C.editor(n, editorCB))
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
	id    int32
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
	CW = iota
	SSB
	FM
	AM
	RTTY
	OTHER
)

/*
 QSO構造体の周波数帯の列挙子です。
*/
const (
	K1900 = iota
	K3500
	M7
	M10
	M14
	M18
	M21
	M24
	M28
	M50
	M144
	M430
	M1200
	M2400
	M5600
	G10UP
)

/*
 交信を識別する番号を返します。
*/
func (qso *QSO) GetID() int32 {
	return qso.id / 100
}

/*
 交信が行われた時刻を返します。
*/
func (qso *QSO) GetTime() time.Time {
	t := math.Abs(qso.time)
	h := time.Duration((t - float64(int(t))) * 24)
	d := time.Date(1899, 12, 30, 0, 0, 0, 0, zone)
	return d.Add(h*time.Hour).AddDate(0, 0, int(t))
}

/*
 呼出符号のポータブル表記を除く部分を返します。
*/
func (qso *QSO) GetCallSign() string {
	return strings.Split(qso.GetCall(), "/")[0]
}

/*
 交信相手の呼出符号を返します。
*/
func (qso *QSO) GetCall() string {
	return zylo_decode_string(qso.call[:])
}

/*
 送信したコンテストナンバーを返します。
*/
func (qso *QSO) GetSent() string {
	return zylo_decode_string(qso.sent[:])
}

/*
 受信したコンテストナンバーを返します。
*/
func (qso *QSO) GetRcvd() string {
	return zylo_decode_string(qso.rcvd[:])
}

/*
 運用者名を返します。
*/
func (qso *QSO) GetName() string {
	return zylo_decode_string(qso.name[:])
}

/*
 備考を返します。
*/
func (qso *QSO) GetNote() string {
	return zylo_decode_string(qso.note[:])
}

/*
 主マルチプライヤを返します。
*/
func (qso *QSO) GetMul1() string {
	return zylo_decode_string(qso.mul1[:])
}

/*
 マルチプライヤを返します。
*/
func (qso *QSO) GetMul2() string {
	return zylo_decode_string(qso.mul2[:])
}

/*
 副マルチプライヤを返します。
*/
func (qso *QSO) SetCall(value string) {
	copy(qso.call[:], zylo_encode_string(value))
}

/*
 送信したコンテストナンバーを設定します。
*/
func (qso *QSO) SetSent(value string) {
	copy(qso.sent[:], zylo_encode_string(value))
}

/*
 受信したコンテストナンバーを設定します。
*/
func (qso *QSO) SetRcvd(value string) {
	copy(qso.rcvd[:], zylo_encode_string(value))
}

/*
 運用者名を設定します。
*/
func (qso *QSO) SetName(value string) {
	copy(qso.name[:], zylo_encode_string(value))
}

/*
 備考を設定します。
*/
func (qso *QSO) SetNote(value string) {
	copy(qso.note[:], zylo_encode_string(value))
}

/*
 主マルチプライヤを設定します。
*/
func (qso *QSO) SetMul1(value string) {
	copy(qso.mul1[:], zylo_encode_string(value))
}

/*
 副マルチプライヤを設定します。
*/
func (qso *QSO) SetMul2(value string) {
	copy(qso.mul2[:], zylo_encode_string(value))
}

func zylo_decode_string(f []byte) string {
	v := string(f[1 : int(f[0])+1])
	return strings.TrimSpace(v)
}

func zylo_encode_string(v string) []byte {
	return append([]byte{byte(len(v))}, v...)
}

/*
 重複交信を検査します。
 有効な交信の場合は真を返します。
*/
func (qso *QSO) VerifyDupe() bool {
	return !qso.Dupe || zylo_dupes
}

/*
 周波数帯を検査します。
 有効な交信の場合は真を返します。
*/
func (qso *QSO) VerifyBand() bool {
	return zylo_bands[qso.Band]
}

/*
 通信方式を検査します。
 有効な交信の場合は真を返します。
*/
func (qso *QSO) VerifyMode() bool {
	return zylo_modes[qso.Mode]
}

/*
 交信相手の呼出符号を検査します。
 有効な交信の場合は真を返します。
*/
func (qso *QSO) VerifyCall() bool {
	return zylo_calls.MatchString(qso.GetCall())
}

/*
 コンテストナンバーを検査します。
 有効な交信の場合は真を返します。
*/
func (qso *QSO) VerifyRcvd() bool {
	return zylo_rcvds.MatchString(qso.GetRcvd())
}

/*
 コンテストナンバーを検査します。
 有効な交信の場合は真を返します。
*/
func (qso *QSO) VerifySent() bool {
	return zylo_sents.MatchString(qso.GetSent())
}

/*
 コンテストナンバーを正規表現でグループ化します。
*/
func (qso *QSO) GetRcvdGroups() []string {
	return zylo_rcvds.FindStringSubmatch(qso.GetRcvd())
}

/*
 コンテストナンバーを正規表現でグループ化します。
*/
func (qso *QSO) GetSentGroups() []string {
	return zylo_sents.FindStringSubmatch(qso.GetSent())
}

/*
 QSOを無効な交信とします。
*/
func (qso *QSO) Invalid() {
	qso.Score = 0
	qso.SetMul1("")
	qso.SetMul2("")
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
 指定された交信記録をzLog側に追加します。
*/
func (qso *QSO) Insert() {
	C.insert(unsafe.Pointer(qso), insertCB)
}

/*
 指定された交信記録をzLog側で削除します。
*/
func (qso *QSO) Delete() {
	C.delete(unsafe.Pointer(qso), deleteCB)
}

/*
 指定されたzLog側の交信記録を更新します。
*/
func (qso *QSO) Update() {
	C.update(unsafe.Pointer(qso), updateCB)
}

/*
 指定された設定を取得します。
*/
func GetINI(section, key string) string {
	init, _ := ini.LooseLoad(SettingsFileName)
	return init.Section(section).Key(key).String()
}

/*
 指定された設定を保存します。
*/
func SetINI(section, key, value string) {
	init, _ := ini.LooseLoad(SettingsFileName)
	init.Section(section).Key(key).SetValue(value)
	init.SaveTo(SettingsFileName)
}

/*
 指定された文字列をダイアログで表示します。
*/
func DisplayModal(msg string, args ...interface{}) {
	text := C.CString(fmt.Sprintf(msg, args...))
	defer C.free(unsafe.Pointer(text))
	C.dialog(text, dialogCB)
}

/*
 指定された文字列を通知バナーに表示します。
*/
func DisplayToast(msg string, args ...interface{}) {
	text := C.CString(fmt.Sprintf(msg, args...))
	defer C.free(unsafe.Pointer(text))
	C.notify(text, notifyCB)
}

/*
 指定されたクエリを問い合わせます。
*/
func Query(text string) string {
	buf := make([]byte, ResponseCapacity+1)
	copy(buf[:ResponseCapacity], text[:])
	C.access(unsafe.Pointer(&buf[0]), accessCB)
	return string(buf[:bytes.IndexByte(buf, 0)])
}

/*
 指定されたウィンドウハンドルを取得します。
 例:

  GetUI("MainForm.FileOpenItem")
  GetUI("MenuForm.CancelButton")
*/
func GetUI(expression string) uintptr {
	e := C.CString(expression)
	defer C.free(unsafe.Pointer(e))
	return uintptr(C.handle(e, handleCB))
}

/*
 指定されたスクリプトを実行します。
*/
func RunDelphi(exp string, args ...interface{}) int {
	e := C.CString(fmt.Sprintf(exp, args...))
	defer C.free(unsafe.Pointer(e))
	return int(C.script(e, scriptCB))
}

var zylo_dupes = false
var zylo_bands = make(map[byte]bool)
var zylo_modes = make(map[byte]bool)
var zylo_calls = regexp.MustCompile(`^.+$`)
var zylo_rcvds = regexp.MustCompile(`^.+$`)
var zylo_sents = regexp.MustCompile(`^.+$`)

/*
 重複交信を許容します。
*/
func AllowDupe() {
	zylo_dupes = true
}

/*
 指定された周波数帯を許容します。
*/
func AllowBand(bands ...byte) {
	for _, band := range bands {
		zylo_bands[band] = true
	}
}

/*
 指定された通信方式を許容します。
*/
func AllowMode(modes ...byte) {
	for _, mode := range modes {
		zylo_modes[mode] = true
	}
}

/*
 指定された周波数帯を許容します。
*/
func AllowBandRange(lo, hi byte) {
	for b := lo; b <= hi; b++ {
		AllowBand(b)
	}
}

/*
 指定された通信方式を許容します。
*/
func AllowModeRange(lo, hi byte) {
	for m := lo; m <= hi; m++ {
		AllowMode(m)
	}
}

/*
 交信相手の呼出符号の正規表現を設定します。
*/
func AllowCall(pattern string) {
	zylo_calls = regexp.MustCompile(pattern)
}

/*
 コンテストナンバーの正規表現を設定します。
*/
func AllowRcvd(pattern string) {
	zylo_rcvds = regexp.MustCompile(pattern)
}

/*
 コンテストナンバーの正規表現を設定します。
*/
func AllowSent(pattern string) {
	zylo_sents = regexp.MustCompile(pattern)
}

/*
 I/O拡張機能が利用する専用の変数です。
*/
var FileExtFilter string

/*
 DATファイルを内蔵するための変数です。
*/
var CityMultiList string

/*
 拡張機能の起動時に呼び出されます。
*/
var OnLaunchEvent = func() {}

/*
 拡張機能の終了時に呼び出されます。
*/
var OnFinishEvent = func() {}

/*
 ウィンドウメッセージを受信します。
*/
var OnWindowEvent = func(msg uintptr) {}

/*
 交信記録をZLOファイルに変換する要求を処理します。
*/
var OnImportEvent = func(source, target string) error {
	return nil
}

/*
 ZLOファイルを他の書式に変換する要求を処理します。
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
 無効な交信の場合は、マルチプライヤを空の文字列にします。
*/
var OnVerifyEvent = func(qso *QSO) {
	ok := true
	ok = ok && qso.VerifyDupe()
	ok = ok && qso.VerifyBand()
	ok = ok && qso.VerifyMode()
	ok = ok && qso.VerifyRcvd()
	ok = ok && qso.VerifySent()
	if ok {
		OnAcceptEvent(qso)
	} else {
		qso.Invalid()
	}
}

/*
 有効と判定した交信に得点やマルチプライヤを設定します。
*/
var OnAcceptEvent = func(qso *QSO) {
	qso.SetMul1(qso.GetRcvd())
}

/*
 総合得点を計算します。
 引数は交信の合計得点と主マルチプライヤの異なり数です。
*/
var OnPointsEvent = func(score, mul1s int) int {
	return score * mul1s
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
