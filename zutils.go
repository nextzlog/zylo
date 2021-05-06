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
	"gopkg.in/go-toast/toast.v1"
	"github.com/amenzhinsky/go-memexec"
)

/*
 defines a QSO data frame in zLog binary format.
 */
type QSO struct {
	time float64
	call [13]byte
	sent [31]byte
	rcvd [31]byte
	void byte
	sRST uint16
	rRST uint16
	seID uint32
	Mode byte
	Band byte
	pow1 byte
	mul1 [31]byte
	mul2 [31]byte
	new1 bool
	new2 bool
	mark byte
	name [15]byte
	note [65]byte
	isCQ bool
	dupe bool
	rsv1 byte
	txID byte
	pow2 uint32
	rsv2 uint32
	rsv3 uint32
}

/*
 QSO list.
 */
type Log []QSO

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
 reference to qxsl.exe
 */
type QxSL struct {
	exe *memexec.Exec
}

/*
 a bridge function of Delphi InputBox.
 */
var Ibox func(string, string) (string, bool)

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
 converts a raw pointer into a QSO list.
 */
func ToLog(ptr uintptr, size int) (list Log) {
	for i := 0; i < size; i++ {
		list = append(list, *ToQSO(ptr))
		ptr += QBYTES
	}
	return
}

/*
 converts a QSO list into a binary data.
 */
func (log *Log) Dump(zone *time.Location) []byte {
	_, off := time.Now().In(zone).Zone()
	min := int16(-off / 60)
	buf := new(bytes.Buffer)
	buf.Write(make([]byte, 0x54))
	binary.Write(buf, binary.LittleEndian, min)
	buf.Write(make([]byte, 0xAA))
	binary.Write(buf, binary.LittleEndian, log)
	return buf.Bytes()
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
 extracts the contest QSO transmitted serial number.
 */
func (qso *QSO) GetSent() string {
	return getString(qso.sent[:])
}

/*
 extracts the contest QSO received serial number.
 */
func (qso *QSO) GetRcvd() string {
	return getString(qso.rcvd[:])
}

/*
 extracts the logging operator's name.
 */
func (qso *QSO) GetName() string {
	return getString(qso.name[:])
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
 extracts whether the 1st multiplier appears for the first time.
 */
func (qso *QSO) IsNewMul1() bool {
	return qso.new1
}

/*
 extracts whether the 2nd multiplier appears for the first time.
 */
func (qso *QSO) IsNewMul2() bool {
	return qso.new2
}

/*
 extracts the contacted station's callsign.
 */
func (qso *QSO) SetCall(value string) {
	setString(qso.call[:], value)
}

/*
 extracts the contest QSO transmitted serial number.
 */
func (qso *QSO) SetSent(value string) {
	setString(qso.sent[:], value)
}

/*
 extracts the contest QSO received serial number.
 */
func (qso *QSO) SetRcvd(value string) {
	setString(qso.rcvd[:], value)
}

/*
 extracts the logging operator's name.
 */
func (qso *QSO) SetName(value string) {
	setString(qso.name[:], value)
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
 sets whether the 1st multiplier appears for the 1st time.
 */
func (qso *QSO) SetNewMul1(value bool) {
	qso.new1 = value
}

/*
 sets whether the 2nd multiplier appears for the 1st time.
 */
func (qso *QSO) SetNewMul2(value bool) {
	qso.new2 = value
}

/*
 loads qxsl.exe from the byte array.
 */
func NewQxSL(bytes []byte) (qxsl *QxSL, err error) {
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
func (qxsl *QxSL) Close() error {
	return qxsl.exe.Close()
}

/*
 calls qxsl.exe to obtain filter string for a file dialog.
 */
func (qxsl *QxSL) Filter() (filter string, err error) {
	cmd := qxsl.exe.Command("filter")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	return string(out), err
}

/*
 calls qxsl.exe to format the specified log into another format.
 */
func (qxsl *QxSL) Format(source, target, format string) error {
	cmd := qxsl.exe.Command("format", source, target, format)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	_, err := cmd.Output()
	return err
}

/*
 displays a modal dialog asking the user to enter a string value.
 */
func Prompt(label, value string) (ret string, ok bool) {
	return Ibox(label, value);
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
