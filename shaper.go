package main

import "math"

var coordList = []plotterCoordinate{
	plotterCoordinate{0, 0},
	plotterCoordinate{0, 100},
	plotterCoordinate{100, 100},
	plotterCoordinate{100, 0},
	plotterCoordinate{0, 0},
}

func (ps plotterStatus) generatePlotMessage(newCoord plotterCoordinate) plotMessage {
	pc := ps.Config
	oldCoord := ps.CurrentCoord
	distPerDegree := math.Pi * pc.SpoolDiameter / 360
	oldLengthLeft := math.Sqrt(math.Pow(pc.StartCoord.X+oldCoord.X, 2) + math.Pow(pc.StartCoord.Y+oldCoord.Y, 2))
	oldLengthRight := math.Sqrt(math.Pow((pc.AnchorDistance-(pc.StartCoord.X+oldCoord.X)), 2) + math.Pow(pc.StartCoord.Y+oldCoord.Y, 2))
	newLengthLeft := math.Sqrt(math.Pow(pc.StartCoord.X+newCoord.X, 2) + math.Pow(pc.StartCoord.Y+newCoord.Y, 2))
	newLengthRight := math.Sqrt(math.Pow(pc.AnchorDistance-(pc.StartCoord.X+newCoord.X), 2) + math.Pow(pc.StartCoord.Y+newCoord.Y, 2))
	leftRotDelta := (newLengthLeft - oldLengthLeft) / distPerDegree
	rightRotDelta := (newLengthRight - oldLengthRight) / distPerDegree
	return plotMessage{
		stepperCommand{int(math.Abs(leftRotDelta)), leftRotDelta > 0},
		stepperCommand{int(math.Abs(rightRotDelta)), rightRotDelta > 0},
	}
}
