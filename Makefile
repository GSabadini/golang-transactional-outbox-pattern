up:
	docker-compose up -d

down:
	docker-compose down --remove-orphans --volumes

run:
	go run main.go