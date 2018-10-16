package main

import "math"

var coordList = []stepperCoordinate{
	stepperCoordinate{0, 0},
	stepperCoordinate{0, 10},
	stepperCoordinate{10, 10},
	stepperCoordinate{10, 0},
	stepperCoordinate{0, 0},
}

func (s stepperConfig) generatePlotMessage(oldCoord stepperCoordinate, newCoord stepperCoordinate) plotMessage {
	distPerDegree := math.Pi * s.SpoolDiameter / 360
	oldLengthLeft := math.Sqrt(math.Pow(oldCoord.X, 2) + math.Pow(oldCoord.Y, 2))
	oldLengthRight := math.Sqrt(math.Pow((s.AnchorDistance-oldCoord.X), 2) + math.Pow(oldCoord.Y, 2))
	newLengthLeft := math.Sqrt(math.Pow(newCoord.X, 2) + math.Pow(newCoord.Y, 2))
	newLengthRight := math.Sqrt(math.Pow(newCoord.X, 2) + math.Pow(newCoord.Y, 2))
	leftRotDelta := (newLengthLeft - oldLengthLeft) / distPerDegree
	rightRotDelta := (newLengthRight - oldLengthRight) / distPerDegree
	return plotMessage{
		stepperCommand{int(math.Abs(leftRotDelta)), leftRotDelta > 0},
		stepperCommand{int(math.Abs(rightRotDelta)), rightRotDelta > 0},
	}
}
