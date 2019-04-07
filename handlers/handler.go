package handlers

import (
	"github.com/labstack/echo"
	"github.com/roboncode/go-urlshortener/stores"
	"github.com/speps/go-hashids"
	"net/http"
	"strconv"
)

const (
	// :: Internal ::
	MissingRequiredUrlMsg = `Missing required property "url"`
)

type Handler struct {
	Store  stores.Store
	HashID *hashids.HashID
}

func (h *Handler) CreateLink(c echo.Context) error {
	var body = new(struct {
		Url string `json:"url"`
	})

	if err := c.Bind(&body); err != nil {
		return err
	}

	if body.Url == "" {
		return c.JSON(http.StatusBadRequest, MissingRequiredUrlMsg)
	}

	counter := h.Store.IncCount()
	if code, err := h.HashID.Encode([]int{counter}); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else if link, err := h.Store.Create(code, body.Url); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	} else {
		return c.JSON(http.StatusCreated, link)
	}
}

func (h *Handler) GetLinks(c echo.Context) error {
	skip, _ := strconv.Atoi(c.QueryParam("s"))
	limit, _ := strconv.Atoi(c.QueryParam("l"))
	links := h.Store.List(limit, skip)
	return c.JSON(http.StatusOK, links)
}

func (h *Handler) GetLink(c echo.Context) error {
	if link, err := h.Store.Read(c.Param("code")); err != nil {
		return c.NoContent(http.StatusNotFound)
	} else {
		return c.JSON(http.StatusOK, link)
	}
}

func (h *Handler) DeleteLink(c echo.Context) error {
	if count := h.Store.Delete(c.Param("code")); count == 0 {
		return c.NoContent(http.StatusNotFound)
	} else {
		return c.NoContent(http.StatusOK)
	}
}

func (h *Handler) RedirectToUrl(c echo.Context) error {
	if link, err := h.Store.Read(c.Param("code")); err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/404")
		//return nil
	} else {
		return c.Redirect(http.StatusMovedPermanently, link.LongUrl)
	}
}
