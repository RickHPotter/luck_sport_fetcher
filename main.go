package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"example.com/m/models"
	"example.com/m/style"
	_ "modernc.org/sqlite"
)

const (
	workerCount  = 3
	workerCount2 = 2
	workerCount3 = 1
)

func worker(dates <-chan string, database *sql.DB, wg *sync.WaitGroup) {
	defer wg.Done()

	for date := range dates {
		records, fetchResultMessage, err := models.FetchData(date)
		if err != nil {
			fmt.Println(err.Error())
			continue
		} else {
			if databaseResultMessage, err := models.SaveToDatabase(database, date, records); err != nil {
				fmt.Println(fetchResultMessage + err.Error())
			} else {
				fmt.Println(fetchResultMessage + databaseResultMessage)
			}
		}

	}
}

func main() {
	startDateStr := ""
	endDateStr := ""

	switch len(os.Args) {
	case 1:
		{
			fmt.Println("Argumentos de datas não fornecidos. Formato: yyyy-mm-dd")
			return
		}
	case 2:
		{
			database := models.SetUpDatabase()
			models.RunQuery(database, os.Args[1])
			models.CloseDatabase(database)
			return
		}
	case 3:
		{
			startDateStr = os.Args[1]
			endDateStr = os.Args[2]
		}
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		log.Fatalf("Data Invalida: %v\n", err)
	}

	endDate := startDate
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			log.Fatalf("Data Invalida: %v\n", err)
		}
	}

	style.CheckColourSupport()

	start := time.Now()

	database := models.SetUpDatabase()

	models.FetchMatchDateCounts(database, startDateStr, endDateStr)

	if err != nil {
		log.Fatal(err)
	}

	firstTryDates := make(chan string)
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(firstTryDates, database, &wg)
	}

	for currentDate := startDate; !currentDate.After(endDate); currentDate = currentDate.AddDate(0, 0, 1) {
		firstTryDates <- currentDate.Format("2006-01-02")
	}

	close(firstTryDates)
	wg.Wait()

	secondTryDates := make(chan string)

	for i := 0; i < workerCount2; i++ {
		wg.Add(1)
		go worker(secondTryDates, database, &wg)
	}

	fmt.Println(style.Colour("FALHAS DA PRIMEIRA TENTATIVA", style.Yellow, style.Bold))
	fmt.Println("Falhas de API: [" + strings.Join(models.TheOnesThatGotAwayAPI, ", ") + "]")
	fmt.Println("Falhas de JSON: [" + strings.Join(models.TheOnesThatGotAwayJSON, ", ") + "]")
	fmt.Println("Falhas de BANCO: [" + strings.Join(models.TheOnesThatGotAwayDB, ", ") + "]")

	if len(models.TheOnesThatGotAwayJSON) != 0 {
		fmt.Println(style.Colour("Reportando datas que houveram falhas de leitura ou sanitização de dados:", style.Red, style.Bold))
		for _, currentDate := range models.TheOnesThatGotAwayJSON {
			fmt.Print(currentDate + " ")
		}
		fmt.Println()
	}

	if len(models.TheOnesThatGotAwayAPI) != 0 {
		fmt.Println(style.Colour("Tentando novamente dias que houveram falhas:", style.Red, style.Bold))
		items := models.TheOnesThatGotAwayAPI
		models.TheOnesThatGotAwayAPI = []string{}
		for _, currentDate := range items {
			secondTryDates <- currentDate
		}
	}

	if len(models.TheOnesThatGotAwayDB) != 0 {
		fmt.Println(style.Colour("Tentando novamente dias que houveram falhas de gravação no Banco de Dados:", style.Red, style.Bold))
		items := models.TheOnesThatGotAwayDB
		models.TheOnesThatGotAwayDB = []string{}
		for _, currentDate := range items {
			secondTryDates <- currentDate
		}
	}

	close(secondTryDates)
	wg.Wait()

	if len(models.TheOnesThatGotAwayAPI) != 0 && len(models.TheOnesThatGotAwayDB) != 0 {
		thirdTryDates := make(chan string)

		for i := 0; i < workerCount3; i++ {
			wg.Add(1)
			go worker(thirdTryDates, database, &wg)
		}
		fmt.Println(style.Colour("FALHAS DA SEGUNDA TENTATIVA", style.Yellow, style.Bold))
		fmt.Println("Falhas de API: [" + strings.Join(models.TheOnesThatGotAwayAPI, ", ") + "]")
		fmt.Println("Falhas de JSON: [" + strings.Join(models.TheOnesThatGotAwayJSON, ", ") + "]")
		fmt.Println("Falhas de BANCO: [" + strings.Join(models.TheOnesThatGotAwayDB, ", ") + "]")

		if len(models.TheOnesThatGotAwayJSON) != 0 {
			fmt.Println(style.Colour("Reportando datas que houveram falhas de leitura ou sanitização de dados:", style.Red, style.Bold))
			for _, currentDate := range models.TheOnesThatGotAwayJSON {
				fmt.Print(currentDate + " ")
			}
			fmt.Println()
		}

		if len(models.TheOnesThatGotAwayAPI) != 0 {
			fmt.Println(style.Colour("Tentando novamente dias que houveram falhas:", style.Red, style.Bold))
			for _, currentDate := range models.TheOnesThatGotAwayAPI {
				thirdTryDates <- currentDate
			}
		}

		if len(models.TheOnesThatGotAwayDB) != 0 {
			fmt.Println(style.Colour("Tentando novamente dias que houveram falhas de gravação no Banco de Dados:", style.Red, style.Bold))
			for _, currentDate := range models.TheOnesThatGotAwayDB {
				thirdTryDates <- currentDate
			}
		}

		close(thirdTryDates)
		wg.Wait()
	}

	models.CloseDatabase(database)

	end := time.Now()

	duration := end.Sub(start)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	fmt.Printf(style.Colour("Processo finalizado em %02d:%02d:%02d (hh:mm:ss)\n", style.Green, style.Bold), hours, minutes, seconds)

	fmt.Println()
	fmt.Println(style.Colour("Pressione ENTER para fechar ou espere 60 segundos.", style.Italic, style.Bold))

	// Create a channel to wait for user input
	inputChan := make(chan struct{})

	// Start a goroutine to wait for user input
	go func() {
		bufio.NewReader(os.Stdin).ReadString('\n')
		inputChan <- struct{}{}
	}()

	// Wait for either user input or 5 seconds
	select {
	case <-inputChan:
		fmt.Println("Fechando agora...")
	case <-time.After(60 * time.Second):
		fmt.Println("Tempo esgotado. Fechando...")
	}
}
