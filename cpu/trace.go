package cpu

import "fmt"

func TraceCPU(c *CPU) string {
	var ret string
	op_code, addr := c.GetNextOpCode()
	switch op_code.mode {
	case IMPLIED:
		ret = fmt.Sprintf("%02X  %02X        %s                             A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, op_code.name,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ACCUMULATOR:
		ret = fmt.Sprintf("%02X  %02X        %s A                           A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, op_code.name,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case IMMEDIATE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		ret = fmt.Sprintf("%02X  %02X %02X     %s #$%02X                        A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, p1,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case RELATIVE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		ret = fmt.Sprintf("%02X  %02X %02X     %s $%02X                       A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, addr,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ZEROPAGE, ZEROPAGEX, ZEROPAGEY:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		mem_val := c.MemRead(addr)
		ret = fmt.Sprintf("%02X  %02X %02X     %s $%02X = %02X                    A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, addr, mem_val,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ABSOLUTE, ABSOLUTEX, ABSOLUTEY, INDIRECTX, INDIRECT, INDIRECTY:
		if shouldReturnAddress(op_code.name) {
			p1 := c.MemRead(c.ProgramCounter() + 1)
			p2 := c.MemRead(c.ProgramCounter() + 2)
			mem_val := c.MemRead(addr)
			ret = fmt.Sprintf("%02X  %02X %02X %02X  %s $%04X = %02X                  A:%02X X:%02X Y:%02X P:%02X SP:%02X",
				c.ProgramCounter(), op_code.code, p1, p2, op_code.name, addr, mem_val,
				c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
			)
		} else {
			p1 := c.MemRead(c.ProgramCounter() + 1)
			p2 := c.MemRead(c.ProgramCounter() + 2)
			ret = fmt.Sprintf("%02X  %02X %02X %02X  %s $%04X                       A:%02X X:%02X Y:%02X P:%02X SP:%02X",
				c.ProgramCounter(), op_code.code, p1, p2, op_code.name, addr,
				c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
			)
		}
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
