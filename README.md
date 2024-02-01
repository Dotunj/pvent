##  Pvent

Pvent is a CLI tool for sending events across different message brokers. At the moment, only Google Pub/Sub,
Kafka and Amazon SQS are supported.

### Installation
To install Pvent using homebrew, you can run the following commands:

```bash
brew tap dotunj/tools
brew install dotunj/tools/pvent
```
You can also check the [releases section](https://github.com/Dotunj/pvent/releases) to download the binary for your OS.

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
  - type string
        Message Brokers Type (sqs, google, kafka)
  - rate int
        Total number of events to send
  - target string
        Path to JSON payload to dispatch
  - config
       Path to Pvent config file
```

### License
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.