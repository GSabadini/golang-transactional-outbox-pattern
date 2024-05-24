#!/bin/bash

## SQS ## (Just to see the messages to debug)
aws sqs create-queue --queue-name EventSubscriber --endpoint-url http://localhost:4566 --region sa-east-1

## SNS ##
aws --endpoint-url=http://localhost:4566 sns create-topic --name Events --region sa-east-1

## SNS SUBSCRIBER ##
aws --endpoint-url=http://localhost:4566 sns subscribe --topic-arn arn:aws:sns:sa-east-1:000000000000:Events --protocol sqs --notification-endpoint arn:aws:sqs:sa-east-1:000000000000:EventSubscriber --region sa-east-1

