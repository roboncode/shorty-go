package handlers

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/roboncode/go-urlshortener/consts"
	"github.com/roboncode/go-urlshortener/helpers"
	"github.com/roboncode/go-urlshortener/models"
	"github.com/roboncode/go-urlshortener/stores/mocks"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	mockDB = map[int]models.Link{
		1: {
			ID:       1,
			Code:     "abc",
			LongUrl:  "https://roboncode.com",
			ShortUrl: "https://ac.me/abc",
		},
		2: {
			ID:       2,
			Code:     "def",
			LongUrl:  "https://google.com",
			ShortUrl: "https://ac.me/def",
		},
	}
	createJSON = `{"url":"https://roboncode.com"}`
)

func TestHandler_CreateLink(t *testing.T) {
	viper.SetDefault(consts.HashMin, 5)

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(createJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/shorten")

	ID := 1
	firstMockDB := mockDB[1]

	mockStore := mocks.Store{}
	mockStore.On("IncCount").Return(int64(ID))
	mockStore.On("Create", "lejRe", "https://roboncode.com").Return(&firstMockDB, nil)

	h := &Handler{
		Store:  &mockStore,
		HashID: helpers.NewHashID(),
	}

	// Assertions
	if assert.NoError(t, h.CreateLink(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, string(firstMockDB.EncodeLink()), strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandler_GetLink(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/links/:code")
	c.SetParamNames("code")
	c.SetParamValues("abc")

	firstMockDB := mockDB[1]

	mockStore := mocks.Store{}
	mockStore.On("Read", c.Param("code")).Return(&firstMockDB, nil)

	h := &Handler{
		Store:  &mockStore,
		HashID: helpers.NewHashID(),
	}

	// Assertions
	if assert.NoError(t, h.GetLink(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(firstMockDB.EncodeLink()), strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandler_GetLinks(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/links")

	links := []models.Link{mockDB[1], mockDB[2]}
	limit := int64(0)
	skip := int64(0)

	mockStore := mocks.Store{}
	mockStore.On("List", limit, skip).Return(links)

	h := &Handler{
		Store:  &mockStore,
		HashID: helpers.NewHashID(),
	}

	strLinks, _ := json.Marshal(links)

	// Assertions
	if assert.NoError(t, h.GetLinks(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(strLinks), strings.TrimSpace(rec.Body.String()))
	}
}

func TestHandler_DeleteLink(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/links/:code")
	c.SetParamNames("code")
	c.SetParamValues("abc")

	firstMockDB := mockDB[1]

	mockStore := mocks.Store{}
	mockStore.On("Read", c.Param("code")).Return(&firstMockDB, nil)
	mockStore.On("Delete", c.Param("code")).Return(nil)

	h := &Handler{
		Store:  &mockStore,
		HashID: helpers.NewHashID(),
	}

	// Assertions
	if assert.NoError(t, h.GetLink(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
