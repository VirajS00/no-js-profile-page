package main

import (
	"encoding/json"
	"fmt"
	"net/textproto"
	"os"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/google/uuid"
)

type ImageData struct {
	ImageName string               `json:"image_name"`
	ImageUrl  string               `json:"image_url"`
	Header    textproto.MIMEHeader `json:"header"`
	Size      int64                `json:"size"`
}

type Data struct {
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	ShowAddress bool      `json:"show_address"`
	Address     string    `json:"address"`
	DarkMode    bool      `json:"dark_mode"`
	ImageData   ImageData `json:"image_data"`
}

func main() {
	var port string

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = "3000"
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/assets", "./views/assets")
	app.Static("/images", "./images")

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
		dataByte, err := os.ReadFile(dataFile)

		if err != nil {
			return err
		}

		var data Data

		if err := json.Unmarshal(dataByte, &data); err != nil {
			return err
		}

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
			DarkMode:    data.DarkMode,
			ImageData:   data.ImageData,
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
			ImageData:   data.ImageData,
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

	app.Post("/profilepic", func(c *fiber.Ctx) error {
		file, err := c.Context().FormFile("ImageInput")

		if err != nil {
			return err
		}

		dataByte, err := os.ReadFile(dataFile)

		if err != nil {
			fmt.Println("ihdidhidhidhidhdi")
			return err
		}

		var data Data

		if err := json.Unmarshal(dataByte, &data); err != nil {
			return err
		}

		value := reflect.ValueOf(data)
		field := value.FieldByName("ImageData")

		imgName := data.ImageData.ImageName

		if field.IsValid() && imgName != "" {
			if err := os.Remove(fmt.Sprintf("./images/%s", imgName)); err != nil {
				return err
			}
		}

		uniqueId := uuid.New()

		filename := strings.Replace(uniqueId.String(), "-", "", -1)

		fileExt := strings.Split(file.Filename, ".")[1]

		image := fmt.Sprintf("%s.%s", filename, fileExt)

		if err = c.SaveFile(file, fmt.Sprintf("./images/%s", image)); err != nil {
			return err
		}

		imageUrl := fmt.Sprintf("%s/images/%s", c.BaseURL(), image)

		imageData := ImageData{
			ImageName: image,
			ImageUrl:  imageUrl,
			Header:    file.Header,
			Size:      file.Size,
		}

		newData := Data{
			Name:        data.Name,
			Email:       data.Email,
			ShowAddress: data.ShowAddress,
			Address:     data.Address,
			DarkMode:    data.DarkMode,
			ImageData:   imageData,
		}

		jsonStr, err := json.Marshal(newData)

		if err != nil {
			return err
		}

		if err := os.WriteFile(dataFile, []byte(""), os.ModeAppend); err != nil {
			return err
		}

		if err := os.WriteFile(dataFile, jsonStr, os.ModeAppend); err != nil {
			fmt.Println("writefile error")
			return err
		}

		return c.Redirect("/")
	})

	app.Listen("0.0.0.0:" + port)
}
