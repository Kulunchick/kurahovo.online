package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/jasonlvhit/gocron"
	"github.com/kulunchick/kurahovo.online/utils"
)

func test() {
	insert, err := db.Query("INSERT INTO Test(value, date) VALUES ( ?, NOW())", 28.8)

	if err != nil {
		log.Fatal(err.Error())
	}

	defer insert.Close()

	fmt.Println("Test")
}

var db *sql.DB
var f *os.File

func SetupLogger() {
	var err error
	f, err = os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		panic(err)
	}

	wrt := io.MultiWriter(os.Stdout, f)

	log.SetOutput(wrt)
}

func SetubJobs() {
	gocron.Every(145).Seconds().Do(test)
}

func main() {
	var err error
	app := fiber.New()

	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?parseTime=true", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")))

	if err != nil {
		panic(err)
	}

	app.Get("/api/latest", func(c *fiber.Ctx) error {
		data, errors := utils.GetJatestData(db)

		if len(errors) != 0 {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"errors": errors,
			})
		}

		return c.Status(fiber.StatusOK).JSON(data)
	})

	SetubJobs()
	SetupLogger()
	test()

	defer f.Close()
	defer db.Close()

	go gocron.Start()
	panic(app.Listen(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))))
}
