# Jaak Pub Sub
A little pub subexercise

## Build & Run the server
Build and run the application with
```Makefile
make build && ./build/jaakpubsub
```

Alternatively simply run
```
go run main.go
```

## Publishing and Subscribing
Once the server is up an running clients can subscribe by going to
http://localhost:9090

Messages can be published by sending a POST request to '/'.
The request JSON payload must contain a message with the following structure: ` {message: "my test message"}`.
If using curl a request example could look like:
`curl -X POST -d '{"message" : "my test message"}' http://localhost:9090`

## To Do
[] topics
[] test all the things
[] better persistence?

## License
MIT.