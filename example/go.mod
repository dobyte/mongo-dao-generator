module example

go 1.17

require (
	github.com/dobyte/gen-mongo-dao v0.0.1
	go.mongodb.org/mongo-driver v1.11.0
)

replace github.com/dobyte/gen-mongo-dao => ../
