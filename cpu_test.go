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
		t.Error(`Register not set to correct value`)
	}
}

// LDA
func TestLDAImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x05)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0010)
}

func TestLDAImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0000)
}

func TestLDAZeroPageLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA5, 0x10, 0x00}
	c.mem_write(0x10, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAZeroPageXLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xa2, 0x0F, 0xB5, 0x80, 0x00}
	c.mem_write(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAZeroPageXLoadDataWhenOverflow(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0xFF and A to 0x80
	// This should fetch from memory location 0x8F due to overflow
	vec := []uint8{0xa2, 0xFF, 0xB5, 0x80, 0x00}
	c.mem_write(0x7F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xAD, 0x05, 0x80, 0x00}
	c.mem_write(0x8005, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x92, 0xBD, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x92, 0xB9, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDAIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa0, 0x04, 0xB1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 10)
	assert_status(t, c.status, 0b0000_0000)
}

// LDX

func TestLDXImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x05)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDXImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x00, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x00)
	assert_status(t, c.status, 0b0000_0010)
}

func TestLDXImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0000)
}

func TestLDXZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA6, 0x10, 0x00}
	c.mem_write(0x10, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDXZeroPageY(t *testing.T) {
	c := InitCPU()
	// Sets y register to 0x0F and x to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xa0, 0x0F, 0xB6, 0x80, 0x00}
	c.mem_write(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDXZeroPageYLoadDataWhenOverflow(t *testing.T) {
	c := InitCPU()
	// Sets y register to 0xFF and x to 0x80
	// This should fetch from memory location 0x8F due to overflow
	vec := []uint8{0xA0, 0xFF, 0xB6, 0x80, 0x00}
	c.mem_write(0x7F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDXAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xAE, 0x05, 0x80, 0x00}
	c.mem_write(0x8005, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDXAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x92, 0xBE, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 10)
	assert_status(t, c.status, 0b0000_0000)
}

// LDY
func TestLDYImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x05)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDYImmediateLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x00, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x00)
	assert_status(t, c.status, 0b0000_0010)
}

func TestLDYImmediateLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0b_1100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0000)
}

func TestLDYZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA4, 0x10, 0x00}
	c.mem_write(0x10, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDYZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xa2, 0x0F, 0xB4, 0x80, 0x00}
	c.mem_write(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDYZeroPageXLoadDataWhenOverflow(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0xFF and A to 0x80
	// This should fetch from memory location 0x8F due to overflow
	vec := []uint8{0xa2, 0xFF, 0xB4, 0x80, 0x00}
	c.mem_write(0x7F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDYAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xAC, 0x05, 0x80, 0x00}
	c.mem_write(0x8005, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLDYAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x92, 0xBC, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 10)
	assert_status(t, c.status, 0b0000_0000)
}

// TAX
func TestTAXLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0xAA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x05)
	assert_status(t, c.status, 0b0000_0000)
}

func TestTAXLoadDataWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0xAA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x00)
	assert_status(t, c.status, 0b0000_0010)
}

func TestTAXLoadDataWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b_1100_0000, 0xAA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0b_1100_0000)
	assert_status(t, c.status, 0b1000_0000)
}

// INX
func TestInxAdd1(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xe8, 0x00})
	assert_register(t, c.register_x, 1)
	assert_status(t, c.status, 0b0000_0000)
}

func TestInxOverflowTo0(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xff, 0xAA, 0xe8, 0x00})
	assert_register(t, c.register_x, 0)
	assert_status(t, c.status, 0b0000_0010)
}

func TestInxOverflow(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xff, 0xAA, 0xe8, 0xe8, 0x00})
	assert_register(t, c.register_x, 1)
	assert_status(t, c.status, 0b0000_0000)
}

func TestInxWhenBit7Set(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 200, 0xAA, 0xe8, 0x00})
	assert_register(t, c.register_x, 201)
	assert_status(t, c.status, 0b1000_0000)
}

// ADC
func TestAdcImmediateWithoutCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x05, 0x69, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x07)
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0000)
}

func TestAdcImmediateWithOutgoingCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0xFF, 0x69, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x01)
	assert_status(t, c.status, 0b0000_0001)
}

func TestAdcImmediateWithOverflowFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x70, 0x69, 0x70, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0xE0)
	assert_status(t, c.status, 0b1100_0000)
}

func TestAdcZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x01, 0x65, 0x15, 0x00}
	c.mem_write(0x15, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0000)
}

func TestAdcZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x01
	// Runs adc instr with 0x80
	// This should fetch from memory location 0x8F
	// and add the current value of a register to it (0x01)
	vec := []uint8{0xa9, 0x01, 0xa2, 0x0F, 0x75, 0x80, 0x00}
	c.mem_write(0x8F, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0000)
}

func TestAdcAbsolute(t *testing.T) {
	c := InitCPU()
	// Sets x register a to 0x01
	// Runs adc instr with 0x10 and 0x80
	// This should fetch from memory location 0x8010
	// and add the current value of a register to it (0x01)
	vec := []uint8{0xa9, 0x01, 0x6D, 0x10, 0x80, 0x00}
	c.mem_write(0x8010, 10)
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
	c.mem_write(0x2092, 10)
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
	c.mem_write(0x2092, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0000)
}

func TestAdcIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x01, 0xa2, 0x04, 0x61, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0000)
}

func TestAdcIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x01, 0xa0, 0x04, 0x71, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 10)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 11)
	assert_status(t, c.status, 0b0000_0000)
}

//And
func TestANDImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0001, 0x29, 0b0000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

func TestANDImmediateWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x29, 0b0000_0010, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0010)
}

func TestANDImmediateWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x29, 0b1000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1000_0001)
	assert_status(t, c.status, 0b1000_0000)
}

func TestANDZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x25, 0xF8, 0x00}
	c.mem_write(0xF8, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

func TestANDZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0011, 0xA2, 0x0F, 0x35, 0x80, 0x00}
	c.mem_write(0x8F, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

func TestANDAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x2D, 0x05, 0x90, 0x00}
	c.mem_write(0x9005, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

func TestANDAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xa2, 0x92, 0x3D, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

func TestANDAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xA0, 0x92, 0x39, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

func TestANDIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xa2, 0x04, 0x21, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

func TestANDIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0xa0, 0x04, 0x31, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0b1000_0001)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0000)
}

// Combination tests
func TestFiveOpsWorkingTogether(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	assert_register(t, c.register_x, 0xc1)
	assert_status(t, c.status, 0b1000_0000)
}
