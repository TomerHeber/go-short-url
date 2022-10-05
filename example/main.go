package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/TomerHeber/go-short-url"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/pkg/browser"
)

type CreateUrlRequest struct {
	Url string `json:"url" validate:"required,url"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"shortUrl"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func main() {
	s, err := short.NewShortener()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	e.File("/", "example/public/index.html")
	e.File("/index.html", "example/public/index.html")

	e.POST("/create", func(c echo.Context) error {
		cur := new(CreateUrlRequest)
		if err := c.Bind(cur); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := c.Validate(cur); err != nil {
			return err
		}

		surl, err := s.CreateShortenedUrl(c.Request().Context(), cur.Url)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, &CreateUrlResponse{
			ShortUrl: surl.GetUrl(),
		})
	})

	e.GET("/:id", func(c echo.Context) error {
		id := c.Param("id")

		if len(id) > 7 {
			return echo.NewHTTPError(http.StatusNotFound)
		}

		url, err := s.GetUrlFromShortenedUrlId(c.Request().Context(), id)
		if err != nil {
			if errors.Is(err, &short.IdNotFoundError{}) {
				return echo.NewHTTPError(http.StatusNotFound, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.Redirect(http.StatusMovedPermanently, url)
	})

	//nolint
	go browser.OpenURL("http://localhost:8080")
	e.Logger.Fatal(e.Start(":8080"))
}
