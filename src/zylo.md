ZyLOを利用すれば[zLog](https://zlog.org)の拡張機能をGo言語で開発できます。
例えば、

- 独自のユーザインタフェースや機能の追加
- 従来対応できなかった複雑な規約への対応
- 他のソフトウェアやハードウェアとの連携

利用者はzLogのマーケットプレイス機能を通じて多彩な拡張機能を入手できます。
無限の可能性を切り開きましょう。

## 具体例

- `format.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/format))
- `latest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/latest))
- `toasty.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/toasty))
- `hstest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/hstest))
- `tmtest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/tmtest))
- `yltest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/yltest))
- `rttest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/rttest))

## プラグイン管理機能

- zLogのエンドユーザは、zLogのプラグイン管理機能を通じて、拡張機能を簡単にインストールできます。

![marketplace](market.png)

## 拡張機能のビルド例

- 通知欄に文字列を表示するだけの簡単な拡張機能として、`toasty.dll`のソースファイルを入手します。

```sh
> git clone https://github.com/nextzlog/zylo
> cd zylo/utils/toasty
```

- 開発環境を用意して[zbuild](https://github.com/nextzlog/zylo/releases/tag/zbuild)を実行します。

```bat
> zbuild-windows.exe setup
> zbuild-windows.exe compile
```

- これで`toasty.dll`が生成されます。

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

- 適当なマーケット管理者に依頼して、その管理者が公開する`market.toml`に、DLLの詳細を追記します。

```toml
[pkg.toasty]
tag = "title"
msg = "description"
web = "website URL"
use = ["cfg.toasty", "dat.toasty", "dll.toasty"]

[cfg.toasty]
url = "https://example.com/releases/toasty.cfg"

[dat.toasty]
url = "https://example.com/releases/toasty.dat"

[dll.toasty]
url = "https://example.com/releases/toasty.dll"
```

## クローラの定期巡回

- 毎週末にクローラが巡回し、`market.toml`の内容を検査して、DLLをマーケットプレイスに公開します。

{{.EmitUsage}}
