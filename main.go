package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vchatchai/refresh/db"
	"github.com/vchatchai/refresh/model"
	"github.com/vchatchai/refresh/tcp"
	"gopkg.in/yaml.v2"
)

var config model.Config
var database db.DB
var client tcp.HttpClient

func main() {

	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal(err)
	}

	d, err := sql.Open("sqlserver", dataSource())
	if err != nil {
		log.Fatal(err)
	}
	d.SetMaxOpenConns(0)
	d.SetMaxIdleConns(10)
	defer d.Close()

	database = db.NewDB(d)

	// bookings := csv.ReadCSVBooking(config)
	bookings, err := database.GetBooking(config)
	if err != nil {
		log.Fatal("GetBooking", err)
	}

	// fmt.Println("requestBody", string(requestBody))

	// jsonStr, _ := json.Marshal(bookings)

	// ioutil.WriteFile("data.json", jsonStr, 0644)

	resp, err := client.Post(config.URLRefresh+"/booking", "application/json", bookings)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response GetBooking")
	log.Println(string(body))

	// client.ApiToken = "test"

	// users := csv.ReadLadenUser(config)
	users, err := database.GetUser()
	if err != nil {
		log.Fatal("GetUser", err)
	}

	// log.Println("requestBody", string(requestBody))

	resp, err = client.Post(config.URLRefresh+"/user", "application/json", users)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("response GetUser")
	log.Println(string(body))
	/*
		// containers := csv.ReadLadenContainer(config)
		containers, err := database.GetContainer()

		if err != nil {
			log.Fatal("GetContainer", err)
		}

		// log.Println("requestBody", string(requestBody))

		resp, err = client.Post(config.URLRefresh+"/container", "application/json", containers)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("response GetContainer")
		log.Println(string(body))

	*/

}

func dataSource() string {
	// host := "localhost"
	// pass := "manager1"
	// if os.Getenv("profile") == "prod" {
	// 	host = "db"
	// 	pass = os.Getenv("db_pass")
	// }

	if config.Database.Debug {
		// fmt.Printf(" password:%s\n", config.Database.Password)
		fmt.Printf(" port:%d\n", config.Database.DatabasePort)
		fmt.Printf(" server:%s\n", config.Database.Server)
		fmt.Printf(" user:%s\n", config.Database.User)
		fmt.Printf(" database:%s\n", config.Database.Database)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=disable", config.Database.Server, config.Database.User, config.Database.Password, config.Database.DatabasePort, config.Database.Database)
	if config.Database.Debug {
		fmt.Printf(" connString:%s\n", connString)
	}

	return connString

	// return "goxygen:" + pass + "@tcp(" + host + ":3306)/goxygen"
}
