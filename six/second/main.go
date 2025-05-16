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

type Position struct {
	x         int
	y         int
	watchLine int // 0 up 1 right 2 down 3 left
}

func generateMap(strCh chan string) ([][]int, Position) {
	newMap := make([][]int, 0)
	pos := Position{}
	for line := range strCh {
		symbols := strings.Split(line, "")

		wallLine := make([]int, 0)
		for ind, sym := range symbols {
			isFound := false

			switch sym {
			case "#":
				wallLine = append(wallLine, 1)
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
			default:
				wallLine = append(wallLine, 0)
			}

			if isFound {
				pos.x = ind
				pos.y = len(newMap) - 1
				if pos.y < 0 {
					pos.y = 0
				}
				wallLine = append(wallLine, 0)
			}
		}
	}

	return newMap, pos
}

func getCanBePlace(generatedMap [][]int, x, y int) {

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
