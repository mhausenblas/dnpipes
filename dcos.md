# A Mesos reference implementation

This reference implementation of `dnpipes` is based on [DC/OS](https://dcos.io) using [Apache Kafka](http://kafka.apache.org/).

## Install

From source:

```bash
$ go get github.com/mhausenblas/dnpipes
$ go build
$ sudo mv dnpipes /usr/local/dnpipes
```

From binaries, for Linux:

```bash
$ curl -s -L https://github.com/mhausenblas/dnpipes/releases/download/0.1.0/linux-dnpipes -o dnpipes
$ sudo mv dnpipes /usr/local/dnpipes
$ sudo chmod +x /usr/local/bin/dnpipes
```

From binaries, for macOS:

```bash
$ curl -s -L https://github.com/mhausenblas/dnpipes/releases/download/0.1.0/macos-dnpipes -o dnpipes
$ sudo mv dnpipes /usr/local/dnpipes
$ sudo chmod +x /usr/local/bin/dnpipes
```

## Example session

To try it out yourself, you first need a to [install a DC/OS cluster](https://dcos.io/install/) and then Apache Kafka like so:

```bash
$ dcos package install kafka
```

Note that if you are unfamiliar with Kafka and its terminology, you can check out the respective [Kafka 101 example](https://github.com/dcos/examples/tree/master/1.8/kafka) now.

Next, figure out where the brokers are (in my case I started Kafka with one broker):

```bash
$ dcos kafka connection

{
  "address": [
    "10.0.2.94:9951"
  ],
  "zookeeper": "master.mesos:2181/dcos-service-kafka",
  "dns": [
    "broker-0.kafka.mesos:9951"
  ],
  "vip": "broker.kafka.l4lb.thisdcos.directory:9092"
}
```
Now, an example session using the `dnpipes` reference implementation looks as follows.
I've set up two terminals, in one I'm starting the `dnpipes` in publisher mode:

```bash
$ ./dnpipes --mode=publisher --broker=broker-0.kafka.mesos:9951 --topic=test
> hello!
> bye
> RESET
2016/12/11 13:38:57 Connected to 10.0.6.130:2181
2016/12/11 13:38:57 Authenticated: id=97087250213175381, timeout=4000
[0] - &{Czxid:295 Mzxid:295 Ctime:1481463523762 Mtime:1481463523762 Version:0 Cversion:1 Aversion:0 EphemeralOwner:0 DataLength:0 NumChildren:1 Pzxid:296}
reset this dnpipes
> ^C
```

The second terminal has `dnpipes` in subscriber mode running:

```bash
$ ./dnpipes --mode=subscriber --broker=broker-0.kafka.mesos:9951 --topic=test 2>/dev/null
hello!
bye
^C
```

And here's a screen shot of the whole thing:

![screen shot of example dnpipes session](img/session.png)
