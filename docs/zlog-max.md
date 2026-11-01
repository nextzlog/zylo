zLog Max の機能を zLog 令和版から使うためのプラグインです。
高校コンテストのように zLog 本体での対応が難しいコンテストをサポートします。
また、zLog 令和版のインポートとエクスポートに LG8 や HDB などの形式を加えます。

#### コンテスト定義

CFG ファイルの `exit` の前に `dll zlog-max.dll` と書くと、zLog はこのプラグインに採点を任せます。
zLog Max で扱う Lua スクリプト内蔵型の CFG ファイルを zLog 令和版でもそのまま使えます。

```lua
dll    zlog-max.dll
exit
script lua

function init()
    return {}
end

function add(state, qso)
    return state, {}
end

function total(state, sums)
    return sums.score * sums.mults
end
```

#### 対応する形式

| 形式 | 拡張子 |
|---|---|
| QXML | xml, qxml |
| ADIF | adi, adx |
| Cabrillo | cbr |
| JARL 電子ログ R2.1 | txt |
| zLog | all, txt |
| CTESTWIN | txt, lg8 |
| Turbo HAMLOG | hdb |
