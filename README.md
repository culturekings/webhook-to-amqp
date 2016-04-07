# Webhook to RabbitMQ / AMQP

This simple project logs webhook events in to a RabbitMQ queue. It is designed to work with Heroku and require no installation.

## Setup

1. Create app in Heroku
2. Deploy this repository
3. Set a "AMQP_SERVER" environmental variable in Heroku to point to your Queue Server

## Usage

The service is designed to handle any url, method, headers and body content input and log it against the queue task in a simple json format. An example is below

    {
        "method" : "POST",
        "url" : "/webhook/example",
        "headers" : "header",
        "body" : "Your body content"
    }

This allows workers to consume the queue and handle the allocation of each task to its relevant processor.