package cpu

import "fmt"

type AddressingMode uint8

const STACK_RESET uint8 = 0xFD
const PROGRAM_START uint16 = 0xC000

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
	INDIRECT
)

type OpCode struct {
	code   uint8
	name   string
	mode   AddressingMode
	bytes  uint8
	cycles uint8
	f_call func(*CPU, OpCode)
}

var OPTABLE = map[uint8]OpCode{
	0xA9: {0xA9, "LDA", IMMEDIATE, 2, 2, (*CPU).lda},
	0xA5: {0xA5, "LDA", ZEROPAGE, 2, 3, (*CPU).lda},
	0xB5: {0xB5, "LDA", ZEROPAGEX, 2, 4, (*CPU).lda},
	0xAD: {0xAD, "LDA", ABSOLUTE, 3, 4, (*CPU).lda},
	0xBD: {0xBD, "LDA", ABSOLUTEX, 3, 4, (*CPU).lda}, // plus 1 cycle if page crossed
	0xB9: {0xB9, "LDA", ABSOLUTEY, 3, 4, (*CPU).lda}, // plus 1 cycle if page crossed
	0xA1: {0xA1, "LDA", INDIRECTX, 2, 6, (*CPU).lda},
	0xB1: {0xB1, "LDA", INDIRECTY, 2, 5, (*CPU).lda}, // plus 1 cycle if page crossed
	0xA2: {0xA2, "LDX", IMMEDIATE, 2, 2, (*CPU).ldx},
	0xA6: {0xA6, "LDX", ZEROPAGE, 2, 3, (*CPU).ldx},
	0xB6: {0xB6, "LDX", ZEROPAGEY, 2, 4, (*CPU).ldx},
	0xAE: {0xAE, "LDX", ABSOLUTE, 3, 4, (*CPU).ldx},
	0xBE: {0xBE, "LDX", ABSOLUTEY, 3, 4, (*CPU).ldx}, // plus 1 cycle if page crossed
	0xA0: {0xA0, "LDY", IMMEDIATE, 2, 2, (*CPU).ldy},
	0xA4: {0xA4, "LDY", ZEROPAGE, 2, 3, (*CPU).ldy},
	0xB4: {0xB4, "LDY", ZEROPAGEX, 2, 4, (*CPU).ldy},
	0xAC: {0xAC, "LDY", ABSOLUTE, 3, 4, (*CPU).ldy},
	0xBC: {0xBC, "LDY", ABSOLUTEX, 3, 4, (*CPU).ldy}, // plus 1 cycle if page crossed
	0xAA: {0xAA, "TAX", IMPLIED, 2, 2, (*CPU).tax},
	0xE8: {0xE8, "INX", IMPLIED, 2, 2, (*CPU).inx},
	0x69: {0x69, "ADC", IMMEDIATE, 2, 2, (*CPU).adc},
	0x65: {0x65, "ADC", ZEROPAGE, 2, 3, (*CPU).adc},
	0x75: {0x75, "ADC", ZEROPAGEX, 2, 4, (*CPU).adc},
	0x6D: {0x6D, "ADC", ABSOLUTE, 3, 4, (*CPU).adc},
	0x7D: {0x7D, "ADC", ABSOLUTEX, 3, 4, (*CPU).adc},
	0x79: {0x79, "ADC", ABSOLUTEY, 3, 4, (*CPU).adc},
	0x61: {0x61, "ADC", INDIRECTX, 3, 4, (*CPU).adc},
	0x71: {0x71, "ADC", INDIRECTY, 3, 4, (*CPU).adc},
	0x29: {0x29, "AND", IMMEDIATE, 2, 2, (*CPU).and},
	0x25: {0x25, "AND", ZEROPAGE, 2, 3, (*CPU).and},
	0x35: {0x35, "AND", ZEROPAGEX, 2, 4, (*CPU).and},
	0x2D: {0x2D, "AND", ABSOLUTE, 3, 4, (*CPU).and},
	0x3D: {0x3D, "AND", ABSOLUTEX, 3, 4, (*CPU).and},
	0x39: {0x39, "AND", ABSOLUTEY, 3, 4, (*CPU).and},
	0x21: {0x21, "AND", INDIRECTX, 3, 4, (*CPU).and},
	0x31: {0x31, "AND", INDIRECTY, 3, 4, (*CPU).and},
	0x0A: {0x0A, "ASL", ACCUMULATOR, 1, 2, (*CPU).asl},
	0x06: {0x06, "ASL", ZEROPAGE, 2, 5, (*CPU).asl},
	0x16: {0x16, "ASL", ZEROPAGEX, 2, 6, (*CPU).asl},
	0x0E: {0x0E, "ASL", ABSOLUTE, 2, 6, (*CPU).asl},
	0x1E: {0x1E, "ASL", ABSOLUTEX, 2, 7, (*CPU).asl},
	0x90: {0x90, "BCC", RELATIVE, 2, 2, (*CPU).bcc}, // plus 1 if branch succeeds, plus 2 if new page
	0xB0: {0xB0, "BCS", RELATIVE, 2, 2, (*CPU).bcs}, // plus 1 if branch succeeds, plus 2 if new page
	0xF0: {0xF0, "BEQ", RELATIVE, 2, 2, (*CPU).beq}, // plus 1 if branch succeeds, plus 2 if new page
	0x24: {0x24, "BIT", ZEROPAGE, 2, 3, (*CPU).bit},
	0x2C: {0x2C, "BIT", ABSOLUTE, 2, 3, (*CPU).bit},
	0x30: {0x30, "BMI", RELATIVE, 2, 2, (*CPU).bmi}, // plus 1 if branch succeeds, plus 2 if new page
	0xD0: {0xD0, "BNE", RELATIVE, 2, 2, (*CPU).bne}, // plus 1 if branch succeeds, plus 2 if new page
	0x10: {0x10, "BPL", RELATIVE, 2, 2, (*CPU).bpl}, // plus 1 if branch succeeds, plus 2 if new page
	0x00: {0x00, "BRK", IMPLIED, 1, 7, (*CPU).brk},
	0x50: {0x50, "BVC", RELATIVE, 2, 2, (*CPU).bvc}, // plus 1 if branch succeeds, plus 2 if new page
	0x70: {0x70, "BVS", RELATIVE, 2, 2, (*CPU).bvs}, // plus 1 if branch succeeds, plus 2 if new page
	0x18: {0x18, "CLC", IMPLIED, 1, 2, (*CPU).clc},
	0xD8: {0xD8, "CLD", IMPLIED, 1, 2, (*CPU).cld},
	0x58: {0x58, "CLI", IMPLIED, 1, 2, (*CPU).cli},
	0xB8: {0xB8, "CLV", IMPLIED, 1, 2, (*CPU).clv},
	0xC9: {0xC9, "CMP", IMMEDIATE, 2, 2, (*CPU).cmp},
	0xC5: {0xC5, "CMP", ZEROPAGE, 2, 2, (*CPU).cmp},
	0xD5: {0xD5, "CMP", ZEROPAGEX, 2, 2, (*CPU).cmp},
	0xCD: {0xCD, "CMP", ABSOLUTE, 2, 2, (*CPU).cmp},
	0xDD: {0xDD, "CMP", ABSOLUTEX, 2, 2, (*CPU).cmp},
	0xD9: {0xD9, "CMP", ABSOLUTEY, 2, 2, (*CPU).cmp},
	0xC1: {0xC1, "CMP", INDIRECTX, 2, 2, (*CPU).cmp},
	0xD1: {0xD1, "CMP", INDIRECTY, 2, 2, (*CPU).cmp},
	0xE0: {0xE0, "CPX", IMMEDIATE, 2, 2, (*CPU).cpx},
	0xE4: {0xE4, "CPX", ZEROPAGE, 2, 2, (*CPU).cpx},
	0xEC: {0xEC, "CPX", ABSOLUTE, 2, 2, (*CPU).cpx},
	0xC0: {0xC0, "CPY", IMMEDIATE, 2, 2, (*CPU).cpy},
	0xC4: {0xC4, "CPY", ZEROPAGE, 2, 2, (*CPU).cpy},
	0xCC: {0xCC, "CPY", ABSOLUTE, 2, 2, (*CPU).cpy},
	0xC6: {0xC6, "DEC", ZEROPAGE, 2, 2, (*CPU).dec},
	0xD6: {0xD6, "DEC", ZEROPAGEX, 2, 2, (*CPU).dec},
	0xCE: {0xCE, "DEC", ABSOLUTE, 2, 2, (*CPU).dec},
	0xDE: {0xDE, "DEC", ABSOLUTEX, 2, 2, (*CPU).dec},
	0xCA: {0xCA, "DEX", IMPLIED, 2, 2, (*CPU).dex},
	0x88: {0x88, "DEY", IMPLIED, 2, 2, (*CPU).dey},
	0x49: {0x49, "EOR", IMMEDIATE, 2, 2, (*CPU).eor},
	0x45: {0x45, "EOR", ZEROPAGE, 2, 2, (*CPU).eor},
	0x55: {0x55, "EOR", ZEROPAGEX, 2, 2, (*CPU).eor},
	0x4D: {0x4D, "EOR", ABSOLUTE, 2, 2, (*CPU).eor},
	0x5D: {0x5D, "EOR", ABSOLUTEX, 2, 2, (*CPU).eor},
	0x59: {0x59, "EOR", ABSOLUTEY, 2, 2, (*CPU).eor},
	0x41: {0x41, "EOR", INDIRECTX, 2, 2, (*CPU).eor},
	0x51: {0x51, "EOR", INDIRECTY, 2, 2, (*CPU).eor},
	0xE6: {0xE6, "INC", ZEROPAGE, 2, 2, (*CPU).inc},
	0xF6: {0xF6, "INC", ZEROPAGEX, 2, 2, (*CPU).inc},
	0xEE: {0xEE, "INC", ABSOLUTE, 2, 2, (*CPU).inc},
	0xFE: {0xFE, "INC", ABSOLUTEX, 2, 2, (*CPU).inc},
	0xC8: {0xC8, "INY", IMPLIED, 2, 2, (*CPU).iny},
	0x4C: {0x4C, "JMP", ABSOLUTE, 3, 3, (*CPU).jmp},
	0x6C: {0x6C, "JMP", INDIRECT, 3, 5, (*CPU).jmp},
	0x4A: {0x4A, "LSR", ACCUMULATOR, 1, 2, (*CPU).lsr},
	0x46: {0x46, "LSR", ZEROPAGE, 2, 5, (*CPU).lsr},
	0x56: {0x56, "LSR", ZEROPAGEX, 2, 6, (*CPU).lsr},
	0x4E: {0x4E, "LSR", ABSOLUTE, 2, 6, (*CPU).lsr},
	0x5E: {0x5E, "LSR", ABSOLUTEX, 2, 7, (*CPU).lsr},
	0xEA: {0xEA, "NOP", IMPLIED, 2, 7, (*CPU).nop},
	0x09: {0x09, "ORA", IMMEDIATE, 2, 2, (*CPU).ora},
	0x05: {0x05, "ORA", ZEROPAGE, 2, 2, (*CPU).ora},
	0x15: {0x15, "ORA", ZEROPAGEX, 2, 2, (*CPU).ora},
	0x0D: {0x0D, "ORA", ABSOLUTE, 2, 2, (*CPU).ora},
	0x1D: {0x1D, "ORA", ABSOLUTEX, 2, 2, (*CPU).ora},
	0x19: {0x19, "ORA", ABSOLUTEY, 2, 2, (*CPU).ora},
	0x01: {0x01, "ORA", INDIRECTX, 2, 2, (*CPU).ora},
	0x11: {0x11, "ORA", INDIRECTY, 2, 2, (*CPU).ora},
	0x2A: {0x2A, "ROL", ACCUMULATOR, 1, 2, (*CPU).rol},
	0x26: {0x26, "ROL", ZEROPAGE, 2, 5, (*CPU).rol},
	0x36: {0x36, "ROL", ZEROPAGEX, 2, 6, (*CPU).rol},
	0x2E: {0x2E, "ROL", ABSOLUTE, 2, 6, (*CPU).rol},
	0x3E: {0x3E, "ROL", ABSOLUTEX, 2, 7, (*CPU).rol},
	0x6A: {0x6A, "ROR", ACCUMULATOR, 1, 2, (*CPU).ror},
	0x66: {0x66, "ROR", ZEROPAGE, 2, 5, (*CPU).ror},
	0x76: {0x76, "ROR", ZEROPAGEX, 2, 6, (*CPU).ror},
	0x6E: {0x6E, "ROR", ABSOLUTE, 2, 6, (*CPU).ror},
	0x7E: {0x7E, "ROR", ABSOLUTEX, 2, 7, (*CPU).ror},
	0xE9: {0xE9, "SBC", IMMEDIATE, 2, 2, (*CPU).sbc},
	0xE5: {0xE5, "SBC", ZEROPAGE, 2, 2, (*CPU).sbc},
	0xF5: {0xF5, "SBC", ZEROPAGEX, 2, 2, (*CPU).sbc},
	0xED: {0xED, "SBC", ABSOLUTE, 2, 2, (*CPU).sbc},
	0xFD: {0xFD, "SBC", ABSOLUTEX, 2, 2, (*CPU).sbc},
	0xF9: {0xF9, "SBC", ABSOLUTEY, 2, 2, (*CPU).sbc},
	0xE1: {0xE1, "SBC", INDIRECTX, 2, 2, (*CPU).sbc},
	0xF1: {0xF1, "SBC", INDIRECTY, 2, 2, (*CPU).sbc},
	0x38: {0x38, "SEC", IMPLIED, 2, 2, (*CPU).sec},
	0xF8: {0xF8, "SED", IMPLIED, 2, 2, (*CPU).sed},
	0x78: {0x78, "SEI", IMPLIED, 2, 2, (*CPU).sei},
	0x85: {0x85, "STA", ZEROPAGE, 2, 2, (*CPU).sta},
	0x95: {0x95, "STA", ZEROPAGEX, 2, 2, (*CPU).sta},
	0x8D: {0x8D, "STA", ABSOLUTE, 2, 2, (*CPU).sta},
	0x9D: {0x9D, "STA", ABSOLUTEX, 2, 2, (*CPU).sta},
	0x99: {0x99, "STA", ABSOLUTEY, 2, 2, (*CPU).sta},
	0x81: {0x81, "STA", INDIRECTX, 2, 2, (*CPU).sta},
	0x91: {0x91, "STA", INDIRECTY, 2, 2, (*CPU).sta},
	0x86: {0x86, "STX", ZEROPAGE, 2, 2, (*CPU).stx},
	0x96: {0x96, "STX", ZEROPAGEY, 2, 2, (*CPU).stx},
	0x8E: {0x8E, "STX", ABSOLUTE, 2, 2, (*CPU).stx},
	0x84: {0x84, "STY", ZEROPAGE, 2, 2, (*CPU).sty},
	0x94: {0x94, "STY", ZEROPAGEX, 2, 2, (*CPU).sty},
	0x8C: {0x8C, "STY", ABSOLUTE, 2, 2, (*CPU).sty},
	0xA8: {0xA8, "TAY", IMPLIED, 2, 2, (*CPU).tay},
	0x8A: {0x8A, "TXA", IMPLIED, 2, 2, (*CPU).txa},
	0x98: {0x98, "TYA", IMPLIED, 2, 2, (*CPU).tya},
	0x20: {0x20, "JSR", ABSOLUTE, 2, 2, (*CPU).jsr},
	0x48: {0x48, "PHA", IMPLIED, 2, 2, (*CPU).pha},
	0x08: {0x08, "PHP", IMPLIED, 2, 2, (*CPU).php},
	0x68: {0x68, "PLA", IMPLIED, 2, 2, (*CPU).pla},
	0x28: {0x28, "PLP", IMPLIED, 2, 2, (*CPU).plp},
	0x40: {0x40, "RTI", IMPLIED, 2, 2, (*CPU).rti},
	0x60: {0x60, "RTS", IMPLIED, 2, 2, (*CPU).rts},
	0xBA: {0xBA, "TSX", IMPLIED, 2, 2, (*CPU).tsx},
	0x9A: {0x9A, "TXS", IMPLIED, 2, 2, (*CPU).txs},
	// Custom instruction to quit and leave emulator in current state
	0x02: {0x02, "HLT", IMPLIED, 2, 2, (*CPU).hlt},
}

type CPU struct {
	register_a      uint8
	register_x      uint8
	register_y      uint8
	status          uint8
	program_counter uint16
	stack_pointer   uint8
	memory          [0x10000]uint8
}

func (c *CPU) ProgramCounter() uint16 {
	return c.program_counter
}

func (c *CPU) GetRegisterX() uint8 {
	return c.register_x
}

func (c *CPU) GetRegisterY() uint8 {
	return c.register_y
}

func (c *CPU) GetRegisterA() uint8 {
	return c.register_a
}

func (c *CPU) GetStatus() uint8 {
	return c.status
}

func (c *CPU) GetStackPointer() uint8 {
	return c.stack_pointer
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
	c.register_y = 0
	c.register_x = 0
	c.status = 0b100100

	c.program_counter = c.MemRead16(0xFFFC)
	c.stack_pointer = STACK_RESET
}

func (c *CPU) RunWithCallback(f_call func()) {
	var op OpCode
	var ok bool
	for {
		f_call()
		opcode := c.MemRead(c.program_counter)
		c.program_counter++
		if op, ok = OPTABLE[opcode]; !ok {
			panic(fmt.Sprintf("No instr found for %x", opcode))
		}
		op.f_call(c, op)
		if opcode == 0x00 || opcode == 0x02 {
			return
		}
	}
}

func (c *CPU) Run() {
	c.RunWithCallback(func() {})
}

func (c *CPU) Load(program []uint8) {
	copy(c.memory[PROGRAM_START:], program)

	c.MemWrite16(0xFFFC, PROGRAM_START)
}

func (c *CPU) Step(f_call func()) bool {
	f_call()
	opcode := c.MemRead(c.program_counter)
	c.program_counter++
	op, ok := OPTABLE[opcode]
	if !ok {
		panic(fmt.Sprintf("Unknown opcode: %x", opcode))
	}
	op.f_call(c, op)
	return opcode != 0x00 && opcode != 0x02
}

func (c *CPU) GetNextOpCode() (OpCode, uint16) {
	var addr uint16
	var op OpCode
	var ok bool
	opcode := c.MemRead(c.program_counter)
	if op, ok = OPTABLE[opcode]; !ok {
		panic(fmt.Sprintf("No instr found for %x", opcode))
	}
	if op.mode == IMPLIED || op.mode == ACCUMULATOR {
		return op, 0
	}
	// This is pretty hacked together and should be fixed
	c.program_counter++
	c.interpret_mode(op.mode, &addr, false)
	c.program_counter--
	return op, addr
}

func (c *CPU) push(val uint8) {
	c.memory[0x0100+uint16(c.stack_pointer)] = val
	c.stack_pointer--
}

func (c *CPU) pull() uint8 {
	c.stack_pointer++
	return c.MemRead(0x0100 + uint16(c.stack_pointer))
}

func (c *CPU) push_16(val uint16) {
	lo := uint8(val & 0xFF)
	hi := uint8(val >> 8)
	c.push(hi)
	c.push(lo)
}

func (c *CPU) brk(op OpCode) {
	c.push_16(c.program_counter + 2)
	l_status := c.status
	l_status |= 0b0011_0000
	c.push(l_status)
	c.status |= 0b0000_0100
	c.program_counter = c.MemRead16(0xFFFE)
}

func (c *CPU) MemRead(addr uint16) uint8 {
	return c.memory[addr]
}

func (c *CPU) MemWrite(addr uint16, v uint8) {
	c.memory[addr] = v
}

func (c *CPU) MemRead16(addr uint16) uint16 {
	lo := c.memory[addr]
	hi := c.memory[addr+1]
	return make_16_bit(hi, lo)
}

func (c *CPU) MemWrite16(addr uint16, v uint16) {
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

func (c *CPU) set_decimal_bit() {
	c.status |= 0b0000_1000
}

func (c *CPU) set_interrupt_bit() {
	c.status |= 0b0000_0100
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
	if ((res ^ b) & (res ^ a) & 0x80) != 0 {
		c.status |= 0b0100_0000
	} else {
		c.status &= 0b1011_1111
	}
}

func (c *CPU) copy_overflow_flag(v uint8) {
	if (v & 0b0100_0000) > 0 {
		c.status |= 0b0100_0000
	} else {
		c.status &= 0b1011_1111
	}
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
	c.set_negative_flag(reg - val)
}

func (c *CPU) jmp(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	c.program_counter++
	c.program_counter = addr
}

func (c *CPU) jsr(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	ret_addr := c.program_counter
	c.push_16(ret_addr)
	c.program_counter = addr
	/*
		addr := c.mem_read_16(c.program_counter)
		c.push_16(c.program_counter + 2 - 1)
		c.program_counter = addr */
}

func (c *CPU) pha(op OpCode) {
	c.push(c.register_a)
}

func (c *CPU) pla(op OpCode) {
	c.register_a = c.pull()
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) plp(op OpCode) {
	t_status := c.pull()
	t_status &= 0b1110_1111
	t_status |= 0b0010_0000
	c.status = t_status
}

func (c *CPU) hlt(op OpCode) {
}

func (c *CPU) rti(op OpCode) {
	t_status := c.pull()
	t_status &= 0b1110_1111
	t_status |= 0b0010_0000
	c.status = t_status
	lo := c.pull()
	hi := c.pull()
	c.program_counter = make_16_bit(hi, lo)
}

func (c *CPU) rts(op OpCode) {
	lo := c.pull()
	hi := c.pull()
	c.program_counter = make_16_bit(hi, lo) + 1
}

func (c *CPU) tsx(op OpCode) {
	c.register_x = c.stack_pointer
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) txs(op OpCode) {
	c.stack_pointer = c.register_x
}

func (c *CPU) php(op OpCode) {
	push_status := c.status | 0b0011_0000
	c.push(push_status)
}

func (c *CPU) sec(op OpCode) {
	c.set_carry_bit()
}

func (c *CPU) sed(op OpCode) {
	c.set_decimal_bit()
}

func (c *CPU) sei(op OpCode) {
	c.set_interrupt_bit()
}

func (c *CPU) sta(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	c.program_counter++
	c.MemWrite(addr, c.register_a)
}

func (c *CPU) stx(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	c.program_counter++
	c.MemWrite(addr, c.register_x)
}

func (c *CPU) sty(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	c.program_counter++
	c.MemWrite(addr, c.register_y)
}

func (c *CPU) cmp(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.do_compare(val, c.register_a)
}

func (c *CPU) cpx(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.do_compare(val, c.register_x)
}

func (c *CPU) cpy(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true)
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
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if c.is_zero_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) bmi(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if !c.is_negative_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) bvs(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if !c.is_overflow_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) bpl(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if c.is_negative_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) bvc(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if c.is_overflow_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) bit(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.set_zero_flag(val & c.register_a)
	c.copy_overflow_flag(val)
	c.set_negative_flag(val)
}

func (c *CPU) and(op OpCode) {
	c.register_a &= c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) eor(op OpCode) {
	c.register_a ^= c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) ora(op OpCode) {
	c.register_a |= c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) sbc(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true)
	// TODO: Understand what the below actually does properly
	// The carry flag in 6502 is inverted when subtracting:
	// if carry=1 => borrowIn=0, if carry=0 => borrowIn=1.
	// So we take the CPU’s carry bit and flip it.
	old_a := c.register_a
	carry_bit := c.status & 0b00000001
	borrow_in := uint16(1 - carry_bit) // 1 if carry=0, 0 if carry=1

	// Do a 16-bit subtraction so we can detect borrow
	temp := uint16(c.register_a) - uint16(val) - borrow_in
	result := byte(temp & 0xFF) // 8-bit final

	// Store the result back
	c.register_a = result

	// Determine if a borrow occurred by seeing if the 16-bit result is < 0x100
	// If temp < 0x100, it fits in 8 bits => no borrow => set carry
	// If temp >= 0x100 (wrap-around), that means a borrow => clear carry
	if temp < 0x100 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}

	c.compute_overflow_bit(old_a, ^val, result)

	c.set_zero_and_negative_flag(result)

	c.program_counter++
}

func (c *CPU) asl(op OpCode) {
	var pre_val uint8
	var val uint8
	if op.mode == ACCUMULATOR {
		pre_val = c.register_a
		c.register_a <<= 1
		val = c.register_a
	} else {
		var addr uint16
		val = c.interpret_mode(op.mode, &addr, true)
		pre_val = val
		val <<= 1
		c.MemWrite(addr, val)
		c.program_counter++
	}
	// Check if the bit 7 is set and set carry flag if it is
	if pre_val&0b1000_0000 > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) lsr(op OpCode) {
	var pre_val uint8
	var val uint8
	if op.mode == ACCUMULATOR {
		pre_val = c.register_a
		c.register_a >>= 1
		val = c.register_a
	} else {
		var addr uint16
		val = c.interpret_mode(op.mode, &addr, true)
		pre_val = val
		val >>= 1
		c.MemWrite(addr, val)
		c.program_counter++
	}
	// Check if the bit 0 is set and set carry flag if it is
	if pre_val&0b0000_0001 > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) rol(op OpCode) {
	var pre_val uint8
	var val uint8
	if op.mode == ACCUMULATOR {
		pre_val = c.register_a
		val = c.register_a
		val <<= 1
		// Copy carry bit to bit 0
		val = (c.status & 0b0000_0001) | (val & 0b1111_1110)
		c.register_a = val
	} else {
		var addr uint16
		val = c.interpret_mode(op.mode, &addr, true)
		pre_val = val
		val <<= 1
		// Copy carry bit to bit 0
		val = (c.status & 0b0000_0001) | (val & 0b1111_1110)
		c.MemWrite(addr, val)
		c.program_counter++
	}

	// Check if the bit 7 is set and set carry flag if it is
	if pre_val&0b1000_0000 > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) ror(op OpCode) {
	var pre_val uint8
	var val uint8
	var carry_and uint8
	if c.is_carry_set() {
		carry_and = 0b1000_0000
	}
	if op.mode == ACCUMULATOR {
		pre_val = c.register_a
		val = c.register_a
		val >>= 1
		val = carry_and | (val & 0b0111_1111)
		c.register_a = val
	} else {
		var addr uint16
		val = c.interpret_mode(op.mode, &addr, true)
		pre_val = val
		val >>= 1
		val = carry_and | (val & 0b0111_1111)
		c.MemWrite(addr, val)
		c.program_counter++
	}

	// Check if the bit 0 is set and set carry flag if it is
	if pre_val&0b0000_0001 > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) nop(op OpCode) {
}

func (c *CPU) bcc(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if c.is_carry_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) bcs(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if !c.is_carry_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) beq(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true)
	if !c.is_zero_set() {
		c.program_counter++
		return
	}
	c.program_counter = addr
}

func (c *CPU) adc(op OpCode) {
	mem_val := c.interpret_mode(op.mode, nil, true)
	val := c.add_carry_bit(mem_val)
	result := val + c.register_a
	c.decide_carry_bit(result, c.register_a)
	c.compute_overflow_bit(mem_val, c.register_a, result)
	c.register_a = result
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) lda(op OpCode) {

	c.register_a = c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) ldy(op OpCode) {
	c.register_y = c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_y)
}

func (c *CPU) ldx(op OpCode) {
	c.register_x = c.interpret_mode(op.mode, nil, true)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) tax(op OpCode) {
	c.register_x = c.register_a
	c.set_zero_and_negative_flag(c.register_x)
}

func (c *CPU) txa(op OpCode) {
	c.register_a = c.register_x
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) tya(op OpCode) {
	c.register_a = c.register_y
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) tay(op OpCode) {
	c.register_y = c.register_a
	c.set_zero_and_negative_flag(c.register_y)
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
	val := c.interpret_mode(op.mode, &addr, true)
	c.program_counter++
	val--
	c.MemWrite(addr, val)
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) inc(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr, true)
	c.program_counter++
	val++
	c.MemWrite(addr, val)
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

func (c *CPU) interpret_mode(m AddressingMode, read_adr *uint16, incr_pc bool) uint8 {
	var val uint8
	var addr uint16
	var incr_count uint16
	next_val := c.MemRead(c.program_counter)
	switch m {
	case IMMEDIATE:
		val = next_val
	case RELATIVE:
		val = next_val
		addr = c.program_counter + uint16(int16(int8(val))) + 1
	case ZEROPAGE:
		addr = uint16(next_val)
		val = c.MemRead(addr)
	case ZEROPAGEX:
		addr = uint16(next_val + c.register_x)
		val = c.MemRead(addr)
	case ZEROPAGEY:
		addr = uint16(next_val + c.register_y)
		val = c.MemRead(addr)
	case ABSOLUTE:
		addr = c.MemRead16(c.program_counter)
		incr_count++
		val = c.MemRead(addr)
	case ABSOLUTEX:
		in := c.MemRead16(c.program_counter)
		incr_count++
		addr = in + uint16(c.register_x)
		val = c.MemRead(addr)
	case ABSOLUTEY:
		in := c.MemRead16(c.program_counter)
		incr_count++
		addr = in + uint16(c.register_y)
		val = c.MemRead(addr)
	case INDIRECTX:
		in := next_val + c.register_x
		target := c.MemRead16(uint16(in))
		incr_count++
		val = c.MemRead(target)
		addr = target
		incr_count++
	case INDIRECTY:
		in := next_val + c.register_y
		target := c.MemRead16(uint16(in))
		incr_count++
		val = c.MemRead(target)
		addr = target
		incr_count++
	case INDIRECT:
		in := c.MemRead16(c.program_counter)
		target := c.MemRead16(in)
		incr_count++
		val = c.MemRead(target)
		addr = target
		incr_count++
	default:
		panic("Unknown addresing mode")
	}
	if read_adr != nil {
		*read_adr = addr
	}
	if incr_pc {
		c.program_counter += incr_count
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
