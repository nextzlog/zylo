/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22
 * License: The MIT License since 2021 October 28 (see LICENSE)
 * Author: Journal of Hamradio Informatics (http://pafelog.net)
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
		Thre: 0.03,
		STFT: &stft.STFT{
			FrameShift: RATE / 100,
			FrameLen:   4096,
			Window:     window.CreateHanning(4096),
		},
	}
	tone := encoder.Tone(TextToCode(TEXT))
	for _, text := range decoder.Read(tone) {
		if CodeToText(text) != TEXT {
			t.Errorf("%s != %s", text, TEXT)
		} else {
			return
		}
	}
	t.Error("no text decoded successfully")
}
