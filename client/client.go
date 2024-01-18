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
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080", nil)
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

	if _, err := os.Stat("./cotacao.txt"); err != nil {

		f, err := os.Create("./cotacao.txt")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("Dolar: %f \n", quote.Bid))

	} else {
		f, err := os.OpenFile("./cotacao.txt", os.O_APPEND|os.O_WRONLY, 644)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		f.WriteString(fmt.Sprintf("Dolar: %f \n", quote.Bid))

	}

	io.Copy(os.Stdout, res.Body)
}
