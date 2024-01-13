package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		result, err := performTask(ctx)
		if err != nil {
			fmt.Print("cancelado")
			return
		}

		w.Write([]byte("Deu Certo - " + result))
	})

	http.ListenAndServe(":8080", nil)

}

func performTask(ctx context.Context) (string, error) {
	select {
	case <-time.After(3 * time.Millisecond):
		return "O request foi finalizado no tempo certo", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
