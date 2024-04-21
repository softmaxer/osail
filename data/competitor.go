package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"

	ollama "github.com/ollama/ollama/api"
	"gorm.io/gorm"
)

type Competitor struct {
	gorm.Model
	Name         string `gorm:"primaryKey"`
	Host         string
	Port         int32
	Stream       bool
	Rating       float64
	ExperimentId string
}

type CompetitorRequest struct {
	Name string `form:"name"`
	Host string `form:"host"`
	Port string `form:"port"`
}

func (c *Competitor) GetUrl() (*url.URL, error) {
	stringify := fmt.Sprintf("%s:%d", c.Host, c.Port)
	url, err := url.Parse(stringify)
	return url, err
}

func (c *Competitor) GetCompletion(prompt string) ollama.GenerateResponse {
	genRequest := ollama.GenerateRequest{Model: c.Name, Prompt: prompt, Stream: &c.Stream}
	requestBody, err := json.Marshal(genRequest)
	if err != nil {
		log.Printf("Error marshalling Chat request")
	}

	baseUrl, err := c.GetUrl()
	if err != nil {
		log.Println(err.Error())
	}

	response, err := http.Post(
		baseUrl.String()+"/api/generate",
		"application/json",
		bytes.NewBuffer(requestBody),
	)

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err.Error())
	}
	defer response.Body.Close()
	var genereationResponse ollama.GenerateResponse
	err = json.Unmarshal(responseData, &genereationResponse)
	if err != nil {
		log.Println(err.Error())
	}

	return genereationResponse
}

func (c *Competitor) expectedScore(opponent *Competitor) float64 {
	ratingDiff := (opponent.Rating - c.Rating) / 400
	expScore := 1 / (1 + math.Pow(10.0, float64(ratingDiff)))
	return expScore
}

func (c *Competitor) UpdateRating(opponent *Competitor, actualScore float64) {
	expScore := c.expectedScore(opponent)
	c.Rating += 32 * (actualScore - expScore)
	log.Printf("Updated rating to: %.2f", c.Rating)
}
