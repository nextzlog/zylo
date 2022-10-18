---
title: ZyLO
subtitle: プラグイン開発マニュアル
---

## 拡張機能の利用方法

- zLogのエンドユーザは、zLogに内蔵された設定画面を通じて、拡張機能を簡単にインストールできます。

1. 設定メニューからプラグイン管理メニューを選び、管理画面を開く。
2. 画面上部のリストから好きな拡張機能を選ぶと、詳細が表示される。
3. 画面下部のボタンを押して、拡張機能を有効化・無効化・更新する。

|ボタン |動作                                        |
|-------|--------------------------------------------|
|Install|拡張機能を有効化する。押すと同時に起動する。|
|Disable|拡張機能を無効化する。再起動後に反映される。|
|Upgrade|拡張機能を最新にする。再起動後に反映される。|

- 拡張機能をテストするには、`zlog.ini`に拡張機能のパスをカンマ区切りで記載し、zLogを起動します。

```ini
[zylo]
DLLs=hstest.dll,yltest.dll,rttest.dll
```

## 拡張機能の開発環境

- まず、拡張機能の開発に利用する環境に合わせて、以下のパッケージ管理システムを準備しておきます。

|OS      |package manager                |
|--------|-------------------------------|
|Windows |[choco](https://chocolatey.org)|
|macOS   |[brew](https://brew.sh)        |
|Ubuntu  |[apt](https://debian.org)      |

- 開発環境に適合する[zbuild](https://github.com/nextzlog/zylo/releases/tag/zbuild)を準備します。

|OS      |zbuild            |
|--------|------------------|
|Windows |zbuild-windows.exe|
|macOS   |zbuild-macos      |
|Ubuntu  |zbuild-linux      |

- 最後に、以下のコマンドを実行すると、Go言語の開発環境が自動的にインストールされ、準備完了です。

```bat
> zbuild-windows.exe setup
```

## 拡張機能のビルド例

- 以下に拡張機能の雛形を示します。専用のディレクトリを作成し、適当な名前で雛形を保存しましょう。

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
	// when calculate score
	return score * mults
}
```

- そのフォルダで`compile`コマンドを実行すると、フォルダと同じ名前で、拡張機能のDLLが完成します。

```bat
> zbuild-windows.exe compile --help
USAGE:
    zbuild compile [VER]

ARGS:
    <VER>    [default: 2.8]

> zbuild-windows.exe compile
> zbuild-windows.exe compile 2.8.3.0
```

## 得点計算の移譲方法

- CFGファイルの末尾に、以下の項目を追記することで、そのコンテストの得点計算をDLLに委譲できます。

```
exit # required statement
dll rttest.dll # basename
```

- 具体的には、以下の関数が、CFGファイルで指定されたDLLに対してのみ、呼び出されるようになります。

|関数         |機能                          |
|-------------|------------------------------|
|OnAssignEvent|得点計算の委譲時に通知される。|
|OnVerifyEvent|交信の得点と有効性を確定する。|
|OnPointsEvent|交信記録の総合得点を計算する。|

## 高度な拡張機能の例

- 拡張機能でWinAPIを使う場合は、`GetUI`関数で、zLogのGUI部品のウィンドウハンドルを取得できます。

```go
handle := reiwa.GetUI("MainForm.FileNewItem")
```

- 拡張機能は、zLogのボタン及びメニュー項目のクリックや、記入欄のキーボードの入力を監視できます。

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

- 拡張機能は、WinAPIのGo実装を利用することで、zLogに独自の部品を追加し、イベントを処理できます。

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

## プロパティ操作の例

- `RunDelphi`関数は、zLogのGUI部品のプロパティやメソッドをDelphiの式で参照もしくは実行できます。

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

### 演算子

|演算子|意味  |返り値の型|
|------|------|----------|
| +    |加算  |Extended  |
| -    |減算  |Extended  |
| *    |乗算  |Extended  |
| /    |除算  |Extended  |
| =    |等値  |Boolean   |
| <>   |不等  |Boolean   |
| <    |未満  |Boolean   |
| >    |超過  |Boolean   |
| <=   |以下  |Boolean   |
| >=   |以上  |Boolean   |
| and  |論理積|Boolean   |
| or   |論理和|Boolean   |
| not  |否定  |Boolean   |

### 組み込み関数

|関数|意味            |
|----|----------------|
|Int |実数を整数に変換|
|Put |プロパティを設定|

```pascal
op.Int(real_value) // returns int value
op.Put(object, property_name, property_value) // returns the object
```

## クエリを実行する例

- zLogのマクロの内容を取得できます。CWキーボードのマクロに加え、表に掲載する変数が利用可能です。

```go
fmt.Println(reiwa.Query("$B,$X,$R,$F,$Z,$I,$Q,$V,$O,$S,$P,$A,$N,$L,$C,$E,$M"))
```

|変数|内容          |
|----|--------------|
|{V} |バージョン番号|
|{F} |編集ファイル名|
|{C} |自局の呼出符号|
|{B} |運用中のバンド|
|{M} |運用中のモード|

## 拡張機能の頒布方法

- 最初に、適当なTOMLファイルを作成し、DLLの名称に続け、リリースとチェックサムのURLを記載します。

```toml
# dll.DLL名
[dll.sample]
url = "https://example.com/releases/sample.dll"
sum = "https://example.com/releases/sample.dll.md5"
```

- また、DLLに付随してCFGファイルやDATファイル等を記載できます。DLLと同様に必要事項を記載します。

```toml
[cfg.sample]
url = "https://example.com/releases/sample.cfg"

[dat.sample]
url = "https://example.com/releases/sample.dat"
```

- 最後に、拡張機能の詳細な情報を記述します。この内容が、zLogの拡張機能の管理画面に表示されます。

```toml
[pkg.sample]
tag = "title"
msg = "description"
web = "https://example.com/sample/index.html"
doc = "https://example.com/sample/details.md"
use = ["cfg.sample", "dat.sample", "dll.sample"] # dependency
exp = "unstable" # or "stable"
```

- 適当な管理者を選び、前掲のTOMLファイルの公開を依頼して、マーケットプレイスへの反映を待ちます。

## 自動リリースの設定

- 拡張機能をGitHubで開発している場合は、以下に示すワークフローにより、リリースを自動化できます。

{% raw %}

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

{% endraw %}

## オブジェクトの便覧

- 以下にzLogのメイン画面のオブジェクトを列挙します。`GetUI`関数や`RunDelphi`関数で参照できます。

```pascal
{% include form.dfm %}
```

## 組み込み関数と変数


