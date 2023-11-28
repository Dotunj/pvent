##  Pvent

Pvent is a CLI tool for sending events across different message brokers. At the moment, only Google Pub/Sub,
Kafka and Amazon SQS are supported.

### Getting Started

To get started with using Pvent, you need to provide your message brokers credentials in a config file or via the CLI directly. Refer [here](https://github.com/Dotunj/pvent#cli-usage-manual) for the CLI usage manual and the 
flags tied to the respective brokers.

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
  - access-key-id string
        AWS Access Key ID
  - secret-access-key string
        AWS Secret Access Key
  - region string
        AWS Region
  - queue-name
        AWS Queue Name
  - address
        Kafka Cluster Address
  - auth string
        Kafka Auth Type (plain, scram)
  - hash string
        Kafka Auth Hash (SHA256, SHA512)
  - ktopic string
        Kafka Topic Name
  - username string
        Kafka Auth Username
  - password string
        Kafka Auth Password
  - gtopic string
        Google Pub/Sub Topic Name
  - project-id string
        Google Pub/Sub Project ID
  - type string
        Message Brokers Type (sqs, google, kafka)
  - rate int
        Total number of events to send
  - target string
        Path to JSON payload to dispatch
  - config
       Path to Config file
```

### License
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.