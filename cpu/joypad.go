package cpu

type JoypadButton uint8

const (
	ButtonA JoypadButton = 0b0000_0001
	ButtonB JoypadButton = 0b0000_0010
	Select  JoypadButton = 0b0000_0100
	Start   JoypadButton = 0b0000_1000
	Up      JoypadButton = 0b0001_0000
	Down    JoypadButton = 0b0010_0000
	Left    JoypadButton = 0b0100_0000
	Right   JoypadButton = 0b1000_0000
)

type Joypad struct {
	strobe     bool
	button_idx uint8
	Button     JoypadButton
}

func NewJoypad() *Joypad {
	return &Joypad{}
}

func (j *Joypad) WriteData(v uint8) {
	j.strobe = v&1 == 1
	if j.strobe {
		j.button_idx = 0
	}
}

func (j *Joypad) ReadData() uint8 {
	if j.button_idx > 7 {
		return 1
	}
	response := (j.Button & (1 << j.button_idx)) >> j.button_idx
	if !j.strobe && j.button_idx <= 7 {
		j.button_idx++
	}
	return uint8(response)
}

func (j *Joypad) SetButtonPressedStatus(button JoypadButton, pressed bool) {
	if pressed {
		j.Button |= button
	} else {
		j.Button &^= button
	}
}
