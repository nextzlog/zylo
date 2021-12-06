Tutorial
====

ZyLOを利用すれば[zLog](https://zlog.org)の拡張機能をGo言語で開発できます。
例えば、

- 独自のユーザインタフェースや機能の追加
- 従来対応できなかった複雑な規約への対応
- 他のソフトウェアやハードウェアとの連携

無限の可能性を切り開きましょう。

## 具体例

- [公開された拡張機能のリスト](https://zylo.pafelog.net/market.html)

## 拡張機能の利用方法

### 利用者向けの情報

- zLogのエンドユーザは、zLogに内蔵された設定画面を通じて、拡張機能を簡単にインストールできます。

1. 設定メニューからプラグイン管理メニューを選び、管理画面を開く。
2. 画面上部のリストから好きな拡張機能を選ぶと、詳細が表示される。
3. 以下に示すボタンを押して、拡張機能を有効化・無効化・更新する。

|ボタン |動作                                        |
|-------|--------------------------------------------|
|Install|拡張機能を有効化する。押すと同時に起動する。|
|Disable|拡張機能を無効化する。再起動後に反映される。|
|Upgrade|拡張機能を最新にする。再起動後に反映される。|

### 開発者向けの情報

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

func init() {
	// when plugin loaded
	OnLaunchEvent = onLaunchEvent
	OnFinishEvent = onFinishEvent
	OnAttachEvent = onAttachEvent
	OnAssignEvent = onAssignEvent
	OnDetachEvent = onDetachEvent
	OnInsertEvent = onInsertEvent
	OnDeleteEvent = onDeleteEvent
	OnVerifyEvent = onVerifyEvent
	OnPointsEvent = onPointsEvent
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

func onInsertEvent(qso *QSO) {
	// when insert this QSO
}

func onDeleteEvent(qso *QSO) {
	// when delete this QSO
}

func onVerifyEvent(qso *QSO) {
	// score and multiplier
}

func onPointsEvent(score, mults int) int {
	// when calculate score
	return score * mults
}
```

- そのフォルダで`compile`コマンドを実行すると、フォルダと同じ名前で、拡張機能のDLLが完成します。

```bat
> zbuild-windows.exe compile
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

### ウィンドウハンドルの取得

- 拡張機能でWinAPIを使う場合は、`GetUI`関数で、zLogのGUI部品のウィンドウハンドルを取得できます。
- 部品は、[main.dfm](https://github.com/jr8ppg/zLog/blob/master/zlog/main.dfm)等で確認できます。

```go
handle := GetUI("MainForm.FileNewItem")
```

### ボタンやタイピングの監視

- 拡張機能は、zLogのボタン及びメニュー項目のクリックや、記入欄のキーボードの入力を監視できます。

```go
package main

func init() {
	OnLaunchEvent = onLaunchEvent
}

func onLaunchEvent() {
	HandleButton("MainForm.CWPlayButton", onButton)
	HandleEditor("MainForm.CallsignEdit", onEditor)
}

func onButton(num int) {
	DisplayToast("click CWPlayButton")
}

func onEditor(key int) {
	DisplayToast(Query("QSO with $B"))
}
```

### 独自のメニュー項目の追加

- 拡張機能は、WinAPIのGo実装を利用することで、zLogに独自の部品を追加し、イベントを処理できます。

```go
package main

import (
	"fmt"
	"unsafe"
	"github.com/gonutz/w32"
)

func init() {
	OnLaunchEvent = onLaunchEvent
	OnWindowEvent = onWindowEvent
}

func onLaunchEvent() {
	h := w32.HMENU(GetUI("MainForm.MainMenu"))
	w32.AppendMenu(h, w32.MF_STRING, 810, "GO!")
	w32.DrawMenuBar(w32.HWND(GetUI("MainForm")))
}

func onWindowEvent(msg uintptr) {
	m := (*w32.MSG)(unsafe.Pointer(msg))
	fmt.Printf("Window Message %v\n", m)
}
```

## 拡張機能の頒布方法

### パッケージの作成

- 最初に、適当なTOMLファイルを作成し、その冒頭に、DLLの名称と最新版の配布場所のURLを記載します。

```toml
# dll.DLL名
[dll.sample]
url = "https://example.com/releases/sample.dll"
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
web = "website URL"
use = ["cfg.sample", "dat.sample", "dll.sample"] # dependency
exp = "unstable" # or "stable"
```

### パッケージの公開

- 適当な管理者を選び、前掲のTOMLファイルの公開を依頼して、マーケットプレイスへの反映を待ちます。

### 自動ビルドの設定

- 拡張機能をGitHubで開発している場合は、便利なActionを使って、拡張機能のビルドを自動化できます。
- 以下のワークフローは、`main`ブランチが更新されると、DLLをビルドして`nightly`タグで公開します。

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

    # optional: request market update
    - uses: peter-evans/repository-dispatch@v1
      with:
        token: YOUR_PRIVATE_ACCESS_TOKEN
        repository: nextzlog/zylo
        event-type: store
```
