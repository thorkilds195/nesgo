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
		t.Error("Chr rom return value not correct")
	}
}

func TestWritesToOAMAddr(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToOAMAddr(0x25)
	if !(p.oam_addr_reg == 0x25) {
		t.Error("OAM addr value not correct")
	}
}
func TestWritesToPPUCtrl(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToPPUCtrl(0x25)
	if !(p.ctrl.value == 0x25) {
		t.Error("CTRL value not correct")
	}
}
func TestWritesToScrollInXMode(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToScroll(0x25)
	if !(p.scroll.x_value == 0x25) {
		t.Error("Scroll x_value not correct")
	}
}

func TestWritesToScrollInYMode(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToScroll(0x15)
	p.WriteToScroll(0x25)
	if !(p.scroll.x_value == 0x15 && p.scroll.y_value == 0x25) {
		t.Error("Scroll values not correct")
	}
}
func TestReadsStatusRegister(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.status.value = 0b1100_1000
	if !(p.ReadStatusRegister() == 0b1100_1000 &&
		!p.status.isVblankSet() &&
		p.addr.is_hi &&
		p.scroll.is_x) {
		t.Error("Wrong status register value")
	}
}

func TestInteractsWithOAMRegister(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteOAMData(0x25)
	if !(p.oam_data[0] == 0x25) {
		t.Error("Wrong value for OAM data write")
	}
	p.oam_data[1] = 0x15
	if !(p.ReadOAMData() == 0x15) {
		t.Error("Wrong value for OAM data read")
	}
}

func TestWritesToMask(t *testing.T) {
	p := setupTestPPU(VERTICAL)
	p.WriteToMask(0b0001_1000)
	if !(p.mask.value == 0b0001_1000) {
		t.Error("Wrong value for mask data write")
	}
}
