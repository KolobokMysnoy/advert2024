package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	defaultFunc "github.com/KolobokMysnoy/advert2024/defaultFunc"
	"github.com/joho/godotenv"
)

var LINK string = "https://adventofcode.com/2024/day/7/input"
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
var OPERATORS = []string{"+", "*"}

func checkOperators(index int, currentSum, neededSum int, numbers []int, operator string) []int {
	if index+1 >= len(numbers) {
		return []int{currentSum}
	}

	var localSum int
	rightValue := numbers[index+1]

	switch operator {
	case "+":
		localSum = currentSum + rightValue
	case "*":
		localSum = currentSum * rightValue
	}

	if localSum > neededSum {
		return []int{}
	}

	sums := make([]int, 0)
	plusSum := checkOperators(index+1, localSum, neededSum, numbers, "+")
	productSum := checkOperators(index+1, localSum, neededSum, numbers, "*")

	sums = append(sums, plusSum...)
	sums = append(sums, productSum...)

	return sums
}

func checkLine(inputLine string) (int, bool) {
	splitString := strings.Split(inputLine, ":")
	if len(splitString) != 2 {
		fmt.Printf("error at checkLine: incorrect line: line %s \n", inputLine)
		return 0, false
	}

	neededValue, err := strconv.Atoi(splitString[0])
	if err != nil {
		return 0, false
	}

	numbers, err := defaultFunc.ConvertNumbersInString(strings.TrimSpace(splitString[1]), " ")
	if err != nil {
		fmt.Print(err)
		return 0, false
	}

	plusSums := checkOperators(0, numbers[0], neededValue, numbers, "+")
	for _, v := range plusSums {
		if v == neededValue {
			return neededValue, true
		}
	}

	productSums := checkOperators(0, numbers[0], neededValue, numbers, "*")
	for _, v := range productSums {
		if v == neededValue {
			return neededValue, true
		}
	}

	return 0, false
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
		number, ok := checkLine(v)
		if ok {
			sum += number
		}
	}

	fmt.Println(BLUE, "Output: ", sum)
}
