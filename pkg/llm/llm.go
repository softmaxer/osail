package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"

	ollama "github.com/ollama/ollama/api"
)

type Llm interface {
	GetUrl() (*url.URL, error)
	GetCompletion(string) ollama.GenerateResponse
}

type Model struct {
	Name   string  `json:"name"`
	Host   string  `json:"host"`
	Port   int     `json:"port"`
	Rating float64 `json:"rating"`
	Stream bool    `json:"stream"`
}

type Download struct {
	Name string `json:"name"`
}

func (model *Model) expectedScore(opponent *Model) float64 {
	ratingDiff := (opponent.Rating - model.Rating) / 400
	expScore := 1 / (1 + math.Pow(10.0, float64(ratingDiff)))
	return expScore
}

func (model *Model) UpdateRating(opponent *Model, actualScore float64) {
	expScore := model.expectedScore(opponent)
	model.Rating += 32 * (actualScore - expScore)
	log.Printf("Updated rating to: %.2f", model.Rating)
}

func (model *Model) GetUrl() (*url.URL, error) {
	stringify := fmt.Sprintf("%s:%d", model.Host, model.Port)
	url, err := url.Parse(stringify)
	return url, err
}

func (model *Model) GetCompletion(prompt string) ollama.GenerateResponse {
	genRequest := ollama.GenerateRequest{Model: model.Name, Prompt: prompt, Stream: &model.Stream}
	requestBody, err := json.Marshal(genRequest)
	if err != nil {
		log.Printf("Error marshalling Chat request")
	}

	baseUrl, err := model.GetUrl()
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

func ParseJSON(response ollama.GenerateResponse, obj any) error {
	if response.Response == "" {
		return errors.New("Empty response")
	}
	start := strings.Index(response.Response, "{")
	end := strings.LastIndex(response.Response, "}") + 1
	jsonPart := response.Response[start:end]
	err := json.Unmarshal([]byte(jsonPart), &obj)
	return err
}
