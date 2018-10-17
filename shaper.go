package main

import "math"

var coordList = []plotterCoordinate{
	plotterCoordinate{0, 0},
	plotterCoordinate{0, 10},
	plotterCoordinate{10, 10},
	plotterCoordinate{10, 0},
	plotterCoordinate{0, 0},
}

func (pc plotterConfig) generatePlotMessage(oldCoord plotterCoordinate, newCoord plotterCoordinate) plotMessage {
	distPerDegree := math.Pi * pc.SpoolDiameter / 360
	oldLengthLeft := math.Sqrt(math.Pow(oldCoord.X, 2) + math.Pow(oldCoord.Y, 2))
	oldLengthRight := math.Sqrt(math.Pow((pc.AnchorDistance-oldCoord.X), 2) + math.Pow(oldCoord.Y, 2))
	newLengthLeft := math.Sqrt(math.Pow(newCoord.X, 2) + math.Pow(newCoord.Y, 2))
	newLengthRight := math.Sqrt(math.Pow(pc.AnchorDistance-newCoord.X, 2) + math.Pow(newCoord.Y, 2))
	leftRotDelta := (newLengthLeft - oldLengthLeft) / distPerDegree
	rightRotDelta := (newLengthRight - oldLengthRight) / distPerDegree
	return plotMessage{
		stepperCommand{int(math.Abs(leftRotDelta)), leftRotDelta > 0},
		stepperCommand{int(math.Abs(rightRotDelta)), rightRotDelta > 0},
	}
}
