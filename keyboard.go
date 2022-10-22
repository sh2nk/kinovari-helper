package main

type Action struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
	Label   string `json:"label"`
}

type Button struct {
	Action Action `json:"action"`
	Color  string `json:"color"`
}

type Keyboard struct {
	OneTime    bool        `json:"one_time"`
	Inline     bool        `json:"inline"`
	ButtonRows []ButtonRow `json:"buttons"`
}

type ButtonRow []Button

func NewKeyboard() *Keyboard {
	return &Keyboard{OneTime: false, Inline: true}
}

func (k *Keyboard) AddRow() {
	k.ButtonRows = append(k.ButtonRows, ButtonRow{})
}

func (k *Keyboard) AddRows(count int) {
	for i := 0; i < count; i++ {
		k.AddRow()
	}
}

func (br *ButtonRow) AddButton(label, color string, payload ...string) {
	var p string
	if len(payload) > 0 {
		p = payload[0]
	} else {
		p = ""
	}
	b := Button{
		Action: Action{
			Type:    "callback",
			Payload: p,
			Label:   label,
		},
		Color: color,
	}
	*br = append(*br, b)
}
