zLog+ ZyLO for Windows
====

![image](https://img.shields.io/badge/Go-1.17-red.svg)
![image](https://img.shields.io/badge/Rust-1.56-red.svg)
![image](https://img.shields.io/badge/Delphi-10.4-red.svg)
![image](https://img.shields.io/badge/license-MIT-darkblue.svg)
![badge](https://github.com/nextzlog/zylo/actions/workflows/build.yaml/badge.svg)

[ZyLO](https://github.com/nextzlog/zylo) is a plugin and marketplace mechanism integrated into [zLog](http://zlog.org).
zLog is a simple yet powerful logging software for ham radio contests, which has been loved by many users for 30 years.

## Documents

- [API](https://nextzlog.github.io/zylo)

## Samples

- `format.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/format))
- `latest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/latest))
- `toasty.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/utils/toasty))
- `hstest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/hstest))
- `tmtest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/tmtest))
- `yltest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/yltest))
- `rttest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/rttest))

## Build DLL

- First, clone the sample project [toasty](https://github.com/nextzlog/zylo/tree/master/utils/toasty) as follows.

```sh
$ git clone https://github.com/nextzlog/zylo
$ cd zylo/rules/toasty
```

- Next, create `toasty.dll` by the [`zbuild`](https://github.com/nextzlog/zylo/releases/tag/zbuild) command as follows, and you will find `toasty.dll` in the directory.

### Build DLL on Linux

```sh
$ ./zbuild-linux setup
$ ./zbuild-linux compile
```

### Build DLL on macOS

```sh
$ ./zbuild-macos setup
$ ./zbuild-macos compile
```

### Build DLL on Windows

```bat
> zbuild-windows.exe setup
> zbuild-windows.exe compile
```

## Install DLL

- To install the DLL manually, place it in the same directory as zLog and add the following lines to `zlog.ini`.

```ini
[zylo]
DLLs=foo.dll,bar.dll,baz.dll
```

- To define a contest that uses the DLL, add the following lines to the CFG file.

```
exit
dll foo.dll
```

## Publish DLL

- Ask one of the [market managers](https://github.com/nextzlog/zylo/blob/master/src/market.list) to add the release URL of the DLL to `market.toml`.

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

- Crawler generates [market.json](https://nextzlog.github.io/zylo/market.json) every Saturday at 0:00 from the TOML files to notify zLog of the update.

## Contribution

Feel free to contact [@nextzlog](https://twitter.com/nextzlog) on Twitter.

## License

### Author

[無線部開発班](https://pafelog.net)

- JG1VPP
- JS2FVO
- JJ1GUJ

### Clauses

[MIT License](LICENSE)
