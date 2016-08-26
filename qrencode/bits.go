package qrencode

import (
	"bytes"
	"image"
	"image/color"
	"io"
)

// The test benchmark shows that encoding with boolBitVector/boolBitGrid is
// twice as fast as byteBitVector/byteBitGrid and uin32BitVector/uint32BitGrid.

type BitVector struct {
	boolBitVector
}

type BitGrid struct {
	boolBitGrid
}

func (v *BitVector) AppendBits(b BitVector) {
	v.boolBitVector.AppendBits(b.boolBitVector)
}

func NewBitGrid(width, height int) *BitGrid {
	return &BitGrid{newBoolBitGrid(width, height)}
}

/*
type BitVector struct {
	byteBitVector
}

type BitGrid struct {
	byteBitGrid
}

func (v *BitVector) AppendBits(b BitVector) {
	v.byteBitVector.AppendBits(b.byteBitVector)
}

func NewBitGrid(width, height int) *BitGrid {
	return &BitGrid{newByteBitGrid(width, height)}
}
*/

/*
type BitVector struct {
	uint32BitVector
}

type BitGrid struct {
	uint32BitGrid
}

func (v *BitVector) AppendBits(b BitVector) {
	v.uint32BitVector.AppendBits(b.uint32BitVector)
}

func NewBitGrid(width, height int) *BitGrid {
	return &BitGrid{newUint32BitGrid(width, height)}
}
*/

func (v *BitVector) String() string {
	b := bytes.Buffer{}
	for i, l := 0, v.Length(); i < l; i++ {
		if v.Get(i) {
			b.WriteString("1")
		} else {
			b.WriteString("0")
		}
	}
	return b.String()
}

func (g *BitGrid) String() string {
	b := bytes.Buffer{}
	for y, w, h := 0, g.Width(), g.Height(); y < h; y++ {
		for x := 0; x < w; x++ {
			if g.Empty(x, y) {
				b.WriteString(" ")
			} else if g.Get(x, y) {
				b.WriteString("#")
			} else {
				b.WriteString("_")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}

// Outputs the Code in UTF8 Block Characters. Each Bit is half a character
func (g *BitGrid) WriteUtf8BlockChars(w io.Writer, inverse bool) {
	var upperHalfBlock []byte
	var lowerHalfBlock []byte
	var completeBlock []byte
	var space []byte

	if inverse {
		upperHalfBlock = []byte("\u2584")
		lowerHalfBlock = []byte("\u2580")
		completeBlock = []byte("\u0020")
		space = []byte("\u2588")

	} else {
		upperHalfBlock = []byte("\u2580")
		lowerHalfBlock = []byte("\u2584")
		completeBlock = []byte("\u2588")
		space = []byte("\u0020")

	}
	newline := []byte("\n")

	w.Write(space)
	for i := 0; i < g.Width(); i++ {
		if g.Get(i, 0) {
			w.Write(lowerHalfBlock)
		} else {
			w.Write(space)
		}
	}
	w.Write(space)
	w.Write(newline)

	for i := 1; i < g.Height()-1; i = i + 2 {
		w.Write(space)
		for j := 0; j < g.Width(); j++ {
			if g.Get(j, i) {
				if g.Get(j, i+1) {
					w.Write(completeBlock)
				} else {
					w.Write(upperHalfBlock)
				}
			} else {
				if g.Get(j, i+1) {
					w.Write(lowerHalfBlock)
				} else {
					w.Write(space)
				}
			}
		}
		w.Write(space)
		w.Write(newline)
	}

	w.Write(lowerHalfBlock)
	for i := 0; i < g.Width(); i++ {
		w.Write(lowerHalfBlock)
	}
	w.Write(lowerHalfBlock)
	w.Write(newline)

}

// Encode the Grid in ANSI escape sequences and set the background according
// to the values in the BitGrid surrounded by a white frame
func (g *BitGrid) WriteTerminalOutput(w io.Writer) {
	white := []byte("\033[47m  \033[0m")
	black := []byte("\033[40m  \033[0m")
	newline := []byte("\n")

	w.Write(white)
	for i := 0; i <= g.Width(); i++ {
		w.Write(white)
	}
	w.Write(newline)

	for i := 0; i < g.Height(); i++ {
		w.Write(white)
		for j := 0; j < g.Width(); j++ {
			if g.Get(j, i) {
				w.Write(black)
			} else {
				w.Write(white)
			}
		}
		w.Write(white)
		w.Write(newline)
	}
	w.Write(white)
	for i := 0; i <= g.Width(); i++ {
		w.Write(white)
	}
	w.Write(newline)
}

// Return an image of the grid, with black blocks for true items and
// white blocks for false items, with the given block size and a
// default margin.
func (g *BitGrid) Image(blockSize int) image.Image {
	return g.ImageWithMargin(blockSize, 4)
}

// Return an image of the grid, with black blocks for true items and
// white blocks for false items, with the given block size and margin.
func (g *BitGrid) ImageWithMargin(blockSize, margin int) image.Image {
	width := blockSize * (2*margin + g.Width())
	height := blockSize * (2*margin + g.Height())
	i := image.NewGray16(image.Rect(0, 0, width, height))
	for y := 0; y < blockSize*margin; y++ {
		for x := 0; x < width; x++ {
			i.Set(x, y, color.White)
			i.Set(x, height-1-y, color.White)
		}
	}
	for y := blockSize * margin; y < height-blockSize*margin; y++ {
		for x := 0; x < blockSize*margin; x++ {
			i.Set(x, y, color.White)
			i.Set(width-1-x, y, color.White)
		}
	}
	for y, w, h := 0, g.Width(), g.Height(); y < h; y++ {
		for x := 0; x < w; x++ {
			x0 := blockSize * (x + margin)
			y0 := blockSize * (y + margin)
			c := color.White
			if g.Get(x, y) {
				c = color.Black
			}
			for dy := 0; dy < blockSize; dy++ {
				for dx := 0; dx < blockSize; dx++ {
					i.Set(x0+dx, y0+dy, c)
				}
			}
		}
	}
	return i
}
