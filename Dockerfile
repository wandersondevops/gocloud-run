# Use uma imagem base oficial do Go
FROM golang:1.17-alpine

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie o go.mod e o go.sum, e baixe as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copie o código fonte da aplicação para o contêiner
COPY . .

# Compile a aplicação
RUN go build -o main .

# Defina a variável de ambiente PORT para 8080
ENV PORT=8080

# Exponha a porta que a aplicação usará
EXPOSE 8080

# Comando para iniciar a aplicação
CMD ["./main"]
