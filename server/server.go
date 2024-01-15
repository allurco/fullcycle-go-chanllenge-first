package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	USDBRL Usdbrl `json:"USDBRL"`
}

type Usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func main() {
	log.Println("Iniciando o Servidor")
	db, err := connectToDb()
	defer db.Close()
	if err != nil {
		log.Println(err)
		panic(err)
	}

	results, err := createTables(db)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	log.Println(results)
	log.Println("Servidor Iniciado")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println("Request Iniciada")
		defer log.Println("Request Finalizada")

		bid, err := performTask(ctx, db)
		if err != nil {
			log.Println(err)
			return
		}

		writeResult := fmt.Sprintf("Deu Certo - %f", bid)

		w.Write([]byte(writeResult))
	})

	http.ListenAndServe(":8080", nil)

}

func performTask(ctx context.Context, db *sql.DB) (float32, error) {
	select {
	case <-time.After(200 * time.Millisecond):
		quote, err := getQuote()
		if err != nil {
			return 0, err
		}

		err = saveData(&quote, db)
		if err != nil {
			return 0, err
		}

		bid, err := getLatestQuote(db)

		return bid, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func connectToDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./exchange.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) (string, error) {
	sqlStmt := `
	drop table exchange;
	create table exchange (id integer not null primary key autoincrement, real_value double(10,2), requested_at timestamp default CURRENT_TIMESTAMP);
	delete from exchange;
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return "", err
	}

	return "", nil
}

func getQuote() (Usdbrl, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return Usdbrl{}, err
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Usdbrl{}, err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Usdbrl{}, err
	}

	var quote Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		return Usdbrl{}, err
	}

	return quote.USDBRL, nil
}

func saveData(quote *Usdbrl, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("insert into exchange(real_value) values(?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(quote.Bid)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func getLatestQuote(db *sql.DB) (float32, error) {

	stmt, err := db.Prepare("select real_value from exchange order by id desc limit 1")
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer stmt.Close()
	var bid float32
	err = stmt.QueryRow().Scan(&bid)
	if err != nil {
		return 0, err
	}

	return bid, nil
}
