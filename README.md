# 🚀 Go URL Shortener (Versão Estruturada)

![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)

Um encurtador de URLs simples e eficiente construído com Go (Golang). Este projeto evoluiu de um script único para uma API RESTful completa, implementada seguindo uma arquitetura profissional em camadas (Clean Architecture) para garantir manutenibilidade, testabilidade e escalabilidade.

## 📄 Descrição

Este projeto é uma API RESTful que permite encurtar uma ou mais URLs longas, gerando um código curto para cada uma. A API armazena os links de forma persistente em um banco de dados SQLite. Ao acessar a URL com o código curto, o usuário é redirecionado para a URL original.

O objetivo principal é demonstrar a construção de uma aplicação Go robusta, cobrindo conceitos como servidores HTTP, roteamento avançado, manipulação de JSON, persistência de dados com SQL e, mais importante, a organização de código em uma estrutura de pacotes desacoplada.

## ✨ Features

- [x] **Arquitetura em Camadas:** Código organizado em pacotes `domain`, `storage`, e `handler`.
- [x] **Injeção de Dependência:** As camadas são desacopladas através de interfaces.
- [x] **Persistência de Dados:** Armazenamento em banco de dados **SQLite**.
- [x] **API RESTful Completa:** Suporte aos métodos `GET`, `POST`, e `DELETE`.
- [x] **Prevenção de Duplicatas:** Reutiliza links curtos para URLs já existentes.
- [x] **Gerenciador de Tarefas:** Comandos simplificados com `Makefile`.
- [x] **Roteamento Limpo:** Utilizando a biblioteca **Chi**.
- [x] **Respostas Estruturadas:** Comunicação via JSON.

## 🛠️ Tecnologias Utilizadas

- **Linguagem:** [Go](https://go.dev/)
- **Banco de Dados:** [SQLite](https://www.sqlite.org/index.html)
- **Roteador HTTP:** [Chi v5](https://github.com/go-chi/chi)
- **Gerenciador de Tarefas:** [Make](https://www.gnu.org/software/make/)
- **Ferramenta de Teste de API:** [Insomnia](https://insomnia.rest/)

## 📂 Estrutura do Projeto

O projeto segue um layout padrão para aplicações Go, separando claramente as responsabilidades:

```
golang-encurtador/
├── cmd/
│   └── api/
│       └── main.go         # Ponto de entrada: inicialização e injeção de dependência.
├── internal/
│   ├── domain/
│   │   └── link.go         # Define as estruturas de dados centrais.
│   ├── handler/
│   │   └── link_handler.go # Camada de apresentação (HTTP/JSON).
│   └── storage/
│       └── sqlite_storage.go # Camada de acesso a dados (SQL).
├── Makefile                # Apelidos para comandos comuns (run, build, test).
├── go.mod
└── storage.db              # (Criado na primeira execução, ignorado pelo Git).
```

## ⚙️ Como Executar o Projeto

1.  **Clone o repositório:**

    ```bash
    git clone [https://github.com/dev-thomaz/golang-encurtador.git](https://github.com/dev-thomaz/golang-encurtador.git)
    ```

2.  **Navegue até a pasta do projeto:**

    ```bash
    cd golang-encurtador
    ```

3.  **Execute a aplicação com Make:**
    ```bash
    make run
    ```
    Alternativamente, sem o Make: `go run ./cmd/api/main.go`

O servidor estará rodando em `http://localhost:8080`.

### Comandos Úteis (Makefile)

- `make run`: Executa a aplicação em modo de desenvolvimento.
- `make build`: Compila o binário otimizado para produção.
- `make test`: Roda todos os testes do projeto.
- `make clean`: Remove o binário construído.
- `make help`: Mostra todos os comandos disponíveis.

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

### Listar Todos os Links

- **Endpoint:** `GET /links`
- **Descrição:** Retorna uma lista de todos os links armazenados no banco de dados.

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
