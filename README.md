# Go Weather App

## Visão Geral

Este é um aplicativo simples escrito em Go que recebe um CEP brasileiro, identifica a cidade correspondente e retorna o clima atual (temperatura em graus Celsius, Fahrenheit e Kelvin) usando a API ViaCEP para buscar a localização e a WeatherAPI para buscar as informações climáticas.

## Requisitos

- Go 1.19 ou superior
- Docker
- Conta no Google Cloud
- API Key da WeatherAPI

## Configuração

1. Clone este repositório:
    ```bash
    git clone https://github.com/seu-usuario/go-weather-app.git
    cd go-weather-app
    ```

2. Configure seu ambiente Go:
    ```bash
    go mod tidy
    ```

3. Adicione sua chave da API WeatherAPI ao ambiente:
    ```bash
    export WEATHER_API_KEY=6cdc84347236418999e50624242507
    ```

## Construção e Execução

### Localmente

Para construir e executar o aplicativo localmente:

1. Execute o aplicativo:
    ```bash
    go run main.go
    ```

2. Teste o endpoint:
    ```bash
    curl "http://localhost:8080/weather?cep=01001000"
    ```

### Docker

Para construir e executar o aplicativo usando Docker:

1. Construa a imagem Docker:
    ```bash
    docker build -t go-weather-app .
    ```

2. Execute o container Docker:
    ```bash
    docker run -p 8080:8080 -e WEATHER_API_KEY=6cdc84347236418999e50624242507 go-weather-app
    ```

3. Teste o endpoint:
    ```bash
    curl "http://localhost:8080/weather?cep=01001000"
    ```

## Implantação no Google Cloud Run

### Configuração Inicial

1. Autentique-se no Google Cloud:
    ```bash
    gcloud auth login
    gcloud config set project gorun
    ```

2. Ative as APIs necessárias:
    ```bash
    gcloud services enable run.googleapis.com
    gcloud services enable containerregistry.googleapis.com
    ```

### Construção e Envio da Imagem Docker

1. Configure o Docker para usar o Google Container Registry:
    ```bash
    gcloud auth configure-docker
    ```

2. Construa a imagem Docker e envie-a para o Google Container Registry:
    ```bash
    docker build -t gcr.io/gorun/go-weather-app .
    docker push gcr.io/gorun/go-weather-app
    ```

### Deploy no Google Cloud Run

1. URL para acessar a aplicação. Teste o endpoint:
    ```bash
    curl "https://gorun-cloud-571351426586.us-east1.run.app/weather?cep=01001000/"
    ```

