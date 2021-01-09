package main

import (
	"log"
	"net/http"

	"github.com/higker/go-session"
	"github.com/labstack/echo/v4"
)

func init() {
	cfg := session.Config{
		CookieName:     "sessionID",
		Path:           "/",
		MaxAge:         30 * 24 * 60 * 60,
		HttpOnly:       true,
		Secure:         false,
		RedisAddr:      "127.0.0.1:6379",
		RedisPassword:  "",
		RedisDB:        0,
		RedisKeyPrefix: session.RedisPrefix,
	}
	err := session.Builder(session.Redis, &cfg)
	if err != nil {
		log.Fatal(err)
	}
}

type CartData struct {
	ProductID   int64
	ProductName string
	Qty         int64
	Price       float64
}

type SessionData struct {
	ID      int64
	Name    string
	Email   string
	Address string
	Cart    []CartData
}

func main() {
	e := echo.New()

	e.GET("/", index)
	e.POST("/set", set)
	e.GET("/get", get)
	e.GET("/del", del)
	e.GET("/clean", clean)

	e.Logger.Fatal(e.Start(":1323"))
}

func set(c echo.Context) error {
	var sessionData SessionData
	err := c.Bind(&sessionData)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	ctx, err := session.Ctx(c.Response().Writer, c.Request())
	if err != nil {
		c.Logger().Error(err.Error())
	}

	err = ctx.Set("K1", sessionData)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.HTML(http.StatusOK, "set value ok")
}

func get(c echo.Context) error {
	ctx, err := session.Ctx(c.Response().Writer, c.Request())
	if err != nil {
		c.Logger().Error(err.Error())
	}
	bytes, err := ctx.Get("K1")
	if err != nil {
		c.Logger().Error(err.Error())
	}
	sd := new(SessionData)
	//Deserialize data into objects
	err = session.DeSerialize(bytes, sd)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.JSON(http.StatusOK, sd)
}

func clean(c echo.Context) error {
	ctx, err := session.Ctx(c.Response().Writer, c.Request())
	if err != nil {
		c.Logger().Error(err.Error())
	}

	// clean session all data by session
	ctx.Clean(c.Response().Writer)

	return c.HTML(http.StatusOK, "clean data ok")
}

func del(c echo.Context) error {
	ctx, err := session.Ctx(c.Response().Writer, c.Request())
	if err != nil {
		c.Logger().Error(err.Error())
	}
	err = ctx.Del("K1")
	if err != nil {
		c.Logger().Error(err.Error())
	}
	return c.HTML(http.StatusOK, "delete v1 successful")
}

func index(c echo.Context) error {
	c.Response().Header().Add("Content-Type", "text/html")

	return c.HTML(http.StatusOK, `
	Go session storage example:<br><br>
	<a href="/set">Store key in session</a><br>
	<a href="/get">Get key value from session</a><br>
	<a href="/del">Destroy session</a><br>
	<a href="/clean">Clean session</a>
	<a href="https://github.com/higker/go-session">to github</a><br>`)
}