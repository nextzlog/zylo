/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/
package morse

import (
	"bufio"
	_ "embed"
	"strings"
)

//go:embed latin.dat
var morse string
var reverse = make(map[string]rune)
var forward = make(map[rune]string)

func init() {
	reader := strings.NewReader(morse)
	stream := bufio.NewScanner(reader)
	for stream.Scan() {
		val := stream.Text()
		reverse[val[1:]] = rune(val[0])
		forward[rune(val[0])] = val[1:]
	}
}

/*
 モールス信号の文字列を欧文に変換します。
*/
func CodeToText(code string) (result string) {
	for _, s := range strings.Split(code, " ") {
		if val, ok := reverse[s]; ok {
			result += string(val)
		} else {
			result += "?"
		}
	}
	return
}

/*
 欧文をモールス信号の文字列に変換します。
*/
func TextToCode(text string) (result string) {
	for _, s := range text {
		result += " " + forward[s]
	}
	if result != "" {
		return result[1:]
	} else {
		return
	}
}
