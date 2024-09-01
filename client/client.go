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

type Cotacao struct {
	Bid string `json:"bid"`
}

func buscaCotacao() (*Cotacao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cotacao Cotacao
	err = json.Unmarshal(data, &cotacao)
	if err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func salvaCotacao(cotacao *Cotacao) error {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", cotacao.Bid))
	return err
}

func main() {
	cotacao, err := buscaCotacao()
	if err != nil {
		panic(err)
	}

	err = salvaCotacao(cotacao)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Arquivo criado com sucesso!\n")
}
