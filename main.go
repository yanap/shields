package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
)

type Config struct {
	GoPath      string `envconfig:"GOPATH"`
	Go111Module string `envconfig:"GO111MODULE"`
	GoEnvShell  string `envconfig:"GOENV_SHELL"`
}

type CustomContext struct {
	echo.Context
	Config
}

func ConfigMiddleware(config Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(*CustomContext)
			cc.Config = config
			return next(cc)
		}
	}
}

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}

func main() {
	var config Config
	err := envconfig.Process("config", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c, config}
			return next(cc)
		}
	})
	e.Use(ConfigMiddleware(config))
	e.GET("/env", Env)
	e.Start(":1234")
}

func Env(c echo.Context) error {
	cc := c.(*CustomContext)
	c.String(200, cc.GetEnv())
	return nil
}

func (c Config) GetEnv() string {
	return fmt.Sprintf(
		"%s:%s:%s",
		c.GoPath,
		c.Go111Module,
		c.GoEnvShell,
	)
}
