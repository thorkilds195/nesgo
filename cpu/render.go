package cpu

const WITDH = 256
const HEIGHT = 240

type Frame struct {
	Data []uint8
}

func NewFrame() *Frame {
	return &Frame{
		Data: make([]uint8, WITDH*HEIGHT*4),
	}
}

func (f *Frame) SetPixel(x uint32, y uint32, rgb RGB) {
	base := y*4*WITDH + x*4
	if base+3 < uint32(len(f.Data)) {
		f.Data[base] = rgb.R
		f.Data[base+1] = rgb.G
		f.Data[base+2] = rgb.B
		f.Data[base+3] = 0xFF
	}
}
