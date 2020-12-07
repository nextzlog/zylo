ZyLO: Go Extension Mechanism for zLog+
====

![image](https://img.shields.io/badge/Go-1.15-red.svg)
![image](https://img.shields.io/badge/Delphi-10.4-red.svg)
![image](https://img.shields.io/badge/license-GPL3-darkblue.svg)

[zLog](http://zlog.org) is a simple logging software for amateur radio contests, originally developed by JJ1MED at the [University of Tokyo Amateur Radio Club](http://ja1zlo.u-tokyo.org), loved by many contesters for 30 years.

## Features

- imports and exports the QSO data in arbitrary formats including QXML, [ADIF](http://adif.org), [Cabrillo](https://wwrof.org/cabrillo/), [CTESTWIN](http://e.gmobb.jp/ctestwin/Download.html), etc.
- defines event handler specifications for implementing external applications as a DLL that work together with [zLog REIWA edition](https://github.com/jr8ppg/zLog).

## Releases

[Download here](https://github.com/nextzlog/zylo/releases).

## Documents

- Run [GoDoc](https://godoc.org) as follows:

```sh
$ godoc -http=localhost:8000
```

- Then open http://localhost:8000 in a web browser.

## Events

- When zLog loads a user defined contest, zLog will try to load a DLL with the same name as the CFG file except for the file extension.
- The following functions need to be provided to zLog by the DLL and are called by zLog as needed.

### Launch

- When the contest is initialized, zLog calls the `zlaunch` function and provides the DLL with the zLog configurations via text.

```go
func zlaunch(cfg string) {}
```

### Verify

- When zLog needs to update the score, zLog calls the `zrevise` function to extract the multiplier, and calls the `zverify` function to calculate the individual QSO score, and finally invoke the `zresult` function to update the total score.

```go
func zrevise(qso uintptr) {}
func zverify(qso uintptr) (score int) {}
func zresult(log uintptr) (total int) {}
```

### Modify

- zLog calls the `zinsert` function every time zLog imports or appends QSOs, and calls the `zdelete` function when zLog deletes or modifies a QSO.

```go
func zinsert(qso uintptr) {}
func zdelete(qso uintptr) {}
```

### Finish

- zLog calls the `zfinish` function after the contest is closed in the main window.

```go
func zfinish() {}
```

## Build

### zLog

1. Clone this repository.
2. Invoke `setup.bat` and you will find the `zLog` directory.
3. Open `zLog/Zlog.dpr` via Delphi RAD Studio and [build zLog following JR8PPG's instructions](https://github.com/jr8ppg/zLog).

### ZyLO

```go
$ go get github.com/nextzlog/zylo
```

## Dependency

- [TProcess-Delphi](https://github.com/z505/TProcess-Delphi)

## Contribution

Feel free to contact [@nextzlog](https://twitter.com/nextzlog) on Twitter.

## License

### Author

[無線部開発班](https://pafelog.net)

### Clauses

- This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

- This program is distributed in the hope that it will be useful, but **without any warranty**; without even the implied warranty of **merchantability or fitness for a particular purpose**.
See the GNU General Public License for more details.

- You should have received a copy of the GNU General Public License along with this program.
If not, see <http://www.gnu.org/licenses/>.
