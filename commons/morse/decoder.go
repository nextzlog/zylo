/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22
 * License: The MIT License since 2021 October 28 (see LICENSE)
 * Author: Journal of Hamradio Informatics (http://pafelog.net)
*******************************************************************************/
package morse

import (
	"github.com/r9y9/gossp"
	"github.com/r9y9/gossp/stft"
	"github.com/thoas/go-funk"
	"math"
	"strings"
)

const MIN_RELIABLE_DOT = 2

/*
 モールス信号の文字列です。
*/
type Message struct {
	Code string
	Freq int
}

/*
 文章の区切りを検出した場合は真を返します。
*/
func (m *Message) Finish() (finish bool) {
	return strings.HasSuffix(m.Code, " ; ")
}

/*
 この文字列の末尾に次の文字列を結合します。
*/
func (prev *Message) Join(next *Message) {
	best := 1
	save := 0
	skip := 0
	head := strings.Split(prev.Code, " ")
	tail := strings.Split(next.Code, " ")
	for idx := 0; idx < len(head); idx++ {
		rel := len(head)-idx >= len(tail)
		zip := funk.Zip(head[idx:], tail)
		pos := 0
		neg := 0
		for idx, tuple := range zip {
			p := tuple.Element1.(string)
			n := tuple.Element2.(string)
			if p == n {
				pos += 1
			} else if rel {
				neg = idx + 1
			} else {
				neg = idx
			}
		}
		if pos >= best {
			best = pos
			save = idx + neg
			skip = neg
		}
	}
	p := strings.Join(head[:save], " ")
	n := strings.Join(tail[skip:], " ")
	prev.Code = strings.TrimSpace(p + n)
	next.Code = strings.TrimSpace(p + n)
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
	max := funk.MaxFloat64(signal)
	for idx, val := range signal {
		signal[idx] = val * math.Min(d.Gain, max/val)
	}
	gmm := means{X: signal}
	gmm.optimize(d.Iter)
	pre := 0
	for idx, val := range signal {
		cls := gmm.class(val)
		if pre != cls {
			result = append(result, &step{
				time: idx,
				down: cls == 0,
			})
		}
		pre = cls
	}
	return append(result, &step{time: len(signal)})
}

func (d *Decoder) detect(signal []float64) (result Message) {
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
	Decoder Decoder
	samples []float64
}

/*
 規定の設定が適用された解析器を構築します。
*/
func DefaultMonitor(SamplingRateInHz int) (monitor Monitor) {
	shift := int(math.Round(0.02 * float64(SamplingRateInHz)))
	return Monitor{
		MaxHold: 10 * SamplingRateInHz,
		Decoder: Decoder{
			Iter: 10,
			Bias: 10,
			Gain: 2,
			Thre: 0.03,
			STFT: stft.New(shift, 2048),
		},
	}
}

/*
 音声からモールス信号の文字列を抽出します。
*/
func (m *Monitor) Read(signal []float64) (result []Message) {
	m.samples = append(m.samples, signal...)
	result = m.Decoder.Read(m.samples)
	if len(m.samples) > m.MaxHold {
		m.samples = m.samples[len(signal):]
	}
	return
}
