package main

import (
	"encoding/json"
	"log"
)

type Question struct {
	Ques string   `json:"ques"`
	Opts []string `json:"opts"`
	Ans  string   `json:"ans"`
}

type Questions map[string]Question

func getQuestions() (Questions, error) {
	//correct answer is the index of option opts , this json can be fetched from an api
	questionsJSON := `{
		"q1": {"ques":"What is dummy Ques 1", "opts": ["a", "b", "c"], "ans": "1"},
		"q2": {"ques":"What is dummy Ques 2", "opts": ["a", "b", "c"], "ans": "2"},
		"q3": {"ques":"What is dummy Ques 3", "opts": ["a", "b", "c"], "ans": "0"}
	}`

	var questions Questions

	if err := json.Unmarshal([]byte(questionsJSON), &questions); err != nil {
		log.Printf("Error unmarshaling JSON: %v", err)
		return nil, err
	}

	return questions, nil

}
