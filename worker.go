package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func checkStatus(url string, ch chan<- string) {
	response, err := http.Get(url)
	if err != nil {
		ch <- url + " " + err.Error()
		return
	}
	defer response.Body.Close()
	ch <- url + " " + response.Status
}

func main() {
	readFile, err := os.Open("data.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close()
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	ch := make(chan string)

	urls := 0
	for fileScanner.Scan() {
		go checkStatus(fileScanner.Text(), ch)
		urls++
	}
	for i := 0; i <= urls; i++ {
		fmt.Println(<-ch)
	}

}
