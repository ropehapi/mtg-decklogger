# Variáveis
APP_NAME=mtg-decklogger
API_HOST=http://localhost
API_PORT=8080

run:
	go run ./cmd/api

build:
	sudo docker compose build

up :
	sudo docker compose up -d

down:
	sudo docker compose down

# Derruba tudo e remove volume do Postgres
reset:
	sudo docker compose down -v

# Limpa cache e reconstroi
rebuild:
	sudo docker compose build --no-cache

# Testa se a API está respondendo
health:
	curl -i $(API_HOST):$(API_PORT)/health

# Faz um commit
gcm:
	@if [ -z "$(m)" ]; then \
		echo "Erro: é necessário fornecer uma mensagem de commit."; \
		exit 1; \
	fi
	@git add . && git commit -m "$(m)"

fmt:
	go fmt ./...
