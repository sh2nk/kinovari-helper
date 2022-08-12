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
	OneTime bool     `json:"one_time"`
	Inline  bool     `json:"inline"`
	Buttons []Button `json:"buttons"`
}

type KeyboardBuilder struct {
	Keyboard
}
