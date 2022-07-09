package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

type task struct {
	taskID int
	dur    time.Duration
}

func processTask(tasks <-chan task, quit <-chan struct{}, done chan<- struct{}) {
	for {
		select {
		case task := <-tasks:
			log.Printf("start task %d\n", task.taskID)
			time.Sleep(task.dur)
			log.Printf("finish task %d\n", task.taskID)

			done <- struct{}{}
		case <-quit:
			return
		}
	}
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

	var maxProc int
	fmt.Println("Enter max proc count:")
	if _, err := fmt.Scan(&maxProc); err != nil {
		log.Fatal("can't read from stdin")
	}

	runtime.GOMAXPROCS(maxProc)

	tasks := make(chan task, maxProc)
	quit := make(chan struct{}, maxProc)
	done := make(chan struct{}, maxProc)

	for i := 0; i < maxProc; i++ {
		go processTask(tasks, quit, done)
	}

	scanner := bufio.NewScanner(file)
	taskCount := 0

	for scanner.Scan() {
		dur, err := time.ParseDuration(scanner.Text())
		if err != nil {
			log.Printf("error in parse task %d\n", taskCount)
			continue
		}
		taskCount++

		tasks <- task{
			taskID: taskCount,
			dur:    dur,
		}
	}

	for i := 0; i < taskCount; i++ {
		<-done
	}

	for i := 0; i < maxProc; i++ {
		quit <- struct{}{}
	}

}
