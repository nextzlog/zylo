zLog+ ZyLO for Windows
====

ZyLOを利用すれば[zLog](https://zlog.org)の拡張機能をGo言語で開発できます。
例えば、

- 独自のユーザインタフェースや機能の追加
- 従来対応できなかった複雑な規約への対応
- 他のソフトウェアやハードウェアとの連携

利用者はzLogのプラグイン管理機能を通じて多彩な拡張機能にアクセスできます。
無限の可能性を切り開きましょう。

## 具体例

- `format.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/format))
- `latest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/latest))
- `hstest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/hstest))
- `tmtest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/tmtest))
- `yltest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/yltest))
- `rttest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/rttest))

## プラグイン管理機能

- zLogのエンドユーザは、zLogのプラグイン管理機能を通じて、拡張機能を簡単にインストールできます。

> メニューバー &rarr; Settings &rarr; Plugin Manager &rarr; ZyLO Plugin Managerの画面が開く

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

- 最後に、**パッケージ**の詳細を記述します。この内容が、zLogのプラグイン管理画面に表示されます。

```toml
[pkg.sample]
tag = "title"
msg = "description"
web = "website URL"
use = ["cfg.sample", "dat.sample", "dll.sample"] # 依存先
```

- 毎週末にクローラが巡回し、`market.toml`の内容を検査して、DLLをzLogの全ての利用者に公開します。

# ZyLO API

{{.EmitUsage}}
