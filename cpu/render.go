package cpu

import "fmt"

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

func (f *Frame) Render(p *PPU) {
	bank := p.ctrl.BnkdPatternAddress()
	for i := 0; i < 0x3C0; i++ {
		tile_nr := uint16(p.vram[i])
		tile_column := i % 32
		tile_row := i / 32
		tile := p.chr_rom[(bank + tile_nr*16):(bank + tile_nr*16 + 16)]
		if tile_nr > 0 {
			fmt.Println(tile_nr)
		}
		for y := 0; y <= 7; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				value := (1&upper)<<1 | (1 & lower)
				upper = upper >> 1
				lower = lower >> 1
				var rgb RGB
				switch value {
				case 0:
					rgb = SYSTEM_PALLETE[0x01]
				case 1:
					fmt.Println("Hit other")
					rgb = SYSTEM_PALLETE[0x23]
				case 2:
					fmt.Println("Hit other")
					rgb = SYSTEM_PALLETE[0x27]
				case 3:
					fmt.Println("Hit other")
					rgb = SYSTEM_PALLETE[0x30]
				default:
					panic("Not valid rgb rom")
				}
				f.SetPixel(uint32(tile_column*8+x), uint32(tile_row*8+y), rgb)
			}
		}
	}
	return
	for i := len(p.oam_data) - 4; i >= 0; i -= 4 {
		tile_idx := uint16(p.oam_data[i+1])
		tile_x := int(p.oam_data[i+3])
		tile_y := int(p.oam_data[i])

		flip_vertical := false
		if p.oam_data[i+2]>>7&1 == 1 {
			flip_vertical = true
		}
		flip_horizontal := false
		if p.oam_data[i+2]>>6&1 == 1 {
			flip_horizontal = true
		}
		bank := p.ctrl.SprtPatternAddress()

		tile := p.chr_rom[(bank + tile_idx*16):(bank + (tile_idx)*16 + 16)]

		for y := 0; y < 8; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				value := (1&lower)<<1 | (1 & upper)
				upper = upper >> 1
				lower = lower >> 1
				var rgb RGB
				switch value {
				case 0:
					continue
				case 1:
					rgb = SYSTEM_PALLETE[0x23]
				case 2:
					rgb = SYSTEM_PALLETE[0x27]
				case 3:
					rgb = SYSTEM_PALLETE[0x30]
				default:
					panic("Not valid rgb rom")
				}
				if !flip_horizontal && !flip_vertical {
					f.SetPixel(uint32(tile_x+x), uint32(tile_y+y), rgb)
				} else if flip_horizontal && !flip_vertical {
					f.SetPixel(uint32(tile_x+7-x), uint32(tile_y+y), rgb)
				} else if !flip_horizontal && flip_vertical {
					f.SetPixel(uint32(tile_x+x), uint32(tile_y+7-y), rgb)
				} else if flip_horizontal && flip_vertical {
					f.SetPixel(uint32(tile_x+7-x), uint32(tile_y+7-y), rgb)
				}
			}
		}
	}
}
