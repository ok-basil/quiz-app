package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func problemPuller(filename string) ([]problem, error) {
	if fObj, err := os.Open(filename); err == nil {
		csvR := csv.NewReader(fObj)
		if cLines, err := csvR.ReadAll(); err == nil {
			return parseProblem(cLines), nil
		} else {
			return nil, fmt.Errorf("error reading data in CSV"+"format from %s file; %s", filename, err.Error())
		}
	} else {
		return nil, fmt.Errorf("error in opening %s file; %s", filename, err.Error())
	}
}

func main() {
	fName := flag.String("f", "quiz.csv", "File to read")

	timer := flag.Int("t", 30, "Timer for quiz")

	flag.Parse()

	problems, err := problemPuller(*fName)

	if err != nil {
		exit(fmt.Sprintf("Something went wrong: %s", err.Error()))
	}

	correctAnswer := 0

	tObj := time.NewTimer(time.Duration(*timer) * time.Second)
	ansC := make(chan string)

problemLoop:
	for i, p := range problems {
		var answer string
		fmt.Printf("Problem %d: %s=", i+1, p.question)
		go func() {
			n, err := fmt.Scanf("%s", &answer)
			if err != nil || n != 1 {
				answer = ""
				fmt.Println("Invalid input!")
				return
			}
			ansC <- answer
		}()
		select {
		case <-tObj.C:
			fmt.Println("Time's up!")
			break problemLoop
		case iAns := <-ansC:
			if iAns == p.answer {
				correctAnswer++
				fmt.Println("Correct!")
			}
			if i == len(problems)-1 {
				close(ansC)
			}
		}
	}

	fmt.Printf("You got %d out of %d correct!\n", correctAnswer, len(problems))
	fmt.Printf("Press enter to exit...")
	<-ansC
}

func parseProblem(lines [][]string) []problem {
	r := make([]problem, len(lines))
	for i := 0; i < len(lines); i++ {
		r[i] = problem{question: lines[i][0], answer: lines[i][1]}
	}
	return r
}

type problem struct {
	question string
	answer   string
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
