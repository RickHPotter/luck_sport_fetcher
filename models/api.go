package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"example.com/m/style"
)

func FetchData(date string) ([]Record, string, error) {
	resultMessage := fmt.Sprintf("%s == ", date)

	url := fmt.Sprintf("https://px-1x2.7mdt.com/data/history/en/%s/17.js?nocache=%d", date, time.Now().Unix())
	resp, err := http.Get(url)
	if err != nil {
		TheOnesThatGotAwayAPI = append(TheOnesThatGotAwayAPI, date)
		return nil, "", fmt.Errorf(resultMessage+style.Colour("falha ao acessar API. Erro: %v", style.Red, style.Bold), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		TheOnesThatGotAwayAPI = append(TheOnesThatGotAwayAPI, date)
		return nil, "", fmt.Errorf(resultMessage+style.Colour("falha ao trazer dados. Status: %d", style.Red, style.Bold), resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		TheOnesThatGotAwayJSON = append(TheOnesThatGotAwayJSON, date)
		return nil, "", fmt.Errorf(resultMessage+style.Colour("falha ao ler os dados. Erro: %v", style.Red, style.Bold), err)
	}

	bodyStr := sanitizeData(string(body))

	var recordsArray []string
	err = json.Unmarshal([]byte(bodyStr), &recordsArray)
	if err != nil {
		bodyStr = sanitizeUffef(bodyStr)
		err = json.Unmarshal([]byte(bodyStr), &recordsArray)
		if err != nil {
			TheOnesThatGotAwayJSON = append(TheOnesThatGotAwayJSON, date)
			return nil, "", fmt.Errorf(resultMessage+style.Colour("falha ao sanitizar dados. Erro: %v", style.Red, style.Bold), err)
		}
	}

	resultMessage += fmt.Sprintf("%4d JOGOS ENCONTRADOS ", len(recordsArray))

	var records []Record
	for _, recordData := range recordsArray {
		recordData = strings.Trim(recordData, "\"")

		fields := strings.Split(recordData, "|")

		if len(fields) < 17 {
			TheOnesThatGotAwayJSON = append(TheOnesThatGotAwayJSON, date)
			return nil, "", fmt.Errorf(resultMessage + "jogo inválido")
		}

		date, err := time.Parse("2006-01-02", date)
		if err != nil {
			date = time.Now()
		}

		score := fields[8] + "-" + fields[9] + "\n" + "(" + fields[10] + ")"
		record := Record{
			MatchDate:  date.Format("02/01/2006"),
			League:     fields[3],
			HomeTeam:   fields[6],
			AwayTeam:   fields[7],
			EarlyOdds1: fields[11],
			EarlyOddsX: fields[12],
			EarlyOdds2: fields[13],
			FinalOdds1: fields[14],
			FinalOddsX: fields[15],
			FinalOdds2: fields[16],
			Score:      score,
		}

		records = append(records, record)
	}

	return records, resultMessage, nil
}

func sanitizeData(rawData string) string {
	sanitized := strings.TrimPrefix(rawData, "\xef\xbb\xbf") // UTF-8

	sanitized = strings.Replace(sanitized, "var dt = ", "", 1)
	sanitized = strings.TrimSuffix(sanitized, ";")
	sanitized = strings.ReplaceAll(sanitized, "\\'", "'")
	sanitized = strings.ReplaceAll(sanitized, "\t", " ")

	return sanitized
}

func sanitizeUffef(rawData string) string {
	return strings.TrimPrefix(rawData, "\ufeff")
}
