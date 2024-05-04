ZyLO
====

![image](https://img.shields.io/badge/Go-1.22-red.svg)
![image](https://img.shields.io/badge/license-MIT-darkblue.svg)
![badge](https://github.com/nextzlog/zylo/actions/workflows/build.yaml/badge.svg)
![badge](https://github.com/nextzlog/zylo/actions/workflows/clean.yaml/badge.svg)

ZyLO is the official plugin SDK for [zLog v2.8 or later](https://use.zlog.org).
zLog is powerful and extensible logging software for amateur radio contests, and users can easily install plugins via plugin manager dialog.

## Documents

- [Slides](https://nextzlog.dev/zylo.pdf)
- [Manual](https://nextzlog.github.io/zylo)

## Get Started

### SDK

First, prepare the package management system (`brew`, `choco`, and `apt`) and install Go.

```sh
$ brew install go
$ choco install golang
$ sudo apt install -y golang-go
```

Then, run the following commands to set up the SDK.

```sh
$ go install github.com/nextzlog/zylo/zbuild@HEAD
$ zbuild setup
```

### Template

Create a directory and the source file shown below in it in the same way as creating a Go package.

```go
package main

import "zylo/reiwa"

func init() {
  // when plugin loaded
  reiwa.OnLaunchEvent = onLaunchEvent
  reiwa.OnFinishEvent = onFinishEvent
  reiwa.OnAttachEvent = onAttachEvent
  reiwa.OnAssignEvent = onAssignEvent
  reiwa.OnDetachEvent = onDetachEvent
  reiwa.OnInsertEvent = onInsertEvent
  reiwa.OnDeleteEvent = onDeleteEvent
  reiwa.OnVerifyEvent = onVerifyEvent
  reiwa.OnPointsEvent = onPointsEvent
}

func onLaunchEvent() {
  // when zLog launched
}

func onFinishEvent() {
  // when zLog finished
}

func onAttachEvent(contest, configs string) {
  // when contest attached
}

func onAssignEvent(contest, configs string) {
  // if CFG calls this DLL
}

func onDetachEvent(contest, configs string) {
  // when contest detached
}

func onInsertEvent(qso *reiwa.QSO) {
  // when insert this QSO
}

func onDeleteEvent(qso *reiwa.QSO) {
  // when delete this QSO
}

func onVerifyEvent(qso *reiwa.QSO) {
  // score and multiplier
}

func onPointsEvent(score, mults int) int {
  return score * mults
}
```

### Build

Run `zbuild` inside the directory to build a plugin DLL with the same name as the directory.

```sh
$ zbuild build
$ zbuild build --version 2.8.3.0
```

### Launch

To test the plugin, put the relative path to the DLL into `zlog.ini` as follows and launch zLog.

```ini
[zylo]
DLLs=hstest.dll,yltest.dll,rttest.dll
```

## Publish

First, create a TOML file and include the download URL and md5 checksum in the `dll` table.
You can include the checksum directly or include the URL of the md5 file.

```toml
# dll.<name>
[dll.sample]
url = "https://example.com/releases/sample.dll"
sum = "https://example.com/releases/sample.dll.md5"
```

You can also publish other files along with the DLL, such as CFG and DAT files.

```toml
# cfg.<name>
[cfg.sample]
url = "https://example.com/releases/sample.cfg"

# dat.<name>
[dat.sample]
url = "https://example.com/releases/sample.dat"
```

Finally, write the plugin meta information, which will be displayed in the plugin manager.

```toml
# pkg.<name>
[pkg.sample]
tag = "title"
msg = "description"
web = "https://example.com/sample/index.html"
doc = "https://example.com/sample/details.md"
use = ["cfg.sample", "dat.sample", "dll.sample"] # dependency
exp = "unstable" # or "stable"
```

Publish the TOML file to a Git repository, and make an issue at [nextzlog/todo](https://github.com/nextzlog/todo) to request crawling.

## GitHub Actions

The following workflow automates plugin releases.

```yaml
name: 'build'
on:
  push:
    branches:
    - 'main'
jobs:
  BuildDLL:
    runs-on: ubuntu-latest
    steps:
    - uses: nextzlog/zylo@master
      with:
        token: ${{secrets.GITHUB_TOKEN}}
        directory: /path/to/golang/files
```

## Scoring

To delegate score calculation for user-defined contests to a plugin, add the following commands to the end of the CFG file.

```
exit
dll sample.dll # basename
```

The following functions are called only for DLLs specified in the CFG file.

- `OnAssignEvent` invoked only once at the start of the scoring delegation.
- `OnVerifyEvent` determines the score and multiplier of the communication.
- `OnPointsEvent` determines the total score of the contest communications.

## Windows API

Plugins can monitor button and menu clicks and keyboard input in zLog application.

```go
package main

import "zylo/reiwa"

func init() {
  reiwa.OnLaunchEvent = onLaunchEvent
}

func onLaunchEvent() {
  reiwa.HandleButton("MainForm.CWPlayButton", onButton)
  reiwa.HandleEditor("MainForm.CallsignEdit", onEditor)
}

func onButton(num int) {
  reiwa.DisplayToast("CWPlayButton was clicked")
}

func onEditor(key int) {
  reiwa.DisplayToast(reiwa.Query("QSO with $B"))
}
```

Plugins can obtain the window handle of a zLog component with the `GetUI` function.
Plugins can also use Go binding of Windows API to add their own components and handle events.

```go
package main

import (
  "fmt"
  "unsafe"
  "github.com/gonutz/w32"
  "zylo/reiwa"
)

func init() {
  reiwa.OnLaunchEvent = onLaunchEvent
  reiwa.OnWindowEvent = onWindowEvent
}

func onLaunchEvent() {
  h := w32.HMENU(reiwa.GetUI("MainForm.MainMenu"))
  w32.AppendMenu(h, w32.MF_STRING, 810, "My Menu")
  w32.DrawMenuBar(w32.HWND(reiwa.GetUI("MainForm")))
}

func onWindowEvent(msg uintptr) {
  m := (*w32.MSG)(unsafe.Pointer(msg))
  fmt.Printf("Window Message %v\n", m)
}
```

## Delphi API

Plugins can get and set properties and execute methods of zLog components in the form of Delphi expressions.

```go
package main

import "zylo/reiwa"

func init() {
  reiwa.OnLaunchEvent = onLaunchEvent
}

func onLaunchEvent() {
  reiwa.RunDelphi(`PluginMenu.Add(op.Put(MainMenu.CreateMenuItem(), "Name", "MyMenu"))`)
  reiwa.RunDelphi(`op.Put(MainMenu.FindComponent("MyMenu"), "Caption", "Special Menu")`)
}
```

The following built-in functions are available.

```pascal
op.Int(number) // converts real value to int value and returns it
op.Put(obj, key, value) // sets properties and returns the object
```

## Query

Plugins can retrieve the variables of the CW keyboard in zLog.

```go
fmt.Println(reiwa.Query("$B,$X,$R,$F,$Z,$I,$Q,$V,$O,$S,$P,$A,$N,$L,$C,$E,$M"))
```

The following variables are also available.

```go
fmt.Println(reiwa.Query("{V}")) // version number
fmt.Println(reiwa.Query("{F}")) // ZLO(ZLOX) file
fmt.Println(reiwa.Query("{C}")) // your call sign
fmt.Println(reiwa.Query("{B}")) // operating band
fmt.Println(reiwa.Query("{M}")) // operating mode
```

## Contribution

Feel free to make issues at [nextzlog/todo](https://github.com/nextzlog/todo).
Follow [@nextzlog](https://twitter.com/nextzlog) on Twitter.

## License

### Author

[無線部開発班](https://nextzlog.dev)

- JG1VPP
- JS2FVO
- JJ1GUJ

### Clauses

[MIT License](LICENSE)
