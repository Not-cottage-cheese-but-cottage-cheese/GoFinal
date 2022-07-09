package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/Not-cottage-cheese-but-cottage-cheese/final-go/server"
)

const workers = 10

func TestTime() {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, err := http.Post(
				"http://localhost:3000/add?mod=async",
				"text/plain",
				bytes.NewBuffer([]byte("1s")),
			)

			if err != nil {
				log.Panic(err)
			}

		}()
	}

	wg.Wait()

	resp, err := http.Get("http://localhost:3000/time")
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	expected := (time.Second * (workers - 1)).String()
	if !reflect.DeepEqual(expected, string(body)) {
		log.Panicf("time test: expected %s, actual: %s", expected, string(body))
	}
}

func TestSchedule() {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, err := http.Post(
				"http://localhost:3000/add?mod=async",
				"text/plain",
				bytes.NewBuffer([]byte("1s")),
			)

			if err != nil {
				log.Panic(err)
			}

		}()
	}

	wg.Wait()

	resp, err := http.Get("http://localhost:3000/schedule")
	if err != nil {
		log.Panic(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	var tasks []server.Task
	if err := json.Unmarshal(body, &tasks); err != nil {
		log.Panic(err)
	}

	expectedTask := make([]server.Task, workers-1)
	for i := range expectedTask {
		expectedTask[i] = server.Task{
			Dur: time.Second,
		}
	}

	if !reflect.DeepEqual(expectedTask, tasks) {
		log.Panicf("schedule test: expected %v, actual: %v", expectedTask, tasks)
	}
}

func TestSync() {
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			_, err := http.Post(
				"http://localhost:3000/add?mod=async",
				"text/plain",
				bytes.NewBuffer([]byte("1s")),
			)

			if err != nil {
				log.Panic(err)
			}

		}()
	}

	wg.Wait()

	start := time.Now()
	_, err := http.Post(
		"http://localhost:3000/add?mod=sync",
		"text/plain",
		bytes.NewBuffer([]byte("0s")),
	)
	if err != nil {
		log.Panic(err)
	}

	elapsed := time.Since(start)

	eps := float64(time.Second*workers) * 0.1

	if !(math.Abs(float64(elapsed-time.Second*workers)) < eps) {
		log.Panicf("sync test: expected: %s +- %.6fs, actual: %v", (time.Second * workers).String(), eps, elapsed)
	}
}

func main() {
	tests := []func(){
		TestTime,
		TestSchedule,
		TestSync,
	}

	for i, test := range tests {
		log.Printf("starting test %d", i+1)
		test()
		log.Printf("test %d done\n", i+1)
		if i != len(tests)-1 {
			log.Printf("sleeping until next test\n")
			time.Sleep(2 * workers * time.Second)
		}
	}
}
