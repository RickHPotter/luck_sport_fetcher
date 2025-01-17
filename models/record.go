package models

type Record struct {
	MatchDate  string `json:"matchDate"`
	League     string `json:"league"`
	HomeTeam   string `json:"homeTeam"`
	AwayTeam   string `json:"awayTeam"`
	EarlyOdds1 string `json:"earlyOdds1"`
	EarlyOddsX string `json:"earlyOddsX"`
	EarlyOdds2 string `json:"earlyOdds2"`
	FinalOdds1 string `json:"finalOdds1"`
	FinalOddsX string `json:"finalOddsX"`
	FinalOdds2 string `json:"finalOdds2"`
	Score      string `json:"score"`
}

var (
	TheOnesThatGotAwayAPI  = []string{}
	TheOnesThatGotAwayJSON = []string{}
	TheOnesThatGotAwayDB   = []string{}
	DateCountsMap          = map[string]int{}
)
