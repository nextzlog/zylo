/*******************************************************************************
 * Amateur Radio Operational Logging Software 'ZyLO' since 2020 June 22
 * License: The MIT License since 2021 October 28 (see LICENSE)
 * Author: Journal of Hamradio Informatics (http://pafelog.net)
*******************************************************************************/
package morse

import (
	"encoding/binary"
	"github.com/r9y9/gossp"
	"github.com/r9y9/gossp/stft"
	"github.com/thoas/go-funk"
	"math"
	"sort"
	"strings"
)

const MIN_RELIABLE_DOT = 2

type Samples []float64

type means struct {
	X Samples
	m Samples
}

func (b *means) optimize(iterations int) {
	b.m = append(b.m, funk.MinFloat64(b.X))
	b.m = append(b.m, funk.MaxFloat64(b.X))
	for i := 0; i < iterations; i++ {
		newN := make(Samples, len(b.m))
		newM := make(Samples, len(b.m))
		for _, x := range b.X {
			k := b.class(x)
			newN[k] += 1
			newM[k] += x
		}
		for k, m := range newM {
			b.m[k] = m / newN[k]
		}
	}
	sort.Sort(sort.Float64Slice(b.m))
}

func (b *means) class(x float64) (k int) {
	lo := math.Abs(x - b.m[0])
	hi := math.Abs(x - b.m[1])
	if lo < hi {
		return 0
	} else {
		return 1
	}
}

func (b *means) extra(x float64) (k int) {
	hi := math.Abs(x - b.m[1]*1)
	ex := math.Abs(x - b.m[1]*3)
	if hi < ex {
		return b.class(x)
	} else {
		return 2
	}
}

type step struct {
	time int
	down bool
	span float64
}

func (s *step) tone(class int) string {
	switch class {
	case 0:
		return "."
	case 1:
		return "_"
	default:
		return "_"
	}
}

func (s *step) mute(class int) string {
	switch class {
	case 0:
		return ""
	case 1:
		return " "
	default:
		return " ; "
	}
}

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

func (d *Decoder) binary(signal Samples) (result []*step) {
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

func (d *Decoder) detect(signal Samples) (result Message) {
	steps := d.binary(signal)
	tones := make(Samples, 0)
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

func (d *Decoder) search(spectrum Samples) (result []int) {
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
func (d *Decoder) Read(signal Samples) (result []Message) {
	spec, _ := gossp.SplitSpectrogram(d.STFT.STFT(signal))
	dist := make(Samples, d.STFT.FrameLen/2)
	for _, s := range spec {
		for idx, val := range s[d.Bias:len(dist)] {
			dist[idx] += val * val
		}
	}
	buff := make(Samples, len(spec))
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
	samples Samples
}

/*
 規定の設定が適用された解析器を構築します。
*/
func DefaultMonitor(SampleRateInHz int) (monitor Monitor) {
	shift := int(math.Round(0.02 * float64(SampleRateInHz)))
	return Monitor{
		MaxHold: 10 * SampleRateInHz,
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
func (m *Monitor) Read(signal Samples) (result []Message) {
	m.samples = append(m.samples, signal...)
	result = m.Decoder.Read(m.samples)
	if len(m.samples) > m.MaxHold {
		m.samples = m.samples[len(signal):]
	}
	return
}

/*
 音声のバイト表現から音声信号を取得します。
*/
func Read32BitSignedInt(signal []byte) (result []float64) {
	for _, buffer := range funk.Chunk(signal, 4).([][]byte) {
		v := binary.LittleEndian.Uint32(buffer)
		result = append(result, float64(int32(v)))
	}
	return
}
