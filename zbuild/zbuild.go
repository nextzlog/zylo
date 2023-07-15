/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/
package main

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

//go:embed INSTALL.yaml
var install string

//go:embed REPLACE.yaml
var replace string

//go:embed _build.go
var version string

var installCmd map[bool][]map[string][]string
var replaceCmd []string

func init() {
	yaml.UnmarshalStrict([]byte(install), &installCmd)
	yaml.UnmarshalStrict([]byte(replace), &replaceCmd)
}

func main() {
	setupFlags := []cli.Flag{
		cli.BoolFlag{
			Name: "sudo",
		},
	}
	buildFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "version",
			Value: "2.8",
		},
	}
	setupCmd := cli.Command{
		Name:   "setup",
		Flags:  setupFlags,
		Action: setup,
	}
	buildCmd := cli.Command{
		Name:   "build",
		Flags:  buildFlags,
		Action: build,
	}
	commands := []cli.Command{
		setupCmd,
		buildCmd,
	}
	app := cli.App{
		Name:     "zbuild",
		Commands: commands,
	}
	app.Run(os.Args)
}

func call(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stdout
	cmd.Stdout = os.Stderr
	return cmd.Run()
}

func setup(c *cli.Context) error {
	cmds := installCmd[c.Bool("sudo")]
	for _, cmd := range cmds {
		for name, arg := range cmd {
			if call(name, arg...) == nil {
				return nil
			}
		}
	}
	return errors.New("failed")
}

func build(c *cli.Context) error {
	name, _ := os.Getwd()
	name = filepath.Base(name)
	dllName := fmt.Sprintf("%s.dll", name)
	md5Name := fmt.Sprintf("%s.md5", name)
	os.Setenv("GOOS", "windows")
	os.Setenv("GOARCH", "amd64")
	os.Setenv("CGO_ENABLED", "1")
	os.Setenv("GOPROXY", "direct")
	os.Setenv("CC", "x86_64-w64-mingw32-gcc")
	main, _ := os.Create("main.go")
	defer main.Close()
	main.WriteString(fmt.Sprintf(version, c.String("version")))
	call("go", "mod", "init", name)
	call("go", replaceCmd...)
	call("go", "get", "-u", "all")
	call("go", "mod", "tidy")
	call("go", "build", "-o", dllName, "-buildmode=c-shared")
	call("upx", dllName)
	file, _ := os.Open(dllName)
	save, _ := os.Create(md5Name)
	defer file.Close()
	defer save.Close()
	hash := md5.New()
	io.Copy(hash, file)
	save.WriteString(hex.EncodeToString(hash.Sum(nil)[:]))
	return nil
}
