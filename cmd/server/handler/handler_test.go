package handler_test

import (
	"net/http"
	"net/http/httptest"
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
