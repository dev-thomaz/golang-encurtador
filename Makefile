# Define o nome do nosso binário final.
BINARY_NAME=shortener-api

# Define o ponto de entrada principal da nossa aplicação.
MAIN_PATH=./cmd/api/main.go

## ====================================================================================
## Comandos de Desenvolvimento
## ====================================================================================

# O .PHONY garante que o 'make' execute o comando mesmo que já exista um arquivo com o mesmo nome.
.PHONY: run
run:
	@go run $(MAIN_PATH)

.PHONY: tidy
tidy:
	@go mod tidy

## ====================================================================================
## Comandos de Build e Teste
## ====================================================================================

.PHONY: build
build: tidy
	@echo "Construindo o binário..."
	@go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Binário $(BINARY_NAME) construído com sucesso."

.PHONY: test
test:
	@go test -v ./...

.PHONY: clean
clean:
	@echo "Limpando o binário construído..."
	@if [ -f $(BINARY_NAME) ]; then rm $(BINARY_NAME); fi
	@echo "Limpeza concluída."

## ====================================================================================
## Ajuda
## ====================================================================================

.PHONY: help
help:
	@echo "Comandos disponíveis:"
	@echo "  run        - Compila e executa a aplicação em modo de desenvolvimento."
	@echo "  build      - Compila o binário da aplicação para produção."
	@echo "  test       - Roda todos os testes do projeto."
	@echo "  tidy       - Organiza as dependências do go.mod."
	@echo "  clean      - Remove o binário construído."