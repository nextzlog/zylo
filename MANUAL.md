zLog+ ZyLO
====

ZyLOは[zLog](http://zlog.org)の拡張機能を提供するDLLを作るための仕組みです。
DLLのプログラミング次第でzLogに無限の可能性を付加します。
例えば、

- 独自のユーザインタフェースや機能の追加
- 従来対応できなかった複雑な規約への対応
- 他のソフトウェアやハードウェアとの連携

## 動作条件

- Window 10 (64bit版)
- [zLog 2.7.0.0+](https://github.com/jr8ppg/zlog) (64bit版)

## インストール方法

ZyLOは最新のzLogに統合されており、DLLをzLogと同じ場所に配置して、設定を変更するだけで利用できます。
例えば、`foo.dll`と`bar.dll`を利用する場合は、`zlog.ini`を開いて、

```ini
[zylo]
DLLs=foo.dll,bar.dll
```

設定項目`[zylo]`を作成し、カンマ区切りでDLLの名前を並べて、保存します。
この状態でzLogを起動すると、DLLの機能と連携し始めます。

## 得点計算の移譲方法

ZyLOに対応したCFGファイルの末尾には、以下の2行があります。

```
exit
dll foo.dll
```

この場合は、得点計算が`foo.dll`側で行われます。
なお、`foo.dll`は事前にインストールが必要です。

## 対応済みコンテスト

従来zLogで対応困難だったコンテストを中心に、DLLを実装しています。

- [YLコンテスト](https://github.com/nextzlog/zylo/tree/master/rules/yltest)
- [高校コンテスト](https://github.com/nextzlog/zylo/tree/master/rules/hstest)
- [リアルタイムコンテスト](https://github.com/nextzlog/zylo/tree/master/rules/rttest)

## その他の拡張機能

ネットワークプログラミングを得意とするGo言語の強みと生産性を活かした便利機能を提供しています。

- [電子ログ変換機能](https://github.com/nextzlog/zylo/tree/master/rules/format)
- [更新自動通知機能](https://github.com/nextzlog/zylo/tree/master/rules/latest)

## 拡張機能の開発方法

ZyLOではGoogleが開発した[Go言語](https://golang.org)によりDLLを開発します。

### ビルド環境

- `x86_64-w64-mingw32-gcc`
- Go 1.16

### ビルド方法

[zbuild](https://github.com/nextzlog/zylo/releases/tag/zbuild)を入手して、DLLのソースコードと同じ場所で`zbuild`を実行します。

Windowsで開発している場合は、

```bat
> build-windows.exe
```

Ubuntuで開発している場合は、

```sh
$ ./zbuild-ubuntu
```

macOSで開発している場合は、

```sh
$ ./zbuild-macos
```

必要に応じて`zbuild`がGoプロジェクトを初期化し、`go.mod`を生成すると同時に、ライブラリとして`zutils.go`を生成し、DLLをビルドします。

### イベントハンドラ

DLLでは、得点計算やQSOの追加・削除を受信するためのイベントハンドラ([詳細](https://nextzlog.github.io/zylo))を適宜実装します。
