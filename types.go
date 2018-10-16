package main

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
		stepperCommand{360, true},
		stepperCommand{360, false},
	},
	plotMessage{
		stepperCommand{360, false},
		stepperCommand{360, true},
	},
	plotMessage{
		stepperCommand{360, true},
		stepperCommand{360, false},
	},
	plotMessage{
		stepperCommand{360, false},
		stepperCommand{360, true},
	},
}

type wsMessage struct {
	MessageType string      `json:"msgType"`
	Payload     interface{} `json:"payload"`
}
