// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package image

import (
	"image/color"
	"testing"
)

type image interface {
	Image
	Opaque() bool
	Set(int, int, color.Color)
	SubImage(Rectangle) Image
}

func cmp(cm color.Model, c0, c1 color.Color) bool {
	r0, g0, b0, a0 := cm.Convert(c0).RGBA()
	r1, g1, b1, a1 := cm.Convert(c1).RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1 && a0 == a1
}

var testImages = []struct {
	name  string
	image func() image
}{
	{"rgba", func() image { return NewRGBA(Rect(0, 0, 10, 10)) }},
	{"rgba64", func() image { return NewRGBA64(Rect(0, 0, 10, 10)) }},
	{"nrgba", func() image { return NewNRGBA(Rect(0, 0, 10, 10)) }},
	{"nrgba64", func() image { return NewNRGBA64(Rect(0, 0, 10, 10)) }},
	{"alpha", func() image { return NewAlpha(Rect(0, 0, 10, 10)) }},
	{"alpha16", func() image { return NewAlpha16(Rect(0, 0, 10, 10)) }},
	{"gray", func() image { return NewGray(Rect(0, 0, 10, 10)) }},
	{"gray16", func() image { return NewGray16(Rect(0, 0, 10, 10)) }},
	{"paletted", func() image {
		return NewPaletted(Rect(0, 0, 10, 10), color.Palette{
			Transparent,
			Opaque,
		})
	}},
}

func TestImage(t *testing.T) {
	for _, tc := range testImages {
		m := tc.image()
		if !Rect(0, 0, 10, 10).Eq(m.Bounds()) {
			t.Errorf("%T: want bounds %v, got %v", m, Rect(0, 0, 10, 10), m.Bounds())
			continue
		}
		if !cmp(m.ColorModel(), Transparent, m.At(6, 3)) {
			t.Errorf("%T: at (6, 3), want a zero color, got %v", m, m.At(6, 3))
			continue
		}
		m.Set(6, 3, Opaque)
		if !cmp(m.ColorModel(), Opaque, m.At(6, 3)) {
			t.Errorf("%T: at (6, 3), want a non-zero color, got %v", m, m.At(6, 3))
			continue
		}
		if !m.SubImage(Rect(6, 3, 7, 4)).(image).Opaque() {
			t.Errorf("%T: at (6, 3) was not opaque", m)
			continue
		}
		m = m.SubImage(Rect(3, 2, 9, 8)).(image)
		if !Rect(3, 2, 9, 8).Eq(m.Bounds()) {
			t.Errorf("%T: sub-image want bounds %v, got %v", m, Rect(3, 2, 9, 8), m.Bounds())
			continue
		}
		if !cmp(m.ColorModel(), Opaque, m.At(6, 3)) {
			t.Errorf("%T: sub-image at (6, 3), want a non-zero color, got %v", m, m.At(6, 3))
			continue
		}
		if !cmp(m.ColorModel(), Transparent, m.At(3, 3)) {
			t.Errorf("%T: sub-image at (3, 3), want a zero color, got %v", m, m.At(3, 3))
			continue
		}
		m.Set(3, 3, Opaque)
		if !cmp(m.ColorModel(), Opaque, m.At(3, 3)) {
			t.Errorf("%T: sub-image at (3, 3), want a non-zero color, got %v", m, m.At(3, 3))
			continue
		}
		// Test that taking an empty sub-image starting at a corner does not panic.
		m.SubImage(Rect(0, 0, 0, 0))
		m.SubImage(Rect(10, 0, 10, 0))
		m.SubImage(Rect(0, 10, 0, 10))
		m.SubImage(Rect(10, 10, 10, 10))
	}
}

func Test16BitsPerColorChannel(t *testing.T) {
	testColorModel := []color.Model{
		color.RGBA64Model,
		color.NRGBA64Model,
		color.Alpha16Model,
		color.Gray16Model,
	}
	for _, cm := range testColorModel {
		c := cm.Convert(color.RGBA64{0x1234, 0x1234, 0x1234, 0x1234}) // Premultiplied alpha.
		r, _, _, _ := c.RGBA()
		if r != 0x1234 {
			t.Errorf("%T: want red value 0x%04x got 0x%04x", c, 0x1234, r)
			continue
		}
	}
	testImage := []image{
		NewRGBA64(Rect(0, 0, 10, 10)),
		NewNRGBA64(Rect(0, 0, 10, 10)),
		NewAlpha16(Rect(0, 0, 10, 10)),
		NewGray16(Rect(0, 0, 10, 10)),
	}
	for _, m := range testImage {
		m.Set(1, 2, color.NRGBA64{0xffff, 0xffff, 0xffff, 0x1357}) // Non-premultiplied alpha.
		r, _, _, _ := m.At(1, 2).RGBA()
		if r != 0x1357 {
			t.Errorf("%T: want red value 0x%04x got 0x%04x", m, 0x1357, r)
			continue
		}
	}
}

func BenchmarkAt(b *testing.B) {
	for _, tc := range testImages {
		b.Run(tc.name, func(b *testing.B) {
			m := tc.image()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.At(4, 5)
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	c := color.Gray{0xff}
	for _, tc := range testImages {
		b.Run(tc.name, func(b *testing.B) {
			m := tc.image()
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m.Set(4, 5, c)
			}
		})
	}
}

func BenchmarkRGBAAt(b *testing.B) {
	m := NewRGBA(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.RGBAAt(4, 5)
	}
}

func BenchmarkRGBASetRGBA(b *testing.B) {
	m := NewRGBA(Rect(0, 0, 10, 10))
	c := color.RGBA{0xff, 0xff, 0xff, 0x13}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetRGBA(4, 5, c)
	}
}

func BenchmarkRGBA64At(b *testing.B) {
	m := NewRGBA64(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.RGBA64At(4, 5)
	}
}

func BenchmarkRGBA64SetRGBA64(b *testing.B) {
	m := NewRGBA64(Rect(0, 0, 10, 10))
	c := color.RGBA64{0xffff, 0xffff, 0xffff, 0x1357}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetRGBA64(4, 5, c)
	}
}

func BenchmarkNRGBAAt(b *testing.B) {
	m := NewNRGBA(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.NRGBAAt(4, 5)
	}
}

func BenchmarkNRGBASetNRGBA(b *testing.B) {
	m := NewNRGBA(Rect(0, 0, 10, 10))
	c := color.NRGBA{0xff, 0xff, 0xff, 0x13}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetNRGBA(4, 5, c)
	}
}

func BenchmarkNRGBA64At(b *testing.B) {
	m := NewNRGBA64(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.NRGBA64At(4, 5)
	}
}

func BenchmarkNRGBA64SetNRGBA64(b *testing.B) {
	m := NewNRGBA64(Rect(0, 0, 10, 10))
	c := color.NRGBA64{0xffff, 0xffff, 0xffff, 0x1357}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetNRGBA64(4, 5, c)
	}
}

func BenchmarkAlphaAt(b *testing.B) {
	m := NewAlpha(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.AlphaAt(4, 5)
	}
}

func BenchmarkAlphaSetAlpha(b *testing.B) {
	m := NewAlpha(Rect(0, 0, 10, 10))
	c := color.Alpha{0x13}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetAlpha(4, 5, c)
	}
}

func BenchmarkAlpha16At(b *testing.B) {
	m := NewAlpha16(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Alpha16At(4, 5)
	}
}

func BenchmarkAlphaSetAlpha16(b *testing.B) {
	m := NewAlpha16(Rect(0, 0, 10, 10))
	c := color.Alpha16{0x13}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetAlpha16(4, 5, c)
	}
}

func BenchmarkGrayAt(b *testing.B) {
	m := NewGray(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.GrayAt(4, 5)
	}
}

func BenchmarkGraySetGray(b *testing.B) {
	m := NewGray(Rect(0, 0, 10, 10))
	c := color.Gray{0x13}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetGray(4, 5, c)
	}
}

func BenchmarkGray16At(b *testing.B) {
	m := NewGray16(Rect(0, 0, 10, 10))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.Gray16At(4, 5)
	}
}

func BenchmarkGraySetGray16(b *testing.B) {
	m := NewGray16(Rect(0, 0, 10, 10))
	c := color.Gray16{0x13}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.SetGray16(4, 5, c)
	}
}
