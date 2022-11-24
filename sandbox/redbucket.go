package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	LANGVAR = "BLEVELANG"
	CONFIG  = "./sandbox/config/redconfig.json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	curDir, _ := os.Executable()
	fmt.Printf("current directory %s\n", curDir)
	file, err := os.Open(CONFIG)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	check(err)
}
