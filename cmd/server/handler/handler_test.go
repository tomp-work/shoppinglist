package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/require"
	"github.com/tomp-work/shoppinglist/cmd/server/handler"
)

func TestGetItemList(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(nil, rec)

	h := &handler.Handler{
		Items: map[string]*handler.Item{
			"1": {Id: "1", Name: "Apple", Quantity: 5},
			"2": {Id: "2", Name: "Orange", Quantity: 3},
		},
	}

	expectedJSON := `[{"id":"1","name":"Apple","quantity":5},{"id":"2","name":"Orange","quantity":3}]`

	require.NoError(t, h.GetItemList(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, expectedJSON, rec.Body.String())
}

func TestCreateItem(t *testing.T) {
	const itemJSON = `{"id":"1","name":"Apple","quantity":5}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(itemJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handler.Handler{
		Items: map[string]*handler.Item{},
	}

	require.NoError(t, h.CreateItem(c))
	require.Equal(t, http.StatusCreated, rec.Code)
	require.JSONEq(t, itemJSON, rec.Body.String())
	require.Equal(t, h.Items, map[string]*handler.Item{"1": {Id: "1", Name: "Apple", Quantity: 5}})
}

func TestDeleteItemNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/item/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "999"}})

	h := &handler.Handler{
		Items: map[string]*handler.Item{
			"1": {Id: "1", Name: "Apple", Quantity: 5},
			"2": {Id: "2", Name: "Orange", Quantity: 3},
			"3": {Id: "3", Name: "Bread", Quantity: 2},
		},
	}

	require.NoError(t, h.DeleteItem(c))
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "id (999) not found", rec.Body.String())
}

func TestDeleteItem(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/item/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "2"}})

	h := &handler.Handler{
		Items: map[string]*handler.Item{
			"1": {Id: "1", Name: "Apple", Quantity: 5},
			"2": {Id: "2", Name: "Orange", Quantity: 3},
			"3": {Id: "3", Name: "Bread", Quantity: 1},
		},
	}

	expectedItems := map[string]*handler.Item{
		"1": {Id: "1", Name: "Apple", Quantity: 5},
		"3": {Id: "3", Name: "Bread", Quantity: 1},
	}

	require.NoError(t, h.DeleteItem(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())
	require.Equal(t, h.Items, expectedItems)
}
