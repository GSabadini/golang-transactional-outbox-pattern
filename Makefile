docker-up:
	chmod +x ./_scripts/aws/localstack.sh
	docker-compose up -d

docker-down:
	docker-compose down --remove-orphans --volumes

run:
	go run main.go

list-queues:
	aws sqs --endpoint-url http://localhost:4566 list-queues

list-topics:
	aws sns --endpoint-url http://localhost:4566 list-topics

list-events:
	aws sqs receive-message --endpoint-url http://localhost:4566 --queue-url http://localhost:4566/000000000000/EventSubscriber --attribute-names All --message-attribute-names All --max-number-of-messages 10