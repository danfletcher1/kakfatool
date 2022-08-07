/*

Date: 23/04/18
Codebase:
Description:

To use the package you need to first set the connection string
it takes a host and messaging topic
then you can log message or errors

USAGE:

kafka.Connect("localhost")
go func() {
	kafka.Receive(conf["kafkaTopic"], receiveChannel)
		if err != nil {
		fmt.Println(err)
	}
}()

for {
	fmt.Println( <- receiveChannel)
}

kafka.Send("testtopic", "hello everyone")

*/

package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

// Local variables
var producer sarama.AsyncProducer // the Kafka producer connection
var kafkaServer string
var config *sarama.Config

// Connect to kafak. You must connect before any other function
func Connect(Server string) error {
	// save this as a function variable
	kafkaServer = Server

	// Create Kafka connection string
	config = sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages using snappy
	//config.Producer.CompressionLevel =  // don't know the compression options leaving to default
	config.Producer.Flush.Frequency = time.Millisecond * 100 // Bunches messages and max send every ....
	//config.Producer.Flush.Bytes = 1000 // Bunch message and send when message size is maximum of ....
	// Don't go over 1Mb without updating the consumers default on the server ....

	var err error
	// Connect to produce on Kafka just error and return, dont want the main program failing
	producer, err = sarama.NewAsyncProducer([]string{kafkaServer}, config)
	if err != nil {
		return errors.New("Failed to connect to kafka server. " + err.Error())
	}
	return nil
}

// Receive will output messages from a kafka topic as they arrive on the channel provided
// run as a goroutine, the function will hang waiting for messages, and exits with <-signals
func Receive(topic string, receiveChannel chan []byte, offset int64) error {

	// Connect to consume Kafka just error and return, dont want the main program failing
	consumer, err := sarama.NewConsumer([]string{kafkaServer}, config)
	if err != nil {
		return errors.New("Failed to open Kafka consumer " + topic + ". " + err.Error())
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			errors.New("ERROR: Failed to close Kafka producer. " + err.Error())
		}
	}()

	partitions, err := consumer.Partitions(topic)
	if err != nil {
		return errors.New("ERROR: Failed to find a list of kafka partitions for this topic. " + err.Error())
	}

	// Start a listener for all partitions
	for _, partition := range partitions {
		go func(partition int32) {

			// Connect to partition
			fmt.Println("Listening to Kafka topic", topic, "partition", partition)
			partitionConsumer, err := consumer.ConsumePartition(topic, partition, offset)
			if err != nil {
				errors.New("ERROR: Failed to open Kafka consumer partition. " + err.Error())
			}

			for {
				select {
				// We will receive the message
				case msg := <-partitionConsumer.Messages():
					receiveChannel <- msg.Value
					// If there is a terminate signal
				}

			}
			// Close partition
			defer func() {
				if err := partitionConsumer.Close(); err != nil {
					errors.New("ERROR: Failed to close Kafka consumer partition. " + err.Error())
				}

			}()
		}(partition)

	}

	// End operation but don't close
	signals := make(chan bool)
	<-signals

	return nil
}

// Close should be run as a defered process right after connect
func Close() error {
	if err := producer.Close(); err != nil {
		return errors.New("Failed to disconnect Kafka producer. " + err.Error())
	}
	return nil
}

// Send any message to the kafka topic specified
func Send(topic, msg string) error {
	producer.Input() <- &sarama.ProducerMessage{Topic: topic, Value: sarama.StringEncoder(msg)}

	// we are doing a wait here for a reply
	select {
	case <-producer.Successes():	
		return nil
	case err := <-producer.Errors():
		// if we get an error
		return errors.New("Failed to send kafka message to " + topic + ". " + err.Error())
	}
}
