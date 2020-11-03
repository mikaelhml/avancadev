package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type CEP struct {
	numero string
}

var ceps = []string{"1", "2", "4", "8", "16", "32", "64", "128"}

func (c CEP) Check(cep string) string {

	for _, item := range ceps {
		if cep == item {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status   string
	Mensagem string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	numero := r.PostFormValue("cep")

	cep := CEP{
		numero: numero,
	}

	var valid = CEP.Check(cep, numero)
	result := Result{Status: valid}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Fprintf(w, string(jsonResult))
}
