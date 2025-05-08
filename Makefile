COMPOSE=docker-compose

.PHONY: server client1

all: server client1

server: docker-compose.yml
	$(COMPOSE) up --build -d server

client1: docker-compose.yml
	$(COMPOSE) up --build -d client1 

down:
	$(COMPOSE) down

config:
	$(COMPOSE) config 
