/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22nd
 * Released under the MIT License (or GPL v3 until 2021 Oct 28th) (see LICENSE)
 * Univ. Tokyo Amateur Radio Club Development Task Force (https://nextzlog.dev)
*******************************************************************************/
package morse

import (
	"github.com/r9y9/gossp"
	"github.com/r9y9/gossp/stft"
	"github.com/thoas/go-funk"
	"math"
)

const MIN_RELIABLE_DOT = 2

/*
 モールス信号の文字列です。
*/
type Message struct {
	Data []float64
	Code string
	Freq int
	Life int
	Miss int
}

/*
 モールス信号の解析器です。
*/
type Decoder struct {
	Iter int
	Bias int
	Gain float64
	Thre float64
	STFT *stft.STFT
}

func (d *Decoder) binary(signal []float64) (result []*step) {
	key := make([]float64, len(signal))
	max := funk.MaxFloat64(signal)
	for idx, val := range signal {
		key[idx] = val * math.Min(d.Gain, max/val)
	}
	gmm := means{X: key}
	gmm.optimize(d.Iter)
	pre := 0
	for idx, val := range key {
		cls := gmm.class(val)
		if pre != cls {
			result = append(result, &step{
				time: idx,
				down: cls == 0,
			})
		}
		pre = cls
	}
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
		if funk.MinFloat64(gmm.m) > MIN_RELIABLE_DOT {
			for _, s := range steps[1:] {
				if s.down {
					result.Code += s.tone(gmm.class(s.span))
				} else {
					result.Code += s.mute(gmm.extra(s.span))
				}
			}
		}
	}
	return
}

func (d *Decoder) search(spectrum []float64) (result []int) {
	total := funk.SumFloat64(spectrum)
	value := 0.0
	index := 0
	for idx, val := range spectrum {
		if val > value {
			value = val
			index = idx
		} else if value > d.Thre*total {
			result = append(result, index)
			value = 0.0
			index = 0
		}
	}
	return
}

/*
 音声からモールス信号の文字列を抽出します。
 複数の周波数のモールス信号を分離できます。
*/
func (d *Decoder) Read(signal []float64) (result []Message) {
	spec, _ := gossp.SplitSpectrogram(d.STFT.STFT(signal))
	dist := make([]float64, d.STFT.FrameLen/2)
	for _, s := range spec {
		for idx, val := range s[d.Bias:len(dist)] {
			dist[idx] += val * val
		}
	}
	buff := make([]float64, len(spec))
	for _, idx := range d.search(dist) {
		for t, s := range spec {
			buff[t] = s[d.Bias+idx]
		}
		if m := d.detect(buff); m.Code != "" {
			m.Freq = int(d.Bias + idx)
			result = append(result, m)
		}
	}
	return
}

/*
 モールス信号の逐次的な解析器です。
*/
type Monitor struct {
	MaxHold int
	MaxMiss int
	Decoder Decoder
	samples []float64
	targets []Message
}

/*
 規定の設定が適用された解析器を構築します。
*/
func DefaultMonitor(SamplingRateInHz int) (monitor Monitor) {
	shift := int(math.Round(0.01 * float64(SamplingRateInHz)))
	return Monitor{
		MaxHold: 5 * SamplingRateInHz,
		MaxMiss: 5,
		Decoder: Decoder{
			Iter: 5,
			Bias: 5,
			Gain: 2,
			Thre: 0.01,
			STFT: stft.New(shift, 2048),
		},
	}
}

/*
 音声からモールス信号の文字列を抽出します。
*/
func (m *Monitor) Read(signal []float64) (result []Message) {
	shift := m.Decoder.STFT.FrameShift
	m.samples = append(m.samples, signal...)
	if len(m.samples) > m.MaxHold {
		m.samples = m.samples[len(signal)/shift*shift:]
	}
	for _, next := range m.Decoder.Read(m.samples) {
		for _, prev := range m.targets {
			if next.Freq == prev.Freq {
				drop := len(next.Data) - (len(signal) / shift)
				data := append(prev.Data, next.Data[drop:]...)
				next = m.Decoder.detect(data)
				next.Freq = prev.Freq
				next.Life = prev.Life
			}
		}
		next.Life += 1
		result = append(result, next)
	}
	for _, prev := range m.targets {
		miss := true
		for _, next := range result {
			if next.Freq == prev.Freq {
				miss = false
			}
		}
		if miss && prev.Miss < m.MaxMiss {
			prev.Miss += 1
			result = append(result, prev)
		}
	}
	m.targets = result
	return
}
