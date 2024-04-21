package board

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"

	"github.com/softmaxer/osail/data"
	"github.com/softmaxer/osail/pkg/llm"
)

type Board struct {
	Competitors []data.Competitor
}

func (board *Board) GenerateCompetitors() (*data.Competitor, *data.Competitor) {
	if len(board.Competitors) == 2 {
		return &board.Competitors[0], &board.Competitors[1]
	}
	player1 := rand.Intn(len(board.Competitors))
	player2 := rand.Intn(len(board.Competitors))
	for player2 == player1 {
		player2 = rand.Intn(len(board.Competitors))
	}

	return &board.Competitors[player1], &board.Competitors[player2]
}

func (board *Board) Init() {
	for _, competitor := range board.Competitors {

		body, err := json.Marshal(llm.Download{Name: competitor.Name})
		if err != nil {
			log.Fatal("Error marshalling json: ", err.Error())
		}
		competitorUrl, err := competitor.GetUrl()
		if err != nil {
			log.Printf("Error parsing Ollama client url: %s\n", competitorUrl)
		}

		response, err := http.Post(
			competitorUrl.String()+"/api/pull",
			"application/json",
			bytes.NewBuffer(body),
		)
		if err != nil {
			log.Fatal("Couldn't sent request")
		}
		defer response.Body.Close()

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf(err.Error())
		}
		log.Printf(string(responseData))
	}
}

func (board *Board) SortRatings() {
	for i := range len(board.Competitors) {
		for j := 1; j < i; j++ {
			if board.Competitors[j-1].Rating < board.Competitors[j].Rating {
				intermediate := board.Competitors[j]
				board.Competitors[j] = board.Competitors[j-1]
				board.Competitors[j-1] = intermediate
			}
		}
	}
}
