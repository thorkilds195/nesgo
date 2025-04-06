package cpu

import (
	"strings"
	"testing"
)

func TestBasicTrace(t *testing.T) {
	vec := []uint8{0x4C, 0xF5, 0xC5}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:00 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
