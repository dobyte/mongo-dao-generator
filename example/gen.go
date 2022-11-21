package main

import (
	"example/model/mail"
	"example/model/user"
	"github.com/dobyte/gen-mongo-dao"
	"log"
)

func main() {
	g := gen.NewGenerator(&gen.Options{
		OutputDir:    "./dao",
		OutputPkg:    "example/dao",
		EnableSubPkg: true,
	})

	g.AddModels(
		&mail.Mail{},
		&user.User{},
	)

	err := g.MakeDao()
	if err != nil {
		log.Fatal(err)
	}
}
