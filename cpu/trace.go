package cpu

import "fmt"

func TraceCPU(c *CPU) string {
	var ret string
	op_code, addr := c.GetNextOpCode()
	switch op_code.mode {
	case IMMEDIATE, RELATIVE:
		ret = fmt.Sprintf("%x  %x %x", c.ProgramCounter(), op_code.code, c.MemRead(c.ProgramCounter()+1))
	case ZEROPAGE, ZEROPAGEX, ZEROPAGEY:
		ret = fmt.Sprintf("%x  %x %x", c.ProgramCounter(), op_code.code, c.MemRead(c.ProgramCounter()+1))
	case ABSOLUTE, ABSOLUTEX, ABSOLUTEY, INDIRECTX, INDIRECT, INDIRECTY:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		p2 := c.MemRead(c.ProgramCounter() + 2)
		ret = fmt.Sprintf("%02X  %02X %02X %02X  %s $%02X                       A:%02X X:%02X Y:%02X P:00 SP:%02X",
			c.ProgramCounter(), op_code.code, p1, p2, op_code.name, addr,
			c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStackPointer(),
		)
	}
	return ret
}
