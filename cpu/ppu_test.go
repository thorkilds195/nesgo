package cpu

import "testing"

func setupTestPPU(m Mirroring) *PPU {
	chrrom := make([]uint8, 0x2000)
	return NewPPU(chrrom, m)
}

func TestWritesValueToPPUAddr(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToPPUAddr(0x21)
	p.WriteToPPUAddr(0x15)
	if !(p.addr.get() == 0x2115) {
		t.Error("PPU addr not correct")
	}
}
func TestWritesValueToPPUVramInVertical(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToPPUAddr(0x21)
	p.WriteToPPUAddr(0x15)
	p.WriteToData(0x25)
	if !(p.vram[0x0115] == 0x25) {
		t.Error("PPU vram value not correct")
	}
}

func TestWritesValueToPPUVramInHorizontal(t *testing.T) {
	p := setupTestPPU(HORIZONTAL)
	p.WriteToPPUAddr(0x21)
	p.WriteToPPUAddr(0x15)
	p.WriteToData(0x25)
	if !(p.vram[0x0115] == 0x25) {
		t.Error("PPU vram value not correct")
	}
}
func TestWritesValueToPPUVramInVerticalAnd2NameTable(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToPPUAddr(0x28)
	p.WriteToPPUAddr(0x15)
	p.WriteToData(0x25)
	if !(p.vram[0x0015] == 0x25) {
		t.Error("PPU vram value not correct")
	}
}
func TestWritesValueToPPUVramInHorizontalAnd2NameTable(t *testing.T) {
	p := setupTestPPU(HORIZONTAL)
	p.WriteToPPUAddr(0x28)
	p.WriteToPPUAddr(0x15)
	p.WriteToData(0x25)
	if !(p.vram[0x0415] == 0x25) {
		t.Error("PPU vram value not correct")
	}
}
func TestReadsValueFromPPUVramInVertical(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.vram[0x0115] = 0x25
	p.WriteToPPUAddr(0x21)
	p.WriteToPPUAddr(0x15)
	p.ReadData()
	if !(p.ReadData() == 0x25) {
		t.Error("PPU vram value not correct")
	}
}

func TestReadsValueFromPPUVramInHorizontal(t *testing.T) {
	p := setupTestPPU(HORIZONTAL)
	p.vram[0x0115] = 0x25
	p.WriteToPPUAddr(0x21)
	p.WriteToPPUAddr(0x15)
	p.ReadData()
	if !(p.ReadData() == 0x25) {
		t.Error("PPU vram value not correct")
	}
}
func TestReadsValueFromPPUVramInVerticalAnd2NameTable(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.vram[0x0015] = 0x25
	p.WriteToPPUAddr(0x28)
	p.WriteToPPUAddr(0x15)
	p.ReadData()
	if !(p.ReadData() == 0x25) {
		t.Error("PPU vram value not correct")
	}
}
func TestReadsValueFromPPUVramInHorizontalAnd2NameTable(t *testing.T) {
	p := setupTestPPU(HORIZONTAL)
	p.vram[0x0415] = 0x25
	p.WriteToPPUAddr(0x28)
	p.WriteToPPUAddr(0x15)
	p.ReadData()
	if !(p.ReadData() == 0x25) {
		t.Error("PPU vram value not correct")
	}
}

func TestReadsValueFromChrSpace(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.chr_rom[0x0015] = 0x25
	p.WriteToPPUAddr(0x00)
	p.WriteToPPUAddr(0x15)
	p.ReadData()
	if !(p.ReadData() == 0x25) {
		t.Error("PPU vram value not correct")
	}
}
