package main

import (
	"sync"
)

var (
	uid       int
	userScore = make(map[int]int)
	mu        sync.Mutex
)

func init() {
	uid = 1
}

func CreateUser() int {
	mu.Lock()
	defer mu.Unlock()

	userID := uid
	uid++

	userScore[userID] = 0
	return userID
}

func GetScore(userID int) (int, bool) {
	mu.Lock()
	defer mu.Unlock()

	score, exists := userScore[userID]
	return score, exists
}

func IncrementScore(userID, increment int) bool {
	mu.Lock()
	defer mu.Unlock()

	userScore[userID] += increment
	return true
}
