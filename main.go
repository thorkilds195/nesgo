package main

import (
	"log"
	"math/rand"
	"nesgo/cpu"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 32
	screenHeight = 32
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
	for i := 0; i < 100; i++ {
		alive := e.cpu.Step(func() {
			e.cpu.MemWrite(0xFE, uint8(rand.Intn(15)+1))
		})

		if !alive {
			panic("No longer alive")
		}
	}
	if readScreenState(e.cpu, e) {
		e.texture.WritePixels(e.framebuffer)
	}
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(e.texture, op)
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewEmulator(c *cpu.CPU) *Emulator {
	fb := make([]byte, screenWidth*screenHeight*4)
	c.Reset()
	texture := ebiten.NewImage(screenWidth, screenHeight)

	return &Emulator{
		framebuffer: fb,
		cpu:         c,
		texture:     texture,
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth*10, screenHeight*10)
	ebiten.SetWindowTitle("NES Snake Emulator")
	dat, err := os.ReadFile("./snake.nes")
	if err != nil {
		panic(err)
	}
	rom := cpu.InitRom(dat)
	bus := cpu.InitBus(rom)
	cpu := cpu.InitCPU(bus)
	game := NewEmulator(cpu)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
