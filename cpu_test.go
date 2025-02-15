package cpu

import "testing"

// LDA
func TestLDAImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0x00}
	c.Interpet(vec)
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
	c.Interpet(vec)
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
	c.Interpet(vec)
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

// TAX
func TestTAXLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0xAA, 0x00}
	c.Interpet(vec)
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
	c.Interpet(vec)
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
	c.Interpet(vec)
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
	c.register_x = 0
	c.Interpet([]uint8{0xe8, 0x00})
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
	c.register_x = 0xff
	c.Interpet([]uint8{0xe8, 0x00})
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
	c.register_x = 0xff
	c.Interpet([]uint8{0xe8, 0xe8, 0x00})
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
	c.register_x = 200
	c.Interpet([]uint8{0xe8, 0x00})
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
	c.Interpet([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	if !(c.register_x == 0xc1) {
		t.Error(`Register not set to correct value`)
	}
}
