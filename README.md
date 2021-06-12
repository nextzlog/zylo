ZyLO: Go Extension Mechanism for zLog+
====

![image](https://img.shields.io/badge/Go-1.16-red.svg)
![image](https://img.shields.io/badge/Rust-1.51-red.svg)
![image](https://img.shields.io/badge/Delphi-10.4-red.svg)
![image](https://img.shields.io/badge/license-GPL3-darkblue.svg)

[zLog](http://zlog.org) is a simple logging software for amateur radio contests, originally developed by JJ1MED at the [University of Tokyo Amateur Radio Club](http://ja1zlo.u-tokyo.org), loved by many contesters for 30 years.

## Features

- helps develop DLLs that work together with zLog to realize flexible, dynamic definition of amateur radio contests.
- supports importing and exporting QSO data in any format, including QXML, [ADIF](http://adif.org), [Cabrillo](https://wwrof.org/cabrillo/), [CTESTWIN](http://e.gmobb.jp/ctestwin/Download.html), etc.

## Documents

- [GoDoc](https://nextzlog.github.io/zylo)
- [Wiki](https://github.com/nextzlog/zylo/wiki)

## Events

- The following functions need to be provided to zLog by the DLL and are called by zLog as needed.

### Start & Exit

- zLog calls the `zlaunch` (`zfinish`) function to initialize (terminate) the DLL when zLog is launched (terminated). 

```go
func zlaunch() {}
func zfinish() {}
```

### Open & Close

- zLog calls the `zattach` (`zdetach`) function when the contest is opened (closed).

```go
func zattach(name string, path string) {}
func zdetach() {}
```

### Add & Delete

- zLog calls the `zinsert` (`zdelete`) function every time before zLog appends (deletes) or updates a QSO.

```go
func zinsert(qso *zylo.QSO) {}
func zdelete(qso *zylo.QSO) {}
```

### Validate QSO

- zLog calls the `zverify` function to calculate the score and multiplier for the latest QSO, and the mutiplier must be a member of the city list provided by the `zcities` function.

```go
func zverify(qso *zylo.QSO) {}
func zcities() (dat string) {}
```

### Update Score

- zLog calls the `zpoints` function to calculate the current total score.

```go
func zpoints(score, mults int) (total int) {}
```

### Key & Button

- zLog calls the `zeditor` (`zbutton`) function every time the user presses a key (button) to enter a QSO or send a Morse code.

```go
func zeditor(key int, name string) (block bool) {}
func zbutton(btn int, name string) (block bool) {}
```

## Build Tool

- ZyLO provides a [zbuild](https://github.com/nextzlog/zylo/releases/tag/zbuild) tool that compiles Go files in the working directory and calls an external Go compiler to create a DLL.

```sh
$ zbuild
```

- Note that `zbuild` creates some files which includes low-level functions exported to zLog.

## Build

- If you need to reinitialize the `zylo` module, run the following command:

```sh
$ go mod init github.com/nextzlog/zylo
```

## Contribution

Feel free to contact [@nextzlog](https://twitter.com/nextzlog) on Twitter.

## License

### Author

[無線部開発班](https://pafelog.net)

- JG1VPP
- JS2FVO

### Clauses

- This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

- This program is distributed in the hope that it will be useful, but **without any warranty**; without even the implied warranty of **merchantability or fitness for a particular purpose**.
See the GNU General Public License for more details.

- You should have received a copy of the GNU General Public License along with this program.
If not, see <http://www.gnu.org/licenses/>.
