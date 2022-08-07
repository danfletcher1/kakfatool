package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
        "github.com/Shopify/sarama"	
)

func main() {
	server := flag.String("conn", "localhost:9092", "Kafka server and port")
	topic := flag.String("topic", "test", "Kafka topic to listen on")
	old := flag.Bool("oldest", false, "Start reading from oldest record")
	delim := flag.String("delim", "\n", "Delimiter charactor between messages")

	flag.Parse()

	offset := sarama.OffsetNewest
	if *old {
		offset = sarama.OffsetOldest
	}

	// Connect to kafka to read topic contents
	err := Connect(*server)
	if err != nil {
		fmt.Println("I have a problem connecting to kafka, so we can't get any work done. ", err)
		os.Exit(1)
	}
	defer Close()

	// create channel for topic messages
	msg := make(chan []byte, 1)

	// read topic messages from kafka
	go func(){
		Receive(*topic, msg, offset)
	}()

	// print topic messages
	go func(){
		for v := range msg {
			fmt.Printf(string(v)+*delim)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
}
