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

func (f *Frame) Render(p *PPU) {
	// Render background
	bank := p.ctrl.BnkdPatternAddress()
	for i := 0; i < 0x3C0; i++ {
		tile_nr := uint16(p.vram[i])
		tile_column := i % 32
		tile_row := i / 32
		tile := p.chr_rom[(bank + tile_nr*16):(bank + tile_nr*16 + 16)]
		palette := bgPallete(p, uint(tile_column), uint(tile_row))
		for y := 0; y <= 7; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				value := (1&lower)<<1 | (1 & upper)
				upper = upper >> 1
				lower = lower >> 1
				var rgb RGB
				switch value {
				case 0:
					rgb = SYSTEM_PALLETE[palette[0]]
				case 1:
					rgb = SYSTEM_PALLETE[palette[1]]
				case 2:
					rgb = SYSTEM_PALLETE[palette[2]]
				case 3:
					rgb = SYSTEM_PALLETE[palette[3]]
				default:
					panic("Not valid rgb rom")
				}
				f.SetPixel(uint32(tile_column*8+x), uint32(tile_row*8+y), rgb)
			}
		}
	}
	// Render sprites
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
		paletteIdx := p.oam_data[i+2] & 0b11
		sprPallete := spritePallete(p, paletteIdx)
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
				skip := false
				switch value {
				case 0:
					skip = true
				case 1:
					rgb = SYSTEM_PALLETE[sprPallete[1]]
				case 2:
					rgb = SYSTEM_PALLETE[sprPallete[2]]
				case 3:
					rgb = SYSTEM_PALLETE[sprPallete[3]]
				default:
					panic("Not valid rgb rom")
				}
				if skip {
					continue
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

func bgPallete(p *PPU, tile_col uint, tile_row uint) [4]uint8 {
	tableIdx := tile_row/4*8 + tile_col/4
	attrByte := p.vram[0x3C0+tableIdx]
	col_idx := tile_col % 4 / 2
	row_idx := tile_row % 4 / 2
	var palletIdx uint8
	if col_idx == 0 && row_idx == 0 {
		palletIdx = attrByte & 0b11
	} else if col_idx == 1 && row_idx == 0 {
		palletIdx = (attrByte >> 2) & 0b11
	} else if col_idx == 0 && row_idx == 1 {
		palletIdx = (attrByte >> 4) & 0b11
	} else if col_idx == 1 && row_idx == 1 {
		palletIdx = (attrByte >> 6) & 0b11
	}
	palleteStart := 1 + uint(palletIdx)*4
	return [4]uint8{p.palette_table[0], p.palette_table[palleteStart], p.palette_table[palleteStart+1], p.palette_table[palleteStart+2]}
}

func spritePallete(p *PPU, palette_idx uint8) [4]uint8 {
	start := 0x11 + uint(palette_idx*4)
	return [4]uint8{0, p.palette_table[start], p.palette_table[start+1], p.palette_table[start+2]}
}
