package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type CheckCorrectFunc func(int) bool

func calculateUpFromDif(dif int) bool {
	return dif < 0
}

func getGlobalIsUp(values []int) bool {
	difOfFirstElements := calculateUpFromDif(values[0] - values[1])
	difOfCenterElements := calculateUpFromDif(values[len(values)/2] - values[len(values)/2+1])
	difOfEndElements := calculateUpFromDif(values[len(values)-2] - values[len(values)-1])

	cntTrue := 0

	if difOfFirstElements {
		cntTrue++
	}

	if difOfCenterElements {
		cntTrue++
	}

	if difOfEndElements {
		cntTrue++
	}

	return cntTrue >= 2
}

func getIndexToRemove(sequence []int, indexOfElement int, checkCorrect CheckCorrectFunc) int {
	i := indexOfElement

	getDiff := func(a, b int) int {
		return sequence[a] - sequence[b]
	}

	// remove last element
	if i == len(sequence)-2 {
		return len(sequence)
	}

	// remove from first element
	if i == 0 {
		isWithoutLeftCorrect := checkCorrect(getDiff(i+1, i+2))
		isWithoutCenterCorrect := checkCorrect(getDiff(i, i+2))

		if isWithoutLeftCorrect {
			return i
		}

		if isWithoutCenterCorrect {
			return i + 1
		}

		return -1
	}

	isWithoutCenterCorrect := checkCorrect(getDiff(i-1, i+1))
	isWithoutRightCorrect := checkCorrect(getDiff(i, i+2))

	if isWithoutCenterCorrect {
		return i
	}

	if isWithoutRightCorrect {
		return i + 1
	}

	return -1
}

func isSequenceCorrect(sequence, correctDifferences []int) bool {
	globalIsUp := getGlobalIsUp(sequence)
	tolerance := 0

	isDifferenceCorrect := func(dif int) bool {
		return findElement(correctDifferences, abs(dif)) && calculateUpFromDif(dif) == globalIsUp
	}

	for i := 0; i < len(sequence)-1; i++ {
		dif := sequence[i] - sequence[i+1]

		if !isDifferenceCorrect(dif) {
			if tolerance > 0 {
				return false
			}

			tolerance++

			removeIndex := getIndexToRemove(sequence, i, isDifferenceCorrect)
			if removeIndex == -1 {
				return false
			}

			if removeIndex-1 == i {
				sequence[removeIndex] = sequence[i]
			}
		}
	}
	return true
}

var LINK string = "https://adventofcode.com/2024/day/2/input"
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

	correctLines := 0
	for str := range readCh {
		numbs, err := convertNumbersInString(str, " ")
		if err != nil {
			fmt.Println(err)
			return
		}

		correctNumbers := []int{1, 2, 3}

		correct := isSequenceCorrect(numbs, correctNumbers)
		if correct {
			correctLines += 1
		}
	}

	reset := "\033[0m"
	blue := "\033[34m"
	fmt.Println(blue, "Correct lines: ", reset, correctLines)
}
