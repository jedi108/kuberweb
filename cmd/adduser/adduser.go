package main

import (
	"fmt"
	"log"
	"os"

	"git.betfavorit.cf/vadim.tsurkov/kuberweb/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func newConfig() (*viper.Viper, *sqlx.DB, error) {
	c := viper.New()

	c.SetEnvPrefix("PG")
	c.BindEnv("USER")
	pguser := c.Get("USER")
	c.BindEnv("PASSWORD")
	pgpass := c.Get("PASSWORD")
	c.BindEnv("DATABASE")
	pgdatabase := c.Get("DATABASE")

	dsn := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?sslmode=disable", pguser, pgpass, pgdatabase)
	c.SetDefault("cookie_secret", "qaBzlTixkx2c9S6i")
	c.AutomaticEnv()

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, db, err
	}

	return c, db, nil
}
func getUserInput(input interface{}) (int, error) {
	var err error

	switch t := input.(type) {
	case *int:
		_, err = fmt.Scanf("%d", input)
	case *string:
		_, err = fmt.Scanf("%s", input)
	default:
		fmt.Printf("unexpected type %T", t)
	}

	return 0, err
}

func UserConsoleRead() (email string, passd string) {
	fmt.Print("Enter user email:")
	getUserInput(&email)
	if len(email) == 0 {
		fmt.Println("Email not entered")
		os.Exit(1)
	}

	var passwd string
	fmt.Print("Enter user passwd:")
	getUserInput(&passwd)
	if len(email) == 0 {
		fmt.Println("Password not entered")
		os.Exit(1)
	}
	return email, passwd
}

func main() {
	_, db, err := newConfig()
	if err != nil {
		log.Fatalf("error in config:", err)
	}

	email, passwd := UserConsoleRead()

	_, err = models.NewUser(db).Signup(nil, email, passwd, passwd)

	if err != nil {
		log.Fatalf("error in add new user:", err)
	}

	fmt.Printf("User with email %v now is created", email)
}
