/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/
package morse

import (
	"github.com/r9y9/gossp/stft"
	"github.com/r9y9/gossp/window"
	"testing"
)

const (
	RATE = 48000
	TEXT = "CQ DE JA1ZLO"
)

func TestEnDe(t *testing.T) {
	encoder := Encoder{
		Freq: 600,
		WPMs: 10,
		Rate: RATE,
	}
	decoder := Decoder{
		Iter: 50,
		Bias: 2,
		Gain: 2,
		Thre: 0.03,
		STFT: &stft.STFT{
			FrameShift: RATE / 100,
			FrameLen:   4096,
			Window:     window.CreateHanning(4096),
		},
	}
	tone := encoder.Tone(TextToCode(TEXT))
	for _, msg := range decoder.Read(tone) {
		if CodeToText(msg.Code) != TEXT {
			t.Errorf("%s != %s", msg.Code, TEXT)
		} else {
			return
		}
	}
	t.Error("no text decoded successfully")
}
