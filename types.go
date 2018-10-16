package main

type stepperConfig struct {
	SpoolDiameter  float64
	AnchorDistance float64
}

type stepperCoordinate struct {
	X float64
	Y float64
}
type stepperCommand struct {
	Degrees   int  `json:"deg"`
	Clockwise bool `json:"dir"`
}

type plotMessage struct {
	LeftStepper  stepperCommand `json:"left"`
	RightStepper stepperCommand `json:"right"`
}

var commandList = []plotMessage{
	plotMessage{
		stepperCommand{360 * 4, true},
		stepperCommand{360 * 4, true},
	},
	plotMessage{
		stepperCommand{360 * 4, false},
		stepperCommand{360 * 4, false},
	},
}

type wsMessage struct {
	MessageType string      `json:"msgType"`
	Payload     interface{} `json:"payload"`
}
