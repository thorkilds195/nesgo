package cpu

const PRG_ROM_PG_SIZE = 16384
const CHR_ROM_PG_SIZE = 8192

var NESTAG string = string([]byte{0x4E, 0x45, 0x53, 0x1A})

type Mirroring uint8

const (
	VERTICAL Mirroring = iota
	HORIZONTAL
	FOUR_SCREEN
)

type Rom struct {
	prg_rom          []uint8
	chr_rom          []uint8
	mapper           uint8
	screen_mirroring Mirroring
}

func InitRom(data []uint8) *Rom {
	if string(data[:4]) != NESTAG {
		panic("Wrong start of header")
	}
	mapper := (data[7] & 0b1111_0000) | (data[6] >> 4)
	ines_vers := ((data[7] & 0b0000_1100) >> 2)
	if ines_vers > 0 {
		panic("Currently only ines version 1.0 is supported")
	}
	rom_size := uint(data[4]) * PRG_ROM_PG_SIZE
	chr_size := uint(data[5]) * CHR_ROM_PG_SIZE
	skip_trainer := (data[6] & 0b100) != 0
	var prg_rom_start uint
	if skip_trainer {
		prg_rom_start = 16 + 512
	} else {
		prg_rom_start = 16
	}
	chr_rom_start := prg_rom_start + rom_size
	for _, v := range data[9:15] {
		if v != 0x00 {
			panic("Wrong reserved header, all should be 0")
		}
	}
	is_vertical := (data[6] & 0x01) > 0
	is_four_screen := (data[6] & 0b0000_1000) > 0
	var mirroring Mirroring
	if is_four_screen {
		mirroring = FOUR_SCREEN
	} else if is_vertical {
		mirroring = VERTICAL
	} else {
		mirroring = HORIZONTAL
	}
	return &Rom{
		prg_rom:          data[prg_rom_start:(prg_rom_start + rom_size)],
		chr_rom:          data[chr_rom_start:(chr_rom_start + chr_size)],
		mapper:           mapper,
		screen_mirroring: mirroring,
	}
}

func (r *Rom) GetCHRRom() []uint8 {
	return r.chr_rom
}
