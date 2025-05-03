package cpu

import (
	"errors"
	"fmt"
)

func TraceCPU(c *CPU) string {
	var ret string
	var instr_part, assem_part string
	instr_len := 10
	assem_len := 32
	op_code, addr := c.GetNextOpCode()
	is_unsupported := op_code.name[0] == '*'
	if is_unsupported {
		instr_len = 9
		assem_len = 33
	}
	switch op_code.mode {
	case IMPLIED:
		instr_part = fmt.Sprintf("%02X", op_code.code)
		assem_part = op_code.name
	case ACCUMULATOR:
		instr_part = fmt.Sprintf("%02X", op_code.code)
		assem_part = fmt.Sprintf("%s A", op_code.name)
	case IMMEDIATE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		instr_part = fmt.Sprintf("%02X %02X", op_code.code, p1)
		assem_part = fmt.Sprintf("%s #$%02X", op_code.name, p1)
	case RELATIVE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		instr_part = fmt.Sprintf("%02X %02X", op_code.code, p1)
		assem_part = fmt.Sprintf("%s $%02X", op_code.name, addr)
	case ZEROPAGE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		mem_val := c.MemRead(addr)
		instr_part = fmt.Sprintf("%02X %02X", op_code.code, p1)
		assem_part = fmt.Sprintf("%s $%02X = %02X", op_code.name, addr, mem_val)
	case ZEROPAGEX, ZEROPAGEY:
		reg := 'X'
		if op_code.mode == ZEROPAGEY {
			reg = 'Y'
		}
		p1 := c.MemRead(c.ProgramCounter() + 1)
		mem_val := c.MemRead(addr)
		instr_part = fmt.Sprintf("%02X %02X", op_code.code, p1)
		assem_part = fmt.Sprintf("%s $%02X,%c @ %02X = %02X", op_code.name, p1, reg, addr, mem_val)
	case ABSOLUTE:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		p2 := c.MemRead(c.ProgramCounter() + 2)
		mem_val := c.MemRead(addr)
		instr_part = fmt.Sprintf("%02X %02X %02X", op_code.code, p1, p2)
		if shouldReturnAddress(op_code.name) {
			assem_part = fmt.Sprintf("%s $%04X = %02X", op_code.name, addr, mem_val)
		} else {
			assem_part = fmt.Sprintf("%s $%04X", op_code.name, addr)
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
		instr_part = fmt.Sprintf("%02X %02X %02X", op_code.code, p1, p2)
		assem_part = fmt.Sprintf("%s $%04X,%c @ %04X = %02X", op_code.name, in, reg, addr, mem_val)
	case INDIRECTX:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		offset := c.GetRegisterX() + p1
		mem_val := c.MemRead(addr)
		instr_part = fmt.Sprintf("%02X %02X", op_code.code, p1)
		assem_part = fmt.Sprintf("%s ($%02X,X) @ %02X = %04X = %02X", op_code.name, p1, offset, addr, mem_val)
	case INDIRECTY:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		base_ptr := c.mem_read_16_zero(p1)
		mem_val := c.MemRead(addr)
		instr_part = fmt.Sprintf("%02X %02X", op_code.code, p1)
		assem_part = fmt.Sprintf("%s ($%02X),Y = %04X @ %04X = %02X", op_code.name, p1, base_ptr, addr, mem_val)
	case INDIRECT:
		p1 := c.MemRead(c.ProgramCounter() + 1)
		p2 := c.MemRead(c.ProgramCounter() + 2)
		base_ptr := make_16_bit(p2, p1)
		instr_part = fmt.Sprintf("%02X %02X %02X", op_code.code, p1, p2)
		assem_part = fmt.Sprintf("%s ($%04X) = %04X", op_code.name, base_ptr, addr)
	}
	instr_str, _ := padSpaces(instr_part, instr_len)
	assem_str, _ := padSpaces(assem_part, assem_len)
	ret = fmt.Sprintf("%04X  %s%sA:%02X X:%02X Y:%02X P:%02X SP:%02X CYC:%d",
		c.ProgramCounter(), instr_str, assem_str,
		c.GetRegisterA(), c.GetRegisterX(), c.GetRegisterY(), c.GetStatus(), c.GetStackPointer(),
		c.bus.cycles,
	)
	return ret
}

func padSpaces(s string, final_len int) (string, error) {
	if len(s) > final_len {
		return s, errors.New("input string is longer than final length")
	}
	if len(s) == final_len {
		return s, nil
	}
	bytes := make([]byte, final_len)
	for i := 0; i < final_len; i++ {
		if i < len(s) {
			bytes[i] = s[i]
		} else {
			bytes[i] = ' '
		}
	}
	return string(bytes), nil
}

func shouldReturnAddress(instr string) bool {
	switch instr {
	case "JMP", "JSR", "BNE":
		return false
	default:
		return true
	}
}
