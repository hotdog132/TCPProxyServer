package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	//Enter your code here. Read input from STDIN. Print output to STDOUT
	L := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		LStr := scanner.Text()
		L = strings.Split(LStr, ",")
	} else {
		fmt.Println(0)
		return
	}
	maxStreak := 0
	currentStreak := 0
	for _, value := range L {
		value = strings.Replace(value, " ", "", -1)
		if v, err := strconv.Atoi(value); err == nil {
			currentStreak += v
			if currentStreak > maxStreak {
				maxStreak = currentStreak
			}

			if currentStreak < 0 {
				currentStreak = 0
			}
		} else {
			fmt.Println(0)
			return
		}
	}

	fmt.Println(maxStreak)
}
