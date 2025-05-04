package main

import (
	"fmt"
	"log"
	"nesgo/cpu"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 256
	screenHeight = 240
)

var keyMap map[ebiten.Key]cpu.JoypadButton = map[ebiten.Key]cpu.JoypadButton{
	ebiten.KeyArrowDown:  cpu.Down,
	ebiten.KeyArrowLeft:  cpu.Left,
	ebiten.KeyArrowRight: cpu.Right,
	ebiten.KeyArrowUp:    cpu.Up,
	ebiten.KeySpace:      cpu.Select,
	ebiten.KeyEnter:      cpu.Start,
	ebiten.KeyA:          cpu.ButtonA,
	ebiten.KeyS:          cpu.ButtonB,
}

func handleUserInput(c *cpu.Joypad) {
	for key, button := range keyMap {
		c.SetButtonPressedStatus(button, inpututil.KeyPressDuration(key) > 0)
	}
}

type Emulator struct {
	cpu         *cpu.CPU
	texture     *ebiten.Image
	frame       *cpu.Frame
	drawTime    *bool
	frameCount  int
	lastSecond  time.Time
	internalFPS float64
	cpuCycles   uint
}

func (e *Emulator) Update() error {
	handleUserInput(e.cpu.Bus.Joypad)
	prevCycles := e.cpu.GetCycles()
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
	e.cpuCycles = e.cpu.GetCycles() - prevCycles
	copyToBuffer(e.frame, e)
	e.frameCount++
	now := time.Now()
	if now.Sub(e.lastSecond) >= time.Second {
		e.internalFPS = float64(e.frameCount)
		e.frameCount = 0
		e.lastSecond = now
	}
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(e.texture, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Host FPS: %0.2f\nEmulated FPS: %0.2f\nCPU Cycles: %d", ebiten.ActualFPS(), e.internalFPS, e.cpuCycles))
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewEmulator(c *cpu.CPU, f *cpu.Frame, callTrack *bool) *Emulator {
	c.Reset()
	texture := ebiten.NewImage(screenWidth, screenHeight)
	return &Emulator{
		cpu:        c,
		texture:    texture,
		frame:      f,
		drawTime:   callTrack,
		lastSecond: time.Now(),
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
	ebiten.SetVsyncEnabled(true)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowDecorated(true)
	ebiten.SetTPS(60)
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
