package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
)

type ViaCEPResponse struct {
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	UF         string `json:"uf"`
	CEP        string `json:"cep"`
}

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func getCEPInfo(cep string) (*ViaCEPResponse, error) {
	logrus.Info("Consultando API ViaCEP para CEP: ", cep)
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	logrus.Info("Resposta da API ViaCEP: ", bodyString)

	var viaCEPResponse ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResponse); err != nil {
		return nil, err
	}

	if viaCEPResponse.CEP == "" {
		return nil, fmt.Errorf("CEP not found")
	}

	return &viaCEPResponse, nil
}

func getWeather(location string) (*WeatherAPIResponse, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API key is missing")
	}
	logrus.Info("Weather API Key: ", apiKey)

	locationEncoded := url.QueryEscape(location)
	resp, err := http.Get(fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, locationEncoded))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	logrus.Info("WeatherAPI Response Body: ", bodyString)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WeatherAPI returned status code %d: %s", resp.StatusCode, bodyString)
	}

	var weatherResponse WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	return &weatherResponse, nil
}

func main() {
	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Query().Get("cep")
		if len(cep) != 8 {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		cepInfo, err := getCEPInfo(cep)
		if err != nil {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			logrus.Error("Erro ao encontrar o CEP: ", err)
			return
		}

		weatherInfo, err := getWeather(cepInfo.Localidade)
		if err != nil {
			http.Error(w, "error getting weather information", http.StatusInternalServerError)
			logrus.Error("Erro ao obter informações meteorológicas: ", err)
			return
		}

		weatherResponse := map[string]float64{
			"temp_C": weatherInfo.Current.TempC,
			"temp_F": weatherInfo.Current.TempC*1.8 + 32,
			"temp_K": weatherInfo.Current.TempC + 273,
		}

		w.Header().Set("Content-Type", "application/json")
		compactJSON, err := json.Marshal(weatherResponse)
		if err != nil {
			http.Error(w, "error generating JSON response", http.StatusInternalServerError)
			return
		}
		w.Write(compactJSON)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
