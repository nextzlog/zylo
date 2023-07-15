ZyLO
====

![image](https://img.shields.io/badge/Go-1.17-red.svg)
![image](https://img.shields.io/badge/license-MIT-darkblue.svg)
![badge](https://github.com/nextzlog/zylo/actions/workflows/build.yaml/badge.svg)

ZyLO is an official plugin SDK integrated into [zLog v2.8 or later](https://use.zlog.org).
zLog is a powerful amateur radio logging software for contests and has been loved by millions for 30 years.

## Documents

- [Slides](https://nextzlog.dev/zylo.pdf)
- [Manual](https://nextzlog.github.io/zylo)

## Install

```sh
$ go install github.com/nextzlog/zylo/zbuild@HEAD
```

## Usage

To install Go, MinGW and UPX:

```sh
$ zbuild setup
```

To create a DLL for zLog 2.8:

```sh
$ zbuild build --version 2.8
```

## Contribution

Feel free to make issues at [nextzlog/todo](https://github.com/nextzlog/todo).
Follow [@nextzlog](https://twitter.com/nextzlog) on Twitter.

## License

### Author

[無線部開発班](https://nextzlog.dev)

- JG1VPP
- JS2FVO
- JJ1GUJ

### Clauses

[MIT License](LICENSE)
