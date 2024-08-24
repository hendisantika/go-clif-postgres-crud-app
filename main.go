package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/viper"
	"github.com/ukautz/clif"
	"log"
	"os"
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

func updateUser(c *clif.Command) error {
	if c.ParameterCount() < 3 {
		return fmt.Errorf("usage: update <id> <name> <email>")
	}
	id, err := strconv.Atoi(c.Parameter(0))
	if err != nil {
		return fmt.Errorf("invalid user ID")
	}
	name := c.Parameter(1)
	email := c.Parameter(2)

	conn, err := connectDB()
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())

	_, err = conn.Exec(context.Background(), "UPDATE users SET name=$1, email=$2 WHERE id=$3", name, email, id)
	if err != nil {
		return err
	}

	fmt.Println("User updated successfully!")
	return nil
}

func deleteUser(c *clif.Command) error {
	if c.ParameterCount() < 1 {
		return fmt.Errorf("usage: delete <id>")
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

	_, err = conn.Exec(context.Background(), "DELETE FROM users WHERE id=$1", id)
	if err != nil {
		return err
	}

	fmt.Println("User deleted successfully!")
	return nil
}

func main() {
	initConfig()

	cli := clif.New("CRUD App", "A simple CLI to perform CRUD operations", "0.1.0")

	cli.New("create", "Create a new user", createUser)
	cli.New("read", "Read user information", readUser)
	cli.New("update", "Update user information", updateUser)
	cli.New("delete", "Delete a user", deleteUser)

	if err := cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
