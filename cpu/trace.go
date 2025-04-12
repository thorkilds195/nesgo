package cpu

import "fmt"

func TraceCPU(c *CPU) string {
	var ret string
	op_code, addr := c.GetNextOpCode()
	switch op_code.mode {
	case IMPLIED:
		ret = fmt.Sprintf("%04X  %02X        %s                             A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, op_code.name,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ACCUMULATOR:
		ret = fmt.Sprintf("%04X  %02X        %s A                           A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, op_code.name,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case IMMEDIATE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		ret = fmt.Sprintf("%04X  %02X %02X     %s #$%02X                        A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, p1,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case RELATIVE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		ret = fmt.Sprintf("%04X  %02X %02X     %s $%02X                       A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, addr,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ZEROPAGE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		mem_val := c.MemRead(addr)
		ret = fmt.Sprintf("%04X  %02X %02X     %s $%02X = %02X                    A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, addr, mem_val,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ZEROPAGEX, ZEROPAGEY:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		mem_val := c.MemRead(addr)
		reg := 'X'
		if op_code.mode == ZEROPAGEY {
			reg = 'Y'
		}
		ret = fmt.Sprintf("%04X  %02X %02X     %s $%02X,%c @ %02X = %02X             A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, p1, reg, addr, mem_val,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ABSOLUTE:
		if shouldReturnAddress(op_code.name) {
			p1 := c.MemRead(c.ProgramCounter() + 1)
			p2 := c.MemRead(c.ProgramCounter() + 2)
			mem_val := c.MemRead(addr)
			ret = fmt.Sprintf("%04X  %02X %02X %02X  %s $%04X = %02X                  A:%02X X:%02X Y:%02X P:%02X SP:%02X",
				c.ProgramCounter(), op_code.code, p1, p2, op_code.name, addr, mem_val,
				c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
			)
		} else {
			p1 := c.MemRead(c.ProgramCounter() + 1)
			p2 := c.MemRead(c.ProgramCounter() + 2)
			ret = fmt.Sprintf("%04X  %02X %02X %02X  %s $%04X                       A:%02X X:%02X Y:%02X P:%02X SP:%02X",
				c.ProgramCounter(), op_code.code, p1, p2, op_code.name, addr,
				c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
			)
		}
	case ABSOLUTEY, ABSOLUTEX:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		p2 := c.MemRead(c.ProgramCounter() + 2)
		reg := 'X'
		if op_code.mode == ABSOLUTEY {
			reg = 'Y'
		}
		in := make_16_bit(p2, p1)
		mem_val := c.MemRead(addr)
		ret = fmt.Sprintf("%04X  %02X %02X %02X  %s $%04X,%c @ %04X = %02X         A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, p2, op_code.name, in, reg, addr, mem_val,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)

	case INDIRECTX:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		reg := 'X'
		offset := c.GetRegisterX() + p1
		mem_val := c.MemRead(addr)
		ret = fmt.Sprintf("%04X  %02X %02X     %s ($%02X,%c) @ %02X = %04X = %02X    A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, p1, reg, offset, addr, mem_val,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case INDIRECTY:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		base_ptr := c.mem_read_16_zero(p1)
		reg := 'Y'
		mem_val := c.MemRead(addr)
		ret = fmt.Sprintf("%04X  %02X %02X     %s ($%02X),%c = %04X @ %04X = %02X  A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, p1, reg, base_ptr, addr, mem_val,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case INDIRECT:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		p2 := c.MemRead(c.ProgramCounter() + 2)
		base_ptr := make_16_bit(p2, p1)
		ret = fmt.Sprintf("%04X  %02X %02X %02X  %s ($%04X) = %04X              A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, p2, op_code.name, base_ptr, addr,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)

	}

	return ret
}

func shouldReturnAddress(instr string) bool {
	switch instr {
	case "JMP", "JSR", "BNE":
		return false
	default:
		return true
	}
}
