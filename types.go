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

type clientType int

const (
	plotter clientType = 0
	monitor            = 1
)

type registrationMessage struct {
	ClientType clientType `json:"clientType"`
}

type messageType int

const (
	status         messageType = 0
	registration               = 1
	commandRequest             = 2
	plotCommand                = 3
)

type wsMessage struct {
	MessageType messageType `json:"msgType"`
	Payload     interface{} `json:"payload"`
}
