package cpu

import "testing"

// LDA
func TestLDAImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_a == 0x05) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_a == 0x00) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) != 0) {
		t.Error(`Zero flag not set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_a == 0b_1100_0000) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) != 0) {
		t.Error(`Negative flag not set`)
	}
}

func TestLDAZeroPageLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA5, 0x10, 0x00}
	c.mem_write(0x10, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAZeroPageXLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xa2, 0x0F, 0xB5, 0x80, 0x00}
	c.mem_write(0x8F, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAZeroPageXLoadDataWhenOverflow(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0xFF and A to 0x80
	// This should fetch from memory location 0x8F due to overflow
	vec := []uint8{0xa2, 0xFF, 0xB5, 0x80, 0x00}
	c.mem_write(0x7F, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xAD, 0x05, 0x80, 0x00}
	c.mem_write(0x8005, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x92, 0xBD, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x92, 0xB9, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDAIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa0, 0x04, 0xB1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	if !(c.register_a == 10) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

// LDX

func TestLDXImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x05, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_x == 0x05) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDXImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x00, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_x == 0x00) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) != 0) {
		t.Error(`Zero flag not set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDXImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_x == 0b_1100_0000) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) != 0) {
		t.Error(`Negative flag not set`)
	}
}

// LDY
func TestLDYImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x05, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_y == 0x05) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDYImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x00, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_y == 0x00) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) != 0) {
		t.Error(`Zero flag not set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestLDYImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_y == 0b_1100_0000) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) != 0) {
		t.Error(`Negative flag not set`)
	}
}

// TAX
func TestTAXLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0xAA, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_x == 0x05) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestTAXLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0xAA, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_x == 0x00) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) != 0) {
		t.Error(`Zero flag not set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestTAXLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b_1100_0000, 0xAA, 0x00}
	c.LoadAndRun(vec)
	if !(c.register_x == 0b_1100_0000) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) != 0) {
		t.Error(`Negative flag not set`)
	}
}

// INX
func TestInxAdd1(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xe8, 0x00})
	if !(c.register_x == 1) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestInxOverflowTo0(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xff, 0xAA, 0xe8, 0x00})
	if !(c.register_x == 0) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) != 0) {
		t.Error(`Zero flag not set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestInxOverflow(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xff, 0xAA, 0xe8, 0xe8, 0x00})
	if !(c.register_x == 1) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) == 0) {
		t.Error(`Negative flag set`)
	}
}

func TestInxWhenBit7Set(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 200, 0xAA, 0xe8, 0x00})
	if !(c.register_x == 201) {
		t.Error(`Register not set to correct value`)
	}
	if !((c.status & 0b0000_0010) == 0) {
		t.Error(`Zero flag set`)
	}
	if !((c.status & 0b1000_0000) != 0) {
		t.Error(`Negative flag not set`)
	}
}

// Combination tests
func TestFiveOpsWorkingTogether(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	if !(c.register_x == 0xc1) {
		t.Error(`Register not set to correct value`)
	}
}
