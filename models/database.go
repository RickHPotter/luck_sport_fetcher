package models

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"example.com/m/style"
)

func SetUpDatabase() *sql.DB {
	execPath, err := os.Executable()
	if err != nil {
		panic("Falha ao achar o diretório do executável. Processo finalizado.")
	}

	execDir := filepath.Dir(execPath)
	dbPath := filepath.Join(execDir, "odds.sqlite")

	database, err := sql.Open("sqlite", dbPath)
	if err != nil {
		panic("Falha ao abrir o banco de dados. Processo finalizado.")
	}

	database.Exec("PRAGMA journal_mode=WAL;")

	if _, err = database.Exec(schemaSQL); err != nil {
		panic("Erro ao criar o banco de dados do Odds. Processo finalizado:" + err.Error())
	}

	return database
}

func CloseDatabase(database *sql.DB) {
	database.Close()
}

func RunQuery(db *sql.DB, query string) {
	_, err := db.Exec(query)
	if err != nil {
		panic("Erro ao executar a query. Processo finalizado.")
	}
}

func FetchMatchDateCounts(db *sql.DB, startDate string, endDate string) {
	startDate = strings.ReplaceAll(startDate, "-", "")
	endDate = strings.ReplaceAll(endDate, "-", "")

	query := `
  SELECT substr(MatchDate, 7, 4) || '-' || substr(MatchDate, 4, 2) || '-' || substr(MatchDate, 1, 2) As MatchDate, COUNT(*)
  FROM odds
  WHERE CAST(substr(MatchDate, 7, 4) || substr(MatchDate, 4, 2) || substr(MatchDate, 1, 2) AS INTEGER) BETWEEN ` + startDate + ` AND ` + endDate + `
  GROUP BY MatchDate;`

	rows, err := db.Query(query)
	if err != nil {
		panic("Erro ao buscar os dados. Processo finalizado.")
	}
	defer rows.Close()

	dateCounts := make(map[string]int)

	for rows.Next() {
		var date string
		var count int
		if err := rows.Scan(&date, &count); err != nil {
			panic("Erro ao buscar os dados. Processo finalizado.")
		}

		dateCounts[date] = count
	}

	DateCountsMap = dateCounts
}

func SaveToDatabase(database *sql.DB, date string, records []Record) (string, error) {
	resultMessage := ""

	tx, err := database.Begin()
	if err != nil {
		TheOnesThatGotAwayDB = append(TheOnesThatGotAwayDB, date)
		return "", fmt.Errorf(resultMessage+style.Colour("erro ao incializar o banco de dados, erro: %v", style.Red, style.Bold), err)
	}

	if count, exists := DateCountsMap[date]; exists {
		if len(records) == count {
			resultMessage += style.Colour("Dados já haviam sido salvos no Banco de Dados!", style.Blue, style.Bold)
			return resultMessage, nil
		} else if count > 0 {
			database.Exec("DELETE FROM odds WHERE MatchDate = ?", date)
		}
	}

	stmt, err := tx.Prepare(insertSQL)
	if err != nil {
		TheOnesThatGotAwayDB = append(TheOnesThatGotAwayDB, date)
		return "", fmt.Errorf(resultMessage+style.Colour("erro ao preparar inserção de jogo no banco de dados, erro: %v", style.Red, style.Bold), err)
	}

	defer stmt.Close()

	for _, record := range records {
		_, err = stmt.Exec(record.MatchDate, record.League, record.LeagueName, record.HomeTeam, record.AwayTeam, record.EarlyOdds1, record.EarlyOddsX, record.EarlyOdds2,
			record.FinalOdds1, record.FinalOddsX, record.FinalOdds2, record.Score)
		if err != nil {
			TheOnesThatGotAwayDB = append(TheOnesThatGotAwayDB, date)
			return "", fmt.Errorf(resultMessage+style.Colour("erro ao inserir jogo no banco de dados, erro: %v", style.Red, style.Bold), err)
		}
	}

	resultMessage += style.Colour("Jogos salvos no Banco de Dados com sucesso!", style.Green)

	return resultMessage, tx.Commit()
}
