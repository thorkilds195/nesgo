package main

import (
	"log"
	"nesgo/cpu"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 256
	screenHeight = 240
)

func colorFromByte(b uint8) (r, g, b2 byte) {
	switch b {
	case 0:
		return 0x00, 0x00, 0x00
	case 1:
		return 0xFF, 0xFF, 0xFF
	case 2, 9:
		return 0x80, 0x80, 0x80
	case 3, 10:
		return 0xFF, 0x00, 0x00
	case 4, 11:
		return 0x00, 0xFF, 0x00
	case 5, 12:
		return 0x00, 0x00, 0xFF
	case 6, 13:
		return 0xFF, 0x00, 0xFF
	case 7, 14:
		return 0xFF, 0xFF, 0x00
	default:
		return 0x00, 0xFF, 0xFF
	}
}

func readScreenState(c *cpu.CPU, g *Emulator) bool {
	idx := 0
	update := false
	buf := g.framebuffer
	for i := 0x0200; i < 0x0600; i++ {
		color_idx := c.MemRead(uint16(i))
		r, g, b := colorFromByte(color_idx)
		if buf[idx] != r || buf[idx+1] != g || buf[idx+2] != b {
			buf[idx] = r
			buf[idx+1] = g
			buf[idx+2] = b
			buf[idx+3] = 0xFF
			update = true
		}
		idx += 4
	}
	return update
}

func handleUserInput(cpu *cpu.CPU) {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		cpu.MemWrite(0xFF, 0x77)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		cpu.MemWrite(0xFF, 0x73)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		cpu.MemWrite(0xFF, 0x61)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		cpu.MemWrite(0xFF, 0x64)
	}
}

type Emulator struct {
	framebuffer []byte
	cpu         *cpu.CPU
	texture     *ebiten.Image
	frame       *cpu.Frame
	drawTime    *bool
}

func dumpFramebuffer(fb []byte) {
	for y := 0; y < screenHeight; y++ {
		for x := 0; x < screenWidth; x++ {
			idx := (y*screenWidth + x) * 4
			r := fb[idx]
			if r > 0 {
				print("█")
			} else {
				print(" ")
			}
		}
		println()
	}
	println("-------------------------------------------------")
}

func (e *Emulator) Update() error {
	handleUserInput(e.cpu)
	for {
		alive := e.cpu.Step(func() {})

		if !alive {
			panic("No longer alive")
		}
		if *e.drawTime {
			*e.drawTime = false
			break
		}
	}
	copyToBuffer(e.frame, e)
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(e.texture, op)
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewEmulator(c *cpu.CPU, f *cpu.Frame, callTrack *bool) *Emulator {
	fb := make([]byte, screenWidth*screenHeight*4)
	c.Reset()
	texture := ebiten.NewImage(screenWidth, screenHeight)
	return &Emulator{
		framebuffer: fb,
		cpu:         c,
		texture:     texture,
		frame:       f,
		drawTime:    callTrack,
	}
}

func showTileBank(chr_rom []uint8, bank uint32) *cpu.Frame {
	frame := cpu.NewFrame()
	bank = (bank * 0x1000)
	tile_y := 0
	tile_x := 0
	for tile_nr := 0; tile_nr < 255; tile_nr++ {
		if tile_nr != 0 && tile_nr%20 == 0 {
			tile_y += 10
			tile_x = 0
		}
		tile := chr_rom[(bank + uint32(tile_nr)*16):(bank + uint32(tile_nr)*16 + 16)]
		for y := 0; y <= 7; y++ {
			upper := tile[y]
			lower := tile[y+8]
			for x := 7; x >= 0; x-- {
				value := (1&upper)<<1 | (1 & lower)
				upper = upper >> 1
				lower = lower >> 1
				var rgb cpu.RGB
				switch value {
				case 0:
					rgb = cpu.SYSTEM_PALLETE[0x01]
				case 1:
					rgb = cpu.SYSTEM_PALLETE[0x23]
				case 2:
					rgb = cpu.SYSTEM_PALLETE[0x27]
				case 3:
					rgb = cpu.SYSTEM_PALLETE[0x30]
				default:
					panic("Not valid rgb rom")
				}
				frame.SetPixel(uint32(tile_x+x), uint32(tile_y+y), rgb)
			}
		}
		tile_x += 10
	}
	return frame
}

func copyToBuffer(f *cpu.Frame, e *Emulator) {
	e.texture.WritePixels(f.Data)
}

func main() {
	ebiten.SetWindowSize(screenWidth*10, screenHeight*10)
	ebiten.SetWindowTitle("NES Emulator")
	dat, err := os.ReadFile("./pacman.nes")
	if err != nil {
		panic(err)
	}
	var callTrack bool
	frame := cpu.NewFrame()
	rom := cpu.InitRom(dat)
	bus := cpu.InitBus(rom, func(p *cpu.PPU) {
		frame.Render(p)
		callTrack = true
	},
	)
	cpu := cpu.InitCPU(bus)
	game := NewEmulator(cpu, frame, &callTrack)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
