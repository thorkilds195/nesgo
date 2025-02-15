package cpu

import "testing"

// LDA Tests

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


func Test0xa9LdaImmediateLoadDataWhenBit7Set(t *testing.T) {
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

// TAX Tests
