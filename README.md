# kakfatool
Interact with Kafak using this command line tool. You can read/write a file or use stdin/stdout.
I find this very useful in debugging to see the messages that are in the queue, and to send in messages.
I know there are versions like this already provided by the kafka team, but its in java an thusly heavyweight. This compiles to 8.2mb binaries, and you don't need to install anything. You can build using containers if you do not wish to install go. Go can cross compile to another OS. 


# Build
Most will want to build for your current enviroment, you don't need to use the build, which will build using docker containers. You will need to install golang, and this was desined on linux, other OS you may need to make alterations. If your happy with docker, build using the docker file. You can compile direct on the OS with. 

go get ./...
go build -o kafkaListen kafka.go listen.go
go build -o kafkaSend kafka.go send.go


# Use
./kafakSend -?

./kafkaListen -?

cat thisfile | ./kafakSend -conn remote-kafka-server:9092 -topic mytopic

./kafkaListen -conn remote-kafak-server:9092 -topic mytopic (messages will be dumped to screen)
