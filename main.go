package main

import (
	"encoding/json"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

type Data struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	ShowAddress bool   `json:"show_address"`
	Address     string `json:"address"`
	DarkMode    bool   `json:"dark_mode"`
}

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/assets", "./views/assets")

	dataFile := "./data/data.json"

	app.Get("/", func(c *fiber.Ctx) error {
		dataByte, err := os.ReadFile(dataFile)

		if err != nil {
			return err
		}

		var data Data

		if err := json.Unmarshal(dataByte, &data); err != nil {
			return err
		}

		return c.Render("index", data)
	})

	app.Post("/", func(c *fiber.Ctx) error {
		name := string(c.Context().FormValue("Name"))
		email := string(c.Context().FormValue("Email"))
		address := string(c.Context().FormValue("Address"))
		showAddressStr := string(c.Context().FormValue("ShowAddress"))

		var showAddress bool

		if showAddressStr == "showAddess" {
			showAddress = true
		}

		d := Data{
			Name:        name,
			Email:       email,
			Address:     address,
			ShowAddress: showAddress,
		}

		jsonStr, err := json.Marshal(d)

		if err != nil {
			return err
		}

		if err := os.WriteFile(dataFile, []byte(""), os.ModeAppend); err != nil {
			return err
		}

		if err := os.WriteFile(dataFile, jsonStr, os.ModeAppend); err != nil {
			return err
		}

		return c.Redirect("/")
	})

	app.Post("/darkmode", func(c *fiber.Ctx) error {
		dataByte, err := os.ReadFile(dataFile)

		if err != nil {
			return err
		}

		var data Data

		if err := json.Unmarshal(dataByte, &data); err != nil {
			return err
		}

		darkModeToggle := string(c.Context().FormValue("DarkMode"))

		darkMode := !(darkModeToggle == "darkMode")

		newData := Data{
			Name:        data.Name,
			Email:       data.Email,
			ShowAddress: data.ShowAddress,
			Address:     data.Address,
			DarkMode:    darkMode,
		}

		jsonStr, err := json.Marshal(newData)

		if err != nil {
			return err
		}

		if err := os.WriteFile(dataFile, []byte(""), os.ModeAppend); err != nil {
			return err
		}

		if err := os.WriteFile(dataFile, jsonStr, os.ModeAppend); err != nil {
			return err
		}

		return c.Redirect("/")

	})

	app.Listen(":3000")
}
