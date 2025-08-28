# 🚀 Go URL Shortener

![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Um encurtador de URLs simples e eficiente construído com Go (Golang) como parte de um projeto de estudo e portfólio.

## 📄 Descrição

Este projeto é uma API RESTful que permite encurtar uma ou mais URLs longas, gerando um código curto para cada uma. Ao acessar a URL com o código curto, o usuário é redirecionado para a URL original.

O objetivo principal é demonstrar conceitos fundamentais de desenvolvimento de backend em Go, incluindo servidores HTTP, roteamento, manipulação de JSON e, futuramente, persistência de dados.

## ✨ Features

- [x] Encurtamento de uma ou mais URLs em uma única requisição (batch).
- [x] Respostas estruturadas em formato JSON.
- [x] Roteamento de API limpo utilizando a biblioteca Chi.
- [x] Redirecionamento de códigos curtos para as URLs originais.
- [ ] Persistência de dados com SQLite (próxima etapa).

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** [Go](https://go.dev/)
- **Roteador HTTP:** [Chi v5](https://github.com/go-chi/chi)
- **Ferramenta de Teste de API:** [Insomnia](https://insomnia.rest/)

## ⚙️ Como Executar o Projeto

1.  **Clone o repositório:**

    ```bash
    git clone [https://github.com/dev-thomaz/golang-encurtador.git](https://github.com/dev-thomaz/golang-encurtador.git)
    ```

2.  **Navegue até a pasta do projeto:**

    ```bash
    cd golang-encurtador
    ```

3.  **Execute a aplicação:**
    ```bash
    go run main.go
    ```

O servidor estará rodando em `http://localhost:8080`.

## 📡 Como Usar a API

### Encurtar URLs

- **Endpoint:** `POST /shorten`
- **Descrição:** Envia um array de URLs para serem encurtadas.
- **Corpo da Requisição (JSON):**
  ```json
  {
    "urls": [
      "[https://www.google.com](https://www.google.com)",
      "[https://github.com](https://github.com)"
    ]
  }
  ```
- **Resposta de Sucesso (201 Created):**
  ```json
  {
    "results": [
      {
        "original_url": "[https://www.google.com](https://www.google.com)",
        "short_url": "http://localhost:8080/aBcDeF"
      },
      {
        "original_url": "[https://github.com](https://github.com)",
        "short_url": "http://localhost:8080/gHiJkL"
      }
    ]
  }
  ```

### Redirecionar

- **Endpoint:** `GET /{shortCode}`
- **Descrição:** Acessar esta URL no navegador redirecionará para a URL original correspondente.
- **Exemplo:** `http://localhost:8080/aBcDeF`
