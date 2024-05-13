package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func readInput(filename string) [][]string {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Dosya açılamadı:", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var groups [][]string
	var group []string

	space := " "

	group = append(group, space)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			group = append(group, line)
		} else {
			groups = append(groups, group)
			group = nil
		}
	}
	groups = append(groups, group)

	if err := scanner.Err(); err != nil {
		fmt.Println("Dosya okunurken bir hata oluştu:", err)
		os.Exit(1)
	}

	return groups
}

func FindLength(s string) int {
	return len(s)
}

func PrintAscii(runeArr []int, groups [][]string, space int, align string) {
	var lenArr int
	spaceCount := 0

	for _, val := range runeArr {
		if val == 32 {
			lenArr -= 6
			spaceCount++
		}

		lenArr += FindLength(groups[val-31][0])
	}

	for index := 0; index < 8; index++ {

		if align == "right" {
			fmt.Print(strings.Repeat(" ", (space - lenArr)))
		}
		if align == "center" {
			fmt.Print(strings.Repeat(" ", (space-lenArr)/2))
		}

		if align == "justify" && spaceCount == 0 {
			fmt.Print(strings.Repeat(" ", (space-lenArr)/2))
		}

		for i, val := range runeArr {

			fmt.Print(groups[val-31][index])

			if val == 32 {
				fmt.Print(" ")
				if align == "justify" && i == 0 {
				} else if align == "justify" && i == len(runeArr)-1 {
					// val align right
					fmt.Print(strings.Repeat(" ", (space - lenArr)))
				} else if align == "justify" && spaceCount == 0 {
					fmt.Print(strings.Repeat(" ", ((space - lenArr) / 2)))
				} else if align == "justify" {
					fmt.Print(strings.Repeat(" ", ((space-lenArr)/spaceCount)-8))
				}

			}

		}

		fmt.Println()
	}
}

func stringReplace(input string) string {
	// \\n karakterini \n olarak değiştir
	replaced := strings.Replace(input, "\\n", "\n", -1)

	// Eğer başında 'n' karakteri varsa, bu bir kaçış karakteridir, bunu gerçek '\n' karakterine dönüştür
	if strings.HasPrefix(replaced, "n") {
		replaced = replaced[1:]
	}

	// Eğer sonunda 'n' karakteri varsa, bu bir kaçış karakteridir, bunu gerçek '\n' karakterine dönüştür
	if strings.HasSuffix(replaced, "n") {
		replaced = replaced[:]
	}

	return replaced
}

func printTT(temp2StrArr [][]string, groups [][]string, space int, align string) {
	if len(temp2StrArr) > 0 {
		for _, val := range temp2StrArr {
			for _, val2 := range val {
				convertedRuneArr := []int{}
				for _, val3 := range val2 {
					convertedRuneArr = append(convertedRuneArr, int(val3))
				}
				if len(convertedRuneArr) > 0 {
					PrintAscii(convertedRuneArr, groups, space, align)
				}
				if len(convertedRuneArr) == 0 {
					fmt.Println()
				}

			}
		}
	}
}

func ArgsLength() {
	if len(os.Args) < 3 || len(os.Args) > 4 {

		fmt.Println("Usage: go run . [OPTION] [STRING] [BANNER]")
		fmt.Println("Example: go run . --align=right something standard")

		os.Exit(1)
	}
}

func checkFileName(fileName string) {
	if fileName != "standard" && fileName != "shadow" && fileName != "thinkertoy" {
		fmt.Println("Invalid banner name")
		fmt.Println("Usage: go run . [STRING] [BANNER]")
		fmt.Println("EX: go run . something standard")
		os.Exit(1)
	}
}

func getAlign(align *string) string {
	var alignedText string
	switch *align {
	case "center":
		alignedText = "center"
	case "left":
		alignedText = "left"
	case "right":
		alignedText = "right"
	case "justify":
		alignedText = "justify"
	default:
		return ""
	}
	return alignedText
}

func getTerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	width := strings.Split(string(out), " ")[1]

	widthInt, _ := strconv.Atoi(strings.TrimSpace(width))

	return widthInt
}

func main() {
	// Dosyadan okunan grupları al

	// groups := readInput("shadow.txt")
	ArgsLength()

	osArgs := os.Args

	var fileName string
	var input []string

	if len(osArgs) == 3 {
		fileName = "standard"
		input = osArgs[2:3]
	} else {
		fileName = osArgs[3]
		checkFileName(fileName)

		input = osArgs[2:3]
	}

	// if args has \n character, show it as a slash

	align := flag.String("align", "left", "Alignment type (center, left, right, justify)")

	flag.Parse()

	groups := readInput(fileName + ".txt")

	temStrArr := [][]string{}
	temp := ""

	if len(input) > 0 {
		for _, val := range input {
			for index, runeVal := range val {
				if runeVal == '\\' && index < len(val)-1 && val[index+1] == 'n' {
					temStrArr = append(temStrArr, []string{temp})
					temp = ""
				} else {
					temp += string(runeVal)
				}

				if index == len(val)-1 {
					temStrArr = append(temStrArr, []string{temp})
					temp = ""
				}
			}
		}
	}

	temp2StrArr := [][]string{}

	for _, val := range temStrArr {
		for _, val2 := range val {
			temp2StrArr = append(temp2StrArr, []string{stringReplace(val2)})
		}
	}

	// Uygun hizalamayı seç
	space := getTerminalWidth()

	alignedText := getAlign(align)

	printTT(temp2StrArr, groups, space, alignedText)
	// Metni ekrana yazdır
	// fmt.Println(alignedText)
}
