package main

import (
	"flag"
	"fmt"
	"os"

	defaultFunc "github.com/KolobokMysnoy/advert2024/defaultFunc"
	"github.com/joho/godotenv"
)

var LINK string = "https://adventofcode.com/2024/day/5/input"
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

func isCorrectOrder(pages map[int]int, currentPos int, before []int) bool {
	for _, v := range before {
		position, ok := pages[v]
		if !ok {
			continue
		}

		if position >= currentPos {
			return false
		}
	}

	return true
}

func replaceElement(number []int, currentPos int, before []int) bool {
	firstPos := currentPos
	i := currentPos - 1

	for i >= 0 {
		if defaultFunc.FindElement(before, number[i]) {
			tmp := number[firstPos]
			number[firstPos] = number[i]
			number[i] = tmp
			return false
		}
		i--
	}

	return true
}

func isSatisfiesRestriction(pages map[int]int, rules map[int][]int) bool {
	if len(pages) == 0 {
		return false
	}

	for pageNumber := range pages {
		placement := pages[pageNumber]

		beforePages := rules[pageNumber]
		if !isCorrectOrder(pages, placement, beforePages) {
			return false
		}
	}

	return true
}

func getRules(readCh chan string) (map[int][]int, error) {
	rules := make(map[int][]int, 0)
	for i := range readCh {
		if i == "" {
			break
		}

		numbers, err := defaultFunc.ConvertNumbersInString(i, "|")
		if err != nil || len(numbers) != 2 {
			return nil, fmt.Errorf("error at convertNumbersInString: %s", err)
		}

		pageRestrictions, ok := rules[numbers[1]]
		if ok {
			rules[numbers[1]] = append(pageRestrictions, numbers[0])
		} else {
			rules[numbers[1]] = []int{numbers[0]}
		}
	}

	return rules, nil
}

func getSequenceOfPages(numbers []int) (map[int]int, error) {
	sequenceOfPages := make(map[int]int, 0)
	for ind, v := range numbers {
		pagesSeq, ok := sequenceOfPages[v]
		if ok {
			newErr := fmt.Errorf("error page already exist: %d %d", pagesSeq, v)
			return nil, newErr
		}

		sequenceOfPages[v] = ind
	}

	return sequenceOfPages, nil
}

func getCorrectOrder(rules map[int][]int, order []int) ([]int, error) {

	i := len(order) - 1

	lastI := i
	cnt := 0
	fmt.Print(order)
	for i > 0 {
		if cnt == 20 {
			return nil, fmt.Errorf("error here: ", order)
		}
		val := rules[order[i]]

		ok := replaceElement(order, i, val)

		if !ok {
			if lastI == i {
				cnt++
			} else {
				cnt = 0
			}

			continue
		}
		i--
	}

	return order, nil
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

	rules, err := getRules(readCh)
	if err != nil {
		fmt.Println("error happened at getRules: ", err)
		return
	}

	sum := 0
	for i := range readCh {
		numbers, err := defaultFunc.ConvertNumbersInString(i, ",")
		if err != nil {
			fmt.Print("error at convertNumbersInString: ", err)
			return
		}

		sequenceOfPages, err := getSequenceOfPages(numbers)
		if err != nil {
			fmt.Println("error happened at getSequenceOfPages: ", err)
			return
		}

		if !isSatisfiesRestriction(sequenceOfPages, rules) {
			newOrder, err := getCorrectOrder(rules, numbers)
			if err != nil {
				fmt.Println("error happened at getCorrectOrder: ", err)
			}
			sum += newOrder[len(newOrder)/2]
		}
	}

	fmt.Println(BLUE, "Output: ", RESET, sum)
}
