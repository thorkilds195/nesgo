package cpu

import (
	"bufio"
	"fmt"
	"os"
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
	expected := "C000  A9 F5     LDA #$F5                        A:00 X:00 Y:00 P:24 SP:FD"
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
	expected := "C000  A5 F5     LDA $F5 = 00                    A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceZeroPageXInstr(t *testing.T) {
	vec := []uint8{0xB5, 0x33}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  B5 33     LDA $33,X @ 33 = 00             A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceZeroPageYInstr(t *testing.T) {
	vec := []uint8{0xB6, 0x33}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  B6 33     LDX $33,Y @ 33 = 00             A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceAbsolutInstr(t *testing.T) {
	vec := []uint8{0xAD, 0x47, 0x06}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  AD 47 06  LDA $0647 = 00                  A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceAbsoluteYInstr(t *testing.T) {
	vec := []uint8{0xB9, 0x00, 0x03}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  B9 00 03  LDA $0300,Y @ 0300 = 00         A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceAbsoluteXInstr(t *testing.T) {
	vec := []uint8{0xBD, 0x00, 0x03}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  BD 00 03  LDA $0300,X @ 0300 = 00         A:00 X:00 Y:00 P:24 SP:FD"
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

func TestTraceIndirectX(t *testing.T) {
	vec := []uint8{0xA1, 0x80}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  A1 80     LDA ($80,X) @ 80 = 0000 = 00    A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceIndirectY(t *testing.T) {
	vec := []uint8{0xB1, 0x89}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  B1 89     LDA ($89),Y = 0000 @ 0000 = 00  A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceIndirect(t *testing.T) {
	vec := []uint8{0x6C, 0x00, 0x02}
	c := InitCPU()
	c.Load(vec)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	actual := TraceCPU(c)
	expected := "C000  6C 00 02  JMP ($0200) = 0000              A:00 X:00 Y:00 P:24 SP:FD"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestCompareAgainstNesLog(t *testing.T) {
	dat, err := os.ReadFile("../nestest.nes")
	if err != nil {
		panic(err)
	}
	file, err := os.Open("../nestest.log.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var answer []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		answer = append(answer, scanner.Text())
	}
	c := InitCPU()
	c.Load(dat)
	c.MemWrite16(0xFFFC, 0xC000)
	c.Reset()
	idx := 0
	c.RunWithCallback(func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic occurred at log line %d: %v\n", idx, r)
				panic(r) // Re-throw panic so test fails
			}
		}()
		actual := TraceCPU(c)
		expected := formatAns(answer[idx])
		if !(expected == actual) {
			t.Fatalf("Log comparison error on line %d\nGot:\n%s\nExpected:\n%s", idx, actual, expected)
		}
		idx++
	})

}

func formatAns(s string) string {
	cutoff_idx := strings.Index(s, " PPU")
	return s[:cutoff_idx]
}
