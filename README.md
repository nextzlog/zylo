zLog+ ZyLO for Windows
====

![image](https://img.shields.io/badge/Go-1.17-red.svg)
![image](https://img.shields.io/badge/Rust-1.55-red.svg)
![image](https://img.shields.io/badge/Delphi-10.4-red.svg)
![image](https://img.shields.io/badge/license-GPL3-darkblue.svg)
![badge](https://github.com/nextzlog/zylo/actions/workflows/build.yaml/badge.svg)

[ZyLO](https://github.com/nextzlog/zylo) is a plugin and marketplace mechanism integrated into [zLog](http://zlog.org).
zLog is a simple yet powerful logging software for ham radio contests, which has been loved by many users for 30 years.

## Documents

- [API](https://nextzlog.github.io/zylo)

## Samples

- `format.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/format))
- `latest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/latest))
- `toasty.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/toasty))
- `hstest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/hstest))
- `rttest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/rttest))
- `yltest.dll` ([Project](https://github.com/nextzlog/zylo/tree/master/rules/yltest))

## Build DLL

- First, clone the sample project [toasty](https://github.com/nextzlog/zylo/tree/master/rules/toasty) as follows.

```sh
$ git clone https://github.com/nextzlog/zylo
$ cd zylo/rules/toasty
```

- Next, create `toasty.dll` by the [`zbuild`](https://github.com/nextzlog/zylo/releases/tag/zbuild) command as follows, and you will find `toasty.dll` in the directory.

### Build DLL on Linux

```sh
$ apt install gcc-mingw-w64 golang-go
$ ./zbuild-linux compile
```

### Build DLL on macOS

```sh
$ brew install mingw-w64 go
$ ./zbuild-macos compile
```

### Build DLL on Windows

```bat
> choco install mingw golang
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

- Ask one of the [market managers](https://github.com/nextzlog/zylo/blob/master/market.list) to add the release URL of the DLL to `market.toml`.

```toml
[pkg.toasty]
tag = "title"
msg = "description"
web = "website URL"
use = ["dll.toasty"]

[dll.toasty]
url = "release URL"
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

- This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

- This program is distributed in the hope that it will be useful, but **without any warranty**; without even the implied warranty of **merchantability or fitness for a particular purpose**.
See the GNU General Public License for more details.

- You should have received a copy of the GNU General Public License along with this program.
If not, see <http://www.gnu.org/licenses/>.
