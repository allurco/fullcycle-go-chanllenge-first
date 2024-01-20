package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Quote struct {
	Bid float32 `json:"bid"`
}

func main() {

	var filename string = "./cotacao.txt"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var quote Quote

	json.Unmarshal(body, &quote)

	if _, err := os.Stat(filename); err != nil {

		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("Dólar: %f \n", quote.Bid))

	} else {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 644)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("Dólar: %f \n", quote.Bid))

	}

	io.Copy(os.Stdout, res.Body)
}
