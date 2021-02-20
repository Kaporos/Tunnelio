# Tunnelio

Tunnelio is a simple self hosted alternative to ngrok !

## Requirements

- Docker
- Golang ( or download compiled binaries )
- These DNS records to your server:

	\*.\*.yourdomain.com

## Usage

1. First, create a "tunnelio" docker network
	
	docker network create tunnelio 

2. Build the docker image
	docker build . -t tunnelio:test

3. Launch server.go , and you're done !

	go run server.go

4. Clients will only have to do : 
	go run client.go -port LOCAL_PORT -url YOURDOMAIN
Ex: go run client.go -port 1234 -url yourdomain.net
	

## How it's works

Thanks to docker, all forwards ports will be fully isolated from your server, i.e. if an X client wants to expose its local port, its port will be exposed in the docker (containerip: randomport), and the public endpoint will act like a reverse proxy, which will redirect the request to the docker
