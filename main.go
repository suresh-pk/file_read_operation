package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

const (
	filePath = "d:/test.txt"
	gap      = 5
)

func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

func (p *program) run() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	lines := make(chan string, 100) // Buffered channel to prevent blocking
	done := make(chan struct{})

	// Reading the lines
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines <- scanner.Text()
		}
		close(lines)
	}()

	// Printing the lines
	go func() {
		defer close(done)

		lineCount := 0
		for line := range lines {
			fmt.Println(line)
			lineCount++
			if lineCount == 5 {
				fmt.Println("<-- Next 5 Lines --->")
				lineCount = 0
				time.Sleep(gap * time.Second)
			}
		}
	}()
	<-done
}

func main() {
	svcConfig := &service.Config{
		Name:        "ReadMicroservice",
		DisplayName: "ReadMicroservice",
		Description: "ReadMicroservice",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println(err)
	}

	if len(os.Args) > 1 {
		err := service.Control(s, os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		fmt.Println(err)
	}
}
