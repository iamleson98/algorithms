package main

import (
	"fmt"
	"log"
	"strings"
)

const (
	MIN_NUM_OF_TEST = 1
	MAX_NUM_OF_TEST = 20

	MIN_WIDTH = 1
	MAX_WIDTH = 30

	MIN_HEIGHT = 1
	MAX_HEIGHT = 30

	MIN_NUM_OF_FIRES = 1
	MAX_NUM_OF_FIRES = 2000

	SHIP  = byte('#')
	WATER = byte('_')
)

type Match struct {
	PlayerOne PlayerState
	PlayerTwo PlayerState

	FleetWidth  int
	FleetHeight int
	NumOfShots  int

	conclusion string
}

type PlayerState struct {
	FleetMap        [][]byte
	TotalShipsAlive int
}

func evaluateError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func validateMinMax(name string, value, min, max int) error {
	if value < min || value > max {
		return fmt.Errorf("value %q must >= %d and <= %d", name, min, max)
	}
	return nil
}

// CheckHit checks if the provided x, y point to a ship of player.
// If true, then subtracts total number of its alive ships by 1
// and change that position into water '_' to represent sink ship.
// returns boolean indicates that if the enemy can continue shooting.
func (p *PlayerState) CheckHit(x, y, height int) bool {
	if p.FleetMap[height-1-y][x] == SHIP {
		p.TotalShipsAlive--
		p.FleetMap[height-1-y][x] = WATER
		return true
	}
	return false
}

func (m *Match) ParseShots() {
	playerOneTurn := true

	playerOneSteps := 0
	playerTwoSteps := 0

	for i := 0; i < m.NumOfShots; i++ {
		var x, y int
		_, err := fmt.Scanf("%d %d", &x, &y)
		evaluateError(err)
		evaluateError(validateMinMax("X coordinate", x, 0, m.FleetWidth-1))
		evaluateError(validateMinMax("Y coordinate", y, 0, m.FleetHeight-1))

		if playerOneTurn {
			playerOneTurn = m.PlayerTwo.CheckHit(x, y, m.FleetHeight)
			if m.PlayerTwo.TotalShipsAlive == 0 && playerOneSteps > playerTwoSteps {
				playerOneTurn = false
			}
			playerOneSteps++
		} else {
			playerOneTurn = !m.PlayerOne.CheckHit(x, y, m.FleetHeight)
			if m.PlayerOne.TotalShipsAlive == 0 && playerOneSteps < playerTwoSteps {
				playerOneTurn = true
			}
			playerTwoSteps++
		}
	}

	if m.PlayerOne.TotalShipsAlive == 0 && m.PlayerTwo.TotalShipsAlive > 0 {
		m.conclusion = "player two wins"
	} else if m.PlayerOne.TotalShipsAlive > 0 && m.PlayerTwo.TotalShipsAlive == 0 {
		m.conclusion = "player one wins"
	} else {
		m.conclusion = "draw"
	}
}

func (m *PlayerState) ParsePlayerState(fleetWidth, fleetHeight int) {
	for y := 0; y < fleetHeight; y++ {
		var row string
		_, err := fmt.Scanf("%s", &row)
		evaluateError(err)

		var shipCount = strings.Count(row, string(SHIP))
		var waterCount = strings.Count(row, string(WATER))

		if rowLength := len(row); rowLength != fleetWidth || shipCount+waterCount != rowLength {
			evaluateError(fmt.Errorf("row must contain %d character(s) and only '_' and '#'", fleetWidth))
		}

		m.FleetMap = append(m.FleetMap, []byte(row))
		m.TotalShipsAlive += shipCount
	}
}

func parseMatches(numOfMatches int) []Match {
	var matches = make([]Match, numOfMatches)

	for i := 0; i < numOfMatches; i++ {
		var match Match

		_, err := fmt.Scanf("%d %d %d", &match.FleetWidth, &match.FleetHeight, &match.NumOfShots)
		evaluateError(err)
		evaluateError(validateMinMax("fleet width", match.FleetWidth, MIN_WIDTH, MAX_WIDTH))
		evaluateError(validateMinMax("fleet height", match.FleetHeight, MIN_HEIGHT, MAX_HEIGHT))
		evaluateError(validateMinMax("numOfShots", match.NumOfShots, MIN_NUM_OF_FIRES, MAX_NUM_OF_FIRES))

		match.PlayerOne.ParsePlayerState(match.FleetWidth, match.FleetHeight)
		match.PlayerTwo.ParsePlayerState(match.FleetWidth, match.FleetHeight)
		match.ParseShots()

		matches[i] = match
	}

	return matches
}

func printResult(matches []Match) {
	for _, match := range matches {
		fmt.Println(match.conclusion)
	}
}

func main() {
	var numOfMatches int
	_, err := fmt.Scanf("%d", &numOfMatches)
	evaluateError(err)
	evaluateError(validateMinMax("Num of matches", numOfMatches, MIN_NUM_OF_TEST, MAX_NUM_OF_TEST))

	matches := parseMatches(numOfMatches)
	printResult(matches)
}
