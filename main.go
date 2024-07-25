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
	logrus.Infof("Consultando API ViaCEP para CEP: %s", cep)
	resp, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep))
	if err != nil {
		logrus.Errorf("Erro ao consultar API ViaCEP: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	logrus.Infof("Resposta da API ViaCEP: %s", bodyString)

	var viaCEPResponse ViaCEPResponse
	if err := json.Unmarshal(bodyBytes, &viaCEPResponse); err != nil {
		logrus.Errorf("Erro ao decodificar resposta da API ViaCEP: %v", err)
		return nil, err
	}

	if viaCEPResponse.CEP == "" {
		logrus.Errorf("CEP não encontrado na resposta da API ViaCEP")
		return nil, fmt.Errorf("CEP not found")
	}

	return &viaCEPResponse, nil
}

func getWeather(location string) (*WeatherAPIResponse, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		logrus.Error("API key is missing")
		return nil, fmt.Errorf("API key is missing")
	}
	logrus.Infof("Weather API Key: %s", apiKey)

	locationEncoded := url.QueryEscape(location)
	weatherAPIURL := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, locationEncoded)
	logrus.Infof("Consultando WeatherAPI com URL: %s", weatherAPIURL)

	resp, err := http.Get(weatherAPIURL)
	if err != nil {
		logrus.Errorf("Erro ao consultar API WeatherAPI: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	logrus.Infof("Resposta da API WeatherAPI: %s", bodyString)

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("WeatherAPI retornou status code %d: %s", resp.StatusCode, bodyString)
		return nil, fmt.Errorf("WeatherAPI returned status code %d: %s", resp.StatusCode, bodyString)
	}

	var weatherResponse WeatherAPIResponse
	if err := json.Unmarshal(bodyBytes, &weatherResponse); err != nil {
		logrus.Errorf("Erro ao decodificar resposta da API WeatherAPI: %v", err)
		return nil, err
	}

	return &weatherResponse, nil
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	// Log the API key at startup
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		logrus.Fatal("WEATHER_API_KEY environment variable is not set")
	} else {
		logrus.Infof("WEATHER_API_KEY is set")
	}

	// Get the port from the environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logrus.Infof("Defaulting to port %s", port)
	}

	http.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		cep := r.URL.Query().Get("cep")
		logrus.Infof("CEP recebido: %s", cep)
		if len(cep) != 8 {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		cepInfo, err := getCEPInfo(cep)
		if err != nil {
			http.Error(w, "can not find zipcode", http.StatusNotFound)
			logrus.Errorf("Erro ao encontrar o CEP: %v", err)
			return
		}
		logrus.Infof("Informações do CEP: %+v", cepInfo)

		weatherInfo, err := getWeather(cepInfo.Localidade)
		if err != nil {
			http.Error(w, "error getting weather information", http.StatusInternalServerError)
			logrus.Errorf("Erro ao obter informações meteorológicas: %v", err)
			return
		}
		logrus.Infof("Informações meteorológicas: %+v", weatherInfo)

		weatherResponse := map[string]float64{
			"temp_C": weatherInfo.Current.TempC,
			"temp_F": weatherInfo.Current.TempC*1.8 + 32,
			"temp_K": weatherInfo.Current.TempC + 273,
		}

		w.Header().Set("Content-Type", "application/json")
		compactJSON, err := json.Marshal(weatherResponse)
		if err != nil {
			http.Error(w, "error generating JSON response", http.StatusInternalServerError)
			logrus.Errorf("Erro ao gerar resposta JSON: %v", err)
			return
		}
		w.Write(compactJSON)
	})

	logrus.Infof("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
