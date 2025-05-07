SERVER_ENV_FILE_PATH=./server/.env

.PHONY: server

server: docker-compose.yml
	docker-compose --env-file $(SERVER_ENV_FILE_PATH) up --build
