zLog+ ZyLO
====

![image](https://img.shields.io/badge/Go-1.16-red.svg)
![image](https://img.shields.io/badge/Rust-1.51-red.svg)
![image](https://img.shields.io/badge/Delphi-10.4-red.svg)
![image](https://img.shields.io/badge/license-GPL3-darkblue.svg)
![badge](https://github.com/nextzlog/zylo/actions/workflows/build.yaml/badge.svg)

ZyLO is a plugin mechanism for [zLog](http://zlog.org) based on DLLs, that is a simple but powerful logging software for amateur radio contests, originally developed at the [University of Tokyo Amateur Radio Club](http://ja1zlo.u-tokyo.org), and loved by many ham contesters for 30 years.

## Features

- helps realize flexible, dynamic definition of amateur radio contests.
- helps customize the zLog import/export formats.

## Documents

- [GoDoc](https://nextzlog.github.io/zylo)
- [Wiki](https://github.com/nextzlog/zylo/wiki)

## Build DLL

- First, download the `zbuild` tool [here](https://github.com/nextzlog/zylo/releases/tag/zbuild) and place it in the working directory.
- For Windows,

```bat
> choco install mingw
> choco install golang
> zbuild-windows.exe
```

- For Ubuntu,

```sh
$ apt install gcc-mingw-w64
$ apt install golang-go
$ ./zbuild-ubuntu
```

- For macOS,

```sh
$ brew install mingw-w64
$ brew install go
$ ./zbuild-macos
```

- `zbuild` creates `zutils.go` and `zutils.h` and then compiles the Go files to create a DLL.

## Use DLL

- First, download [zLog](https://github.com/jr8ppg/zlog/releases).
- Place the plugin DLL in the same directory as `zlog.exe` and add the following lines to `zlog.ini`.

```ini
[zylo]
DLLs=foo.dll,bar.dll,baz.dll
```

- To define a contest that uses the DLL, add the following lines to the CFG file.

```
exit
dll foo.dll
```

- Start zLog (and select the contest CFG file that uses the DLL).

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
