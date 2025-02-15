package cpu

type CPU struct {
	register_a      uint8
	register_x      uint8
	status          uint8
	program_counter uint16
}

func InitCPU() *CPU {
	return &CPU{
		register_a:      0,
		status:          0,
		program_counter: 0,
	}
}

func (c *CPU) Interpet(program []uint8) {
	for {
		opcode := program[c.program_counter]
		c.program_counter++
		switch opcode {
		case 0xA9:
			param := program[c.program_counter]
			c.program_counter++
			c.lda(param)
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

func (c *CPU) lda(v uint8) {
	c.register_a = v
	c.set_zero_and_negative_flag(c.register_a)
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
