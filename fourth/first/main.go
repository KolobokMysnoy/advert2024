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
	if (x+3) >= len(bigTable[y]) || (y+3) >= len(bigTable) {
		return 0
	}

	if bigTable[y][x] == 'X' && bigTable[y+1][x+1] == 'M' && bigTable[y+2][x+2] == 'A' && bigTable[y+3][x+3] == 'S' {
		return 1
	}

	if bigTable[y+3][x+3] == 'X' && bigTable[y+2][x+2] == 'M' && bigTable[y+1][x+1] == 'A' && bigTable[y][x] == 'S' {
		return 1
	}

	return 0
}

func searchDiagonalSecond(bigTable [][]rune, x, y int) int {
	if x-3 < 0 || y+3 >= len(bigTable) {
		return 0
	}

	if bigTable[y][x] == 'X' && bigTable[y+1][x-1] == 'M' && bigTable[y+2][x-2] == 'A' && bigTable[y+3][x-3] == 'S' {
		return 1
	}

	if bigTable[y+3][x-3] == 'X' && bigTable[y+2][x-2] == 'M' && bigTable[y+1][x-1] == 'A' && bigTable[y][x] == 'S' {
		return 1
	}

	return 0
}

// search only from top to bottom(left corner up, right corner up)
func searchDiagonal(bigTable [][]rune, x, y int) int {
	return searchDiagonalMain(bigTable, x, y) + searchDiagonalSecond(bigTable, x, y)
}

// search only from left to right
func searchHorizontal(bigTable [][]rune, x, y int) int {
	if x+3 >= len(bigTable[0]) || y >= len(bigTable) {
		return 0
	}

	row := bigTable[y]

	if row[x] == 'X' && row[x+1] == 'M' && row[x+2] == 'A' && row[x+3] == 'S' {
		return 1
	}

	if row[x+3] == 'X' && row[x+2] == 'M' && row[x+1] == 'A' && row[x] == 'S' {
		return 1
	}

	return 0
}

// search only from top to bottom
func searchVertical(bigTable [][]rune, x, y int) int {
	if y+3 >= len(bigTable) || x >= len(bigTable[y]) {
		return 0
	}

	if bigTable[y][x] == 'X' && bigTable[y+1][x] == 'M' && bigTable[y+2][x] == 'A' && bigTable[y+3][x] == 'S' {
		return 1
	}

	if bigTable[y+3][x] == 'X' && bigTable[y+2][x] == 'M' && bigTable[y+1][x] == 'A' && bigTable[y][x] == 'S' {
		return 1
	}

	return 0
}

func searchForXmas(bigTable [][]rune) int {
	x := 0
	y := 0

	sum := 0
	for y < len(bigTable) {
		for x < len(bigTable[0]) {
			sum += searchHorizontal(bigTable, x, y)
			sum += searchDiagonal(bigTable, x, y)
			sum += searchVertical(bigTable, x, y)

			x++
		}

		x = 0
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
