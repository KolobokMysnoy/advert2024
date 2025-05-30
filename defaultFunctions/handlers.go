package defaultFunc

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func GetStringsFromFile(nameOfFile string, whereToWrite chan string) {
	defer close(whereToWrite)

	file, err := os.Open(nameOfFile)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		whereToWrite <- scanner.Text()
	}
}

func GetDataFromLink(link string, cookie http.Cookie, whereToWrite chan string) {
	defer close(whereToWrite)

	// prepare cookie
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("cant create cookie jar: ", err)
		return
	}

	cookies := []*http.Cookie{&cookie}

	// set cookies
	myUrl, err := url.Parse(link)
	if err != nil {
		fmt.Println("cant parse url: ", err)
		return
	}

	jar.SetCookies(myUrl, cookies)
	client := &http.Client{
		Jar: jar,
	}

	// get data
	resp, err := client.Get(link)
	if err != nil {
		fmt.Println("cant get data: ", err)
		return
	}
	defer resp.Body.Close()

	// write data to channel
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		whereToWrite <- scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error: ", err)
		return
	}
}

func GetValuesFromLinkToFile(fileName, link string, cookie http.Cookie, isNeedToRewrite bool) {
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) || isNeedToRewrite {
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Print(err)
			return
		}
		defer file.Close()

		linkChannel := make(chan string)
		go GetDataFromLink(link, cookie, linkChannel)

		for {
			val, ok := <-linkChannel
			if !ok {
				break
			}

			stringWithNewLine := val + "\n"
			file.WriteString(stringWithNewLine)
		}
	}
}

func ConvertNumbersInString(str string, separator string) ([]int, error) {
	numbers := make([]int, 0)

	for _, v := range strings.Split(str, separator) {
		if v == "" {
			return []int{}, fmt.Errorf("can't convert empty string")
		}

		num, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}

		numbers = append(numbers, num)
	}

	return numbers, nil
}
