package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (c Coupons) Check(code string) string {
	for _, item := range c.Coupon {
		if code == item.Code {
			return "valid"
		}
	}
	return "invalid"
}

type Result struct {
	Status   string
	Mensagem string
}

var coupons Coupons

func main() {
	coupon := Coupon{
		Code: "abc",
	}

	coupons.Coupon = append(coupons.Coupon, coupon)

	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	coupon := r.PostFormValue("coupon")
	cep := r.PostFormValue("cep")
	valid := coupons.Check(coupon)

	resultFrete := makeHttpCall("http://localhost:9093", cep)

	result := Result{Status: valid}

	fmt.Println(resultFrete.Status)

	if resultFrete.Status == "invalid" && result.Status == "invalid" {
		result.Status = "invalid"
		result.Mensagem = "Cupom invalido e frete grátis indisponivel para sua regiao"
	}

	if resultFrete.Status == "invalid" && result.Status == "valid" {
		result.Status = "valid"
		result.Mensagem = "Cupom validado e frete grátis indisponivel para sua regiao"
	}

	if resultFrete.Status == "valid" && result.Status == "invalid" {
		result.Status = "invalid"
		result.Mensagem = "Cupom invalido e frete grátis disponivel para sua regiao"
	}

	if resultFrete.Status == "valid" && result.Status == "valid" {
		result.Status = "valid"
		result.Mensagem = "Cupom valido e frete grátis disponivel para sua regiao"
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error converting json")
	}

	fmt.Println(string(jsonResult))

	fmt.Fprintf(w, string(jsonResult))

}
func makeHttpCall(urlMicroservice string, cep string) Result {
	values := url.Values{}
	values.Add("cep", cep)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err1 := retryClient.PostForm(urlMicroservice, values)
	if err1 != nil {
		result := Result{Status: "Erro durante a requisicao"}
		return result
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Erro durante o parse ")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result

}
