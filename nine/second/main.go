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

type FileInfo struct {
	IndexLeft int
	Length    int
}

func (d *FileInfo) printFileInfo() {
	fmt.Printf("Index %d Len  %d \n", d.IndexLeft, d.Length)
}

type DiskInfo struct {
	Files       []FileInfo
	EmptySpaces []FileInfo
	DiskLine    []int
}

func (d *DiskInfo) printFileInfo() {
	fmt.Println("Files")
	for _, v := range d.Files {
		v.printFileInfo()
	}

	fmt.Println("Empty")
	for _, v := range d.EmptySpaces {
		v.printFileInfo()
	}

	fmt.Println("Line")
	fmt.Println(d.DiskLine)
	fmt.Println("End of file info")
}

const (
	FREE_SPACE_SYMBOL = -1
)

// stupid way
func getCheckSum(filesSizes []int) (int, error) {
	checksum := 0

	for i, v := range filesSizes {
		if v != FREE_SPACE_SYMBOL {
			checksum += v * i
		}
	}

	return checksum, nil
}

// convert from fileLayout to string with file info
func convertToString(fileLayout string) (DiskInfo, error) {
	fileIndex := 0
	disk := DiskInfo{
		Files:       []FileInfo{},
		EmptySpaces: []FileInfo{},
		DiskLine:    []int{},
	}

	resultString := []int{}

	for i, v := range fileLayout {
		number, err := strconv.Atoi(string(v))
		if err != nil {
			return DiskInfo{}, fmt.Errorf("cant convert to int rune: %w", err)
		}

		if i%2 == 1 {
			if number == 0 {
				continue
			}

			// free space
			newFreeSpace := FileInfo{
				Length:    number,
				IndexLeft: len(resultString),
			}

			disk.EmptySpaces = append(disk.EmptySpaces, newFreeSpace)

			for countOfTimes := 0; countOfTimes < number; countOfTimes++ {
				resultString = append(resultString, FREE_SPACE_SYMBOL)
			}
		} else {
			// fileSize

			newFileSpace := FileInfo{
				Length:    number,
				IndexLeft: len(resultString),
			}

			disk.Files = append(disk.Files, newFileSpace)

			for countOfTimes := 0; countOfTimes < number; countOfTimes++ {
				resultString = append(resultString, fileIndex)
			}

			fileIndex++
		}
	}

	disk.DiskLine = resultString

	return disk, nil
}

func normalizeSpace(fileSizes DiskInfo) DiskInfo {
	normalized := fileSizes.DiskLine

	for indOfFile := len(fileSizes.Files) - 1; indOfFile >= 0; indOfFile-- {
		file := fileSizes.Files[indOfFile]

		for i, v := range fileSizes.EmptySpaces {
			if file.IndexLeft < v.IndexLeft {
				break
			}

			if file.Length <= v.Length {
				fileSizes.EmptySpaces[i].Length -= file.Length

				leftPos := v.IndexLeft
				leftVal := file.IndexLeft

				for ind := 0; ind < file.Length; ind++ {
					normalized[ind+leftPos] = indOfFile
					normalized[ind+leftVal] = FREE_SPACE_SYMBOL
				}

				fileSizes.EmptySpaces[i].IndexLeft += file.Length
				break
			}
		}
	}

	fileSizes.DiskLine = normalized
	return fileSizes

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
		convertedString, err := convertToString(v)
		if err != nil {
			fmt.Print(err)
			return
		}

		normalizedFileSizes := normalizeSpace(convertedString)

		checkSum, err := getCheckSum(normalizedFileSizes.DiskLine)
		if err != nil {
			fmt.Println(err)
		}

		sum += checkSum
	}

	fmt.Println(BLUE, "Output: ", sum)
}
