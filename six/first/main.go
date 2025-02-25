package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	defaultFunc "github.com/KolobokMysnoy/advert2024/defaultFunc"
	"github.com/joho/godotenv"
)

var LINK string = "https://adventofcode.com/2024/day/6/input"
var FILE_NAME string = ""

func setup() error {
	// setup
	currentWorkDirectory, _ := os.Getwd()

	err := godotenv.Load(currentWorkDirectory + "/../../.env")
	if err != nil {
		return fmt.Errorf("error loading .env file: %s", err)
	}

	FILE_NAME = os.Getenv("FILE_NAME")
	COOKIES.Value = os.Getenv("COOKIE_VALUE")

	return nil
}

func generateMap(readCh chan string) ([][]string, Position) {
	mapFloor := make([][]string, 0)
	pos := Position{}
	for i := range readCh {
		splitI := strings.Split(i, "")
		mapFloor = append(mapFloor, splitI)

		for i, v := range splitI {
			isFound := false
			switch v {
			case "V":
				isFound = true
				pos.watchLine = 2
			case "^":
				isFound = true
				pos.watchLine = 0
			case "<":
				isFound = true
				pos.watchLine = 3
			case ">":
				isFound = true
				pos.watchLine = 1
			}
			if isFound {
				pos.x = i
				pos.y = len(mapFloor) - 1
			}
		}
	}

	return mapFloor, pos
}

type Position struct {
	x         int
	y         int
	watchLine int // 0 up 1 right 2 down 3 left
}

func checkIfOutside[T any](arr [][]T, x int, y int) bool {
	return x >= len(arr[0]) || x < 0 || y < 0 || y >= len(arr)
}
func goPathDumb(mapFloor [][]string, pos Position) Position {
	yDif := 0
	xDif := 0

	switch pos.watchLine {
	case 0:
		yDif = -1
	case 1:
		xDif = 1
	case 2:
		yDif = 1
	case 3:
		xDif = -1
	}

	for {
		mapFloor[pos.y][pos.x] = "X"

		if !(checkIfOutside(mapFloor, pos.x+xDif, pos.y+yDif)) {
			sym := mapFloor[pos.y+yDif][pos.x+xDif]
			if sym == "#" {
				pos.watchLine = (pos.watchLine + 1) % 4
				return pos
			}
		}

		pos.y += yDif
		pos.x += xDif

		if checkIfOutside(mapFloor, pos.x, pos.y) {
			return pos
		}
	}
}

func main() {
	err := setup()
	if err != nil {
		fmt.Println(err)
		return
	}

	isNeedUpdate := flag.Bool("update", false, "Used to override content from site")
	flag.Parse()

	defaultFunc.GetValuesFromLinkToFile(FILE_NAME, LINK, COOKIES, *isNeedUpdate)

	readCh := make(chan string)
	go defaultFunc.GetStringsFromFile(FILE_NAME, readCh)

	// NOTE work on task ------------------
	sum := 0

	generatedMap, position := generateMap(readCh)
	for {
		position = goPathDumb(generatedMap, position)
		if checkIfOutside[string](generatedMap, position.x, position.y) {
			break
		}
	}

	for _, v := range generatedMap {

		for _, v2 := range v {
			if v2 == "X" {
				sum += 1
			}
		}
	}

	fmt.Println(BLUE, "Output: ", RESET, sum)
}
