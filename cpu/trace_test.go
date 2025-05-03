package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTraceAbsoluteInstr(t *testing.T) {
	vec := []uint8{0x4C, 0xF5, 0xC5}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceAccumulatorInstr(t *testing.T) {
	vec := []uint8{0x0A}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  0A        ASL A                           A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceImmediateInstr(t *testing.T) {
	vec := []uint8{0xA9, 0xF5}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  A9 F5     LDA #$F5                        A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceZeroPageInstr(t *testing.T) {
	vec := []uint8{0xA5, 0xF5}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  A5 F5     LDA $F5 = 00                    A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceZeroPageXInstr(t *testing.T) {
	vec := []uint8{0xB5, 0x33}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  B5 33     LDA $33,X @ 33 = 00             A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceZeroPageYInstr(t *testing.T) {
	vec := []uint8{0xB6, 0x33}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  B6 33     LDX $33,Y @ 33 = 00             A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceAbsolutInstr(t *testing.T) {
	vec := []uint8{0xAD, 0x47, 0x06}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  AD 47 06  LDA $0647 = 00                  A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceAbsoluteYInstr(t *testing.T) {
	vec := []uint8{0xB9, 0x00, 0x03}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  B9 00 03  LDA $0300,Y @ 0300 = 00         A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}
func TestTraceAbsoluteXInstr(t *testing.T) {
	vec := []uint8{0xBD, 0x00, 0x03}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  BD 00 03  LDA $0300,X @ 0300 = 00         A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestBasicTraceImpliedInstr(t *testing.T) {
	vec := []uint8{0xE8}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  E8        INX                             A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestBasicTraceRelativeInstr(t *testing.T) {
	vec := []uint8{0xB0, 0x04}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  B0 04     BCS $8006                       A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceIndirectX(t *testing.T) {
	vec := []uint8{0xA1, 0x80}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  A1 80     LDA ($80,X) @ 80 = 0000 = 00    A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceIndirectY(t *testing.T) {
	vec := []uint8{0xB1, 0x89}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  B1 89     LDA ($89),Y = 0000 @ 0000 = 00  A:00 X:00 Y:00 P:24 SP:FD CYC:0"
	if !(strings.EqualFold(actual, expected)) {
		t.Errorf("TraceCPU not returning correct output\nGot:\n%s\nExpected:\n%s", actual, expected)
	}
}

func TestTraceIndirect(t *testing.T) {
	vec := []uint8{0x6C, 0x00, 0x02}
	c := InitCPU(setupTestBus(vec))
	c.Reset()
	actual := TraceCPU(c)
	expected := "8000  6C 00 02  JMP ($0200) = 0000              A:00 X:00 Y:00 P:24 SP:FD CYC:0"
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
	r := InitRom(dat)
	b := InitBus(r, func(*PPU) {})
	c := InitCPU(b)
	c.Reset()
	// Hotfix for now to match reset cycle
	b.cycles = 7
	c.program_counter = 0xC000
	idx := 0
	c.RunWithCallback(func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Panic occurred at log line %d: %v\n", idx, r)
				panic(r)
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
	return s[:cutoff_idx] + s[cutoff_idx+12:]
}
