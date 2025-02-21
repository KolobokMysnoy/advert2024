package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var LINK string = "https://adventofcode.com/2024/day/4/input"
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

func searchDiagonalMain(bigTable [][]rune, x, y int) int {
	if bigTable[y+1][x+1] == 'M' && bigTable[y-1][x-1] == 'S' {
		return 1
	}

	if bigTable[y-1][x-1] == 'M' && bigTable[y+1][x+1] == 'S' {
		return 1
	}

	return 0
}

func searchDiagonalSecond(bigTable [][]rune, x, y int) int {
	if bigTable[y+1][x-1] == 'M' && bigTable[y-1][x+1] == 'S' {
		return 1
	}

	if bigTable[y-1][x+1] == 'M' && bigTable[y+1][x-1] == 'S' {
		return 1
	}

	return 0
}

// search only from top to bottom(left corner up, right corner up)
func searchDiagonal(bigTable [][]rune, x, y int) int {
	if bigTable[y][x] != 'A' {
		return 0
	}

	if x+1 >= len(bigTable[0]) ||
		x-1 < 0 ||
		y+1 >= len(bigTable) ||
		y-1 < 0 {
		return 0
	}

	if searchDiagonalMain(bigTable, x, y)+searchDiagonalSecond(bigTable, x, y) == 2 {
		return 1
	}
	return 0
}

func searchForXmas(bigTable [][]rune) int {
	x := 1
	y := 1

	sum := 0
	for y < len(bigTable)-1 {
		for x < len(bigTable[0])-1 {
			sum += searchDiagonal(bigTable, x, y)

			x++
		}

		x = 1
		y++
	}

	return sum
}

func main() {
	err := setup()
	if err != nil {
		fmt.Println(err)
		return
	}

	isNeedUpdate := flag.Bool("update", false, "Used to override content from site")
	flag.Parse()

	getValuesFromLinkToFile(FILE_NAME, LINK, COOKIES, *isNeedUpdate)

	readCh := make(chan string)
	go getStringsFromFile(FILE_NAME, readCh)

	// NOTE work on task ------------------

	bigTable := make([][]rune, 0)
	for i := range readCh {
		bigTable = append(bigTable, []rune(i))
	}

	sum := searchForXmas(bigTable)

	fmt.Println(BLUE, "Output: ", RESET, sum)
}
