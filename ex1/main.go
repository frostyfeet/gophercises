package main

import (
	"strings"
	"flag"
	"os"
	"encoding/csv"
	"bufio"
	"fmt"
	"log"
	"time"
)


type Contest struct {
	question[] string
	result[] string
	counter int 
}

func main(){
	csvFlag := flag.String ("csv", "problems.csv", "CSV file") 
	timeFlag := flag.Int ("limit", 5, "Time limit (default 30")
	flag.Parse()

	c, err := loadCSV(*csvFlag)
	if err != nil {
		print (err)
		return
	}

	startTest(c, timeFlag)
}

func startTest (c *Contest, limit *int){
	timer:= time.NewTimer(time.Duration(*limit) * time.Second)
	answerCh := make(chan string)
	fmt.Println("Starting test:")
	for i := 0; i < len(c.question); i++ {
		fmt.Printf("%s = ", c.question[i])
		go func(){
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()
		select {
		case <-timer.C: 
			fmt.Printf("\nCorrect responses: %v\n", c.counter)
			fmt.Printf("Incorrect responses %v\n", len(c.question) - c.counter)
			return
		case guess := <-answerCh:
			if sanitize(guess) == c.result[i]{
				c.counter++
			}
		}
	}	
	fmt.Printf("Correct responses: %v\n", c.counter)
	fmt.Printf("Incorrect responses %v\n", len(c.question) - c.counter)
}

func sanitize(s string) string{
	return strings.ToLower(strings.TrimSpace(s))
}

func loadCSV(path string) (*Contest, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer f.Close()

	c := Contest{nil, nil, 0}
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		if (len(record) == 2){
			c.question = append(c.question, record[0])
			c.result = append(c.result, record[1])
		}
	}
	return &c, nil
}
