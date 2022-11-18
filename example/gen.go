package main

import (
	"example/model/mail"
	"github.com/dobyte/gen-mongo-dao"
	"log"
)

func main() {
	g := gen.NewGenerator()

	g.AddModel(&mail.Mail{}, "./dao", "example/dao")

	err := g.MakeDao()
	if err != nil {
		log.Fatal(err)
	}
}
