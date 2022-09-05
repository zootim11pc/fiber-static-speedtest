package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"image"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"net/url"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	app := fiber.New()
	app.Use(favicon.New(), logger.New())
	app.Use(cors.New())

	app.Get("/", getImg)
	app.Get("img/logo.abbc38cf.png", getImg)

	app.Listen(":3001")
}

func randomColor() color.NRGBA {
	return color.NRGBA{R: uint8(rand.Intn(255)), G: uint8(rand.Intn(255)), B: uint8(rand.Intn(255)), A: 255}
}

func base64img(width, height, blocks int) string {
	var cellSizeX int = width / blocks
	var cellSizeY int = height / blocks

	colors := make([][]color.NRGBA, blocks)
	for i, _ := range colors {
		colors[i] = make([]color.NRGBA, blocks)
		for j := 0; j < blocks; j++ {
			colors[i][j] = randomColor()
		}
	}

	// 32 x 4  (2)
	// cellSizeX 16
	// cellSizeY 2
	// color[16][2]

	m := image.NewNRGBA(image.Rectangle{Min: image.Point{0, 0}, Max: image.Point{width, height}})
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var xx int = x / cellSizeX
			var yy int = y / cellSizeY
			// fmt.Println(x, y, xx, yy)
			var c = colors[xx][yy]
			m.SetNRGBA(x, y, c)
		}
	}

	buf := bytes.NewBuffer([]byte{})
	if err := png.Encode(buf, m); err != nil {
		fmt.Println(err)
	}
	b64str := base64.StdEncoding.EncodeToString(buf.Bytes())
	return b64str
}

func getImg(c *fiber.Ctx) error {
	min := 0
	max := 5

	u, err := url.Parse(c.BaseURL())
	if err != nil {
		log.Fatal(err)
	}

	c.Set(fiber.HeaderCacheControl, "no-cache")

	switch u.Hostname() {
	case os.Getenv("DOMAIN_A"):
		min = 1
		max = 6
	case os.Getenv("DOMAIN_B"):
		min = 500
		max = 800
	case os.Getenv("DOMAIN_C"):
		min = 1500
		max = 1800
	}

	sleepTime := rand.Intn(max-min) + min

	time.Sleep(time.Duration(sleepTime) * time.Millisecond)

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	c.Set(fiber.HeaderAccessControlAllowOrigin, "*")
	return c.SendString("<html><body><img src=\"data:image/png;base64," + base64img(256, 256, 16) + "\" /></body></html>")

}
