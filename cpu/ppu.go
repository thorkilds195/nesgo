package cpu

type PPU struct {
	chr_rom       []uint8
	palette_table [32]uint8
	vram          [2048]uint8
	oam_data      [256]uint8
	oam_addr_reg  uint8
	mirroring     Mirroring
	addr          *AddressRegister
	ctrl          *ControlRegister
	mask          *MaskRegister
	status        *StatusRegister
	scroll        *ScrollRegister
	data_buffer   uint8
}

func NewPPU(chr_rom []uint8, mirroring Mirroring) *PPU {
	return &PPU{
		chr_rom:       chr_rom,
		palette_table: [32]uint8{},
		vram:          [2048]uint8{},
		oam_data:      [256]uint8{},
		mirroring:     mirroring,
		addr:          NewAddressRegister(),
		ctrl:          NewControlRegister(),
		mask:          NewMaskRegister(),
		status:        NewStatusRegister(),
		scroll:        NewScrollRegister(),
	}
}

func (p *PPU) WriteToPPUAddr(v uint8) {
	p.addr.Update(v)
}

func (p *PPU) WriteToOAMAddr(v uint8) {
	p.oam_addr_reg = v
}

func (p *PPU) WriteToPPUCtrl(v uint8) {
	p.ctrl.Update(v)
}

func (p *PPU) WriteToScroll(v uint8) {
	p.scroll.write(v)
}

func (p *PPU) WriteToOAMDMA(data *[256]uint8) {
	for i := range data {
		p.oam_data[p.oam_addr_reg] = data[i]
		p.oam_addr_reg++
	}
}

func (p *PPU) ReadStatusRegister() uint8 {
	val := p.status.value
	p.status.resetVblank()
	p.addr.reset_latch()
	p.scroll.resetLatch()
	return val
}

func (p *PPU) incrVramAddr() {
	p.ctrl.VramAddIncrement()
}

func (p *PPU) WriteToData(v uint8) {
	addr := p.addr.get()
	if addr >= 0 && addr <= 0x1FFF {
		panic("Attempt to write to chr rom space")
	} else if addr >= 0x2000 && addr <= 0x2FFF {
		p.vram[p.mirrorVramAddr(addr)] = v
	} else if addr >= 0x3000 && addr <= 0x3EFF {
		panic("0x3000 to 0x3EFF shouldnt be used")
	} else if addr == 0x3F10 || addr == 0x3F14 ||
		addr == 0x3F18 || addr == 0x3F1C {
		addr_mirr := addr - 0x10
		p.palette_table[addr_mirr-0x3F00] = v
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		p.palette_table[addr-0x3F00] = v
	} else {
		panic("Invalid address passed")
	}
	p.ctrl.VramAddIncrement()
}

func (p *PPU) ReadData() uint8 {
	addr := p.addr.get()
	p.incrVramAddr()

	if addr >= 0x0000 && addr <= 0x1FFF {
		ret := p.data_buffer
		p.data_buffer = p.chr_rom[addr]
		return ret
	} else if addr >= 0x2000 && addr <= 0x2FFF {
		ret := p.data_buffer
		p.data_buffer = p.vram[p.mirrorVramAddr(addr)]
		return ret
	} else if addr >= 0x3000 && addr <= 0x3EFF {
		panic("Space 0x3000 to 0x3EFF is not expected to be used")
	} else if addr >= 0x3F00 && addr <= 0x3FFF {
		return p.palette_table[(addr - 0x3F00)]
	}
	panic("Unexpected access to mirrored space")
}

func (p *PPU) WriteOAMData(v uint8) {
	p.oam_data[p.oam_addr_reg] = v
	p.oam_addr_reg++
}

func (p *PPU) ReadOAMData() uint8 {
	return p.oam_data[p.oam_addr_reg]
}

func (p *PPU) WriteToMask(v uint8) {
	p.mask.update(v)
}

func (p *PPU) mirrorVramAddr(addr uint16) uint16 {
	/*
			There exists 1kb of vram in address 0x0000 to 0x400
			and another one ine 0x401 to 0x800
			Each of these have two screen states, where each screen state has a 1 additionl screen state mapped to these
		    So we have screenstate A and B being addresbale from vram space, with each screen state a and b also being addresable

			This results in the below mapping depending on the mirroring type

			HORIZONTAL:
			[A, a]
			[B, b]
			VERTICAL:
			[A, B]
			[a, b]

			This allows either smooth horizontal scrolling when using vertical mapping or
			smooth vertical scrolling when using horizontal mapping
	*/
	mirrored_vram := addr & 0b10111111111111
	vram_idx := mirrored_vram - 0x2000
	name_tbl := vram_idx / 0x400
	if p.mirroring == VERTICAL && (name_tbl == 2 || name_tbl == 3) {
		return vram_idx - 0x800
	} else if p.mirroring == HORIZONTAL {
		if name_tbl == 2 || name_tbl == 3 {
			return vram_idx - 0x400
		} else if name_tbl == 1 {
			return vram_idx - 0x800
		}
	}
	return vram_idx
}

type AddressRegister struct {
	value [2]uint8 // First entry is high byte
	is_hi bool
}

func NewAddressRegister() *AddressRegister {
	return &AddressRegister{
		is_hi: true,
		value: [2]uint8{},
	}
}

func (a *AddressRegister) set(v uint16) {
	hi := uint8(v >> 8)
	lo := uint8(v & 0xFF)
	a.value[0] = hi
	a.value[1] = lo
}

func (a *AddressRegister) Update(v uint8) {
	if a.is_hi {
		a.value[0] = v
	} else {
		a.value[1] = v
	}
	if a.get() > 0x3FFF {
		a.set(a.get() % 0b11111111111111)
	}
	a.is_hi = !a.is_hi
}

func (a *AddressRegister) get() uint16 {
	return (uint16(a.value[1]) | (uint16(a.value[0]) << 8))
}

func (a *AddressRegister) increment() {
	lo := a.value[1]
	a.value[1]++
	if lo > a.value[1] {
		a.value[0]++
	}
	if a.get() > 0x3FFF {
		a.set(a.get() % 0b11111111111111)
	}
}

func (a *AddressRegister) reset_latch() {
	a.is_hi = true
}

type ControlRegister struct {
	value uint8
}

func NewControlRegister() *ControlRegister {
	return &ControlRegister{value: 0b0000_0000}
}

func (c *ControlRegister) VramAddIncrement() uint8 {
	if c.value&0b0000_0100 > 0 {
		return 32
	} else {
		return 1
	}
}

func (c *ControlRegister) Update(v uint8) {
	c.value = v
}

type MaskRegister struct {
	value uint8
}

func NewMaskRegister() *MaskRegister {
	return &MaskRegister{value: 0}
}

func (m *MaskRegister) update(v uint8) {
	m.value = v
}

func (m *MaskRegister) isGreyScaleSet() bool {
	return (m.value & 0b0000_0001) > 0
}
func (m *MaskRegister) isShowBackgroundSet() bool {
	return (m.value & 0b0000_0010) > 0
}
func (m *MaskRegister) isShowSpritesSet() bool {
	return (m.value & 0b0000_0100) > 0
}

func (m *MaskRegister) isBackgrounRenderingSet() bool {
	return (m.value & 0b0000_1000) > 0
}

func (m *MaskRegister) isSpriteRenderingSet() bool {
	return (m.value & 0b0001_0000) > 0
}

func (m *MaskRegister) isEmphaziseRedSet() bool {
	return (m.value & 0b0010_0000) > 0
}

func (m *MaskRegister) isEmphaziseGreenSet() bool {
	return (m.value & 0b0100_0000) > 0
}

func (m *MaskRegister) isEmphaziseBlueSet() bool {
	return (m.value & 0b1000_0000) > 0
}

type StatusRegister struct {
	value uint8
}

func NewStatusRegister() *StatusRegister {
	return &StatusRegister{}
}

func (m *StatusRegister) isVblankSet() bool {
	return (m.value & 0b1000_0000) > 0
}

func (m *StatusRegister) resetVblank() {
	m.value &= 0b0111_1111
}

func (m *StatusRegister) setVblank() {
	m.value |= 0b1000_0000
}

func (m *StatusRegister) setSprite0Flag() {
	m.value |= 0b0100_0000
}

func (m *StatusRegister) setSpriteOverflow() {
	m.value |= 0b0010_0000
}

type ScrollRegister struct {
	x_value uint8
	y_value uint8
	is_x    bool
}

func NewScrollRegister() *ScrollRegister {
	return &ScrollRegister{x_value: 0, y_value: 0, is_x: true}
}

func (s *ScrollRegister) write(v uint8) {
	if s.is_x {
		s.x_value = v
	} else {
		s.y_value = v
	}
	s.is_x = !s.is_x
}

func (s *ScrollRegister) resetLatch() {
	s.is_x = true
}
