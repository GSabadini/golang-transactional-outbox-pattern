#!/bin/bash

aws sqs create-queue --queue-name Accounts --endpoint-url http://localhost:4566
aws sqs create-queue --queue-name Transactions --endpoint-url http://localhost:4566

aws --endpoint-url=http://localhost:4566 sns create-topic --name Events

aws --endpoint-url=http://localhost:4566 sns subscribe --topic-arn arn:aws:sns:sa-east-1:000000000000:Events --protocol sqs --notification-endpoint http://localhost:4566/000000000000/Accounts
aws --endpoint-url=http://localhost:4566 sns subscribe --topic-arn arn:aws:sns:sa-east-1:000000000000:Events --protocol sqs --notification-endpoint http://localhost:4566/000000000000/Transactions
