package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sangketkit01/personal-financial/api"
	db "github.com/sangketkit01/personal-financial/db/sqlc"
	"github.com/sangketkit01/personal-financial/token"
	"github.com/sangketkit01/personal-financial/util"
)

func main() {
	config, err := util.LoadEnv(".")
	if err != nil{
		log.Fatal("cannot load env",err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	database, err := pgx.Connect(ctx, config.DatabaseSource)
	if err != nil{
		log.Fatal("cannot connect to database", err)
	}

	store := db.NewStore(database) 
	tokenMaker, err := token.NewJWTMaker("12345678901234567890123456789012")
	if err != nil{
		log.Fatal("cannot create token maker", err)
	}

	_, err = api.NewServer(config, store, tokenMaker)
	if err != nil{
		log.Fatal("cannot start server")
	}
	
	log.Println("Server start at: ", config.ServerPort)
}
