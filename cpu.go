package cpu

import "fmt"

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
	IMPLIED
	ACCUMULATOR
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
	0xA5: {0xA5, ZEROPAGE, 2, 3, (*CPU).lda},
	0xB5: {0xB5, ZEROPAGEX, 2, 4, (*CPU).lda},
	0xAD: {0xAD, ABSOLUTE, 3, 4, (*CPU).lda},
	0xBD: {0xBD, ABSOLUTEX, 3, 4, (*CPU).lda}, // plus 1 cycle if page crossed
	0xB9: {0xB9, ABSOLUTEY, 3, 4, (*CPU).lda}, // plus 1 cycle if page crossed
	0xA1: {0xA1, INDIRECTX, 2, 6, (*CPU).lda},
	0xB1: {0xB1, INDIRECTY, 2, 5, (*CPU).lda}, // plus 1 cycle if page crossed
	0xA2: {0xA2, IMMEDIATE, 2, 2, (*CPU).ldx},
	0xA6: {0xA6, ZEROPAGE, 2, 3, (*CPU).ldx},
	0xB6: {0xB6, ZEROPAGEY, 2, 4, (*CPU).ldx},
	0xAE: {0xAE, ABSOLUTE, 3, 4, (*CPU).ldx},
	0xBE: {0xBE, ABSOLUTEY, 3, 4, (*CPU).ldx}, // plus 1 cycle if page crossed
	0xA0: {0xA0, IMMEDIATE, 2, 2, (*CPU).ldy},
	0xA4: {0xA4, ZEROPAGE, 2, 3, (*CPU).ldy},
	0xB4: {0xB4, ZEROPAGEX, 2, 4, (*CPU).ldy},
	0xAC: {0xAC, ABSOLUTE, 3, 4, (*CPU).ldy},
	0xBC: {0xBC, ABSOLUTEX, 3, 4, (*CPU).ldy}, // plus 1 cycle if page crossed
	0xAA: {0xAA, IMPLIED, 2, 2, (*CPU).tax},
	0xE8: {0xE8, IMPLIED, 2, 2, (*CPU).inx},
	0x69: {0x69, IMMEDIATE, 2, 2, (*CPU).adc},
	0x65: {0x65, ZEROPAGE, 2, 3, (*CPU).adc},
	0x75: {0x75, ZEROPAGEX, 2, 4, (*CPU).adc},
	0x6D: {0x6D, ABSOLUTE, 3, 4, (*CPU).adc},
	0x7D: {0x7D, ABSOLUTEX, 3, 4, (*CPU).adc},
	0x79: {0x79, ABSOLUTEY, 3, 4, (*CPU).adc},
	0x61: {0x61, INDIRECTX, 3, 4, (*CPU).adc},
	0x71: {0x71, INDIRECTY, 3, 4, (*CPU).adc},
	0x29: {0x29, IMMEDIATE, 2, 2, (*CPU).and},
	0x25: {0x25, ZEROPAGE, 2, 3, (*CPU).and},
	0x35: {0x35, ZEROPAGEX, 2, 4, (*CPU).and},
	0x2D: {0x2D, ABSOLUTE, 3, 4, (*CPU).and},
	0x3D: {0x3D, ABSOLUTEX, 3, 4, (*CPU).and},
	0x39: {0x39, ABSOLUTEY, 3, 4, (*CPU).and},
	0x21: {0x21, INDIRECTX, 3, 4, (*CPU).and},
	0x31: {0x31, INDIRECTY, 3, 4, (*CPU).and},
	0x0A: {0x0A, ACCUMULATOR, 1, 2, (*CPU).asl},
	0x06: {0x06, ZEROPAGE, 2, 5, (*CPU).asl},
	0x16: {0x16, ZEROPAGEX, 2, 6, (*CPU).asl},
	0x0E: {0x0E, ABSOLUTE, 2, 6, (*CPU).asl},
	0x1E: {0x1E, ABSOLUTEX, 2, 7, (*CPU).asl},
	0x90: {0x90, RELATIVE, 2, 2, (*CPU).bcc}, // plus 1 if branch succeeds, plus 2 if new page
	0xB0: {0xB0, RELATIVE, 2, 2, (*CPU).bcs}, // plus 1 if branch succeeds, plus 2 if new page
	0xF0: {0xF0, RELATIVE, 2, 2, (*CPU).beq}, // plus 1 if branch succeeds, plus 2 if new page
	0x24: {0x24, ZEROPAGE, 2, 3, (*CPU).bit},
	0x2C: {0x2C, ABSOLUTE, 2, 3, (*CPU).bit},
	0x30: {0x30, RELATIVE, 2, 2, (*CPU).bmi}, // plus 1 if branch succeeds, plus 2 if new page
	0xD0: {0xD0, RELATIVE, 2, 2, (*CPU).bne}, // plus 1 if branch succeeds, plus 2 if new page
	0x10: {0x10, RELATIVE, 2, 2, (*CPU).bpl}, // plus 1 if branch succeeds, plus 2 if new page
	0x00: {0x00, IMPLIED, 1, 7, (*CPU).brk},
	0x50: {0x50, RELATIVE, 2, 2, (*CPU).bvc}, // plus 1 if branch succeeds, plus 2 if new page
	0x70: {0x70, RELATIVE, 2, 2, (*CPU).bvs}, // plus 1 if branch succeeds, plus 2 if new page
	0x18: {0x18, IMPLIED, 1, 2, (*CPU).clc},
	0xD8: {0xD8, IMPLIED, 1, 2, (*CPU).cld},
	0x58: {0x58, IMPLIED, 1, 2, (*CPU).cli},
	0xB8: {0xB8, IMPLIED, 1, 2, (*CPU).clv},
	0xC9: {0xC9, IMMEDIATE, 2, 2, (*CPU).cmp},
	0xC5: {0xC5, ZEROPAGE, 2, 2, (*CPU).cmp},
	0xD5: {0xD5, ZEROPAGEX, 2, 2, (*CPU).cmp},
	0xCD: {0xCD, ABSOLUTE, 2, 2, (*CPU).cmp},
	0xDD: {0xDD, ABSOLUTEX, 2, 2, (*CPU).cmp},
	0xD9: {0xD9, ABSOLUTEY, 2, 2, (*CPU).cmp},
	0xC1: {0xC1, INDIRECTX, 2, 2, (*CPU).cmp},
	0xD1: {0xD1, INDIRECTY, 2, 2, (*CPU).cmp},
	0xE0: {0xE0, IMMEDIATE, 2, 2, (*CPU).cpx},
	0xE4: {0xE4, ZEROPAGE, 2, 2, (*CPU).cpx},
	0xEC: {0xEC, ABSOLUTE, 2, 2, (*CPU).cpx},
	0xC0: {0xC0, IMMEDIATE, 2, 2, (*CPU).cpy},
	0xC4: {0xC4, ZEROPAGE, 2, 2, (*CPU).cpy},
	0xCC: {0xCC, ABSOLUTE, 2, 2, (*CPU).cpy},
	0xC6: {0xC6, ZEROPAGE, 2, 2, (*CPU).dec},
	0xD6: {0xD6, ZEROPAGEX, 2, 2, (*CPU).dec},
	0xCE: {0xCE, ABSOLUTE, 2, 2, (*CPU).dec},
	0xDE: {0xDE, ABSOLUTEX, 2, 2, (*CPU).dec},
	0xCA: {0xCA, IMPLIED, 2, 2, (*CPU).dex},
	0x88: {0x88, IMPLIED, 2, 2, (*CPU).dey},
	0x49: {0x49, IMMEDIATE, 2, 2, (*CPU).eor},
	0x45: {0x45, ZEROPAGE, 2, 2, (*CPU).eor},
	0x55: {0x55, ZEROPAGEX, 2, 2, (*CPU).eor},
	0x4D: {0x4D, ABSOLUTE, 2, 2, (*CPU).eor},
	0x5D: {0x5D, ABSOLUTEX, 2, 2, (*CPU).eor},
	0x59: {0x59, ABSOLUTEY, 2, 2, (*CPU).eor},
	0x41: {0x41, INDIRECTX, 2, 2, (*CPU).eor},
	0x51: {0x51, INDIRECTY, 2, 2, (*CPU).eor},
	0xE6: {0xE6, ZEROPAGE, 2, 2, (*CPU).inc},
	0xF6: {0xF6, ZEROPAGEX, 2, 2, (*CPU).inc},
	0xEE: {0xEE, ABSOLUTE, 2, 2, (*CPU).inc},
	0xFE: {0xFE, ABSOLUTEX, 2, 2, (*CPU).inc},
	0xC8: {0xC8, IMPLIED, 2, 2, (*CPU).iny},
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
	c.Load(program)
	c.Reset()
	c.Run()
}

func (c *CPU) Reset() {
	c.register_a = 0
	c.register_x = 0
	c.status = 0

	c.program_counter = c.mem_read_16(0xFFFC)
}

func (c *CPU) Run() {
	var op OpCode
	var ok bool
	for {
		opcode := c.mem_read(c.program_counter)
		c.program_counter++
		if op, ok = OPTABLE[opcode]; !ok {
			panic(fmt.Sprintf("No instr found for %x", opcode))
		}
		op.f_call(c, op)
		if opcode == 0x00 {
			return
		}
	}
}

func (c *CPU) Load(program []uint8) {
	copy(c.memory[0x8000:], program)

	c.mem_write_16(0xFFFC, 0x8000)
}

func (c *CPU) brk(op OpCode) {
	// TODO: Implement this in full
	c.program_counter++
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

func (c *CPU) add_carry_bit(v uint8) uint8 {
	return (c.status & 0b0000_0001) + v
}

func (c *CPU) is_carry_set() bool {
	return (c.status & 0b0000_0001) > 0
}

func (c *CPU) is_zero_set() bool {
	return (c.status & 0b0000_0010) > 0
}

func (c *CPU) is_negative_set() bool {
	return (c.status & 0b1000_0000) > 0
}

func (c *CPU) is_overflow_set() bool {
	return (c.status & 0b0100_0000) > 0
}

func (c *CPU) decide_carry_bit(new_v, old_v uint8) {
	if new_v < old_v {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
}

func (c *CPU) set_carry_bit() {
	c.status |= 0b0000_0001
}

func (c *CPU) clear_carry_bit() {
	c.status &= 0b1111_1110
}

func (c *CPU) clear_decimal_bit() {
	c.status &= 0b1111_0111
}

func (c *CPU) clear_interrupt_bit() {
	c.status &= 0b1111_1011
}

func (c *CPU) clear_overflow_bit() {
	c.status &= 0b1011_1111
}

func (c *CPU) compute_overflow_bit(a, b, res uint8) {
	val_sign := (a & 0x80) != 0
	reg_sign := (b & 0x80) != 0
	res_sign := (res & 0x80) != 0
	if (val_sign == reg_sign) && (val_sign != res_sign) {
		c.status |= 0b0100_0000
	} else {
		c.status &= 0b1011_1111
	}
}

func (c *CPU) copy_overflow_flag(v uint8) {
	c.status |= v & 0b0100_0000
}

func (c *CPU) do_compare(val, reg uint8) {
	if reg >= val {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	if val == reg {
		c.set_zero_flag(0)
	} else {
		c.set_zero_flag(1)
	}
	c.set_negative_flag(reg & val)
}

func (c *CPU) cmp(op OpCode) {
	val := c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.do_compare(val, c.register_a)
}

func (c *CPU) cpx(op OpCode) {
	val := c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.do_compare(val, c.register_x)
}

func (c *CPU) cpy(op OpCode) {
	val := c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.do_compare(val, c.register_y)
}

func (c *CPU) clv(op OpCode) {
	c.clear_overflow_bit()
}

func (c *CPU) cld(op OpCode) {
	c.clear_decimal_bit()
}

func (c *CPU) clc(op OpCode) {
	c.clear_carry_bit()
}

func (c *CPU) cli(op OpCode) {
	c.clear_interrupt_bit()
}

func (c *CPU) bne(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if !c.is_zero_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) bmi(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if c.is_negative_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) bvs(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if !c.is_overflow_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) bpl(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if !c.is_negative_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) bvc(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if c.is_overflow_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) bit(op OpCode) {
	val := c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.set_zero_flag(val & c.register_a)
	c.copy_overflow_flag(val)
	c.set_negative_flag(val)
}

func (c *CPU) and(op OpCode) {
	c.register_a &= c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) eor(op OpCode) {
	c.register_a ^= c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) asl(op OpCode) {
	if op.mode == ACCUMULATOR {
		c.register_a <<= 1
	} else {
		var addr uint16
		val := c.interpret_mode(op.mode, &addr)
		val <<= 1
		c.mem_write(addr, val)
	}
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) bcc(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if c.is_carry_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) bcs(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if !c.is_carry_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) beq(op OpCode) {
	rel := c.interpret_mode(op.mode, nil)
	c.program_counter++
	if !c.is_zero_set() {
		return
	}
	c.program_counter += uint16(int16(int8(rel)))
}

func (c *CPU) adc(op OpCode) {
	val := c.interpret_mode(op.mode, nil)
	val = c.add_carry_bit(val)
	result := val + c.register_a
	c.decide_carry_bit(result, c.register_a)
	c.compute_overflow_bit(val, c.register_a, result)
	c.register_a = result
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) lda(op OpCode) {

	c.register_a = c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) ldy(op OpCode) {
	c.register_y = c.interpret_mode(op.mode, nil)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_y)
}

func (c *CPU) ldx(op OpCode) {
	c.register_x = c.interpret_mode(op.mode, nil)
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

func (c *CPU) iny(op OpCode) {
	c.register_y++
	c.set_zero_and_negative_flag(c.register_y)
}

func (c *CPU) dec(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr)
	c.program_counter++
	val--
	c.mem_write(addr, val)
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) inc(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr)
	c.program_counter++
	val++
	c.mem_write(addr, val)
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) dex(op OpCode) {
	c.register_x--
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) dey(op OpCode) {
	c.register_y--
	c.set_zero_and_negative_flag(c.register_y)
}

func (c *CPU) interpret_mode(m AddressingMode, read_adr *uint16) uint8 {
	var val uint8
	var addr uint16
	next_val := c.mem_read(c.program_counter)
	switch m {
	case IMMEDIATE, RELATIVE:
		val = next_val
	case ZEROPAGE:
		addr = uint16(next_val)
		val = c.mem_read(addr)
	case ZEROPAGEX:
		addr = uint16(next_val + c.register_x)
		val = c.mem_read(addr)
	case ZEROPAGEY:
		addr = uint16(next_val + c.register_y)
		val = c.mem_read(addr)
	case ABSOLUTE:
		addr = c.mem_read_16(c.program_counter)
		c.program_counter++
		val = c.mem_read(addr)
	case ABSOLUTEX:
		in := c.mem_read_16(c.program_counter)
		c.program_counter++
		addr = in + uint16(c.register_x)
		val = c.mem_read(addr)
	case ABSOLUTEY:
		in := c.mem_read_16(c.program_counter)
		c.program_counter++
		addr = in + uint16(c.register_y)
		val = c.mem_read(addr)
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
	if read_adr != nil {
		*read_adr = addr
	}
	return val
}

func (c *CPU) set_zero_and_negative_flag(v uint8) {
	c.set_zero_flag(v)
	c.set_negative_flag(v)
}
func (c *CPU) set_negative_flag(v uint8) {
	// Set negative flag if bit 7 of v is set
	if (v & 0b1000_0000) > 0 {
		c.status |= 0b1000_0000
	} else {
		c.status &= 0b0111_1111
	}
}

func (c *CPU) set_zero_flag(v uint8) {
	// Set zero flag if v is 0 else unset 0 flag
	if v == 0 {
		c.status |= 0b0000_0010
	} else {
		c.status &= 0b1111_1101
	}
}
