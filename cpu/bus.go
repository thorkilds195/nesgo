package cpu

const RAM uint16 = 0x0000
const RAM_MIRRORS_END uint16 = 0x1FFF
const PPU_REGISTERS uint16 = 0x2000
const PPU_REGISTERS_MIRRORS_END uint16 = 0x3FFF

type Bus struct {
	cpu_vram     [2048]uint8
	rom          *Rom
	ppu          *PPU
	cycles       uint
	gameCallback func(*PPU)
}

func InitBus(r *Rom, c func(*PPU)) *Bus {
	p := NewPPU(r.chr_rom, r.screen_mirroring)
	return &Bus{
		rom:          r,
		ppu:          p,
		gameCallback: c,
	}
}

func (b *Bus) Tick(cycles uint8) {
	b.cycles += uint(cycles)
	newFrame := b.ppu.Tick(cycles * 3)
	if newFrame {
		b.gameCallback(b.ppu)
	}
}

func (b *Bus) MemRead(addr uint16) uint8 {
	if addr >= RAM && addr <= RAM_MIRRORS_END {
		mirr_address_down := addr & 0b00000111_11111111
		return b.cpu_vram[mirr_address_down]
	} else if addr >= PPU_REGISTERS && addr <= PPU_REGISTERS_MIRRORS_END {
		switch addr {
		case 0x2000 | 0x2001 | 0x2003 | 0x2005 | 0x2006 | 0x4014:
			panic("Attempt to read from write-only PPU address")
		case 0x2002:
			return b.ppu.ReadStatusRegister()
		case 0x2004:
			return b.ppu.ReadOAMData()
		case 0x2007:
			return b.ppu.ReadData()
		default:
			if addr >= 0x2008 && addr <= PPU_REGISTERS_MIRRORS_END {
				mirror_address_down := addr & 0b00100000_00000111
				b.MemRead(mirror_address_down)
			} else if addr >= 0x8000 && addr <= 0xFFFF {
				b.readPgrRom(addr)
			}
		}
	} else if addr >= 0x8000 && addr <= 0xFFFF {
		return b.readPgrRom(addr)
	}
	return 0
}

func (b *Bus) MemWrite(addr uint16, val uint8) {

	if addr >= RAM && addr <= RAM_MIRRORS_END {
		mirr_address_down := addr & 0b11111111111
		b.cpu_vram[mirr_address_down] = val
	} else if addr >= PPU_REGISTERS && addr <= PPU_REGISTERS_MIRRORS_END {
		switch addr {
		case 0x2000:
			b.ppu.WriteToPPUCtrl(val)
		case 0x2001:
			b.ppu.WriteToMask(val)
		case 0x2002:
			panic("Trying to write to ppu status register")
		case 0x2003:
			b.ppu.WriteToOAMAddr(val)
		case 0x2004:
			b.ppu.WriteOAMData(val)
		case 0x2005:
			b.ppu.WriteToScroll(val)
		case 0x2006:
			b.ppu.WriteToPPUAddr(val)
		case 0x2007:
			b.ppu.WriteToData(val)
		case 0x4014:
			buf := [256]uint8{}
			hi := uint16(val) << 8
			for i := range buf {
				buf[i] = b.MemRead(hi + uint16(i))
			}
			b.ppu.WriteToOAMDMA(&buf)
		default:
			mirror_down_addr := addr & 0b00100000_00000111
			b.MemWrite(mirror_down_addr, val)
		}
	} else if addr >= 0x8000 && addr <= 0xFFFF {
		panic("Attempt to write to rom space")
	}
}

func (b *Bus) PollNMIStatus() *uint8 {
	return b.ppu.PollNMIStatus()
}

func (b *Bus) readPgrRom(addr uint16) uint8 {
	addr -= 0x8000
	if len(b.rom.prg_rom) == 0x4000 && addr >= 0x4000 {
		addr = addr % 0x4000
	}
	return b.rom.prg_rom[addr]
}
