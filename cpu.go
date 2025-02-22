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
)

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
		switch opcode {
		case 0xA9:
			c.lda(IMMEDIATE)
		case 0xA5:
			c.lda(ZEROPAGE)
		case 0xB5:
			c.lda(ZEROPAGEX)
		case 0xAD:
			c.lda(ABSOLUTE)
		case 0xBD:
			c.lda(ABSOLUTEX)
		case 0xB9:
			c.lda(ABSOLUTEY)
		case 0xA1:
			c.lda(INDIRECTX)
		case 0xB1:
			c.lda(INDIRECTY)
		case 0xA2:
			c.ldx(IMMEDIATE)
		case 0xA0:
			c.ldy(IMMEDIATE)
		case 0xAA:
			c.tax()
		case 0xE8:
			c.inx()
		case 0x00:
			return
		default:
			panic("Not implemented")
		}
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

func (c *CPU) lda(m AddressingMode) {
	var val uint8
	next_val := c.mem_read(c.program_counter)
	switch m {
	case IMMEDIATE:
		val = next_val
	case ZEROPAGE:
		val = c.mem_read(uint16(next_val))
	case ZEROPAGEX:
		val = c.mem_read(uint16(next_val + c.register_x))
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
	c.program_counter++
	c.register_a = val
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) ldx(m AddressingMode) {
	var val uint8
	next_val := c.mem_read(c.program_counter)
	switch m {
	case IMMEDIATE:
		val = next_val
	case ZEROPAGE:
		val = c.mem_read(uint16(next_val))

	default:
		panic("Unknown addresing mode")
	}
	c.program_counter++
	c.register_x = val
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) ldy(m AddressingMode) {
	var val uint8
	next_val := c.mem_read(c.program_counter)
	switch m {
	case IMMEDIATE:
		val = next_val
	case ZEROPAGE:
		val = c.mem_read(uint16(next_val))

	default:
		panic("Unknown addresing mode")
	}
	c.program_counter++
	c.register_y = val
	c.set_zero_and_negative_flag(c.register_y)
}

func (c *CPU) tax() {
	c.register_x = c.register_a
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) inx() {
	c.register_x++
	c.set_zero_and_negative_flag(c.register_x)
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
