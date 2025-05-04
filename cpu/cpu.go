package cpu

import (
	"fmt"
)

type AddressingMode uint8

const STACK_RESET uint8 = 0xFD
const PROGRAM_START uint16 = 0x8000

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

// TODO: Check all the cycles and page crossing logic for opcodes

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
	0x61: {0x61, "ADC", INDIRECTX, 3, 6, (*CPU).adc},
	0x71: {0x71, "ADC", INDIRECTY, 3, 5, (*CPU).adc},
	0x29: {0x29, "AND", IMMEDIATE, 2, 2, (*CPU).and},
	0x25: {0x25, "AND", ZEROPAGE, 2, 3, (*CPU).and},
	0x35: {0x35, "AND", ZEROPAGEX, 2, 4, (*CPU).and},
	0x2D: {0x2D, "AND", ABSOLUTE, 3, 4, (*CPU).and},
	0x3D: {0x3D, "AND", ABSOLUTEX, 3, 4, (*CPU).and},
	0x39: {0x39, "AND", ABSOLUTEY, 3, 4, (*CPU).and},
	0x21: {0x21, "AND", INDIRECTX, 3, 6, (*CPU).and},
	0x31: {0x31, "AND", INDIRECTY, 3, 5, (*CPU).and},
	0x0A: {0x0A, "ASL", ACCUMULATOR, 1, 2, (*CPU).asl},
	0x06: {0x06, "ASL", ZEROPAGE, 2, 5, (*CPU).asl},
	0x16: {0x16, "ASL", ZEROPAGEX, 2, 6, (*CPU).asl},
	0x0E: {0x0E, "ASL", ABSOLUTE, 2, 6, (*CPU).asl},
	0x1E: {0x1E, "ASL", ABSOLUTEX, 2, 7, (*CPU).asl},
	0x90: {0x90, "BCC", RELATIVE, 2, 2, (*CPU).bcc}, // plus 1 if branch succeeds, plus 2 if new page
	0xB0: {0xB0, "BCS", RELATIVE, 2, 2, (*CPU).bcs}, // plus 1 if branch succeeds, plus 2 if new page
	0xF0: {0xF0, "BEQ", RELATIVE, 2, 2, (*CPU).beq}, // plus 1 if branch succeeds, plus 2 if new page
	0x24: {0x24, "BIT", ZEROPAGE, 2, 3, (*CPU).bit},
	0x2C: {0x2C, "BIT", ABSOLUTE, 2, 4, (*CPU).bit},
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
	0xC5: {0xC5, "CMP", ZEROPAGE, 2, 3, (*CPU).cmp},
	0xD5: {0xD5, "CMP", ZEROPAGEX, 2, 4, (*CPU).cmp},
	0xCD: {0xCD, "CMP", ABSOLUTE, 2, 4, (*CPU).cmp},
	0xDD: {0xDD, "CMP", ABSOLUTEX, 2, 4, (*CPU).cmp},
	0xD9: {0xD9, "CMP", ABSOLUTEY, 2, 4, (*CPU).cmp},
	0xC1: {0xC1, "CMP", INDIRECTX, 2, 6, (*CPU).cmp},
	0xD1: {0xD1, "CMP", INDIRECTY, 2, 5, (*CPU).cmp},
	0xE0: {0xE0, "CPX", IMMEDIATE, 2, 2, (*CPU).cpx},
	0xE4: {0xE4, "CPX", ZEROPAGE, 2, 3, (*CPU).cpx},
	0xEC: {0xEC, "CPX", ABSOLUTE, 2, 4, (*CPU).cpx},
	0xC0: {0xC0, "CPY", IMMEDIATE, 2, 2, (*CPU).cpy},
	0xC4: {0xC4, "CPY", ZEROPAGE, 2, 3, (*CPU).cpy},
	0xCC: {0xCC, "CPY", ABSOLUTE, 2, 4, (*CPU).cpy},
	0xC6: {0xC6, "DEC", ZEROPAGE, 2, 5, (*CPU).dec},
	0xD6: {0xD6, "DEC", ZEROPAGEX, 2, 6, (*CPU).dec},
	0xCE: {0xCE, "DEC", ABSOLUTE, 2, 6, (*CPU).dec},
	0xDE: {0xDE, "DEC", ABSOLUTEX, 2, 7, (*CPU).dec},
	0xCA: {0xCA, "DEX", IMPLIED, 2, 2, (*CPU).dex},
	0x88: {0x88, "DEY", IMPLIED, 2, 2, (*CPU).dey},
	0x49: {0x49, "EOR", IMMEDIATE, 2, 2, (*CPU).eor},
	0x45: {0x45, "EOR", ZEROPAGE, 2, 3, (*CPU).eor},
	0x55: {0x55, "EOR", ZEROPAGEX, 2, 4, (*CPU).eor},
	0x4D: {0x4D, "EOR", ABSOLUTE, 2, 4, (*CPU).eor},
	0x5D: {0x5D, "EOR", ABSOLUTEX, 2, 4, (*CPU).eor},
	0x59: {0x59, "EOR", ABSOLUTEY, 2, 4, (*CPU).eor},
	0x41: {0x41, "EOR", INDIRECTX, 2, 6, (*CPU).eor},
	0x51: {0x51, "EOR", INDIRECTY, 2, 5, (*CPU).eor},
	0xE6: {0xE6, "INC", ZEROPAGE, 2, 5, (*CPU).inc},
	0xF6: {0xF6, "INC", ZEROPAGEX, 2, 6, (*CPU).inc},
	0xEE: {0xEE, "INC", ABSOLUTE, 2, 6, (*CPU).inc},
	0xFE: {0xFE, "INC", ABSOLUTEX, 2, 7, (*CPU).inc},
	0xC8: {0xC8, "INY", IMPLIED, 2, 2, (*CPU).iny},
	0x4C: {0x4C, "JMP", ABSOLUTE, 3, 3, (*CPU).jmp},
	0x6C: {0x6C, "JMP", INDIRECT, 3, 5, (*CPU).jmp},
	0x4A: {0x4A, "LSR", ACCUMULATOR, 1, 2, (*CPU).lsr},
	0x46: {0x46, "LSR", ZEROPAGE, 2, 5, (*CPU).lsr},
	0x56: {0x56, "LSR", ZEROPAGEX, 2, 6, (*CPU).lsr},
	0x4E: {0x4E, "LSR", ABSOLUTE, 2, 6, (*CPU).lsr},
	0x5E: {0x5E, "LSR", ABSOLUTEX, 2, 7, (*CPU).lsr},
	0xEA: {0xEA, "NOP", IMPLIED, 2, 2, (*CPU).nop},
	0x09: {0x09, "ORA", IMMEDIATE, 2, 2, (*CPU).ora},
	0x05: {0x05, "ORA", ZEROPAGE, 2, 3, (*CPU).ora},
	0x15: {0x15, "ORA", ZEROPAGEX, 2, 4, (*CPU).ora},
	0x0D: {0x0D, "ORA", ABSOLUTE, 2, 4, (*CPU).ora},
	0x1D: {0x1D, "ORA", ABSOLUTEX, 2, 4, (*CPU).ora},
	0x19: {0x19, "ORA", ABSOLUTEY, 2, 4, (*CPU).ora},
	0x01: {0x01, "ORA", INDIRECTX, 2, 6, (*CPU).ora},
	0x11: {0x11, "ORA", INDIRECTY, 2, 5, (*CPU).ora},
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
	0xE5: {0xE5, "SBC", ZEROPAGE, 2, 3, (*CPU).sbc},
	0xF5: {0xF5, "SBC", ZEROPAGEX, 2, 4, (*CPU).sbc},
	0xED: {0xED, "SBC", ABSOLUTE, 2, 4, (*CPU).sbc},
	0xFD: {0xFD, "SBC", ABSOLUTEX, 2, 4, (*CPU).sbc},
	0xF9: {0xF9, "SBC", ABSOLUTEY, 2, 4, (*CPU).sbc},
	0xE1: {0xE1, "SBC", INDIRECTX, 2, 6, (*CPU).sbc},
	0xF1: {0xF1, "SBC", INDIRECTY, 2, 5, (*CPU).sbc},
	0x38: {0x38, "SEC", IMPLIED, 1, 2, (*CPU).sec},
	0xF8: {0xF8, "SED", IMPLIED, 2, 2, (*CPU).sed},
	0x78: {0x78, "SEI", IMPLIED, 2, 2, (*CPU).sei},
	0x85: {0x85, "STA", ZEROPAGE, 2, 3, (*CPU).sta},
	0x95: {0x95, "STA", ZEROPAGEX, 2, 4, (*CPU).sta},
	0x8D: {0x8D, "STA", ABSOLUTE, 2, 4, (*CPU).sta},
	0x9D: {0x9D, "STA", ABSOLUTEX, 2, 5, (*CPU).sta},
	0x99: {0x99, "STA", ABSOLUTEY, 2, 5, (*CPU).sta},
	0x81: {0x81, "STA", INDIRECTX, 2, 6, (*CPU).sta},
	0x91: {0x91, "STA", INDIRECTY, 2, 6, (*CPU).sta},
	0x86: {0x86, "STX", ZEROPAGE, 2, 3, (*CPU).stx},
	0x96: {0x96, "STX", ZEROPAGEY, 2, 4, (*CPU).stx},
	0x8E: {0x8E, "STX", ABSOLUTE, 2, 4, (*CPU).stx},
	0x84: {0x84, "STY", ZEROPAGE, 2, 3, (*CPU).sty},
	0x94: {0x94, "STY", ZEROPAGEX, 2, 4, (*CPU).sty},
	0x8C: {0x8C, "STY", ABSOLUTE, 2, 4, (*CPU).sty},
	0xA8: {0xA8, "TAY", IMPLIED, 2, 2, (*CPU).tay},
	0x8A: {0x8A, "TXA", IMPLIED, 2, 2, (*CPU).txa},
	0x98: {0x98, "TYA", IMPLIED, 2, 2, (*CPU).tya},
	0x20: {0x20, "JSR", ABSOLUTE, 2, 6, (*CPU).jsr},
	0x48: {0x48, "PHA", IMPLIED, 2, 3, (*CPU).pha},
	0x08: {0x08, "PHP", IMPLIED, 2, 3, (*CPU).php},
	0x68: {0x68, "PLA", IMPLIED, 2, 4, (*CPU).pla},
	0x28: {0x28, "PLP", IMPLIED, 2, 4, (*CPU).plp},
	0x40: {0x40, "RTI", IMPLIED, 2, 6, (*CPU).rti},
	0x60: {0x60, "RTS", IMPLIED, 2, 6, (*CPU).rts},
	0xBA: {0xBA, "TSX", IMPLIED, 2, 2, (*CPU).tsx},
	0x9A: {0x9A, "TXS", IMPLIED, 2, 2, (*CPU).txs},
	// Unofficial opcodes
	0x1A: {0x1A, "*NOP", IMPLIED, 2, 2, (*CPU).nop},
	0x3A: {0x3A, "*NOP", IMPLIED, 2, 2, (*CPU).nop},
	0x5A: {0x5A, "*NOP", IMPLIED, 2, 2, (*CPU).nop},
	0x7A: {0x7A, "*NOP", IMPLIED, 2, 2, (*CPU).nop},
	0xDA: {0xDA, "*NOP", IMPLIED, 2, 2, (*CPU).nop},
	0xFA: {0xFA, "*NOP", IMPLIED, 2, 2, (*CPU).nop},
	0x80: {0x80, "*NOP", IMMEDIATE, 2, 2, (*CPU).skb},
	0x82: {0x82, "*NOP", IMMEDIATE, 2, 2, (*CPU).skb},
	0x89: {0x89, "*NOP", IMMEDIATE, 2, 2, (*CPU).skb},
	0xC2: {0xC2, "*NOP", IMMEDIATE, 2, 2, (*CPU).skb},
	0xE2: {0xE2, "*NOP", IMMEDIATE, 2, 2, (*CPU).skb},
	0x0C: {0x0C, "*NOP", ABSOLUTE, 2, 4, (*CPU).ign},
	0x1C: {0x1C, "*NOP", ABSOLUTEX, 2, 4, (*CPU).ign},
	0x3C: {0x3C, "*NOP", ABSOLUTEX, 2, 4, (*CPU).ign},
	0x5C: {0x5C, "*NOP", ABSOLUTEX, 2, 4, (*CPU).ign},
	0x7C: {0x7C, "*NOP", ABSOLUTEX, 2, 4, (*CPU).ign},
	0xDC: {0xDC, "*NOP", ABSOLUTEX, 2, 4, (*CPU).ign},
	0xFC: {0xFC, "*NOP", ABSOLUTEX, 2, 4, (*CPU).ign},
	0x04: {0x04, "*NOP", ZEROPAGE, 2, 3, (*CPU).ign},
	0x44: {0x44, "*NOP", ZEROPAGE, 2, 3, (*CPU).ign},
	0x64: {0x64, "*NOP", ZEROPAGE, 2, 3, (*CPU).ign},
	0x14: {0x14, "*NOP", ZEROPAGEX, 2, 4, (*CPU).ign},
	0x34: {0x34, "*NOP", ZEROPAGEX, 2, 4, (*CPU).ign},
	0x54: {0x54, "*NOP", ZEROPAGEX, 2, 4, (*CPU).ign},
	0x74: {0x74, "*NOP", ZEROPAGEX, 2, 4, (*CPU).ign},
	0xD4: {0xD4, "*NOP", ZEROPAGEX, 2, 4, (*CPU).ign},
	0xF4: {0xF4, "*NOP", ZEROPAGEX, 2, 4, (*CPU).ign},
	0xA3: {0xA3, "*LAX", INDIRECTX, 2, 6, (*CPU).lax},
	0xA7: {0xA7, "*LAX", ZEROPAGE, 2, 3, (*CPU).lax},
	0xAF: {0xAF, "*LAX", ABSOLUTE, 2, 4, (*CPU).lax},
	0xB3: {0xB3, "*LAX", INDIRECTY, 2, 5, (*CPU).lax},
	0xB7: {0xB7, "*LAX", ZEROPAGEY, 2, 4, (*CPU).lax},
	0xBF: {0xBF, "*LAX", ABSOLUTEY, 2, 4, (*CPU).lax},
	0x83: {0x83, "*SAX", INDIRECTX, 2, 6, (*CPU).sax},
	0x87: {0x87, "*SAX", ZEROPAGE, 2, 3, (*CPU).sax},
	0x8F: {0x8F, "*SAX", ABSOLUTE, 2, 4, (*CPU).sax},
	0x97: {0x97, "*SAX", ZEROPAGEY, 2, 4, (*CPU).sax},
	0x9E: {0x9E, "*SHX", ABSOLUTEY, 2, 7, (*CPU).shx},
	0x9C: {0x9C, "*SHY", ABSOLUTEX, 2, 7, (*CPU).shy},
	0x93: {0x93, "*SHA", INDIRECTY, 2, 6, (*CPU).sha},
	0x9F: {0x9F, "*SHA", ABSOLUTEY, 2, 5, (*CPU).sha},
	0xEB: {0xEB, "*SBC", IMMEDIATE, 2, 2, (*CPU).sbc},
	0xC3: {0xC3, "*DCP", INDIRECTX, 2, 8, (*CPU).dcp},
	0xC7: {0xC7, "*DCP", ZEROPAGE, 2, 5, (*CPU).dcp},
	0xCF: {0xCF, "*DCP", ABSOLUTE, 2, 6, (*CPU).dcp},
	0xD3: {0xD3, "*DCP", INDIRECTY, 2, 8, (*CPU).dcp},
	0xD7: {0xD7, "*DCP", ZEROPAGEX, 2, 6, (*CPU).dcp},
	0xDB: {0xDB, "*DCP", ABSOLUTEY, 2, 7, (*CPU).dcp},
	0xDF: {0xDF, "*DCP", ABSOLUTEX, 2, 7, (*CPU).dcp},
	0xE3: {0xE3, "*ISB", INDIRECTX, 2, 8, (*CPU).isc},
	0xE7: {0xE7, "*ISB", ZEROPAGE, 2, 5, (*CPU).isc},
	0xEF: {0xEF, "*ISB", ABSOLUTE, 2, 6, (*CPU).isc},
	0xF3: {0xF3, "*ISB", INDIRECTY, 2, 8, (*CPU).isc},
	0xF7: {0xF7, "*ISB", ZEROPAGEX, 2, 6, (*CPU).isc},
	0xFB: {0xFB, "*ISB", ABSOLUTEY, 2, 7, (*CPU).isc},
	0xFF: {0xFF, "*ISB", ABSOLUTEX, 2, 7, (*CPU).isc},
	0x23: {0x23, "*RLA", INDIRECTX, 2, 8, (*CPU).rla},
	0x27: {0x27, "*RLA", ZEROPAGE, 2, 5, (*CPU).rla},
	0x2F: {0x2F, "*RLA", ABSOLUTE, 2, 6, (*CPU).rla},
	0x33: {0x33, "*RLA", INDIRECTY, 2, 8, (*CPU).rla},
	0x37: {0x37, "*RLA", ZEROPAGEX, 2, 6, (*CPU).rla},
	0x3B: {0x3B, "*RLA", ABSOLUTEY, 2, 7, (*CPU).rla},
	0x3F: {0x3F, "*RLA", ABSOLUTEX, 2, 7, (*CPU).rla},
	0x03: {0x03, "*SLO", INDIRECTX, 2, 8, (*CPU).slo},
	0x07: {0x07, "*SLO", ZEROPAGE, 2, 5, (*CPU).slo},
	0x0F: {0x0F, "*SLO", ABSOLUTE, 2, 6, (*CPU).slo},
	0x13: {0x13, "*SLO", INDIRECTY, 2, 8, (*CPU).slo},
	0x17: {0x17, "*SLO", ZEROPAGEX, 2, 6, (*CPU).slo},
	0x1B: {0x1B, "*SLO", ABSOLUTEY, 2, 7, (*CPU).slo},
	0x1F: {0x1F, "*SLO", ABSOLUTEX, 2, 7, (*CPU).slo},
	0x43: {0x43, "*SRE", INDIRECTX, 2, 8, (*CPU).sre},
	0x47: {0x47, "*SRE", ZEROPAGE, 2, 5, (*CPU).sre},
	0x4F: {0x4F, "*SRE", ABSOLUTE, 2, 6, (*CPU).sre},
	0x53: {0x53, "*SRE", INDIRECTY, 2, 8, (*CPU).sre},
	0x57: {0x57, "*SRE", ZEROPAGEX, 2, 6, (*CPU).sre},
	0x5B: {0x5B, "*SRE", ABSOLUTEY, 2, 7, (*CPU).sre},
	0x5F: {0x5F, "*SRE", ABSOLUTEX, 2, 7, (*CPU).sre},
	0x63: {0x63, "*RRA", INDIRECTX, 2, 8, (*CPU).rra},
	0x67: {0x67, "*RRA", ZEROPAGE, 2, 5, (*CPU).rra},
	0x6F: {0x6F, "*RRA", ABSOLUTE, 2, 6, (*CPU).rra},
	0x73: {0x73, "*RRA", INDIRECTY, 2, 8, (*CPU).rra},
	0x77: {0x77, "*RRA", ZEROPAGEX, 2, 6, (*CPU).rra},
	0x7B: {0x7B, "*RRA", ABSOLUTEY, 2, 7, (*CPU).rra},
	0x7F: {0x7F, "*RRA", ABSOLUTEX, 2, 7, (*CPU).rra},
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
	Bus             *Bus
}

func (c *CPU) GetCycles() uint {
	return c.Bus.cycles
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

func InitCPU(b *Bus) *CPU {
	return &CPU{Bus: b, stack_pointer: STACK_RESET, program_counter: 0x8000, status: 0b100100}
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
		if c.Bus.PollNMIStatus() != nil {
			c.interrupt_nmi()
		}
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
		c.Bus.Tick(op.cycles)
	}
}

func (c *CPU) Run() {
	c.RunWithCallback(func() {})
}

func (c *CPU) Load(program []uint8) {
	for i := range len(program) {
		c.MemWrite(0x0000+uint16(i), program[i])
	}
	//copy(c.memory[PROGRAM_START:], program)

	c.MemWrite16(0xFFFC, PROGRAM_START)
}

func (c *CPU) Step(f_call func()) bool {
	if c.Bus.PollNMIStatus() != nil {
		c.interrupt_nmi()
	}
	f_call()
	opcode := c.MemRead(c.program_counter)
	c.program_counter++
	op, ok := OPTABLE[opcode]
	// fmt.Println(op.name, op.mode)
	if !ok {
		panic(fmt.Sprintf("Unknown opcode: %x", opcode))
	}
	op.f_call(c, op)
	c.Bus.Tick(op.cycles)
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
	c.interpret_mode(op.mode, &addr, false, nil, true)
	c.program_counter--
	return op, addr
}

func (c *CPU) interrupt_nmi() {
	c.push_16(c.program_counter)
	status := c.status | 0b0011_0000
	c.push(status)
	c.status |= 0b0000_0100
	c.program_counter = c.MemRead16(0xFFFA)
	c.Bus.Tick(2)
}

func (c *CPU) push(val uint8) {
	c.Bus.MemWrite(0x0100+uint16(c.stack_pointer), val)
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
	c.doInterrupt(0xFFFE)
}

func (c *CPU) doInterrupt(rd_addr uint16) {
	c.push_16(c.program_counter + 2)
	l_status := c.status
	l_status |= 0b0011_0000
	c.push(l_status)
	c.status |= 0b0000_0100
	c.program_counter = c.MemRead16(rd_addr)
}

func (c *CPU) MemRead(addr uint16) uint8 {
	return c.Bus.MemRead(addr)
}

func (c *CPU) MemWrite(addr uint16, v uint8) {
	c.Bus.MemWrite(addr, v)
}

func (c *CPU) MemRead16(addr uint16) uint16 {
	lo := c.Bus.MemRead(addr)
	hi := c.Bus.MemRead(addr + 1)
	return make_16_bit(hi, lo)
}

func (c *CPU) mem_read_16_zero(addr uint8) uint16 {
	lo := c.MemRead(uint16(addr))
	hi := c.MemRead(uint16((addr + 1) & 0xFF))
	return make_16_bit(hi, lo)
}
func (c *CPU) MemWrite16(addr uint16, v uint16) {
	lo := uint8(v & 0xFF)
	hi := uint8(v >> 8)
	c.Bus.MemWrite(addr, lo)
	c.Bus.MemWrite(addr+1, hi)
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
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	c.program_counter = addr
}

func (c *CPU) jsr(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
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
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	c.MemWrite(addr, c.register_a)
}

func (c *CPU) stx(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	c.MemWrite(addr, c.register_x)
}

func (c *CPU) sty(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	c.MemWrite(addr, c.register_y)
}

func (c *CPU) cmp(op OpCode) {
	crossed := false
	val := c.interpret_mode(op.mode, nil, true, &crossed, false)
	c.program_counter++
	c.do_compare(val, c.register_a)
	if crossed {
		c.Bus.Tick(1)
	}
}

func (c *CPU) cpx(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true, nil, false)
	c.program_counter++
	c.do_compare(val, c.register_x)
}

func (c *CPU) cpy(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true, nil, false)
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
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if c.is_zero_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) will_pg_cross(addr uint16) bool {
	return (c.program_counter & 0xFF00) != (addr & 0xFF00)
}

func (c *CPU) bmi(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if !c.is_negative_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) bvs(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if !c.is_overflow_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) bpl(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if c.is_negative_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) bvc(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if c.is_overflow_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) bit(op OpCode) {
	val := c.interpret_mode(op.mode, nil, true, nil, false)
	c.program_counter++
	c.set_zero_flag(val & c.register_a)
	c.copy_overflow_flag(val)
	c.set_negative_flag(val)
}

func (c *CPU) and(op OpCode) {
	crossed := false
	c.do_and(c.interpret_mode(op.mode, nil, true, &crossed, false))
	if crossed {
		c.Bus.Tick(1)
	}
}

func (c *CPU) do_and(val uint8) {
	c.register_a &= val
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) eor(op OpCode) {
	crossed := false
	c.register_a ^= c.interpret_mode(op.mode, nil, true, &crossed, false)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
	if crossed {
		c.Bus.Tick(1)
	}
}

func (c *CPU) ora(op OpCode) {
	crossed := false
	c.register_a |= c.interpret_mode(op.mode, nil, true, &crossed, false)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
	if crossed {
		c.Bus.Tick(1)
	}
}

func (c *CPU) sbc(op OpCode) {
	crossed := false
	val := c.interpret_mode(op.mode, nil, true, &crossed, false)
	c.do_sbc(val)
	if crossed {
		c.Bus.Tick(1)
	}
}
func (c *CPU) do_sbc(val uint8) {
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
		val = c.interpret_mode(op.mode, &addr, true, nil, false)
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
		val = c.interpret_mode(op.mode, &addr, true, nil, false)
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
		val = c.interpret_mode(op.mode, &addr, true, nil, false)
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
		val = c.interpret_mode(op.mode, &addr, true, nil, false)
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

func (c *CPU) skb(op OpCode) {
	c.interpret_mode(op.mode, nil, true, nil, true)
	c.program_counter++
}

func (c *CPU) ign(op OpCode) {
	crossed := false
	c.interpret_mode(op.mode, nil, true, &crossed, true)
	c.program_counter++
	if crossed {
		c.Bus.Tick(1)
	}
}

func (c *CPU) lax(op OpCode) {
	c.lda(op)
	c.tax(op)
}

func (c *CPU) dcp(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr, true, nil, false)
	c.program_counter++
	val--
	c.MemWrite(addr, val)
	c.set_zero_and_negative_flag(val)
	c.do_compare(val, c.register_a)
}

func (c *CPU) isc(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr, true, nil, false)
	val++
	c.MemWrite(addr, val)
	c.do_sbc(val)
}

func (c *CPU) sax(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	c.MemWrite(addr, c.register_a&c.register_x)
}

func (c *CPU) slo(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr, true, nil, false)
	pre_val := val
	val <<= 1
	c.MemWrite(addr, val)
	if pre_val&0b1000_0000 > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	c.register_a |= val
	c.set_zero_and_negative_flag(c.register_a)
	c.program_counter++
}

func (c *CPU) sre(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr, true, nil, false)
	pre_val := val
	val >>= 1
	c.MemWrite(addr, val)
	c.register_a ^= val
	c.set_zero_and_negative_flag(c.register_a)
	if pre_val&0b0000_0001 > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	c.program_counter++
}

func (c *CPU) rra(op OpCode) {
	var addr uint16
	var carry_and uint8
	if c.is_carry_set() {
		carry_and = 0b1000_0000
	}
	val := c.interpret_mode(op.mode, &addr, true, nil, false)
	pre_val := val
	val >>= 1
	val = carry_and | (val & 0b0111_1111)
	c.MemWrite(addr, val)
	if pre_val&0b0000_0001 > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	val = c.add_carry_bit(val)
	result := val + c.register_a
	c.decide_carry_bit(result, c.register_a)
	c.compute_overflow_bit(val, c.register_a, result)
	c.register_a = result
	c.set_zero_and_negative_flag(c.register_a)
	c.program_counter++
}

func (c *CPU) rla(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr, true, nil, false)

	// ROL
	oldCarry := (c.status >> 0) & 1
	carry := (val >> 7) & 1
	val = ((val << 1) & 0xFF) | oldCarry
	if carry > 0 {
		c.set_carry_bit()
	} else {
		c.clear_carry_bit()
	}
	c.set_zero_and_negative_flag(val)
	c.MemWrite(addr, val)

	// AND
	c.register_a = c.register_a & val
	c.set_zero_and_negative_flag(c.register_a)
	c.program_counter++
}

func (c *CPU) shx(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	hi_b := uint8(addr >> 8)
	wr_data := c.register_x & (hi_b + 1)
	c.program_counter++
	c.MemWrite(addr, wr_data)
}

func (c *CPU) shy(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	hi_b := uint8(addr >> 8)
	wr_data := c.register_y & (hi_b + 1)
	c.program_counter++
	c.MemWrite(addr, wr_data)
}

func (c *CPU) sha(op OpCode) {
}

func (c *CPU) bcc(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if c.is_carry_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) bcs(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if !c.is_carry_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) beq(op OpCode) {
	var addr uint16
	c.interpret_mode(op.mode, &addr, true, nil, true)
	c.program_counter++
	if !c.is_zero_set() {
		return
	}
	c.Bus.Tick(1)
	if c.will_pg_cross(addr) {
		c.Bus.Tick(1)
	}
	c.program_counter = addr
}

func (c *CPU) adc(op OpCode) {
	mem_val := c.interpret_mode(op.mode, nil, true, nil, false)
	val := c.add_carry_bit(mem_val)
	result := val + c.register_a
	c.decide_carry_bit(result, c.register_a)
	c.compute_overflow_bit(mem_val, c.register_a, result)
	c.register_a = result
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
}

func (c *CPU) lda(op OpCode) {
	crossed := false
	c.register_a = c.interpret_mode(op.mode, nil, true, &crossed, false)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_a)
	if crossed {
		c.Bus.Tick(1)
	}
}

func (c *CPU) ldy(op OpCode) {
	crossed := false
	c.register_y = c.interpret_mode(op.mode, nil, true, &crossed, false)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_y)
	if crossed {
		c.Bus.Tick(1)
	}
}

func (c *CPU) ldx(op OpCode) {
	crossed := false
	c.register_x = c.interpret_mode(op.mode, nil, true, &crossed, false)
	c.program_counter++
	c.set_zero_and_negative_flag(c.register_x)
	if crossed {
		c.Bus.Tick(1)
	}
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
	val := c.interpret_mode(op.mode, &addr, true, nil, false)
	c.program_counter++
	val--
	c.MemWrite(addr, val)
	c.set_zero_and_negative_flag(val)
}

func (c *CPU) inc(op OpCode) {
	var addr uint16
	val := c.interpret_mode(op.mode, &addr, true, nil, false)
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

func (c *CPU) interpret_mode(m AddressingMode, read_adr *uint16, incr_pc bool, did_cross *bool, no_read bool) uint8 {
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
		if !no_read {
			val = c.MemRead(addr)
		}
	case ZEROPAGEX:
		addr = uint16(next_val + c.register_x)
		if !no_read {
			val = c.MemRead(addr)
		}
	case ZEROPAGEY:
		addr = uint16(next_val + c.register_y)
		if !no_read {
			val = c.MemRead(addr)
		}
	case ABSOLUTE:
		addr = c.MemRead16(c.program_counter)
		incr_count++
		if !no_read {
			val = c.MemRead(addr)
		}
	case ABSOLUTEX:
		in := c.MemRead16(c.program_counter)
		incr_count++
		addr = in + uint16(c.register_x)
		if did_cross != nil && in&0xFF00 != addr&0xFF00 {
			*did_cross = true
		}
		if !no_read {
			val = c.MemRead(addr)
		}
	case ABSOLUTEY:
		in := c.MemRead16(c.program_counter)
		incr_count++
		addr = in + uint16(c.register_y)
		if did_cross != nil && in&0xFF00 != addr&0xFF00 {
			*did_cross = true
		}
		if !no_read {
			val = c.MemRead(addr)
		}
	case INDIRECTX:
		in := next_val + c.register_x
		target := c.mem_read_16_zero(in)
		if !no_read {
			val = c.MemRead(target)
		}
		addr = target
	case INDIRECTY:
		target := c.mem_read_16_zero(next_val)
		final_target := target + uint16(c.register_y)
		if !no_read {
			val = c.MemRead(final_target)
		}
		addr = final_target
		if did_cross != nil && target&0xFF00 != final_target&0xFF00 {
			*did_cross = true
		}
	case INDIRECT:
		ptr := c.MemRead16(c.program_counter)
		lo := c.MemRead(ptr)
		var hi uint8
		if ptr&0x00FF == 0x00FF {
			// Bug: wrap within the same page
			hi = c.MemRead(ptr & 0xFF00)
		} else {
			hi = c.MemRead(ptr + 1)
		}
		target := make_16_bit(hi, lo)
		val = 0 // JMP does not load a value
		addr = target
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
