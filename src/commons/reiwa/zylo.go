/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/

package reiwa

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
	"github.com/hashicorp/go-version"
	"gopkg.in/ini.v1"
	"io"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strings"
	"time"
	"unsafe"
)

/*
Length limit of query response.
*/
const ResponseCapacity = 256

/*
The name of this plugin.
*/
var PluginName = ""

/*
zLog version where this plugin works.
*/
var MinVersion = "2.8.3.0"

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
	defer DisplayPanic()
	f := C.CString(FileExtFilter)
	defer C.free(unsafe.Pointer(f))
	C.format(f, callback)
}

//export zylo_query_cities
func zylo_query_cities(callback C.CitiesCB) {
	defer DisplayPanic()
	c := C.CString(CityMultiList)
	defer C.free(unsafe.Pointer(c))
	C.cities(c, callback)
}

//export zylo_launch_event
func zylo_launch_event() bool {
	defer DisplayPanic()
	zv, _ := version.NewVersion(Query("{V}"))
	mv, _ := version.NewVersion(MinVersion)
	if !zv.LessThan(mv) {
		OnLaunchEvent()
		return true
	} else {
		return false
	}
}

//export zylo_finish_event
func zylo_finish_event() bool {
	defer DisplayPanic()
	OnFinishEvent()
	return false
}

//export zylo_window_event
func zylo_window_event(msg uintptr) {
	defer DisplayPanic()
	OnWindowEvent(msg)
}

//export zylo_import_event
func zylo_import_event(source, target *C.char) bool {
	defer DisplayPanic()
	src := C.GoString(source)
	tgt := C.GoString(target)
	return OnImportEvent(src, tgt) == nil
}

//export zylo_export_event
func zylo_export_event(target, format *C.char) bool {
	defer DisplayPanic()
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
	defer DisplayPanic()
	t := C.GoString(test)
	c := C.GoString(path)
	OnAttachEvent(t, c)
	bin, _ := os.ReadFile(Query("{F}"))
	for _, qso := range LoadZLO(bin) {
		OnInsertEvent(&qso)
	}
}

//export zylo_detach_event
func zylo_detach_event(test, path *C.char) {
	defer DisplayPanic()
	t := C.GoString(test)
	c := C.GoString(path)
	OnDetachEvent(t, c)
}

//export zylo_assign_event
func zylo_assign_event(test, path *C.char) {
	defer DisplayPanic()
	t := C.GoString(test)
	c := C.GoString(path)
	OnAssignEvent(t, c)
}

//export zylo_insert_event
func zylo_insert_event(ptr uintptr) {
	defer DisplayPanic()
	OnInsertEvent((*QSO)(unsafe.Pointer(ptr)))
}

//export zylo_delete_event
func zylo_delete_event(ptr uintptr) {
	defer DisplayPanic()
	OnDeleteEvent((*QSO)(unsafe.Pointer(ptr)))
}

//export zylo_verify_event
func zylo_verify_event(ptr uintptr) {
	defer DisplayPanic()
	OnVerifyEvent((*QSO)(unsafe.Pointer(ptr)))
}

//export zylo_points_event
func zylo_points_event(pts, muls int) int {
	defer DisplayPanic()
	return OnPointsEvent(pts, muls)
}

//export zylo_button_event
func zylo_button_event(comp, btn int) {
	defer DisplayPanic()
	if h, ok := buttons[comp]; ok {
		h(btn)
	}
}

//export zylo_editor_event
func zylo_editor_event(comp, key int) {
	defer DisplayPanic()
	if h, ok := editors[comp]; ok {
		h(key)
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
zLog QSO structure.
*/
type QSO struct {
	time  float64
	call  [13]byte
	sent  [31]byte
	rcvd  [31]byte
	void  byte
	SRST  int16
	RRST  int16
	id    int32
	Mode  byte
	Band  byte
	Pow1  byte
	mul1  [31]byte
	mul2  [31]byte
	New1  bool
	New2  bool
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
Enum of modes.
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
Enum of bands.
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
Returns the QSO ID.
*/
func (qso *QSO) GetID() int32 {
	return qso.id / 100
}

/*
Returns the time of the QSO.
*/
func (qso *QSO) GetTime() time.Time {
	t := math.Abs(qso.time)
	h := time.Duration((t - float64(int(t))) * 24)
	d := time.Date(1899, 12, 30, 0, 0, 0, 0, zone)
	return d.Add(h*time.Hour).AddDate(0, 0, int(t))
}

/*
Returns the canonical call sign,
excluding the portable designator.
*/
func (qso *QSO) GetCallSign() string {
	return strings.Split(qso.GetCall(), "/")[0]
}

/*
Returns the call sign.
*/
func (qso *QSO) GetCall() string {
	return zylo_decode_string(qso.call[:])
}

/*
Returns the submitted contest number.
*/
func (qso *QSO) GetSent() string {
	return zylo_decode_string(qso.sent[:])
}

/*
Returns the received contest number.
*/
func (qso *QSO) GetRcvd() string {
	return zylo_decode_string(qso.rcvd[:])
}

/*
Returns the operator name.
*/
func (qso *QSO) GetName() string {
	return zylo_decode_string(qso.name[:])
}

/*
Returns remarks.
*/
func (qso *QSO) GetNote() string {
	return zylo_decode_string(qso.note[:])
}

/*
Returns the primary multiplier.
*/
func (qso *QSO) GetMul1() string {
	return zylo_decode_string(qso.mul1[:])
}

/*
Returns the secondary multiplier.
*/
func (qso *QSO) GetMul2() string {
	return zylo_decode_string(qso.mul2[:])
}

/*
Sets the call sign.
*/
func (qso *QSO) SetCall(value string) {
	copy(qso.call[:], zylo_encode_string(value))
}

/*
Sets the submitted contest number.
*/
func (qso *QSO) SetSent(value string) {
	copy(qso.sent[:], zylo_encode_string(value))
}

/*
Sets the received contest number.
*/
func (qso *QSO) SetRcvd(value string) {
	copy(qso.rcvd[:], zylo_encode_string(value))
}

/*
Sets the operator name.
*/
func (qso *QSO) SetName(value string) {
	copy(qso.name[:], zylo_encode_string(value))
}

/*
Sets remarks.
*/
func (qso *QSO) SetNote(value string) {
	copy(qso.note[:], zylo_encode_string(value))
}

/*
Sets the primary multiplier.
*/
func (qso *QSO) SetMul1(value string) {
	copy(qso.mul1[:], zylo_encode_string(value))
}

/*
Sets the secondary multiplier.
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
Checks for duplicate QSOs.
Returns true if the QSO is valid.
*/
func (qso *QSO) VerifyDupe() bool {
	return !qso.Dupe || zylo_dupes
}

/*
Checks the band.
Returns true if the QSO is valid.
*/
func (qso *QSO) VerifyBand() bool {
	return zylo_bands[qso.Band]
}

/*
Checks the mode.
Returns true if the QSO is valid.
*/
func (qso *QSO) VerifyMode() bool {
	return zylo_modes[qso.Mode]
}

/*
Checks the call sign.
Returns true if the QSO is valid.
*/
func (qso *QSO) VerifyCall() bool {
	return zylo_calls.MatchString(qso.GetCall())
}

/*
Checks the submitted contest number.
Returns true if the QSO is valid.
*/
func (qso *QSO) VerifySent() bool {
	return zylo_sents.MatchString(qso.GetSent())
}

/*
Checks the received contest number.
Returns true if the QSO is valid.
*/
func (qso *QSO) VerifyRcvd() bool {
	return zylo_rcvds.MatchString(qso.GetRcvd())
}

/*
Groups submitted contest numbers by regular expression.
*/
func (qso *QSO) GetSentGroups() []string {
	return zylo_sents.FindStringSubmatch(qso.GetSent())
}

/*
Groups received contest numbers by regular expression.
*/
func (qso *QSO) GetRcvdGroups() []string {
	return zylo_rcvds.FindStringSubmatch(qso.GetRcvd())
}

/*
Makes the QSO an invalid QSO.
*/
func (qso *QSO) Invalid() {
	qso.Score = 0
	qso.SetMul1("")
	qso.SetMul2("")
}

/*
Dumps QSO structure without header information.
*/
func (qso *QSO) DumpWithoutHead(w io.Writer) {
	binary.Write(w, binary.LittleEndian, qso)
}

/*
Reads QSO structure without header information.
*/
func (qso *QSO) LoadWithoutHead(r io.Reader) {
	raw := make([]byte, 256)
	binary.Read(r, binary.LittleEndian, raw)
	*qso = *(*QSO)(unsafe.Pointer(&raw[0]))
}

/*
Converts QSOs to byte sequence with header information.
*/
func DumpZLO(qso ...QSO) (bin []byte) {
	_, diff := time.Now().In(zone).Zone()
	log := append(make([]QSO, 1), qso...)
	log[0].SRST = int16(-diff / 60)
	buf := new(bytes.Buffer)
	for _, qso := range log {
		qso.DumpWithoutHead(buf)
	}
	return buf.Bytes()
}

/*
Converts byte sequence with header information to QSOs.
*/
func LoadZLO(bin []byte) (logs []QSO) {
	buf := bytes.NewReader(bin)
	head := make([]byte, 256)
	io.ReadFull(buf, head)
	zlox := make([]byte, 0)
	if string(head[:4]) == "ZLOX" {
		zlox = make([]byte, 128)
	}
	io.ReadFull(buf, zlox)
	for buf.Len() > 0 {
		qso := QSO{}
		qso.LoadWithoutHead(buf)
		logs = append(logs, qso)
		io.ReadFull(buf, zlox)
	}
	return logs
}

/*
Adds the QSO record to zLog.
*/
func (qso *QSO) Insert() {
	C.insert(unsafe.Pointer(qso), insertCB)
}

/*
Deletes the QSO record in zLog.
*/
func (qso *QSO) Delete() {
	C.delete(unsafe.Pointer(qso), deleteCB)
}

/*
Updates the QSO record in zLog.
*/
func (qso *QSO) Update() {
	C.update(unsafe.Pointer(qso), updateCB)
}

/*
zlog.ini.
*/
func PathToINI() string {
	path, _ := os.Executable()
	tail := filepath.Ext(path)
	path = strings.TrimSuffix(path, tail)
	path, _ = filepath.Abs(path + ".ini")
	return path
}

/*
Gets the specified setting.
*/
func GetINI(section, key string) string {
	init, _ := ini.LooseLoad(PathToINI())
	return init.Section(section).Key(key).String()
}

/*
Sets the specified setting.
*/
func SetINI(section, key, value string) {
	init, _ := ini.LooseLoad(PathToINI())
	init.Section(section).Key(key).SetValue(value)
	init.SaveTo(PathToINI())
}

/*
Catch a panic and display it in a dialog.

	defer DisplayPanic()
*/
func DisplayPanic() {
	if err := recover(); err != nil {
		DisplayModal("%s: %s", err, debug.Stack())
	}
}

/*
Displays the specified string in a dialog.
*/
func DisplayModal(msg string, args ...interface{}) {
	text := C.CString(fmt.Sprintf(msg, args...))
	defer C.free(unsafe.Pointer(text))
	C.dialog(text, dialogCB)
}

/*
Displays the specified string in a toast.
*/
func DisplayToast(msg string, args ...interface{}) {
	text := C.CString(fmt.Sprintf(msg, args...))
	defer C.free(unsafe.Pointer(text))
	C.notify(text, notifyCB)
}

/*
Makes the specified query.
*/
func Query(text string) string {
	buf := make([]byte, ResponseCapacity+1)
	copy(buf[:ResponseCapacity], text[:])
	C.access(unsafe.Pointer(&buf[0]), accessCB)
	return string(buf[:bytes.IndexByte(buf, 0)])
}

/*
Gets the specified window handle.

	GetUI("MainForm.FileOpenItem")
	GetUI("MenuForm.CancelButton")
*/
func GetUI(expression string) uintptr {
	e := C.CString(expression)
	defer C.free(unsafe.Pointer(e))
	return uintptr(C.handle(e, handleCB))
}

/*
Calls the specified Delphi script.
Limited to zLog 2.8.3.0 and later.
*/
func RunDelphi(exp string, args ...interface{}) int {
	e := C.CString(fmt.Sprintf(exp, args...))
	defer C.free(unsafe.Pointer(e))
	return int(C.script(e, scriptCB))
}

var zylo_dupes = false
var zylo_bands = make(map[byte]bool)
var zylo_modes = make(map[byte]bool)
var zylo_calls = regexp.MustCompile(`^.*$`)
var zylo_sents = regexp.MustCompile(`^.*$`)
var zylo_rcvds = regexp.MustCompile(`^.*$`)

/*
Allows duplicate QSO.
*/
func AllowDupe() {
	zylo_dupes = true
}

/*
Allows the specified bands.
*/
func AllowBand(bands ...byte) {
	for _, band := range bands {
		zylo_bands[band] = true
	}
}

/*
Allows the specified modes.
*/
func AllowMode(modes ...byte) {
	for _, mode := range modes {
		zylo_modes[mode] = true
	}
}

/*
Allows the specified bands.
*/
func AllowBandRange(lo, hi byte) {
	for b := lo; b <= hi; b++ {
		AllowBand(b)
	}
}

/*
Allows the specified modes.
*/
func AllowModeRange(lo, hi byte) {
	for m := lo; m <= hi; m++ {
		AllowMode(m)
	}
}

/*
Allows call signs of the pattern.
*/
func AllowCall(pattern string) {
	zylo_calls = regexp.MustCompile(pattern)
}

/*
Allows submitted contest numbers of the pattern.
*/
func AllowSent(pattern string) {
	zylo_sents = regexp.MustCompile(pattern)
}

/*
Allows received contest numbers of the pattern.
*/
func AllowRcvd(pattern string) {
	zylo_rcvds = regexp.MustCompile(pattern)
}

/*
File extension filter for I/O plugin.
*/
var FileExtFilter string

/*
Variable for embedding zLog DAT file.
*/
var CityMultiList string

/*
Called when the plugin is launched.
*/
var OnLaunchEvent = func() {}

/*
Called when the plugin is finished.
*/
var OnFinishEvent = func() {}

/*
Receives window messages.
*/
var OnWindowEvent = func(msg uintptr) {}

/*
Converts a file in another format to a ZLO file.
*/
var OnImportEvent = func(source, target string) error {
	return nil
}

/*
Converts a ZLO file to a file in another format.
*/
var OnExportEvent = func(source, format string) error {
	return nil
}

/*
Called just after opening the contest.
*/
var OnAttachEvent = func(contest, configs string) {}

/*
Called just after closing the contest.
*/
var OnDetachEvent = func(contest, configs string) {}

/*
Called when scoring is delegated.
*/
var OnAssignEvent = func(contest, configs string) {}

/*
Called when a QSO record is added.
*/
var OnInsertEvent = func(qso *QSO) {}

/*
Called when a QSO record is deleted.
Deletion is performed before addition.
*/
var OnDeleteEvent = func(qso *QSO) {}

/*
Determines QSO score and multiplier.
Empty multiplier for an invalid QSO.
Called several times before the QSO is recorded.
*/
var OnVerifyEvent = func(qso *QSO) {
	ok := true
	ok = ok && qso.VerifyDupe()
	ok = ok && qso.VerifyBand()
	ok = ok && qso.VerifyMode()
	ok = ok && qso.VerifySent()
	ok = ok && qso.VerifyRcvd()
	if ok {
		OnAcceptEvent(qso)
	} else {
		qso.Invalid()
	}
}

/*
Sets QSO point and multiplier.
*/
var OnAcceptEvent = func(qso *QSO) {
	qso.SetMul1(qso.GetRcvd())
}

/*
Calculates total score.
*/
var OnPointsEvent = func(score, mul1s int) int {
	return score * mul1s
}

/*
Registers an event handler for the button with the given name.
Subsequent registrations after plugin startup will be ignored.
*/
func HandleButton(name string, handler func(int)) {
	buttons[zylo_add_button_handler(name)] = handler
}

/*
Registers an event handler for the editor with the given name.
Subsequent registrations after plugin startup will be ignored.
*/
func HandleEditor(name string, handler func(int)) {
	editors[zylo_add_editor_handler(name)] = handler
}
