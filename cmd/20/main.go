package main

import (
	"bufio"
	"log"
	"os"
	"sync"
	"time"
)

func processTask(taskID int, dur time.Duration) {
	log.Printf("start task %d\n", taskID)
	time.Sleep(dur)
	log.Printf("finish task %d\n", taskID)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("invlaid argument count")
	}
	fileName := os.Args[1]

	if _, err := os.Stat(fileName); err != nil {
		log.Fatal("file does not exist or somethig else wrong")
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("file open error: %s", err.Error())
	}

	scanner := bufio.NewScanner(file)
	taskID := 0

	var wg sync.WaitGroup
	for scanner.Scan() {
		taskID++
		dur, err := time.ParseDuration(scanner.Text())
		if err != nil {
			log.Printf("error in parse task %d\n", taskID)
			continue
		}

		wg.Add(1)
		go func(ID int) {
			processTask(ID, dur)
			wg.Done()
		}(taskID)
	}

	wg.Wait()
}
