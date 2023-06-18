/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/
package morse

import (
	"github.com/r9y9/gossp"
	"github.com/r9y9/gossp/stft"
	"math"
	"sort"
)

/*
モールス信号の文字列です。
*/
type Message struct {
	Data []float64
	Code string
	Freq int
	Life int
}

/*
モールス信号の解析器です。
*/
type Decoder struct {
	Life int
	Iter int
	Bias int
	Gain float64
	Mute float64
	Band int
	Rate int
	Hold int
	STFT *stft.STFT
	prev []Message
	wave []float64
}

func (d *Decoder) binary(signal []float64) (result []*step) {
	key := make([]float64, len(signal))
	max := max64(signal)
	for idx, val := range signal {
		key[idx] = val * math.Min(d.Gain, max/val)
	}
	gmm := means{X: key}
	gmm.optimize(d.Iter)
	result = gmm.steps()
	return
}

func (d *Decoder) detect(signal []float64) (result Message) {
	result.Data = make([]float64, len(signal))
	copy(result.Data, signal)
	steps := d.binary(signal)
	tones := make([]float64, 0)
	if len(steps) >= 1 {
		for idx, s := range steps[1:] {
			s.span = float64(s.time - steps[idx].time)
			if s.down {
				tones = append(tones, s.span)
			}
		}
	}
	if len(tones) >= 1 {
		gmm := &means{X: tones}
		gmm.optimize(d.Iter)
		for _, s := range steps[1:] {
			if s.down {
				result.Code += s.tone(gmm.class(s.span))
			} else {
				result.Code += s.mute(gmm.extra(s.span))
			}
		}
	}
	return
}

func (d *Decoder) search(series [][]float64) (result []int) {
	cut := d.Band * d.STFT.FrameLen / d.Rate
	pow := make([]float64, d.STFT.FrameLen/2)
	for _, sp := range series {
		for idx, val := range sp[:len(pow)] {
			pow[idx] += val * val
		}
	}
	top := 0.0
	pos := 0
	bit := make(map[int]bool)
	lev := d.Mute * sum64(pow[cut:])
	for idx, val := range pow[d.Bias:cut] {
		if val > top {
			top = val
			pos = idx
		} else if val < lev && top > lev {
			bit[d.Bias+pos] = true
			top = 0
			pos = 0
		}
	}
	for _, prev := range d.prev {
		bit[prev.Freq] = true
	}
	for freq := range bit {
		result = append(result, freq)
	}
	sort.Ints(result)
	return
}

func (d *Decoder) scan(signal []float64) (result []Message) {
	spec, _ := gossp.SplitSpectrogram(d.STFT.STFT(signal))
	wave := make([]float64, len(spec))
	for _, idx := range d.search(spec) {
		for t, s := range spec {
			wave[t] = s[idx]
		}
		if m := d.detect(wave); m.Code != "" {
			m.Freq = idx
			result = append(result, m)
		}
	}
	return
}

func (d *Decoder) next(signal []float64) (result []Message) {
	shift := d.STFT.FrameShift
	if len(d.wave) > d.Hold {
		d.wave = d.wave[len(d.wave)-d.Hold:]
	}
	d.wave = append(d.wave, signal...)
	for _, next := range d.scan(d.wave) {
		for _, prev := range d.prev {
			if next.Freq == prev.Freq {
				drop := len(next.Data) - (len(signal) / shift)
				data := append(prev.Data, next.Data[drop:]...)
				next = d.detect(data)
				next.Freq = prev.Freq
				next.Life = prev.Life
			}
		}
		next.Life++
		result = append(result, next)
	}
	d.prev = result
	d.wave = signal
	return
}

/*
音声からモールス信号の文字列を抽出します。
複数の周波数のモールス信号を分離できます。
*/
func (d *Decoder) Read(signal []float64) (result []Message) {
	for _, next := range d.next(signal) {
		if next.Life >= d.Life {
			result = append(result, next)
		}
	}
	return
}

/*
推奨の設定が適用された解析器を構築します。
*/
func DefaultDecoder(SamplingRateInHz int) (decoder Decoder) {
	return Decoder{
		Life: 3,
		Iter: 5,
		Bias: 5,
		Gain: 2,
		Mute: 10,
		Band: 1000,
		Rate: SamplingRateInHz,
		Hold: SamplingRateInHz,
		STFT: stft.New(SamplingRateInHz/100, 2048),
	}
}
