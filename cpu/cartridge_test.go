package cpu

import "testing"

func setupDataArray(header []uint8) []uint8 {
	// Just sets up some large dataarray of 0s that the rom can pickup what it needs from
	mem := make([]uint8, 200000)
	copy(mem, header)
	return mem
}

func TestCreatesRomWithHorizontalNoMapperAnd1Pg(t *testing.T) {
	test_header := []uint8{0x4e, 0x45, 0x53, 0x1a, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	test_data := setupDataArray(test_header)
	actual := InitRom(test_data)
	if !(actual.screen_mirroring == HORIZONTAL) {
		t.Error("Mirroring not set correctly")
	}
	if !(actual.mapper == 0x00) {
		t.Error("Mapper not set correctly")
	}
	if !(len(actual.prg_rom) == 1*PRG_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
	if !(len(actual.chr_rom) == 1*CHR_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
}
func TestCreatesRomWithVerticalNoMapperAnd1Pg(t *testing.T) {
	test_header := []uint8{0x4e, 0x45, 0x53, 0x1a, 0x01, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	test_data := setupDataArray(test_header)
	actual := InitRom(test_data)
	if !(actual.screen_mirroring == VERTICAL) {
		t.Error("Mirroring not set correctly")
	}
	if !(actual.mapper == 0x00) {
		t.Error("Mapper not set correctly")
	}
	if !(len(actual.prg_rom) == 1*PRG_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
	if !(len(actual.chr_rom) == 1*CHR_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
}
func TestCreatesRomWithFourScreenNoMapperAnd1Pg(t *testing.T) {
	test_header := []uint8{0x4e, 0x45, 0x53, 0x1a, 0x01, 0x01, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	test_data := setupDataArray(test_header)
	actual := InitRom(test_data)
	if !(actual.screen_mirroring == FOUR_SCREEN) {
		t.Error("Mirroring not set correctly")
	}
	if !(actual.mapper == 0x00) {
		t.Error("Mapper not set correctly")
	}
	if !(len(actual.prg_rom) == 1*PRG_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
	if !(len(actual.chr_rom) == 1*CHR_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
}
func TestCreatesRomWithHorizontal255MapperAnd1Pg(t *testing.T) {
	test_header := []uint8{0x4e, 0x45, 0x53, 0x1a, 0x01, 0x01, 0xF0, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	test_data := setupDataArray(test_header)
	actual := InitRom(test_data)
	if !(actual.screen_mirroring == HORIZONTAL) {
		t.Error("Mirroring not set correctly")
	}
	if !(actual.mapper == 0xFF) {
		t.Error("Mapper not set correctly")
	}
	if !(len(actual.prg_rom) == 1*PRG_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
	if !(len(actual.chr_rom) == 1*CHR_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
}
func TestCreatesRomWithHorizontal254MapperAnd1Pg(t *testing.T) {
	test_header := []uint8{0x4e, 0x45, 0x53, 0x1a, 0x01, 0x01, 0xE0, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	test_data := setupDataArray(test_header)
	actual := InitRom(test_data)
	if !(actual.screen_mirroring == HORIZONTAL) {
		t.Error("Mirroring not set correctly")
	}
	if !(actual.mapper == 0xFE) {
		t.Error("Mapper not set correctly")
	}
	if !(len(actual.prg_rom) == 1*PRG_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
	if !(len(actual.chr_rom) == 1*CHR_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
}
func TestCreatesRomWithHorizontalNoMapperAnd2PrgPg3ChrPg(t *testing.T) {
	test_header := []uint8{0x4e, 0x45, 0x53, 0x1a, 0x02, 0x03, 0xE0, 0xF0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	test_data := setupDataArray(test_header)
	actual := InitRom(test_data)
	if !(actual.screen_mirroring == HORIZONTAL) {
		t.Error("Mirroring not set correctly")
	}
	if !(actual.mapper == 0xFE) {
		t.Error("Mapper not set correctly")
	}
	if !(len(actual.prg_rom) == 2*PRG_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
	if !(len(actual.chr_rom) == 3*CHR_ROM_PG_SIZE) {
		t.Error("Length of prg rom not correct")
	}
}
