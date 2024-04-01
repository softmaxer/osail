package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	ollama "github.com/ollama/ollama/api"
)

type Judge struct {
	Name           string `json:"name"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	Stream         bool   `json:"stream"`
	systemMessage  string
	promptTemplate string
}

type Judgement struct {
	Result [2]int `json:"result"`
}

func (j *Judge) setSystemMessage() {
	j.systemMessage = `You are tasked to judge two responses from two different LLMs. There are no draws. It is considered a good practice to always respond in JSON format. The responses will be presented to you in the following format:
  Response A: ...
  Response B: ...
  If you think response A is better, you should give:
  {
    "result": [1, 0]
  }

  or if you think response B is better, you should give: 
  {
    "result: [0, 1]"
  }`
}

func (j *Judge) setPromptTempalte() {
	j.promptTemplate = `Here are two responses from two LLMs, tell me which one is better:`
}

func (j *Judge) GetUrl() (*url.URL, error) {
	stringify := fmt.Sprintf("%s:%d", j.Host, j.Port)
	url, err := url.Parse(stringify)
	return url, err
}

func (j *Judge) Init() {
	body, err := json.Marshal(Download{Name: j.Name})
	if err != nil {
		log.Fatal("Error marshalling JSON: ", err.Error())
	}

	judgeUrl, err := j.GetUrl()
	if err != nil {
		log.Printf("Error parsing Ollama url: %s\n", judgeUrl)
	}

	response, err := http.Post(
		judgeUrl.String()+"/api/pull",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		log.Fatal("Error sending request: ", err.Error())
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf(err.Error())
	}
	log.Printf(string(responseData))
}

func (j *Judge) JudgePrompt(firstResponse string, secondResponse string) ollama.GenerateResponse {
	j.setSystemMessage()
	j.setPromptTempalte()
	prompt := fmt.Sprintf(
		"%s\nResponse A: %s\nResponse B: %s\n",
		j.promptTemplate,
		firstResponse,
		secondResponse,
	)

	response := j.GetCompletion(prompt)
	return response
}

func (j *Judge) GetCompletion(prompt string) ollama.GenerateResponse {
	genRequest := ollama.GenerateRequest{Model: j.Name, Prompt: prompt, Stream: &j.Stream}
	genRequest.System = j.systemMessage
	requestBody, err := json.Marshal(genRequest)
	if err != nil {
		log.Printf("Error marshalling Chat request")
	}

	baseUrl, err := j.GetUrl()
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
