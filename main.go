package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Response struct {
	Message string `json:"message"`
}

var (
	scoreFrequency = make(map[int]int)
	totalScores    int
)

func submitScore(userScore int) {
	mu.Lock()
	defer mu.Unlock()

	scoreFrequency[userScore]++
	totalScores++
}

func calculatePercentile(userScore int) int {
	mu.Lock()
	defer mu.Unlock()

	countLess := 0

	for score, freq := range scoreFrequency {
		if score <= userScore {
			countLess += freq
		}
	}

	if totalScores > 0 {
		percentile := (float64(countLess) / float64(totalScores)) * 100
		return int(percentile)
	}

	return 0
}

func main() {
	e := echo.New()

	e.Use(middleware.Recover())
	log.SetLevel(log.DebugLevel)

	t := &Template{
		templates: template.Must(template.ParseGlob("*.html")),
	}
	e.Renderer = t

	//Route for displaying quiz
	e.GET("/", func(c echo.Context) error {
		userID := CreateUser()
		fmt.Printf("Created user with ID: %d\n", userID)

		questions, err := getQuestions()
		if err != nil {
			log.Error("Error retrieving questions:", err)
			return c.String(http.StatusInternalServerError, "Failed to retrieve questions.")
		}

		data := map[string]interface{}{
			"UserID":    userID,
			"Questions": questions,
		}

		err = c.Render(http.StatusOK, "quiz.html", data)
		if err != nil {
			log.Error("Error rendering template:", err)
			return c.String(http.StatusInternalServerError, "Failed to render template.")
		}

		return nil
	})

	//form submission
	e.POST("/submit-answer", func(c echo.Context) error {
		formData, err := c.FormParams()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Error retrieving form data")
		}

		questions, _ := getQuestions()

		userID, _ := strconv.Atoi(c.FormValue("userID"))
		fmt.Print(userID)

		for key, values := range formData {
			if questions[key].Ans == values[0] {
				IncrementScore(userID, 1)
			}
		}

		score, _ := GetScore(userID)
		submitScore(score)
		betterThanPercentage := calculatePercentile(score)

		message := "Your score: " + strconv.Itoa(score) + " , You were better than " + strconv.Itoa(betterThanPercentage) + "% of all quizzers"

		response := Response{
			Message: message,
		}

		return c.JSON(http.StatusOK, response)

	})

	log.Fatal(e.Start(":5547"))
}
