# dnpipes

Distributed Named Pipes, short: dnpipes, are essentially a distributed variant of Unix [named pipes](http://en.wikipedia.org/wiki/Named_pipe). They play a similar role that, for example, [SQS](https://aws.amazon.com/sqs/) plays in AWS or the [Service Bus](https://azure.microsoft.com/en-us/services/service-bus/) plays in Azure. 

## Requirements and interface specification

Interpret the key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "MAY", "MAY NOT", and "OPTIONAL" in the context of this repo document as defined in [RFC 2119](https://tools.ietf.org/html/rfc2119).

A dnpipes implementation MUST support two operations:

- `push(TOPIC, MSG)`
- `pull(TOPIC)`

## Reference implementation

The reference implementation of dnpipes is based on [DC/OS](https://dcos.io) using [Apache Kafka](http://kafka.apache.org/).

### Install

TBD.

### Use

TBD.
