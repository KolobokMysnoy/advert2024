package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func partition(arr []int, low, high int) ([]int, int) {
	pivot := arr[high]
	i := low
	for j := low; j < high; j++ {
		if arr[j] < pivot {
			arr[i], arr[j] = arr[j], arr[i]
			i++
		}
	}
	arr[i], arr[high] = arr[high], arr[i]
	return arr, i
}

func quickSort(arr []int, low, high int) []int {
	if low < high {
		var p int
		arr, p = partition(arr, low, high)
		arr = quickSort(arr, low, p-1)
		arr = quickSort(arr, p+1, high)
	}
	return arr
}

func quickSortStart(arr []int) []int {
	return quickSort(arr, 0, len(arr)-1)
}

func getValueFromFile(ch chan string, name string) {
	f, err := os.Open(name)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		ch <- scanner.Text()
	}

	close(ch)
}

func getValuesFromLink(link string, ch chan string) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Print(err)
		return
	}

	myUrl, err := url.Parse(link)
	if err != nil {
		fmt.Print("incorrect url: ", err)
		return
	}

	cook := []*http.Cookie{&COOKIES}

	jar.SetCookies(myUrl, cook)
	client := &http.Client{
		Jar: jar,
	}

	resp, err := client.Get(link)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		text := scanner.Text()
		ch <- text
	}

	close(ch)

	if err := scanner.Err(); err != nil {
		fmt.Print(err)
		return
	}
}

func getValueToFile(fileName, link string, isNeedToRewrite bool) {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) || isNeedToRewrite {
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Print(err)
			return
		}
		defer f.Close()

		newCh := make(chan string)
		go getValuesFromLink("https://adventofcode.com/2024/day/2/input", newCh)

		for {
			val, ok := <-newCh
			if !ok {
				break
			}

			stringWithEnd := val + "\n"
			f.WriteString(stringWithEnd)
		}

		fmt.Println("Finish writing")
	}

}

func convertToNumbers(str string) ([]int, error) {
	const separator = " "

	arr := make([]int, 0)

	for _, v := range strings.Split(str, separator) {
		numb, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		arr = append(arr, numb)
	}

	return arr, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func findElement(arr []int, target int) bool {
	for _, value := range arr {
		if value == target {
			return true
		}
	}
	return false
}

func checkOnCorrect(values, correctDiff []int) bool {
	var isUp *bool

	for i := 0; i < len(values)-1; i++ {
		dif := values[i] - values[i+1]
		if isUp == nil {
			if dif < 0 {
				isUp = new(bool)
				*isUp = true
			} else {
				isUp = new(bool)
				*isUp = false
			}
		}

		if (*isUp) && dif > 0 || !(*isUp) && dif < 0 {
			return false
		}

		if !findElement(correctDiff, abs(dif)) {
			return false
		}
	}
	return true
}

func main() {
	currentWorkDirectory, _ := os.Getwd()

	err := godotenv.Load(currentWorkDirectory + "/../../.env")
	if err != nil {
		fmt.Print("error loading .env file:", err)
		return
	}

	fileName := os.Getenv("FILE_NAME")
	COOKIES.Value = os.Getenv("COOKIE_VALUE")

	args := os.Args
	getValueToFile(fileName, "https://adventofcode.com/2024/day/2/input", args[len(args)-1] == "get")

	readCh := make(chan string)
	go getValueFromFile(readCh, fileName)

	correctLines := 0
	for str := range readCh {
		numbs, err := convertToNumbers(str)
		if err != nil {
			fmt.Println(err)
			return
		}

		correctNumbers := []int{1, 2, 3}

		correct := checkOnCorrect(numbs, correctNumbers)
		if correct {
			correctLines += 1
		}
	}

	fmt.Print(correctLines)
}
