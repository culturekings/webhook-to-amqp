# Webhook to RabbitMQ / AMQP

This simple project logs webhook events in to a RabbitMQ queue. The benefit of this is to allow webhooks to be executed by a background task rather than syncronously. This allows webhooks to withstand application failure, allows webhooks to scale and allows the application to control processing speed. It is designed to work with Heroku and require no installation.

## Setup

1. Create app in Heroku
2. Deploy this repository
3. Set a ``AMQP_SERVER`` environmental variable in Heroku to point to your Queue Server
4. Optionally the queue routing key can be changed with a ``AMQP_ROUTING_KEY`` variable. The default is ``webhooks``

## Usage

The service is designed to handle any url, method, headers and body content input and log it against the queue task in a simple json format. An example is below

    {
      "Method": "GET",
      "Url": {
        "Scheme": "",
        "Opaque": "",
        "User": null,
        "Host": "",
        "Path": "\/what",
        "RawPath": "",
        "RawQuery": "",
        "Fragment": ""
      },
      "Headers": {
        "Accept": [
          "text\/html,application\/xhtml+xml,application\/xml;q=0.9,image\/webp,*\/*;q=0.8"
        ],
        "Accept-Encoding": [
          "gzip, deflate, sdch"
        ],
        "Accept-Language": [
          "en-US,en;q=0.8"
        ],
        "Cache-Control": [
          "max-age=0"
        ],
        "Connection": [
          "keep-alive"
        ],
        "Upgrade-Insecure-Requests": [
          "1"
        ],
        "User-Agent": [
          "Mozilla\/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit\/537.36 (KHTML, like Gecko) Chrome\/49.0.2623.87 Safari\/537.36"
        ]
      },
      "Body": ""
    }

This allows workers to consume the queue and handle the allocation of each task to its relevant processor.