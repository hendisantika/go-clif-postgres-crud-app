package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"github.com/ukautz/clif"
	"log"
	"strconv"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}
}

var config Config

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}

func connectDB() (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.Database.User, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createUser(c *clif.Command) error {
	if c.ParameterCount() < 2 {
		return fmt.Errorf("usage: create <name> <email>")
	}
	name := c.Parameter(0)
	email := c.Parameter(1)

	conn, err := connectDB()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return err
	}

	fmt.Println("User created successfully!")
	return nil
}

func readUser(c *clif.Command) error {
	if c.ParameterCount() < 1 {
		return fmt.Errorf("usage: read <id>")
	}
	id, err := strconv.Atoi(c.Parameter(0))
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}

	conn, err := connectDB()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	var name, email string
	err = conn.QueryRow(context.Background(), "SELECT name, email FROM users WHERE id=$1", id).Scan(&name, &email)
	if err != nil {
		return err
	}

	fmt.Printf("User: %s, Email: %s\n", name, email)
	return nil
}

func main() {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.
	s := "gopher"
	fmt.Println("Hello and welcome, %s!", s)

	for i := 1; i <= 5; i++ {
		//TIP You can try debugging your code. We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>. To start your debugging session,
		// right-click your code in the editor and select the <b>Debug</b> option.
		fmt.Println("i =", 100/i)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
