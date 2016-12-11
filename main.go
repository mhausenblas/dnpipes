package main

import (
	"flag"
	"fmt"
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/samuel/go-zookeeper/zk"
	"os"
	"time"
)

const (
	VERSION         string = "0.1.0"
	MODE_PUBLISHER         = "publisher"
	MODE_SUBSCRIBER        = "subscriber"
)

var (
	version bool
	// operations mode of this agent:
	omode string
	// FQDN/IP + port of a Kafka broker:
	broker string
	// the active Kafka topic:
	topic string
	// the Kafka producer:
	producer sarama.SyncProducer
)

func init() {
	flag.BoolVar(&version, "version", false, "Display version information")
	flag.StringVar(&omode, "mode", MODE_SUBSCRIBER, fmt.Sprintf("The operations mode of this agent, can be either \"%s\" or \"%s\".", MODE_PUBLISHER, MODE_SUBSCRIBER))
	flag.StringVar(&broker, "broker", "", "The FQDN or IP address and port of a Kafka broker. Example: broker-1.kafka.mesos:9382 or localhost:9092")
	flag.StringVar(&topic, "topic", "", "The topic to publish to or pull from. Example: test")

	flag.Usage = func() {
		fmt.Printf("Usage: %s [args]\n\n", os.Args[0])
		fmt.Println("Arguments:")
		flag.PrintDefaults()
	}
	flag.Parse()
}

func about() {
	fmt.Printf("This is the dnpipes reference implementation in version %s\n", VERSION)
}

func handleReset() {
	zks := []string{"leader.mesos:2181"}
	conn, _, _ := zk.Connect(zks, time.Second)
	partitions, stat, _ := conn.Children("/dcos-service-kafka/brokers/topics/" + topic + "/partitions")
	fmt.Println(fmt.Sprintf("%+v - %+v", partitions, stat))
	for _, p := range partitions {
		if err := conn.Delete("/dcos-service-kafka/brokers/topics/"+topic+"/partitions/"+p+"/state", -1); err != nil {
			log.WithFields(log.Fields{"func": "handleSubscriber"}).Error("There was a problem resetting the topic:", err)
		}
		if err := conn.Delete("/dcos-service-kafka/brokers/topics/"+topic+"/partitions/"+p, -1); err != nil {
			log.WithFields(log.Fields{"func": "handleSubscriber"}).Error("There was a problem resetting the topic:", err)
		}
	}
	if err := conn.Delete("/dcos-service-kafka/brokers/topics/"+topic+"/partitions", -1); err != nil {
		log.WithFields(log.Fields{"func": "handleSubscriber"}).Error("There was a problem resetting the topic:", err)
	}
	if err := conn.Delete("/dcos-service-kafka/brokers/topics/"+topic, -1); err != nil {
		log.WithFields(log.Fields{"func": "handleSubscriber"}).Error("There was a problem resetting the topic:", err)
	}
	fmt.Println("reset this dnpipes")
}

func handlePublisher() {
	if p, err := sarama.NewSyncProducer([]string{broker}, nil); err != nil {
		log.Error(err)
		os.Exit(1)
	} else {
		producer = p
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}()
	imsg := ""
	for {
		fmt.Print("> ")
		fmt.Scanln(&imsg)
		if imsg == "RESET" {
			handleReset()
		} else {
			msg := &sarama.ProducerMessage{Topic: string(topic), Value: sarama.StringEncoder(imsg)}
			if _, _, err := producer.SendMessage(msg); err != nil {
				log.WithFields(log.Fields{"func": "handlePublisher"}).Error("Failed to send message ", err)
			} else {
				log.Debug(fmt.Sprintf("%#v", msg))
			}
		}
	}
}

func handleSubscriber() {
	var consumer sarama.Consumer
	if c, err := sarama.NewConsumer([]string{broker}, nil); err != nil {
		log.WithFields(log.Fields{"func": "handleSubscriber"}).Error(err)
		return
	} else {
		consumer = c
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.WithFields(log.Fields{"func": "handleSubscriber"}).Error(err)
		}
	}()

	if partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest); err != nil {
		log.WithFields(log.Fields{"func": "handleSubscriber"}).Error(err)
		return
	} else {
		defer func() {
			if err := partitionConsumer.Close(); err != nil {
				log.WithFields(log.Fields{"func": "handleSubscriber"}).Error(err)
			}
		}()
		for {
			msg := <-partitionConsumer.Messages()
			log.Debug(fmt.Sprintf("%#v", msg))
			fmt.Println(string(msg.Value))
		}
	}
}

func main() {
	if version {
		about()
		os.Exit(0)
	}
	if broker == "" {
		flag.Usage()
		os.Exit(1)
	}

	switch omode {
	case MODE_PUBLISHER:
		handlePublisher()
	case MODE_SUBSCRIBER:
		handleSubscriber()
	default:
		fmt.Println("Usage error, you provided an unknown mode")
		flag.Usage()
		os.Exit(1)
	}
}
