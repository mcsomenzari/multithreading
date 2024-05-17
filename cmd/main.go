package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	urlApiCep string = "https://brasilapi.com.br/api/cep/v1/%s"
	urlViaCep string = "https://viacep.com.br/ws/%s/json/"
)

type CEPResp struct {
	URL  string
	Body string
}

func main() {
	chanApiCep := make(chan CEPResp)
	chanViaCep := make(chan CEPResp)

	go Worker(urlApiCep, "02712080", chanApiCep, 0)
	go Worker(urlViaCep, "02712080", chanViaCep, 0)

	select {
	case apiCep := <-chanApiCep:
		fmt.Printf("URL: %s\nResposta: %s\n", apiCep.URL, apiCep.Body)
	case viaCep := <-chanViaCep:
		fmt.Printf("URL: %s\nResposta: %s\n", viaCep.URL, viaCep.Body)
	case <-time.After(time.Second):
		log.Fatalln("Tempo de resposta excedido.")
	}
}

func Worker(url string, cep string, bodyChannel chan<- CEPResp, delay time.Duration) {

	time.Sleep(0)

	cr := CEPResp{URL: fmt.Sprintf(url, cep)}
	r, err := http.NewRequest("GET", cr.URL, nil)
	if err != nil {
		close(bodyChannel)
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		close(bodyChannel)
	}

	body, err := io.ReadAll(res.Body)
	_ = res.Body.Close()

	if err != nil {
		close(bodyChannel)
	}

	cr.Body = string(body)

	bodyChannel <- cr
}
