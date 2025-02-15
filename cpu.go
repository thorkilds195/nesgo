package cpu

type CPU struct {
	register_a uint8
	register_x uint8
	status uint8
	program_counter uint16
}

func InitCPU() *CPU {
	return &CPU{
		register_a: 0,
		status:0,
		program_counter:0,
	}
}

func (c *CPU) Interpet(program []uint8) {
	for {
		opcode := program[c.program_counter]
		c.program_counter++
		switch (opcode) {
		case 0xA9:
			param := program[c.program_counter]
			c.program_counter++
			c.register_a = param
			// Set zero flag if A is 0 else unset 0 flag
			if c.register_a == 0 {
				c.status |= 0b0000_0010
			} else {
				c.status &= 0b1111_1101;
			}
			// Set negative flag if bit 7 of A is set
			if (c.register_a & 0b1000_0000) > 0 {
				c.status |= 0b1000_0000
			} else {
				c.status &= 0b0111_1111;
			}
		case 0xAA:
			c.register_x = c.register_a
			// Set zero flag if X is 0 else unset 0 flag
			if c.register_x == 0 {
				c.status |= 0b0000_0010
			} else {
				c.status &= 0b1111_1101;
			}
			// Set negative flag if bit 7 of X is set
			if (c.register_x & 0b1000_0000) > 0 {
				c.status |= 0b1000_0000
			} else {
				c.status &= 0b0111_1111;
			}
		case 0x00:
			return;
		default:
			panic("Not implemented") 
		}
	}
}

func Add(x,y int) int {
	return x+y
}


