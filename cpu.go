package cpu

type AddressingMode uint8

const (
	IMMEDIATE AddressingMode = iota
	ZEROPAGE
	ZEROPAGEX
	ZEROPAGEY
	RELATIVE
	ABSOLUTE
	ABSOLUTEX
	ABSOLUTEY
	INDIRECTX
	INDIRECTY
	NONE
)

type OpCode struct {
	code   uint8
	mode   AddressingMode
	bytes  uint8
	cycles uint8
	f_call func(*CPU, OpCode)
}

var OPTABLE = map[uint8]OpCode{
	0xA9: {0xA9, IMMEDIATE, 2, 2, (*CPU).lda},
	0xA5: {0xA5, ZEROPAGE, 2, 2, (*CPU).lda},
	0xB5: {0xB5, ZEROPAGEX, 2, 2, (*CPU).lda},
	0xAD: {0xAD, ABSOLUTE, 2, 2, (*CPU).lda},
	0xBD: {0xBD, ABSOLUTEX, 2, 2, (*CPU).lda},
	0xB9: {0xB9, ABSOLUTEY, 2, 2, (*CPU).lda},
	0xA1: {0xA1, INDIRECTX, 2, 2, (*CPU).lda},
	0xB1: {0xB1, INDIRECTY, 2, 2, (*CPU).lda},
	0xA2: {0xA2, IMMEDIATE, 2, 2, (*CPU).ldx},
	0xA6: {0xA6, ZEROPAGE, 2, 2, (*CPU).ldx},
	0xB6: {0xB6, ZEROPAGEY, 2, 2, (*CPU).ldx},
	0xAE: {0xAE, ABSOLUTE, 2, 2, (*CPU).ldx},
	0xBE: {0xBE, ABSOLUTEY, 2, 2, (*CPU).ldx},
	0xA0: {0xBE, IMMEDIATE, 2, 2, (*CPU).ldy},
	0xAA: {0xAA, NONE, 2, 2, (*CPU).tax},
	0xE8: {0xE8, NONE, 2, 2, (*CPU).inx},
}

type CPU struct {
	register_a      uint8
	register_x      uint8
	register_y      uint8
	status          uint8
	program_counter uint16
	memory          [0xFFFF]uint8
}

func InitCPU() *CPU {
	return &CPU{}
}

func (c *CPU) LoadAndRun(program []uint8) {
	c.load(program)
	c.reset()
	c.run()
}

func (c *CPU) reset() {
	c.register_a = 0
	c.register_x = 0
	c.status = 0

	c.program_counter = c.mem_read_16(0xFFFC)
}

func (c *CPU) run() {
	for {
		opcode := c.mem_read(c.program_counter)
		c.program_counter++
		if opcode == 0x00 {
			return
		}
		op := OPTABLE[opcode]
		op.f_call(c, op)
	}
}

func (c *CPU) load(program []uint8) {
	copy(c.memory[0x8000:], program)

	c.mem_write_16(0xFFFC, 0x8000)
}

func (c *CPU) mem_read(addr uint16) uint8 {
	return c.memory[addr]
}

func (c *CPU) mem_write(addr uint16, v uint8) {
	c.memory[addr] = v
}

func (c *CPU) mem_read_16(addr uint16) uint16 {
	lo := c.memory[addr]
	hi := c.memory[addr+1]
	return make_16_bit(hi, lo)
}

func (c *CPU) mem_write_16(addr uint16, v uint16) {
	lo := uint8(v & 0xFF)
	hi := uint8(v >> 8)
	c.memory[addr] = lo
	c.memory[addr+1] = hi
}

func make_16_bit(hi, lo uint8) uint16 {
	return (uint16(hi) << 8) | uint16(lo)
}

func (c *CPU) lda(op OpCode) {

	c.register_a = c.interpret_mode(op.mode)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) ldy(op OpCode) {
	c.register_y = c.interpret_mode(op.mode)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_y)
}

func (c *CPU) ldx(op OpCode) {
	c.register_x = c.interpret_mode(op.mode)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) tax(op OpCode) {
	c.register_x = c.register_a
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) inx(op OpCode) {
	c.register_x++
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) interpret_mode(m AddressingMode) uint8 {
	var val uint8
	next_val := c.mem_read(c.program_counter)
	switch m {
	case IMMEDIATE:
		val = next_val
	case ZEROPAGE:
		val = c.mem_read(uint16(next_val))
	case ZEROPAGEX:
		val = c.mem_read(uint16(next_val + c.register_x))
	case ZEROPAGEY:
		val = c.mem_read(uint16(next_val + c.register_y))
	case ABSOLUTE:
		addr := c.mem_read_16(c.program_counter)
		c.program_counter++
		val = c.mem_read(addr)
	case ABSOLUTEX:
		addr := c.mem_read_16(c.program_counter)
		c.program_counter++
		val = c.mem_read(addr + uint16(c.register_x))
	case ABSOLUTEY:
		addr := c.mem_read_16(c.program_counter)
		c.program_counter++
		val = c.mem_read(addr + uint16(c.register_y))
	case INDIRECTX:
		addr := next_val + c.register_x
		target := c.mem_read_16(uint16(addr))
		c.program_counter++
		val = c.mem_read(target)
		c.program_counter++
	case INDIRECTY:
		addr := next_val + c.register_y
		target := c.mem_read_16(uint16(addr))
		c.program_counter++
		val = c.mem_read(target)
		c.program_counter++
	default:
		panic("Unknown addresing mode")
	}
	return val
}

func (c *CPU) set_zero_and_negative_flag(v uint8) {
	// Set zero flag if v is 0 else unset 0 flag
	if v == 0 {
		c.status |= 0b0000_0010
	} else {
		c.status &= 0b1111_1101
	}
	// Set negative flag if bit 7 of v is set
	if (v & 0b1000_0000) > 0 {
		c.status |= 0b1000_0000
	} else {
		c.status &= 0b0111_1111
	}
}
