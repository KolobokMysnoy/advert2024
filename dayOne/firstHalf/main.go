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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

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

func convertToNumbers(ch chan string) ([]int, []int, error) {
	left := make([]int, 0)
	right := make([]int, 0)

	for v := range ch {
		strs := strings.Split(v, "   ")

		leftNumb, err := strconv.Atoi(strs[0])
		if err != nil {
			return nil, nil, err
		}

		rightNumb, err := strconv.Atoi(strs[1])
		if err != nil {
			return nil, nil, err
		}

		left = append(left, leftNumb)
		right = append(right, rightNumb)
	}

	if len(left) != len(right) {
		return nil, nil, fmt.Errorf("not equal len")
	}

	left = quickSortStart(left)
	right = quickSortStart(right)
	return left, right, nil
}

func getSum(left, right []int) int {
	sum := 0
	for i := range left {
		sum = sum + abs(left[i]-right[i])
	}
	return sum
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

	firsitArg := os.Args[1:]

	if len(firsitArg) > 0 && firsitArg[0] == "t" {
		f, err := os.Create(fileName)
		if err != nil {
			fmt.Print(err)
			return
		}
		defer f.Close()

		newCh := make(chan string)
		go getValuesFromLink("https://adventofcode.com/2024/day/1/input", newCh)

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

	readCh := make(chan string)
	go getValueFromFile(readCh, fileName)

	left, right, err := convertToNumbers(readCh)
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Println("sum: ", getSum(left, right))
}
