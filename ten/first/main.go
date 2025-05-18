package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	defaultFunc "github.com/KolobokMysnoy/advert2024/defaultFunc"
	"github.com/joho/godotenv"
)

const (
	DAY = 10
)

var LINK string = fmt.Sprintf("https://adventofcode.com/2024/day/%d/input", DAY)
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

// ----

const (
	START_NUMB  = 0
	FINISH_NUMB = 9
	MAX_DIFF    = 1
)

type FinishPoint struct {
	X int
	Y int
}

type MapTile struct {
	IsVisited          bool
	AchievableFinishes []FinishPoint
	Value              int
}

type MapLine = []MapTile
type Map = []MapLine

func Exist(finishes []FinishPoint, point FinishPoint) bool {
	for _, v := range finishes {
		if v.X == point.X && v.Y == point.Y {
			return true
		}
	}

	return false
}

func Merge(left []FinishPoint, right []FinishPoint) []FinishPoint {
	finishesVal := make([]FinishPoint, len(left))
	copy(finishesVal, left)

	for _, v := range right {
		if !Exist(finishesVal, v) {
			finishesVal = append(finishesVal, v)
		}
	}

	return finishesVal
}

func createMap(mapIsland [][]int) Map {
	newMap := Map{}

	for y := 0; y < len(mapIsland); y++ {
		mapLine := MapLine{}

		for x := 0; x < len(mapIsland[0]); x++ {
			title := MapTile{
				IsVisited:          false,
				AchievableFinishes: []FinishPoint{},
				Value:              mapIsland[y][x],
			}

			mapLine = append(mapLine, title)
		}

		newMap = append(newMap, mapLine)
	}

	return newMap
}

func convertLine(line string) ([]int, []int, error) {
	lineNumber := []int{}
	starts := []int{}

	for i, v := range line {
		number, err := strconv.Atoi(string(v))
		if err != nil {
			return nil, nil, err
		}

		lineNumber = append(lineNumber, number)
		if number == 0 {
			starts = append(starts, i)
		}
	}

	return lineNumber, starts, nil
}

func getFinish(mapOfIsland Map, x, y int) []FinishPoint {
	if mapOfIsland[y][x].IsVisited {
		return mapOfIsland[y][x].AchievableFinishes
	}
	currentVal := mapOfIsland[y][x].Value

	if currentVal == FINISH_NUMB {
		finish := FinishPoint{
			X: x,
			Y: y,
		}

		mapOfIsland[y][x].AchievableFinishes = append(mapOfIsland[y][x].AchievableFinishes, finish)
		mapOfIsland[y][x].IsVisited = true

		return mapOfIsland[y][x].AchievableFinishes
	}

	mapOfIsland[y][x].IsVisited = true

	var sum []FinishPoint

	// left
	if x-1 >= 0 && mapOfIsland[y][x-1].Value-currentVal == MAX_DIFF {
		sum = Merge(getFinish(mapOfIsland, x-1, y), sum)
	}

	// right
	if x+1 < len(mapOfIsland[0]) && mapOfIsland[y][x+1].Value-currentVal == MAX_DIFF {
		sum = Merge(getFinish(mapOfIsland, x+1, y), sum)
	}

	// up
	if y-1 >= 0 && mapOfIsland[y-1][x].Value-currentVal == MAX_DIFF {
		sum = Merge(getFinish(mapOfIsland, x, y-1), sum)
	}

	// down
	if y+1 < len(mapOfIsland) && mapOfIsland[y+1][x].Value-currentVal == MAX_DIFF {
		sum = Merge(getFinish(mapOfIsland, x, y+1), sum)
	}

	// fmt.Printf("Y %d X %d Sum %d\n", y, x, sum)
	mapOfIsland[y][x].AchievableFinishes = Merge(sum, mapOfIsland[y][x].AchievableFinishes)
	return mapOfIsland[y][x].AchievableFinishes

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
	var sum int
	sum = 0

	allStarts := [][]int{}
	mapOfIsland := [][]int{}

	ind := 0
	for v := range readCh {
		currentLine, currentStarts, err := convertLine(v)
		if err != nil {
			fmt.Println(err)
			return
		}

		mapOfIsland = append(mapOfIsland, currentLine)
		for _, v := range currentStarts {
			allStarts = append(allStarts, []int{ind, v})
		}

		ind++
	}

	mapOfIslandWithTitles := createMap(mapOfIsland)

	for _, v := range allStarts {
		paths := getFinish(mapOfIslandWithTitles, v[1], v[0])
		sum += len(paths)
	}

	fmt.Println(BLUE, "Output: ", sum)
}
