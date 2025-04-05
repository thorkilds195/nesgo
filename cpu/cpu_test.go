package cpu

import "testing"

var FLAGNAMES = []string{
	"Carry",
	"Zero",
	"Interrupt Disable",
	"Decimal",
	"No Flag",
	"No Flag",
	"Overflow",
	"Negative",
}

// Helper functions
func assert_status(t *testing.T, actual, expected uint8) {
	if actual == expected {
		// All is good, so return
		return
	}
	// Otherwise find which flags is causing the difference
	diff := actual ^ expected

	idx := 0
	var i uint8
	for i = 0b0000_0001; idx < 8; i <<= 1 {
		if diff&i > 0 {
			t.Errorf(`%s flag not set right`, FLAGNAMES[idx])
		}
		idx++
	}
}

func assert_register(t *testing.T, actual, expected uint8) {
	if !(actual == expected) {
		t.Errorf(`Register not set to correct value, expected %b but got %b`, expected, actual)
	}
}

// LDA
func TestLDAImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x05)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestLDAImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0100)
}

func TestLDAZeroPageLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA5, 0x10, 0x00}
	c.MemWrite(0x10, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAZeroPageXLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xa2, 0x0F, 0xB5, 0x80, 0x00}
	c.MemWrite(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAZeroPageXLoadDataWhenOverflow(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0xFF and A to 0x80
	// This should fetch from memory location 0x8F due to overflow
	vec := []uint8{0xa2, 0xFF, 0xB5, 0x80, 0x00}
	c.MemWrite(0x7F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xAD, 0x05, 0x80, 0x00}
	c.MemWrite(0x8005, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x92, 0xBD, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x92, 0xB9, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDAIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa0, 0x04, 0xB1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0100)
}

// LDX

func TestLDXImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x05)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDXImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x00, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestLDXImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0100)
}

func TestLDXZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA6, 0x10, 0x00}
	c.MemWrite(0x10, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDXZeroPageY(t *testing.T) {
	c := InitCPU()
	// Sets y register to 0x0F and x to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xa0, 0x0F, 0xB6, 0x80, 0x00}
	c.MemWrite(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDXZeroPageYLoadDataWhenOverflow(t *testing.T) {
	c := InitCPU()
	// Sets y register to 0xFF and x to 0x80
	// This should fetch from memory location 0x8F due to overflow
	vec := []uint8{0xA0, 0xFF, 0xB6, 0x80, 0x00}
	c.MemWrite(0x7F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDXAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xAE, 0x05, 0x80, 0x00}
	c.MemWrite(0x8005, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDXAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x92, 0xBE, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0100)
}

// LDY
func TestLDYImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x05)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDYImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x00, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestLDYImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0100)
}

func TestLDYZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA4, 0x10, 0x00}
	c.MemWrite(0x10, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDYZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xa2, 0x0F, 0xB4, 0x80, 0x00}
	c.MemWrite(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDYZeroPageXLoadDataWhenOverflow(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0xFF and A to 0x80
	// This should fetch from memory location 0x8F due to overflow
	vec := []uint8{0xa2, 0xFF, 0xB4, 0x80, 0x00}
	c.MemWrite(0x7F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDYAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xAC, 0x05, 0x80, 0x00}
	c.MemWrite(0x8005, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLDYAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x92, 0xBC, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0100)
}

// TAX
func TestTAXLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0xAA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x05)
	assert_status(t, c.status, 0b0000_0100)
}

func TestTAXLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0xAA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestTAXLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b_1100_0000, 0xAA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0100)
}

// TXA
func TestTXALoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x05, 0x8A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x05)
	assert_status(t, c.status, 0b0000_0100)
}

func TestTXALoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x00, 0x8A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestTXALoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0b_1100_0000, 0x8A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0100)
}

// TAY
func TestTAYLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0xA8, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x05)
	assert_status(t, c.status, 0b0000_0100)
}

func TestTAYLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0xA8, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestTAYLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b_1100_0000, 0xA8, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0100)
}

// TXY
func TestTXYLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa0, 0x05, 0x98, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x05)
	assert_status(t, c.status, 0b0000_0100)
}

func TestTXYLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa0, 0x00, 0x98, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestTXYLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa0, 0b_1100_0000, 0x98, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0100)
}

// INX
func TestInxAdd1(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xe8, 0x00})
	assert_register(t, c.register_x, 1)
	assert_status(t, c.status, 0b0000_0100)
}

func TestInxOverflowTo0(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xff, 0xAA, 0xe8, 0x00})
	assert_register(t, c.register_x, 0)
	assert_status(t, c.status, 0b0000_0110)
}

func TestInxOverflow(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xff, 0xAA, 0xe8, 0xe8, 0x00})
	assert_register(t, c.register_x, 1)
	assert_status(t, c.status, 0b0000_0100)
}

func TestInxWhenBit7Set(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 200, 0xAA, 0xe8, 0x00})
	assert_register(t, c.register_x, 201)
	assert_status(t, c.status, 0b1000_0100)
}

// INY
func TestInyAdd1(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xC8, 0x00})
	assert_register(t, c.register_y, 1)
	assert_status(t, c.status, 0b0000_0100)
}

func TestInyOverflowTo0(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa0, 0xff, 0xAA, 0xC8, 0x00})
	assert_register(t, c.register_y, 0)
	assert_status(t, c.status, 0b0000_0110)
}

func TestInyOverflow(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa0, 0xff, 0xAA, 0xC8, 0xC8, 0x00})
	assert_register(t, c.register_y, 1)
	assert_status(t, c.status, 0b0000_0100)
}

func TestInyWhenBit7Set(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa0, 200, 0xAA, 0xC8, 0x00})
	assert_register(t, c.register_y, 201)
	assert_status(t, c.status, 0b1000_0100)
}

// ADC
func TestAdcImmediateWithoutCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0x69, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x07)
	assert_status(t, c.status, 0b0000_0100)
}

func TestAdcImmediateWithIngoingCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0x69, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	// Manually set the carry flag
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0x08)
	assert_status(t, c.status, 0b0000_0100)
}

func TestAdcImmediateWithOutgoingCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0xFF, 0x69, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x01)
	assert_status(t, c.status, 0b0000_0101)
}

func TestAdcImmediateWithOverflowFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x70, 0x69, 0x70, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0xE0)
	assert_status(t, c.status, 0b1100_0100)
}

func TestAdcZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x01, 0x65, 0x15, 0x00}
	c.MemWrite(0x15, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0100)
}

func TestAdcZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x01
	// Runs adc instr with 0x80
	// This should fetch from memory location 0x8F
	// and add the current value of a register to it (0x01)
	vec := []uint8{0xa9, 0x01, 0xa2, 0x0F, 0x75, 0x80, 0x00}
	c.MemWrite(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0100)
}

func TestAdcAbsolute(t *testing.T) {
	c := InitCPU()
	// Sets x register a to 0x01
	// Runs adc instr with 0x10 and 0x80
	// This should fetch from memory location 0x8010
	// and add the current value of a register to it (0x01)
	vec := []uint8{0xa9, 0x01, 0x6D, 0x10, 0x80, 0x00}
	c.MemWrite(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
}

func TestAdcAbsoluteX(t *testing.T) {
	c := InitCPU()
	// Sets register a to 0x01 and x register to 0x92
	// Runs adc instr with 0x00 and 0x20
	// This should fetch from memory location 0x2092
	// and add the current value of a register to it (0x01)
	vec := []uint8{0xa9, 0x01, 0xa2, 0x92, 0x7D, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
}

func TestAdcAbsoluteY(t *testing.T) {
	c := InitCPU()
	// Sets register a to 0x01 and y register to 0x92
	// Runs adc instr with 0x00 and 0x20
	// This should fetch from memory location 0x2092
	// and add the current value of a register to it (0x01)
	vec := []uint8{0xa9, 0x01, 0xa0, 0x92, 0x79, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0100)
}

func TestAdcIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x01, 0xa2, 0x04, 0x61, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0100)
}

func TestAdcIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x01, 0xa0, 0x04, 0x71, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0100)
}

//And
func TestANDImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0001, 0x29, 0b0000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

func TestANDImmediateWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x29, 0b0000_0010, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestANDImmediateWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x29, 0b1000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1000_0001)
	assert_status(t, c.status, 0b1000_0100)
}

func TestANDZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x25, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

func TestANDZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0011, 0xA2, 0x0F, 0x35, 0x80, 0x00}
	c.MemWrite(0x8F, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

func TestANDAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x2D, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

func TestANDAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xa2, 0x92, 0x3D, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

func TestANDAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xA0, 0x92, 0x39, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

func TestANDIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xa2, 0x04, 0x21, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

func TestANDIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xa0, 0x04, 0x31, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0100)
}

//ASL
func TestASLAccumulator(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x0A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestASLZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x06, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestASLZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x16, 0xF8, 0x00}
	c.MemWrite(0xFA, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xFA), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestASLAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x0E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9005), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestASLAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x1E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9007, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9007), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestASLAccumulatorSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b1000_0011, 0x0A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0101)
}

func TestASLAccumulatorClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b1000_0011, 0x0A, 0x0A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_1100)
	assert_status(t, c.status, 0b0000_0100)
}

func TestASLZeroPageSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x06, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b1000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0101)
}

func TestASLZeroPageClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x06, 0xF8, 0x06, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b1000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_1100)
	assert_status(t, c.status, 0b0000_0100)
}

//BCC
func TestBCCWithCarryFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x90, 0x02, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0101)
}

func TestBCCWithoutCarryFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x90, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

//BCS
func TestBCSWitouthCarryFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB0, 0x02, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBSCWithCarryFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB0, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0101)
}

//BEQ
func TestBEQWithoutZeroFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xF0, 0x02, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBEQWithZeroFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xF0, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0010
	c.Run()
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBEQWithCarryFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB0, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0101)
}

//BIT
func TestBITZeroPageAllStatusZero(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0010, 0x24, 0x10, 0x00}
	c.MemWrite(0x10, 0b0000_0010)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBITZeroPageZeroFlagSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0010, 0x24, 0x10, 0x00}
	c.MemWrite(0x10, 0b0000_0000)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0110)
}

func TestBITZeroPageOverflowSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x24, 0x10, 0x00}
	c.MemWrite(0x10, 0b0100_0000)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0100_0110)
}

func TestBITZeroPageNegativeSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x24, 0x10, 0x00}
	c.MemWrite(0x10, 0b1000_0000)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0110)
}

func TestBITAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0010, 0x2C, 0x10, 0x80, 0x00}
	c.MemWrite(0x8010, 0b0000_0010)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0100)
}

//BMI
func TestBMIWithNegativeFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x30, 0x02, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b1000_0000
	c.Run()
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBMIWithoutNegativeFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x30, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

//BNE
func TestBNEWithoutZeroFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xD0, 0x02, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBNEWithZeroFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0xD0, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

//BPL
func TestBPLWithoutNegativeFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x10, 0x02, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBPLWithNegativeFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x10, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b1000_0000
	c.Run()
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

//BVC
func TestBVCWithOverflowFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x50, 0x02, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0100_0000
	c.Run()
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0100_0100)
}

func TestBVCWithoutOverflowFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x50, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

//BVS
func TestBVSWithoutOverflowFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x70, 0x02, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0100)
}

func TestBVSWithOverflowFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x70, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0100_0000
	c.Run()
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0100_0100)
}

//CLC
func TestCLCWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x18, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

func TestCLCWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x18, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

//CLD
func TestCLDWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xD8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_1000
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

func TestCLDWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xD8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

//CLI
func TestCLIWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x58, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0100
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

func TestCLIWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x58, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

//CLV
func TestCLVWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0100_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

func TestCLVWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

//CMP
func TestCMPImmediateWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC9, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPImmediateWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC9, 0x09, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPImmediateWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0xFF, 0xC9, 0xFF, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCMPZeroPageWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC5, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPZeroPageWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC5, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPZeroPageWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0xFF, 0xC5, 0xF8, 0x00}
	c.MemWrite(0xF8, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCMPZeroPageXWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xD5, 0xF8, 0x00}
	c.MemWrite(0xF9, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPZeroPageXWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xD5, 0xF8, 0x00}
	c.MemWrite(0xF9, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPZeroPageXWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0xFF, 0xD5, 0xF8, 0x00}
	c.MemWrite(0xF9, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCMPAbsoluteWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xCD, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPAbsoluteWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xCD, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPAbsoluteWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0xFF, 0xCD, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCMPAbsoluteXWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xDD, 0x50, 0x80, 0x00}
	c.MemWrite(0x8051, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPAbsoluteXWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xDD, 0x50, 0x80, 0x00}
	c.MemWrite(0x8051, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPAbsoluteXWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0xFF, 0xDD, 0x50, 0x80, 0x00}
	c.MemWrite(0x8051, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCMPAbsoluteYWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x01, 0xA9, 0x09, 0xD9, 0x50, 0x80, 0x00}
	c.MemWrite(0x8051, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPAbsoluteYWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x01, 0xA9, 0x09, 0xD9, 0x50, 0x80, 0x00}
	c.MemWrite(0x8051, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPAbsoluteYWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x01, 0xA9, 0xFF, 0xD9, 0x50, 0x80, 0x00}
	c.MemWrite(0x8051, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCMPIndirectXWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA9, 0x09, 0xC1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPIndirectXWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA9, 0x09, 0xC1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPIndirectXWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA9, 0xFF, 0xC1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCMPIndirectYWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x04, 0xA9, 0x09, 0xD1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCMPIndirectYWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x04, 0xA9, 0x09, 0xD1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCMPIndirectYWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x04, 0xA9, 0xFF, 0xD1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

//CPX
func TestCPXImmediateWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE0, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCPXImmediateWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE0, 0x09, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCPXImmediateWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0xFF, 0xE0, 0xFF, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCPXZeroPageWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE4, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCPXZeroPageWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE4, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCPXZeroPageWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0xFF, 0xE4, 0xF8, 0x00}
	c.MemWrite(0xF8, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCPXAbsoluteWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xEC, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCPXAbsoluteWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xEC, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCPXAbsoluteWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0xFF, 0xEC, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

//CPY
func TestCPYImmediateWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC0, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCPYImmediateWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC0, 0x09, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCPYImmediateWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0xFF, 0xC0, 0xFF, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCPYZeroPageWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC4, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCPYZeroPageWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC4, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCPYZeroPageWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0xFF, 0xC4, 0xF8, 0x00}
	c.MemWrite(0xF8, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

func TestCPYAbsoluteWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xCC, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0101)
}

func TestCPYAbsoluteWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xCC, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0111)
}

func TestCPYAbsoluteWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0xFF, 0xCC, 0x50, 0x80, 0x00}
	c.MemWrite(0x8050, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0111)
}

//DEC
func TestDECZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xC6, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0x01)
	assert_status(t, c.status, 0b0000_0100)
}

func TestDECZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xD6, 0xF8, 0x00}
	c.MemWrite(0xF9, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF9), 0x01)
	assert_status(t, c.status, 0b0000_0100)
}

func TestDECZeroAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xCE, 0x05, 0x80, 0x00}
	c.MemWrite(0x8005, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8005), 0x01)
	assert_status(t, c.status, 0b0000_0100)
}

func TestDECZeroAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xDE, 0x05, 0x80, 0x00}
	c.MemWrite(0x8006, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8006), 0x01)
	assert_status(t, c.status, 0b0000_0100)
}

//DEX
func TestDEX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0xCA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x01)
	assert_status(t, c.status, 0b0000_0100)
}

//DEY
func TestDEY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x02, 0x88, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x01)
	assert_status(t, c.status, 0b0000_0100)
}

//EOR
func TestEORImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0101, 0x49, 0b0000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestEORImmediateWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0001, 0x49, 0b0000_0001, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestEORImmediateWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x49, 0b0100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1000_0001)
	assert_status(t, c.status, 0b1000_0100)
}

func TestEORZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x45, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestEORZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0101, 0xA2, 0x0F, 0x55, 0x80, 0x00}
	c.MemWrite(0x8F, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestEORAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x4D, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestEORAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x92, 0x5D, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestEORAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xA0, 0x92, 0x59, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestEORIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x04, 0x41, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

func TestEORIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa0, 0x04, 0x51, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0100)
}

//INC
func TestINCZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xE6, 0xF8, 0x00}
	c.MemWrite(0xF8, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0x03)
	assert_status(t, c.status, 0b0000_0100)
}

func TestINCZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xF6, 0xF8, 0x00}
	c.MemWrite(0xF9, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF9), 0x03)
	assert_status(t, c.status, 0b0000_0100)
}

func TestINCAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xEE, 0x05, 0x80, 0x00}
	c.MemWrite(0x8005, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8005), 0x03)
	assert_status(t, c.status, 0b0000_0100)
}

func TestINCAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xFE, 0x05, 0x80, 0x00}
	c.MemWrite(0x8006, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8006), 0x03)
	assert_status(t, c.status, 0b0000_0100)
}

// JMP
func TestJMPAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x4C, 0x01, 0xFF, 0x00}
	c.MemWrite(0xFF01, 0xA2)
	c.MemWrite(0xFF02, 0x09)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x09)
	assert_status(t, c.status, 0b0000_0100)
}

func TestJMPIndirect(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x6C, 0x01, 0xFF, 0x00}
	c.mem_write_16(0xFF01, 0xFF10)
	c.MemWrite(0xFF10, 0xA2)
	c.MemWrite(0xFF11, 0x09)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x09)
	assert_status(t, c.status, 0b0000_0100)
}

// LSR
func TestLSRAccumulator(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0110, 0x4A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLSRZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x46, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLSRZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x56, 0xF8, 0x00}
	c.MemWrite(0xFA, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xFA), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLSRAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x4E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9005), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLSRAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x5E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9007, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9007), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLSRAccumulatorSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x4A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0101)
}

func TestLSRAccumulatorClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_1100, 0x4A, 0x4A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestLSRZeroPageSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x46, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0001)
	assert_status(t, c.status, 0b0000_0101)
}

func TestLSRZeroPageClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x46, 0xF8, 0x46, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_1100)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

// NOP
func TestNOP(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x05, 0xEA, 0xA9, 0x08, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x08)
	assert_status(t, c.status, 0b0000_0100)
}

//ORA
func TestORAImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0101, 0x09, 0b0000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestORAImmediateWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0000, 0x09, 0b0000_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0110)
}

func TestORAImmediateWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x09, 0b0100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1100_0001)
	assert_status(t, c.status, 0b1000_0100)
}

func TestORAZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x05, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestORAZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0101, 0xA2, 0x0F, 0x15, 0x80, 0x00}
	c.MemWrite(0x8F, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestORAAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x0D, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestORAAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x92, 0x1D, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestORAAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xA0, 0x92, 0x19, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestORAIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x04, 0x01, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestORAIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa0, 0x04, 0x11, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

// ROL
func TestROLAccumulator(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0110, 0x2A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_1100)
	assert_status(t, c.status, 0b0000_0100)
}

func TestROLZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x26, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_1100)
	assert_status(t, c.status, 0b0000_0100)
}

func TestROLZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x36, 0xF8, 0x00}
	c.MemWrite(0xFA, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xFA), 0b0000_1100)
	assert_status(t, c.status, 0b0000_0100)
}

func TestROLAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x2E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9005), 0b0000_1100)
	assert_status(t, c.status, 0b0000_0100)
}

func TestROLAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x3E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9007, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9007), 0b0000_1100)
	assert_status(t, c.status, 0b0000_0100)
}

func TestROLAccumulatorWhenCarrySet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0110, 0x2A, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_1101)
	assert_status(t, c.status, 0b0000_0100)
}

func TestROLZeroPageWhenCarrySet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x26, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0110)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.MemRead(0xF8), 0b0000_1101)
	assert_status(t, c.status, 0b0000_0100)
}

// ROR
func TestRORAccumulator(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0110, 0x6A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestRORZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x66, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestRORZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x76, 0xF8, 0x00}
	c.MemWrite(0xFA, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xFA), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestRORAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x6E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9005), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestRORAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x7E, 0x05, 0x90, 0x00}
	c.MemWrite(0x9007, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9007), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0100)
}

func TestRORAccumulatorWhenCarrySet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_1100, 0x6A, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestRORZeroPageWhenCarrySet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x66, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_1100)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.MemRead(0xF8), 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

//SBC
func TestSBCImmediateWithoutOverflowAndCarrySet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xE9, 0b0000_0001, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCImmediateWithoutOverflowAndCarryNotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xE9, 0b0000_0001, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_a, 0b0000_0011)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCImmediateWithOverflowAndCarrySet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0001, 0xE9, 0b0000_0010, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b1111_1111)
	assert_status(t, c.status, 0b1100_0100)
}

func TestSBCImmediateWithOverflowAndCarryNotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0001, 0xE9, 0b0000_0010, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_a, 0b1111_1110)
	assert_status(t, c.status, 0b1100_0100)
}

func TestSBCImmediateWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0001, 0xE9, 0b0000_0001, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0000)
	assert_status(t, c.status, 0b0000_0111)
}

func TestSBCZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xE5, 0xF8, 0x00}
	c.MemWrite(0xF8, 0b0000_0001)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0101, 0xA2, 0x0F, 0xF5, 0x80, 0x00}
	c.MemWrite(0x8F, 0b0000_0001)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xED, 0x05, 0x90, 0x00}
	c.MemWrite(0x9005, 0b0000_0001)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x92, 0xFD, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b0000_0001)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xA0, 0x92, 0xF9, 0x00, 0x20, 0x00}
	c.MemWrite(0x2092, 0b0000_0001)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x04, 0xE1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0001)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSBCIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa0, 0x04, 0xF1, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0001)
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_0100)
	assert_status(t, c.status, 0b0000_0101)
}

//SEC
func TestSECWhenNotSetInAdvance(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x38, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0101)
}

func TestSECWhenSetInAdvance(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x38, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_register(t, c.register_a, 0b0000_000)
	assert_status(t, c.status, 0b0000_0101)
}

func TestSECWithOtherInstrs(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x38, 0xA9, 0b0000_0111, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0101)
}

//SED
func TestSEDWhenNotSetInAdvance(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xF8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_1100)
}

func TestSEDWhenSetInAdvance(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xF8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_1000
	c.Run()
	assert_register(t, c.register_a, 0b0000_000)
	assert_status(t, c.status, 0b0000_1100)
}

func TestSEDWithOtherInstrs(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xF8, 0xA9, 0b0000_0111, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_1100)
}

//SEI
func TestSEIWhenNotSetInAdvance(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x78, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0100)
}

func TestSEIWhenSetInAdvance(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x78, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0100
	c.Run()
	assert_register(t, c.register_a, 0b0000_0000)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSEIWithOtherInstrs(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x78, 0xA9, 0b0000_0111, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

// STA
func TestSTAZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0111, 0x85, 0xF8, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), c.register_a)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTAZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0111, 0xA2, 0x0F, 0x95, 0x80, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8F), c.register_a)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTAAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x8D, 0x05, 0x90, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9005), c.register_a)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTAAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x02, 0x9D, 0x00, 0x20, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x2002), c.register_a)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTAAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xA0, 0x02, 0x99, 0x00, 0x20, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x2002), c.register_a)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTAIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x04, 0x81, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8010), c.register_a)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTAIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa0, 0x04, 0x91, 0x20, 0x00}
	c.MemWrite(0x24, 0x10)
	c.MemWrite(0x25, 0x80)
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8010), c.register_a)
	assert_status(t, c.status, 0b0000_0100)
}

// STX
func TestSTXZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0b0000_0111, 0x86, 0xF8, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTXZeroPageY(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA2, 0b0000_0111, 0xA0, 0x0F, 0x96, 0x80, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8F), 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTXAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0b0000_0101, 0x8E, 0x05, 0x90, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9005), 0b0000_0101)
	assert_status(t, c.status, 0b0000_0100)
}

// STY
func TestSTYZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0b0000_0111, 0x84, 0xF8, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0xF8), 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTYZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA0, 0b0000_0111, 0xA2, 0x0F, 0x94, 0x80, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x8F), 0b0000_0111)
	assert_status(t, c.status, 0b0000_0100)
}

func TestSTYAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0b0000_0101, 0x8C, 0x05, 0x90, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.MemRead(0x9005), c.register_y)
	assert_status(t, c.status, 0b0000_0100)
}

// JSR
func TestJSRStackDecrement(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x20, 0x01, 0x82, 0x00}
	c.LoadAndRun(vec)
	// Assumes stack pointer will be referenced to 0xFD first by JSR command
	// since that decrements by 2 and then 3 more decrements by BRK cmd
	assert_register(t, c.stack_pointer, 0xFA)
	if !(c.mem_read_16(0x01FE) == 0x8002) {
		t.Error("Stack pointer return value is wrong")
	}
}

func TestJSRLDA(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x20, 0x01, 0x82, 0x00}
	c.MemWrite(0x8201, 0xA9)
	c.MemWrite(0x8202, 0x09)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x09)
}

// PHA
func TestPHA(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0xFC, 0x48, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0xFC)
	if !(c.mem_read_16(0x01FF) == 0xFC) {
		t.Error("Stack pointer return value is wrong")
	}
}

// PHP
func TestPHP(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x00, 0x08, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0110)
	if !(c.mem_read_16(0x01FF) == 0b0000_0010) {
		t.Error("Stack pointer return value is wrong")
	}
}

// PLA
func TestPLA(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x68, 0x00}
	c.Load(vec)
	c.Reset()
	c.push(0x12)
	c.Run()
	assert_register(t, c.register_a, 0x12)
	assert_status(t, c.status, 0b0000_0100)
}

// PLP
func TestPLP(t *testing.T) {
	flag := uint8(0b0000_0110)
	c := InitCPU()
	vec := []uint8{0x28, 0x00}
	c.Load(vec)
	c.Reset()
	c.push(flag)
	c.Run()
	assert_status(t, c.status, flag)
}

// RTS
func TestRTS(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x20, 0x01, 0x82, 0xA2, 0x09, 0x00}
	c.MemWrite(0x8201, 0xA9)
	c.MemWrite(0x8202, 0x12)
	c.MemWrite(0x8203, 0x60)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x12)
	assert_register(t, c.register_x, 0x09)
}

// RTI
func TestRTI(t *testing.T) {
	flag := uint8(0b0000_0100)
	c := InitCPU()
	vec := []uint8{0x40, 0x02}
	c.Load(vec)
	c.Reset()
	c.MemWrite(0x8005, 0x02)
	c.push_16(0x8005)
	c.push(flag)
	c.Run()
	assert_status(t, c.status, flag)
	// Makes sure it resumes at the halt instruction which is written to
	// location 0x8005 and it will consume that one instr before halting and end up
	// at 0x8006 as the final counter of the program
	if !(c.program_counter == 0x8006) {
		t.Error("Program counter set to wrong value after interrupt")
	}
}

// BRK
func TestBRK(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x00}
	c.mem_write_16(0xFFFE, 0x8002)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0100)
	if !(c.program_counter == 0x8002) {
		t.Error("The program counter is set to the wrong value")
	}
}

// Combination tests
func TestFiveOpsWorkingTogether(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	assert_register(t, c.register_x, 0xc1)
	assert_status(t, c.status, 0b1000_0100)
}
