zLog+ ZyLO for Windows
====

[無線部開発班](https://pafelog.net)

ZyLOを利用すれば[zLog](https://zlog.org)の拡張機能をGo言語で開発できます。
例えば、

- 独自のユーザインタフェースや機能の追加
- 従来対応できなかった複雑な規約への対応
- 他のソフトウェアやハードウェアとの連携

無限の可能性を切り開きましょう。

## 具体例

|拡張機能                                                               |内容                              |
|-----------------------------------------------------------------------|----------------------------------|
|[format.dll](https://github.com/nextzlog/zylo/tree/master/utils/format)|zLogに様々なログ形式を追加します。|
|[latest.dll](https://github.com/nextzlog/zylo/tree/master/utils/latest)|zLogの最新のリリースを通知します。|
|[hstest.dll](https://github.com/nextzlog/zylo/tree/master/rules/hstest)|全国高等学校コンテストの規約です。|
|[rttest.dll](https://github.com/nextzlog/zylo/tree/master/rules/rttest)|リアルタイムコンテストの規約です。|
|[tmtest.dll](https://github.com/nextzlog/zylo/tree/master/rules/tmtest)|東海マラソンコンテストの規約です。|

## 拡張機能の管理画面

- zLogのエンドユーザは、zLogに内蔵された設定画面を通じて、拡張機能を簡単にインストールできます。

1. 設定メニューからプラグイン管理メニューを選び、管理画面を開く。
2. 画面上部のリストから好きな拡張機能を選ぶと、詳細が表示される。
3. 以下に示すボタンを押して、拡張機能を有効化・無効化・更新する。

|ボタン |動作                                        |
|-------|--------------------------------------------|
|Install|拡張機能を有効化する。押すと同時に起動する。|
|Disable|拡張機能を無効化する。再起動後に反映される。|
|Upgrade|拡張機能を最新にする。再起動後に反映される。|

## 拡張機能の開発環境

- 最初に、以下のパッケージ管理システムをインストールして、コマンドを実行できることを確認します。

|OS      |package manager                |
|--------|-------------------------------|
|Windows |[choco](https://chocolatey.org)|
|macOS   |[brew](https://brew.sh)        |
|Ubuntu  |[apt](https://debian.org)      |

- その後、[zbuild](https://github.com/nextzlog/zylo/releases/tag/zbuild)を実行し、準備完了です。

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

- その場で`zbuild`を実行すると、DLLが生成されます。

```bat
> zbuild-windows.exe compile
```

## 拡張機能の起動方法

- 開発時にDLLをzLogと連携させるには、DLLをzLogの場所に置き、`zlog.ini`に以下の項目を追記します。

```ini
[zylo]
DLLs=hstest.dll,yltest.dll,rttest.dll
```

## 得点計算の移譲方法

- CFGファイルに以下の項目を追記することで、そのコンテストの得点計算をDLLに委ねることができます。

```
exit
dll rttest.dll
```

## 高度な拡張機能の例

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

- 適当な管理者に依頼して、その管理者が公開中の`market.toml`に、DLLの名称と配布場所を記載します。

```toml
# dll.DLL名
[dll.sample]
url = "https://example.com/releases/sample.dll"
```

- また、DLLに付随してCFGファイルやDATファイル等を添付できます。DLLと同様に配布場所を記載します。

```toml
[cfg.sample]
url = "https://example.com/releases/sample.cfg"

[dat.sample]
url = "https://example.com/releases/sample.dat"
```

- 最後に、**パッケージ**の詳細を記述します。この内容が、zLogの拡張機能の管理画面に表示されます。

```toml
[pkg.sample]
tag = "title"
msg = "description"
web = "website URL"
use = ["cfg.sample", "dat.sample", "dll.sample"] # 依存先
exp = "unstable"
```

- 毎週末にクローラが巡回し、`market.toml`の内容を検査して、DLLをzLogの全ての利用者に公開します。

# ZyLO API

{{.EmitUsage}}
