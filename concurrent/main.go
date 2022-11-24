package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	DURATION   = 10 * time.Second
	OUTPUTRATE = 1000 * time.Millisecond
	FILE       = "./concurrent/random.txt"
)

func main() {
	fmt.Println("Hi")
	ctx, cancel := context.WithTimeout(context.Background(), DURATION)
	defer cancel()
	content := make(chan string)
	var wg sync.WaitGroup

	wg.Add(1)
	go contentGenerator(ctx, &wg, content)
	wg.Add(1)
	go writeFile(ctx, &wg, FILE, content)
	deleteFile(FILE)
	wg.Wait()

}

func contentGenerator(ctx context.Context, wg *sync.WaitGroup, content chan<- string) {
	defer wg.Done()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	show := func(name string, num int, v1, v2, v3 any) string {
		return fmt.Sprintf("%d\t%s\t%v\t%v\t%v\n", num, name, v1, v2, v3)
	}

	lineNumber := 1
	for {
		select {
		case <-time.Tick(OUTPUTRATE):
			content <- show("Int", lineNumber, r.Int(), r.Int(), r.Int())
			lineNumber++
		case <-ctx.Done():
			fmt.Println("contentGenerator: ", ctx.Err().Error())
			close(content)
			return
		}
	}

}

func writeFile(ctx context.Context, wg *sync.WaitGroup, filename string, content <-chan string) {
	defer wg.Done()
	file, err := os.Create(filename)
	check(err)
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("file.Close(): %s", err.Error())
		}
	}()

	writer := bufio.NewWriter(file)
	defer func() {
		if err := writer.Flush(); err != nil {
			fmt.Printf("writer.Flush(): %s", err.Error())
		}
	}()

	for {
		select {
		case message := <-content:
			fmt.Print(message)
			_, err := writer.WriteString(message)
			check(err)
			if _, err := file.Stat(); err == nil {
				fmt.Printf("%s exits\n", file.Name())
			} else if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("%s doesn't exit\n", file.Name())
			} else {
				fmt.Printf("%s\n", err.Error())
			}
		case <-ctx.Done():
			fmt.Println("writeFile: ", ctx.Err().Error())
			return
		}
	}
}

func deleteFile(filename string) error {
	<-time.Tick(5 * time.Second)
	if err := os.Remove(filename); err != nil {
		fmt.Printf("deleteFile %s", err.Error())
		return err
	}
	fmt.Printf("file %s deleted\n", filename)
	return nil

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
