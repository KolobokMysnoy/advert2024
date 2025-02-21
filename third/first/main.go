package main

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/joho/godotenv"
)

var LINK string = "https://adventofcode.com/2024/day/3/input"
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

func isDigit(a rune) bool {
	// HACK rewrite
	switch a {
	case '0':
		return true
	case '1':
		return true
	case '2':
		return true
	case '3':
		return true
	case '4':
		return true
	case '5':
		return true
	case '6':
		return true
	case '7':
		return true
	case '8':
		return true
	case '9':
		return true
	default:
		return false
	}
}

func parseDigit(runes []rune, from int) ([]int, int, bool) {
	first := make([]rune, 0)
	second := make([]rune, 0)

	i := from

	for ; i < len(runes) && isDigit(runes[i]); i++ {
		first = append(first, runes[i])
	}

	if i >= len(runes) || runes[i] != ',' {
		return nil, i, false
	}
	i++

	for ; i < len(runes) && isDigit(runes[i]); i++ {
		second = append(second, runes[i])
	}

	if i >= len(runes) || runes[i] != ')' {
		return nil, i, false
	}
	i++

	// HACK
	conv := func(r []rune) int {
		f := 0
		sd := 1

		slices.Reverse(r)

		for _, v := range r {
			f += int(v-'0') * sd
			sd *= 10
		}

		return f
	}

	numbers := []int{conv(first), conv(second)}
	return numbers, i, true
}

func isMul(from int, runes []rune) bool {
	if from+3 >= len(runes) {
		return false
	}

	return runes[from] == 'm' && runes[from+1] == 'u' && runes[from+2] == 'l' && runes[from+3] == '('
}

func searchMul(from int, runes []rune) []int {
	i := from

	fmt.Println(runes, 'm', 'u', 'l', '(')

	retNumbs := make([]int, 0)
	for i < len(runes) {

		if runes[i] == 'm' && isMul(i, runes) {
			i += 4
		} else {
			i++
			continue
		}

		numbs, ind, ok := parseDigit(runes, i)
		i = ind

		if !ok {
			i++
			continue
		}
		retNumbs = append(retNumbs, numbs...)
	}

	return retNumbs
}

func getLineSum(str string) int {
	runes := []rune(str)
	numbs := searchMul(0, runes)

	if len(numbs)%2 != 0 {
		fmt.Println("error at getLineSum: ", numbs)
		return 0
	}

	sum := 0
	for i := 0; i < len(numbs); {
		sum += numbs[i] * numbs[i+1]
		i += 2
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
	sum := 0
	for i := range readCh {
		sum += getLineSum(i)
	}

	fmt.Println(BLUE, "Output: ", RESET, sum)
}
