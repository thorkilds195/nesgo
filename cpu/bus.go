package cpu

const RAM uint16 = 0x0000
const RAM_MIRRORS_END uint16 = 0x1FFF
const PPU_REGISTERS uint16 = 0x2000
const PPU_REGISTERS_MIRRORS_END uint16 = 0x3FFF

type Bus struct {
	cpu_vram [2048]uint8
	rom      *Rom
}

func InitBus(r *Rom) *Bus {
	return &Bus{
		rom: r,
	}
}

func (b *Bus) MemRead(addr uint16) uint8 {
	if addr >= RAM && addr <= RAM_MIRRORS_END {
		mirr_address_down := addr & 0b00000111_11111111
		return b.cpu_vram[mirr_address_down]
	} else if addr >= PPU_REGISTERS && addr <= PPU_REGISTERS_MIRRORS_END {
		_ = addr & 0b00100000_00000111
		panic("No PPU supported yet")
	} else if addr >= 0x8000 && addr <= 0xFFFF {
		return b.readPgrRom(addr)
	}
	return 0
}

func (b *Bus) MemWrite(addr uint16, val uint8) {

	if addr >= RAM && addr <= RAM_MIRRORS_END {
		mirr_address_down := addr & 0b00000111_11111111
		b.cpu_vram[mirr_address_down] = val
	} else if addr >= PPU_REGISTERS && addr <= PPU_REGISTERS_MIRRORS_END {
		_ = addr & 0b00100000_00000111
		panic("No PPU supported yet")
	} else if addr >= 0x8000 && addr <= 0xFFFF {
		panic("Attempt to write to rom space")
	}

}

func (b *Bus) readPgrRom(addr uint16) uint8 {
	addr -= 0x8000
	if len(b.rom.prg_rom) == 0x4000 && addr >= 0x4000 {
		addr = addr % 0x4000
	}
	return b.rom.prg_rom[addr]
}
