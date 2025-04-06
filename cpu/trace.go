package cpu

import "fmt"

func TraceCPU(c *CPU) string {
	var ret string
	op_code, addr := c.GetNextOpCode()
	switch op_code.mode {
	case IMPLIED, ACCUMULATOR:
		ret = fmt.Sprintf("%02X  %02X        %s                             A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, op_code.name,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case IMMEDIATE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		ret = fmt.Sprintf("%02X  %02X %02X     %s $%02X                         A:%02X X:%02X Y:%02X P:%02X SP:%02X",
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
		ret = fmt.Sprintf("%02X  %02X %02X     %s $%02X                         A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, op_code.name, addr,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	case ABSOLUTE, ABSOLUTEX, ABSOLUTEY, INDIRECTX, INDIRECT, INDIRECTY:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		p2 := c.MemRead(c.ProgramCounter() + 2)
		ret = fmt.Sprintf("%02X  %02X %02X %02X  %s $%02X                       A:%02X X:%02X Y:%02X P:%02X SP:%02X",
			c.ProgramCounter(), op_code.code, p1, p2, op_code.name, addr,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		)
	}
	return ret
}
