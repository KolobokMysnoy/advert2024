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
	DAY = 9
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
	FREE_SPACE_SYMBOL = -1
)

// stupid way
func getCheckSum(filesSizes []int) (int, error) {
	checksum := 0

	for i, v := range filesSizes {
		if v == FREE_SPACE_SYMBOL {
			break
		}

		checksum += v * i
	}

	return checksum, nil
}

// convert from fileLayout to string with file info
func convertToString(fileLayout string) []int {
	fileIndex := 0
	resultString := []int{}

	for i, v := range fileLayout {
		number, err := strconv.Atoi(string(v))
		if err != nil {
			fmt.Println(fmt.Errorf("cant convert to int rune: %w", err))
			return []int{}
		}

		if i%2 == 1 {
			// free space

			for countOfTimes := 0; countOfTimes < number; countOfTimes++ {
				resultString = append(resultString, FREE_SPACE_SYMBOL)
			}
		} else {
			// fileSize
			for countOfTimes := 0; countOfTimes < number; countOfTimes++ {
				resultString = append(resultString, fileIndex)
			}

			fileIndex++
		}
	}

	return resultString
}

func getIndexFileRight(fileSizes []int, indexOfStart int) int {
	for ; indexOfStart > 0; indexOfStart-- {
		if fileSizes[indexOfStart] != FREE_SPACE_SYMBOL {
			return indexOfStart
		}
	}

	return -1
}

func normalizeSpace(fileSizes []int) []int {
	// normalized := make([]rune, utf8.RuneCountInString(fileSizes))
	normalized := fileSizes

	rightIndex := getIndexFileRight(fileSizes, len(fileSizes)-1)
	if rightIndex == -1 {
		return fileSizes
	}

	for i, v := range fileSizes {
		if rightIndex <= i {
			break
		}
		if v == FREE_SPACE_SYMBOL {
			newRightIndex := getIndexFileRight(fileSizes, rightIndex)
			normalized[i] = normalized[newRightIndex]
			normalized[newRightIndex] = FREE_SPACE_SYMBOL

			rightIndex = newRightIndex - 1
		}
	}

	return normalized

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

	for v := range readCh {
		convertedString := convertToString(v)

		normalizedFileSizes := normalizeSpace(convertedString)
		checkSum, err := getCheckSum(normalizedFileSizes)
		if err != nil {
			fmt.Println(err)
		}

		sum += checkSum
	}

	fmt.Println(BLUE, "Output: ", sum)
}
