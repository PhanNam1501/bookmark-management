package main

import (
	"github.com/PhanNam1501/bookmark-management/pkg/sqldb"
)

type User struct {
	FirstName string
	LastName  string
	Password  string
}

func main() {
	dbClient, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}

	sqldb.MigartePostgresDB(dbClient, "file://./migration", "up", 0)
	// userRepo := repository.NewUser(dbClient)

	// _, _ = userRepo.CreateUser(context.Background(), &model.User{
	// 	UserName:    "John Doe",
	// 	Password:    "123456",
	// 	DisplayName: "John Doe",
	// 	Email:       "johndoe@gmail.com",
	// })
}
