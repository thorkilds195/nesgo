package cpu

import (
	"strings"
	"testing"
)

func TestBasicTraceAbsoluteInstr(t *testing.T) {
	vec := []uint8{0x4C, 0xF5, 0xC5}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestBasicTraceImmediateInstr(t *testing.T) {
	vec := []uint8{0xA9, 0xF5}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  A9 F5     LDA $F5                         A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestBasicTraceZeroPageInstr(t *testing.T) {
	vec := []uint8{0xA5, 0xF5}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  A5 F5     LDA $F5                         A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestBasicTraceImpliedInstr(t *testing.T) {
	vec := []uint8{0xE8}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  E8        INX                             A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestBasicTraceRelativeInstr(t *testing.T) {
	vec := []uint8{0xB0, 0x04}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  B0 04     BCS $C006                       A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
