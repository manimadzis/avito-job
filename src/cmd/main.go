package main

import (
	"avito-job/src/internal/domain"
	"avito-job/src/internal/repository/postgres"
	dbclient "avito-job/src/pkg/dbclient/postgres"
	"avito-job/src/pkg/logging"
	"context"
	"fmt"
)

func main() {
	logger := logging.Get()
	db, err := dbclient.New(dbclient.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "pass",
		Database: "bank_service",
	})
	if err != nil {
		logger.Panic(err)
	}
	defer db.Close()

	repo := postgres.NewRepository(db, logger)
	var i uint
	for i = 1; i < 10; i++ {
		err = repo.ReplenishBalance(context.Background(), i, domain.Float64ToMoney(1000.0))
		fmt.Println(err)
		err = repo.ReserveMoney(context.Background(), i, domain.Float64ToMoney(8.0), 1, 1, "123123")
		fmt.Println(err)
		err = repo.RecognizeRevenue(context.Background(), i, domain.Float64ToMoney(8.0), 1, 1)
		fmt.Println(err)
	}

	//a, err := repo.GetBalance(context.Background(), 2)
	//fmt.Println(a, err)
	//
	//err = repo.ReserveMoney(context.Background(), 2, domain.Float64ToMoney(800.0), 1, 1, "123j")
	//err = repo.RecognizeRevenue(context.Background(), 2, domain.Float64ToMoney(800.0), 1, 1)
	//fmt.Println(err)
	//
	////
	////
	//a, err = repo.GetBalance(context.Background(), 2)
	//fmt.Println(a, err)
}
