# 🚀 Go URL Shortener (v1.0)

![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Um encurtador de URLs simples e eficiente construído com Go (Golang) como parte de um projeto de estudo e portfólio. Esta API RESTful é a primeira versão completa do projeto, com todas as funcionalidades básicas implementadas em um único pacote.

## 📄 Descrição

Este projeto é uma API RESTful que permite encurtar uma ou mais URLs longas, gerando um código curto para cada uma. A API armazena os links de forma persistente em um banco de dados SQLite. Ao acessar a URL com o código curto, o usuário é redirecionado para a URL original.

O objetivo principal é demonstrar conceitos fundamentais de desenvolvimento de backend em Go, incluindo servidores HTTP, roteamento avançado, manipulação de JSON e persistência de dados com SQL.

## ✨ Features

- [x] Encurtamento de uma ou mais URLs em uma única requisição (batch).
- [x] Prevenção de criação de links duplicados para a mesma URL.
- [x] Persistência de dados com banco de dados **SQLite**.
- [x] Roteamento de API limpo utilizando a biblioteca **Chi**.
- [x] Listagem de todos os links encurtados.
- [x] Rota para limpar todos os links do banco de dados.
- [x] Respostas estruturadas em formato JSON.
- [x] Redirecionamento de códigos curtos para as URLs originais.

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** [Go](https://go.dev/)
- **Banco de Dados:** [SQLite](https://www.sqlite.org/index.html)
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

O servidor estará rodando em `http://localhost:8080`. Um arquivo `storage.db` será criado na raiz do projeto.

## 📡 Como Usar a API

### Encurtar URLs

- **Endpoint:** `POST /shorten`
- **Descrição:** Envia um array de URLs para serem encurtadas. Se uma URL já existir, retorna o código existente.
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
  [
    {
      "original_url": "[https://www.google.com](https://www.google.com)",
      "short_url": "http://localhost:8080/aBcDeF"
    },
    {
      "original_url": "[https://github.com](https://github.com)",
      "short_url": "http://localhost:8080/gHiJkL"
    }
  ]
  ```

### Listar Todos os Links

- **Endpoint:** `GET /links`
- **Descrição:** Retorna uma lista de todos os links armazenados no banco de dados.
- **Resposta de Sucesso (200 OK):**
  ```json
  [
    {
      "original_url": "[https://www.google.com](https://www.google.com)",
      "short_url": "http://localhost:8080/aBcDeF"
    },
    {
      "original_url": "[https://github.com](https://github.com)",
      "short_url": "http://localhost:8080/gHiJkL"
    }
  ]
  ```

### Limpar Todos os Links

- **Endpoint:** `DELETE /links`
- **Descrição:** Remove todos os links do banco de dados.
- **Resposta de Sucesso (200 OK):**
  ```json
  {
    "message": "Todos os links foram removidos com sucesso.",
    "links_removidos": 2
  }
  ```

### Redirecionar

- **Endpoint:** `GET /{shortCode}`
- **Descrição:** Acessar esta URL no navegador redirecionará para a URL original correspondente.
- **Exemplo:** `http://localhost:8080/aBcDeF`
