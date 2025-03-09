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

// INY
func TestInyAdd1(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xC8, 0x00})
	assert_register(t, c.register_y, 1)
	assert_status(t, c.status, 0b0000_0000)
}

func TestInyOverflowTo0(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa0, 0xff, 0xAA, 0xC8, 0x00})
	assert_register(t, c.register_y, 0)
	assert_status(t, c.status, 0b0000_0010)
}

func TestInyOverflow(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa0, 0xff, 0xAA, 0xC8, 0xC8, 0x00})
	assert_register(t, c.register_y, 1)
	assert_status(t, c.status, 0b0000_0000)
}

func TestInyWhenBit7Set(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa0, 200, 0xAA, 0xC8, 0x00})
	assert_register(t, c.register_y, 201)
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

//ASL
func TestASLAccumulator(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x0A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestASLZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x06, 0xF8, 0x00}
	c.mem_write(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestASLZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x16, 0xF8, 0x00}
	c.mem_write(0xFA, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xFA), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestASLAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x0E, 0x05, 0x90, 0x00}
	c.mem_write(0x9005, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x9005), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestASLAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x1E, 0x05, 0x90, 0x00}
	c.mem_write(0x9007, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x9007), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestASLAccumulatorSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b1000_0011, 0x0A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0001)
}

func TestASLAccumulatorClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b1000_0011, 0x0A, 0x0A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_1100)
	assert_status(t, c.status, 0b0000_0000)
}

func TestASLZeroPageSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x06, 0xF8, 0x00}
	c.mem_write(0xF8, 0b1000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0b0000_0110)
	assert_status(t, c.status, 0b0000_0001)
}

func TestASLZeroPageClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x06, 0xF8, 0x06, 0xF8, 0x00}
	c.mem_write(0xF8, 0b1000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0b0000_1100)
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0001)
}

func TestBCCWithoutCarryFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x90, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0000)
}

//BCS
func TestBCSWitouthCarryFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB0, 0x02, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0001)
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
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0001)
}

//BIT
func TestBITZeroPageAllStatusZero(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0010, 0x24, 0x10, 0x00}
	c.mem_write(0x10, 0b0000_0010)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0000)
}

func TestBITZeroPageZeroFlagSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0010, 0x24, 0x10, 0x00}
	c.mem_write(0x10, 0b0000_0000)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0010)
}

func TestBITZeroPageOverflowSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x24, 0x10, 0x00}
	c.mem_write(0x10, 0b0100_0000)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0100_0010)
}

func TestBITZeroPageNegativeSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x24, 0x10, 0x00}
	c.mem_write(0x10, 0b1000_0000)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0010)
}

func TestBITAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0010, 0x2C, 0x10, 0x80, 0x00}
	c.mem_write(0x8010, 0b0000_0010)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0000)
}

func TestBMIWithoutNegativeFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x30, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0000)
}

func TestBNEWithZeroFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0x00, 0xD0, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0000)
}

//BPL
func TestBPLWithoutNegativeFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x10, 0x02, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0100_0000)
}

func TestBVCWithoutOverflowFlag(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x50, 0x02, 0xa9, 0x05, 0xA2, 0x02, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_register(t, c.register_x, 0x02)
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0000_0000)
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
	assert_status(t, c.status, 0b0100_0000)
}

//CLC
func TestCLCWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x18, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0001
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

func TestCLCWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x18, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

//CLD
func TestCLDWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xD8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_1000
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

func TestCLDWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xD8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

//CLI
func TestCLIWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x58, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0100
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

func TestCLIWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x58, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

//CLV
func TestCLVWhenSetTo1(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0100_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

func TestCLVWhenSetTo0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xB8, 0x00}
	c.Load(vec)
	c.Reset()
	c.status = 0b0000_0000
	c.Run()
	assert_status(t, c.status, 0b0000_0000)
}

//CMP
func TestCMPImmediateWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC9, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPImmediateWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC9, 0x09, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPImmediateWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0xFF, 0xC9, 0xFF, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCMPZeroPageWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC5, 0xF8, 0x00}
	c.mem_write(0xF8, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPZeroPageWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xC5, 0xF8, 0x00}
	c.mem_write(0xF8, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPZeroPageWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0xFF, 0xC5, 0xF8, 0x00}
	c.mem_write(0xF8, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCMPZeroPageXWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xD5, 0xF8, 0x00}
	c.mem_write(0xF9, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPZeroPageXWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xD5, 0xF8, 0x00}
	c.mem_write(0xF9, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPZeroPageXWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0xFF, 0xD5, 0xF8, 0x00}
	c.mem_write(0xF9, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCMPAbsoluteWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xCD, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPAbsoluteWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x09, 0xCD, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPAbsoluteWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0xFF, 0xCD, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCMPAbsoluteXWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xDD, 0x50, 0x80, 0x00}
	c.mem_write(0x8051, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPAbsoluteXWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0x09, 0xDD, 0x50, 0x80, 0x00}
	c.mem_write(0x8051, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPAbsoluteXWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x01, 0xA9, 0xFF, 0xDD, 0x50, 0x80, 0x00}
	c.mem_write(0x8051, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCMPAbsoluteYWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x01, 0xA9, 0x09, 0xD9, 0x50, 0x80, 0x00}
	c.mem_write(0x8051, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPAbsoluteYWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x01, 0xA9, 0x09, 0xD9, 0x50, 0x80, 0x00}
	c.mem_write(0x8051, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPAbsoluteYWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x01, 0xA9, 0xFF, 0xD9, 0x50, 0x80, 0x00}
	c.mem_write(0x8051, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCMPIndirectXWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA9, 0x09, 0xC1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPIndirectXWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA9, 0x09, 0xC1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPIndirectXWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa2, 0x04, 0xA9, 0xFF, 0xC1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCMPIndirectYWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x04, 0xA9, 0x09, 0xD1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCMPIndirectYWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x04, 0xA9, 0x09, 0xD1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCMPIndirectYWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x04, 0xA9, 0xFF, 0xD1, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

//CPX
func TestCPXImmediateWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE0, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCPXImmediateWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE0, 0x09, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCPXImmediateWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0xFF, 0xE0, 0xFF, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCPXZeroPageWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE4, 0xF8, 0x00}
	c.mem_write(0xF8, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCPXZeroPageWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xE4, 0xF8, 0x00}
	c.mem_write(0xF8, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCPXZeroPageWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0xFF, 0xE4, 0xF8, 0x00}
	c.mem_write(0xF8, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCPXAbsoluteWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xEC, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCPXAbsoluteWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x09, 0xEC, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCPXAbsoluteWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0xFF, 0xEC, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

//CPY
func TestCPYImmediateWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC0, 0x05, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCPYImmediateWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC0, 0x09, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCPYImmediateWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0xFF, 0xC0, 0xFF, 0x00}
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCPYZeroPageWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC4, 0xF8, 0x00}
	c.mem_write(0xF8, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCPYZeroPageWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xC4, 0xF8, 0x00}
	c.mem_write(0xF8, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCPYZeroPageWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0xFF, 0xC4, 0xF8, 0x00}
	c.mem_write(0xF8, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

func TestCPYAbsoluteWhenAGreaterThanM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xCC, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0x05)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0001)
}

func TestCPYAbsoluteWhenAEqualM(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x09, 0xCC, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0x09)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b0000_0011)
}

func TestCPYAbsoluteWhen7BitSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0xFF, 0xCC, 0x50, 0x80, 0x00}
	c.mem_write(0x8050, 0xFF)
	c.LoadAndRun(vec)
	assert_status(t, c.status, 0b1000_0011)
}

//DEC
func TestDECZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xC6, 0xF8, 0x00}
	c.mem_write(0xF8, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0x01)
	assert_status(t, c.status, 0b0000_0000)
}

func TestDECZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xD6, 0xF8, 0x00}
	c.mem_write(0xF9, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF9), 0x01)
	assert_status(t, c.status, 0b0000_0000)
}

func TestDECZeroAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xCE, 0x05, 0x80, 0x00}
	c.mem_write(0x8005, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x8005), 0x01)
	assert_status(t, c.status, 0b0000_0000)
}

func TestDECZeroAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xDE, 0x05, 0x80, 0x00}
	c.mem_write(0x8006, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x8006), 0x01)
	assert_status(t, c.status, 0b0000_0000)
}

//DEX
func TestDEX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0xCA, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x01)
	assert_status(t, c.status, 0b0000_0000)
}

//DEY
func TestDEY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA0, 0x02, 0x88, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_y, 0x01)
	assert_status(t, c.status, 0b0000_0000)
}

//EOR
func TestEORImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0101, 0x49, 0b0000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestEORImmediateWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0001, 0x49, 0b0000_0001, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0010)
}

func TestEORImmediateWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x49, 0b0100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1000_0001)
	assert_status(t, c.status, 0b1000_0000)
}

func TestEORZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x45, 0xF8, 0x00}
	c.mem_write(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestEORZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0101, 0xA2, 0x0F, 0x55, 0x80, 0x00}
	c.mem_write(0x8F, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestEORAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x4D, 0x05, 0x90, 0x00}
	c.mem_write(0x9005, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestEORAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x92, 0x5D, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestEORAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xA0, 0x92, 0x59, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestEORIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x04, 0x41, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

func TestEORIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa0, 0x04, 0x51, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0110)
	assert_status(t, c.status, 0b0000_0000)
}

//INC
func TestINCZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xE6, 0xF8, 0x00}
	c.mem_write(0xF8, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0x03)
	assert_status(t, c.status, 0b0000_0000)
}

func TestINCZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xF6, 0xF8, 0x00}
	c.mem_write(0xF9, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF9), 0x03)
	assert_status(t, c.status, 0b0000_0000)
}

func TestINCAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xEE, 0x05, 0x80, 0x00}
	c.mem_write(0x8005, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x8005), 0x03)
	assert_status(t, c.status, 0b0000_0000)
}

func TestINCAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x01, 0xFE, 0x05, 0x80, 0x00}
	c.mem_write(0x8006, 0x02)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x8006), 0x03)
	assert_status(t, c.status, 0b0000_0000)
}

// JMP
func TestJMPAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x4C, 0x01, 0xFF, 0x00}
	c.mem_write(0xFF01, 0xA2)
	c.mem_write(0xFF02, 0x09)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x09)
	assert_status(t, c.status, 0b0000_0000)
}

func TestJMPIndirect(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x6C, 0x01, 0xFF, 0x00}
	c.mem_write_16(0xFF01, 0xFF10)
	c.mem_write(0xFF10, 0xA2)
	c.mem_write(0xFF11, 0x09)
	c.LoadAndRun(vec)
	assert_register(t, c.register_x, 0x09)
	assert_status(t, c.status, 0b0000_0000)
}

// LSR
func TestLSRAccumulator(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0110, 0x4A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0011)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLSRZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x46, 0xF8, 0x00}
	c.mem_write(0xF8, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLSRZeroPageX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x56, 0xF8, 0x00}
	c.mem_write(0xFA, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xFA), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLSRAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x4E, 0x05, 0x90, 0x00}
	c.mem_write(0x9005, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x9005), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLSRAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA2, 0x02, 0x5E, 0x05, 0x90, 0x00}
	c.mem_write(0x9007, 0b0000_0110)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0x9007), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLSRAccumulatorSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0011, 0x4A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0001)
	assert_status(t, c.status, 0b0000_0001)
}

func TestLSRAccumulatorClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_1100, 0x4A, 0x4A, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0011)
	assert_status(t, c.status, 0b0000_0000)
}

func TestLSRZeroPageSetsCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x46, 0xF8, 0x00}
	c.mem_write(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0b0000_0001)
	assert_status(t, c.status, 0b0000_0001)
}

func TestLSRZeroPageClearCarry(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0x46, 0xF8, 0x46, 0xF8, 0x00}
	c.mem_write(0xF8, 0b0000_1100)
	c.LoadAndRun(vec)
	assert_register(t, c.mem_read(0xF8), 0b0000_0011)
	assert_status(t, c.status, 0b0000_0000)
}

// NOP
func TestNOP(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0x05, 0xEA, 0xA9, 0x08, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x08)
	assert_status(t, c.status, 0b0000_0000)
}

//ORA
func TestORAImmediateLoadDataWhenBit7NotSet(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0101, 0x09, 0b0000_0011, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

func TestORAImmediateWhen0(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b0000_0000, 0x09, 0b0000_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0x00)
	assert_status(t, c.status, 0b0000_0010)
}

func TestORAImmediateWhenBit7Set(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xa9, 0b1100_0001, 0x09, 0b0100_0000, 0x00}
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b_1100_0001)
	assert_status(t, c.status, 0b1000_0000)
}

func TestORAZeroPage(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x05, 0xF8, 0x00}
	c.mem_write(0xF8, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

func TestORAZeroPageX(t *testing.T) {
	c := InitCPU()
	// Sets x register to 0x0F and A to 0x80
	// This should fetch from memory location 0x8F
	vec := []uint8{0xA9, 0b0000_0101, 0xA2, 0x0F, 0x15, 0x80, 0x00}
	c.mem_write(0x8F, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

func TestORAAbsolute(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0x0D, 0x05, 0x90, 0x00}
	c.mem_write(0x9005, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

func TestORAAbsoluteX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x92, 0x1D, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

func TestORAAbsoluteY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xA0, 0x92, 0x19, 0x00, 0x20, 0x00}
	c.mem_write(0x2092, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

func TestORAIndirectX(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa2, 0x04, 0x01, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

func TestORAIndirectY(t *testing.T) {
	c := InitCPU()
	vec := []uint8{0xA9, 0b0000_0101, 0xa0, 0x04, 0x11, 0x20, 0x00}
	c.mem_write(0x24, 0x10)
	c.mem_write(0x25, 0x80)
	c.mem_write_16(0x8010, 0b0000_0011)
	c.LoadAndRun(vec)
	assert_register(t, c.register_a, 0b0000_0111)
	assert_status(t, c.status, 0b0000_0000)
}

// Combination tests
func TestFiveOpsWorkingTogether(t *testing.T) {
	c := InitCPU()
	c.LoadAndRun([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	assert_register(t, c.register_x, 0xc1)
	assert_status(t, c.status, 0b1000_0000)
}
