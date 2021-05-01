package telegram

import (
	"image/color"
	"math"

	"golang.org/x/image/colornames"
)

type Align uint8

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

type Header struct {
	Text  string
	Align Align
	Span  float64
}

func (h Header) GetSpan() float64 {
	return math.Max(h.Span, 1.0)
}

type Headers []Header

func (hs Headers) TotalSpan() float64 {
	total := 0.0
	for _, h := range hs {
		total += h.GetSpan()
	}
	return total
}

func (hs Headers) Spans() []float64 {
	var spans []float64
	for _, h := range hs {
		spans = append(spans, h.GetSpan())
	}
	return spans
}

func (hs Headers) GetAligns() []Align {
	var aligns []Align
	for _, h := range hs {
		aligns = append(aligns, h.Align)
	}
	return aligns
}

type Column struct {
	Data  interface{}
	Color color.Color
}

func (c Column) GetColor() color.Color {
	if c.Color == nil {
		return colornames.White
	}
	return c.Color
}

type Columns []Column

type Row struct {
	Columns Columns
}

type Rows []Row

type Footers []Column

type Table struct {
	Headers Headers
	Rows    Rows
	Footers Footers
	Margin  float64
	Padding float64
}
