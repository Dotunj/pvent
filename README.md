##  Pvent

Pvent is a powerful and flexible CLI tool designed for sending events seamlessly across multiple message brokers. It simplifies event-driven communication by providing a unified interface for interacting with different messaging systems. Currently, Pvent supports Google Pub/Sub, Kafka, Amazon SQS, and RabbitMQ, making it a valuable tool for developers working with distributed systems and microservices.

With Pvent, users can easily publish and manage events across these platforms without needing to switch between different SDKs or APIs. Whether you are integrating with cloud-based messaging services or on-premise brokers, Pvent offers a streamlined experience for efficient event transmission. Future updates aim to expand support for additional message brokers, enhancing its versatility even further.

### Installation
To install Pvent using homebrew, you can run the following commands:

```bash
brew tap dotunj/tools
brew install dotunj/tools/pvent
```
You can also check the [releases section](https://github.com/Dotunj/pvent/releases) to download the binary for your OS.

### Getting Started

To get started with using Pvent, you need to provide your message brokers credentials in a config file.

An example `pvent.json` config file looks like:

```json
{
  "type": "sqs",
  "sqs": {
    "access_key_id": "<access-key-id>",
    "secret_access_key": "<secret-access-key>",
    "region": "<region>",
    "queue_name": "<queue-name>"
  },
  "kafka": {
    "address": "<address>",
    "topic": "<topic>",
    "auth": {
      "tls": true,
      "type": "scram",
      "hash": "SHA256",
      "username": "<username>",
      "password": "<password>"
    }
  },
  "google": {
    "project_id": "<project-id>",
    "topic_name": "<topic-name>"
  },
  "rabbitmq": {
    "dsn": "<dsn>",
    "queue_name": "<queue-name>"
  }
}
```
Using the CLI, you can now run:

```bash
pvent dispatch --target payload.json --rate 1
```
If you're using Google Pub/Sub, you'll need to expose the path to your service account credentials file as an environment variable

```bash
export GOOGLE_APPLICATION_CREDENTIALS="/Users/dotunj/Documents/certs/service.json" 
```


### CLI Usage Manual
```
Usage: pvent <command> [command flags]

dispatch command:
  - type string
        Message Brokers Type (sqs, google, kafka, rabbitmq)
  - rate int
        Total number of events to send
  - target string
        Path to JSON payload to dispatch
  - config
       Path to Pvent config file
```

### License
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.