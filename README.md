#  Go Pub Sub
A little pub sub exercise.

## Build & Run the server
Build and run the application with
```Makefile
make build && ./build/pubsub
```

Alternatively simply run
```
go run main.go
```

## Publishing and Subscribing
Once the server is up an running clients can subscribe by opening a new command line window and run
```
curl http://localhost:9090
```
This will start a new session ready to listen to incoming messages.

Messages can be published by sending a POST request to '/'.
The request JSON payload must contain a message with the following structure: `{message: "my test message"}`.
To broadcast a message open another command line window and issuing the command
```
curl -X POST -d '{"message" : "my test message"}' http://localhost:9090

``` 
will send the message to existing subscribers and log the number of open connections.

## License
MIT.
