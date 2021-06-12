/*
 provides the zLog programming interface to the Go language.
 Copyright (C) 2020 JA1ZLO.
 */
package zylo

import (
	"bytes"
	"fmt"
	"math"
	"runtime"
	"syscall"
	"time"
	"unsafe"
	"encoding/binary"
	"gopkg.in/ini.v1"
	"gopkg.in/go-toast/toast.v1"
	"github.com/amenzhinsky/go-memexec"
)

/*
 a single QSO data.
 */
type QSO struct {
	time  float64
	call  [13]byte
	sent  [31]byte
	rcvd  [31]byte
	void  byte
	sRST  uint16
	rRST  uint16
	ID    uint32
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
	Pow2  uint32
	rsv2  uint32
	rsv3  uint32
}

/*
 mode enumeration.
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
 band enumeration.
 */
const (
	M1_9  =  0
	M3_5  =  1
	M7    =  2
	M10   =  3
	M14   =  4
	M18   =  5
	M21   =  6
	M24   =  7
	M28   =  8
	M50   =  9
	M144  = 10
	M430  = 11
	M1200 = 12
	M2400 = 13
	M5600 = 14
	G10UP = 15
)

/*
 QSO struct size.
 */
const QBYTES = 256

/*
 path to zlog.ini.
 */
const INI = "zlog.ini"

/*
 reference to qxsl.exe
 */
type QxSL struct {
	exe *memexec.Exec
}

/*
 a bridge function to insert a QSO.
 */
var InsertQSO func(qso *QSO)

/*
 a bridge function to delete a QSO.
 */
var DeleteQSO func(qso *QSO)

/*
 a bridge function to update a QSO.
 */
var UpdateQSO func(qso *QSO)

/*
 a bridge function to add editor handler.
 */
var HookEditor func(name string)

/*
 a bridge function to add button handler.
 */
var HookButton func(name string)

/*
 displays a message as a toast.
 */
func Notify(msg string, args ...interface{}) {
	toast := toast.Notification {
		AppID: "ZyLO",
		Title: "ZyLO",
		Message: fmt.Sprintf(msg, args...),
	}
	toast.Push()
}

/*
 converts a raw pointer into a QSO.
 */
func ToQSO(ptr uintptr) (qso *QSO) {
	return (*QSO)(unsafe.Pointer(ptr))
}

/*
 inserts a QSO to zLog.
 */
func (qso *QSO) Insert() {
	InsertQSO(qso);
}

/*
 deletes a QSO in zLog.
 */
func (qso *QSO) Delete() {
	DeleteQSO(qso);
}

/*
 updates a QSO in zLog.
 */
func (qso *QSO) Update() {
	UpdateQSO(qso);
}

/*
 extracts the operation time from the QSO.
 */
func (qso *QSO) GetTime(zone *time.Location) time.Time {
	var t = math.Abs(qso.time)
	var h = time.Duration((t - float64(int(t))) * 24)
	var d = time.Date(1899, 12, 30, 0, 0, 0, 0, zone)
	return d.Add(h * time.Hour).AddDate(0, 0, int(t))
}

/*
 converts the Delphi byte array into a string.
 */
func getString(field []byte) string {
	return string(field[1 : int(field[0]) + 1])
}

/*
 writes the string into the Delphi byte array.
 */
func setString(field []byte, value string) {
	field[0] = byte(len(value))
	copy(field[1:], value)
}

/*
 extracts the contacted station's callsign.
 */
func (qso *QSO) GetCall() string {
	return getString(qso.call[:])
}

/*
 extracts the transmitted serial number.
 */
func (qso *QSO) GetSent() string {
	return getString(qso.sent[:])
}

/*
 extracts the received serial number.
 */
func (qso *QSO) GetRcvd() string {
	return getString(qso.rcvd[:])
}

/*
 extracts the operator's name.
 */
func (qso *QSO) GetName() string {
	return getString(qso.name[:])
}

/*
 extracts the QSO notes.
 */
func (qso *QSO) GetNote() string {
	return getString(qso.note[:])
}

/*
 extracts the 1st multiplier.
 */
func (qso *QSO) GetMul1() string {
	return getString(qso.mul1[:])
}

/*
 extracts the 2nd multiplier.
 */
func (qso *QSO) GetMul2() string {
	return getString(qso.mul2[:])
}

/*
 extracts the contacted station's callsign.
 */
func (qso *QSO) SetCall(value string) {
	setString(qso.call[:], value)
}

/*
 extracts the transmitted serial number.
 */
func (qso *QSO) SetSent(value string) {
	setString(qso.sent[:], value)
}

/*
 extracts the received serial number.
 */
func (qso *QSO) SetRcvd(value string) {
	setString(qso.rcvd[:], value)
}

/*
 extracts the operator's name.
 */
func (qso *QSO) SetName(value string) {
	setString(qso.name[:], value)
}

/*
 extracts the QSO notes.
 */
func (qso *QSO) SetNote(value string) {
	setString(qso.note[:], value)
}

/*
 sets the 1st multiplier.
 */
func (qso *QSO) SetMul1(value string) {
	setString(qso.mul1[:], value)
}

/*
 sets the 2nd multiplier.
 */
func (qso *QSO) SetMul2(value string) {
	setString(qso.mul2[:], value)
}

/*
 converts a QSO into a binary data.
 */
func (qso *QSO) Dump(locale *time.Location) []byte {
	_, off := time.Now().In(locale).Zone()
	min := int16(-off / 60)
	buf := new(bytes.Buffer)
	buf.Write(make([]byte, 0x54))
	binary.Write(buf, binary.LittleEndian, min)
	buf.Write(make([]byte, 0xAA))
	binary.Write(buf, binary.LittleEndian, qso)
	return buf.Bytes()
}

/*
 gets a specified value from zlog.ini
 */
func GetINI(section, key string) (value string) {
	cfg, err := ini.Load(INI)
	if err == nil {
		value = cfg.Section(section).Key(key).String()
	} else {
		value = ""
	}
	return
}

/*
 sets a specified value into zlog.ini
 */
func SetINI(section, key, value string) (err error) {
	cfg, err := ini.Load(INI)
	if err == nil {
		cfg.Section(section).Key(key).SetValue(value)
		err = cfg.SaveTo(INI)
	}
	return
}

/*
 loads qxsl.exe from the byte array.
 */
func LoadQxSL(bytes []byte) (qxsl *QxSL, err error) {
	var exe *memexec.Exec
	exe, err = memexec.New(bytes)
	if err == nil {
		qxsl = &QxSL {
			exe: exe,
		}
	}
	return
}

/*
 releases resources for qxsl.exe.
 */
func (exe *QxSL) Close() error {
	return exe.exe.Close()
}

/*
 calls qxsl.exe to obtain filter string for a file dialog.
 */
func (exe *QxSL) Filter() (filter string, err error) {
	cmd := exe.exe.Command("filter")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	return string(out), err
}

/*
 calls qxsl.exe to format the specified log into another format.
 */
func (exe *QxSL) Format(source, target, format string) error {
	cmd := exe.exe.Command("format", source, target, format)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.Output()
	return err
}

/*
 captures a panic and display the detailed information of the panic.
 */
func CapturePanic() {
	if err := recover(); err != nil {
		_, file, line, _ := runtime.Caller(1)
		Notify("panic occurred at line %d in %s: %s", line, file, err)
	}
}
