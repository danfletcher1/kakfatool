package main

import (
	"fmt"
	"flag"
	"os"
	"os/signal"
	"bufio"
	"io"
	"log"
)

func main() {

	server := flag.String("conn", "localhost:9092", "Kafka server and port")
	topic := flag.String("topic", "test", "Kafka topic to listen on")
	input := flag.String("file", "", "Optional file or stdin" )
	
	flag.Parse()

	// Connect to kafka to read topic contents
	err := Connect(*server)
	if err != nil {
		fmt.Println("I have a problem connecting to kafka, so we can't get any work done. ", err)
		os.Exit(1)
	}
	defer Close()
	threads := make(chan struct{}, 1000)

	// print topic messages
	go func(){
		
		var err error
		// declear stdin as the default
		file := os.Stdin
		
		// open file if one is specified if given
		if *input != "" {
			// Open the file
			file, err = os.Open(*input)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
		}
		
		var d []byte;

		reader := bufio.NewReader(file)

		for {
			// process each line
			b, o, err := reader.ReadLine()
			if o == true {
				d = append(d, b...)
				continue
			}
			if len(d) > 0 {
				b = append(d, b...)
				d = []byte{}
			}
			if err == io.EOF {
				os.Exit(0)
			}
			if err != nil {
				log.Fatal(err)
			}
			threads <- struct{}{}
			go func(topic, data string){
				e := Send(topic, data)
				if e != nil { fmt.Println("kafka send error,", e)}
				<-threads
			}(*topic, string(b))
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	<-signals
	for len(threads) > 0 {}
}
