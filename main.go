package main

import (
	"os"
	"bufio"
	"fmt"
	"bytes"
	"time"
)

func main() {
	args := os.Args[1:]

	initDic(args[0])
	fmt.Printf("dic:\t%v;\t Nums:%d \n", args[0], len(dic))
	fmt.Printf("chars:\t%v \n", args[1])

	sop := make(chan string)
	pop := make(chan []string)
	go search(args[1], sop)
	go packStrings(sop, pop, 10)

	counter := 0
	for pack := range pop {

		fmt.Println("")
		for i, str := range pack {
			fmt.Printf("\t[%d,\t%d]:\t%s\n", counter, i, str)
			counter ++
		}
		time.Sleep(time.Second * 18)
	}

	fmt.Printf("Found: %d", counter)
}


var dic map[string]bool

func initDic(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	dic = make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		dic[scanner.Text()] = true
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func search(chars string, output chan string) {
	visited := make(map[string]bool)
	max := len(chars) - 1
	cursor := []int{0, 0, 0}
	for len(cursor) <= len(chars) {
		cursor = next(cursor, max, false)
		str := toString(cursor, chars)

		if _, ok := visited[str]; ok {
			continue
		}
		visited[str] = true

		if _, ok := dic[str]; ok {
			output <- str
		}
	}

	close(output)
}

func next(cursor []int, max int, allowRepeat bool) []int {
	from := 0;
	for {
		cursor = increment(cursor, max, from)
		if allowRepeat || len(cursor) > max {
			return cursor
		}

		rpi := repeatIndex(cursor)
		if rpi < 0 {
			return cursor
		}
		from = rpi
	}
}

func increment(cursor []int, max int, from int) []int {
	for {
		if len(cursor) <= from {
			cursor = append(cursor, 0)
			return cursor
		}

		cursor[from] ++;
		if cursor[from] > max {
			cursor[from] = 0;
			from ++
		} else {
			break
		}
	}
	return cursor
}

func repeatIndex(cursor []int) int {
	visited := make(map[int]int)
	for i, cur := range cursor {
		if r, ok := visited[cur]; ok {
			return r
		} else {
			visited[cur] = i
		}
	}
	return -1
}

func toString(cursor []int, chars string) string {
	buf := bytes.Buffer{}

	for _, cur := range cursor {
		buf.WriteByte(chars[cur])
	}
	return string(buf.Bytes())
}

func packStrings(input chan string, output chan []string, size int) {
	pack := make([]string, 0, size)
	for {
		str, ok := <- input
		if ok {
			pack = append(pack, str)
		}
		if !ok || len(pack) == size {
			output <- pack
			pack = make([]string, 0, size)
		}
		if !ok {
			break
		}
	}
	close(output)
}