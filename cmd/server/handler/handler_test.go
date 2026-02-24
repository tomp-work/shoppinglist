package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tomp-work/shoppinglist/cmd/server/handler"
	"github.com/tomp-work/shoppinglist/cmd/server/models"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) Send(toEmail, content string) error {
	args := m.Called(toEmail, content)
	return args.Error(0)
}

func TestGetItemList(t *testing.T) {
	e := echo.New()
	rec := httptest.NewRecorder()
	c := e.NewContext(nil, rec)

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", SeqNum: 1, Price: 5},
			"2": {Id: "2", Name: "Orange", SeqNum: 0, Price: 10},
			"3": {Id: "3", Name: "Bread", SeqNum: 2, Price: 15},
		},
	}

	expectedJSON := `[
		{"id":"2","name":"Orange","picked":false,"seqnum": 0,"price":10},
		{"id":"1","name":"Apple","picked":false,"seqnum": 1,"price":5},
		{"id":"3","name":"Bread","picked":false,"seqnum": 2,"price":15}
	]`

	require.NoError(t, h.GetItemList(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, expectedJSON, rec.Body.String())
}

func TestCreateItem(t *testing.T) {
	const itemJSON = `{"id":"2","name":"Orange","picked":false,"seqnum":1,"price":10}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(itemJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	h := &handler.Handler{
		ItemMaxID: 1,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", SeqNum: 0, Price: 5},
		},
		ListDetails: models.ListDetails{
			TotalPrice:    5,
			SpendingLimit: 50,
		},
	}

	require.NoError(t, h.CreateItem(c))
	require.Equal(t, http.StatusCreated, rec.Code)
	require.JSONEq(t, itemJSON, rec.Body.String())
	require.Equal(t, h.Items["1"], &models.Item{Id: "1", Name: "Apple", SeqNum: 0, Price: 5})
	require.Equal(t, h.Items["2"], &models.Item{Id: "2", Name: "Orange", SeqNum: 1, Price: 10})
	require.Equal(t, 15, h.ListDetails.TotalPrice)
}

func TestDeleteItemNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/item/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "999"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", Price: 5},
			"2": {Id: "2", Name: "Orange", Price: 10},
			"3": {Id: "3", Name: "Bread", Price: 15},
		},
		ListDetails: models.ListDetails{
			TotalPrice:    30,
			SpendingLimit: 50,
		},
	}

	require.NoError(t, h.DeleteItem(c))
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "id (999) not found", rec.Body.String())
	require.Equal(t, 30, h.ListDetails.TotalPrice)
}

func TestDeleteItem(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/item/2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "2"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", Price: 5},
			"2": {Id: "2", Name: "Orange", Price: 10},
			"3": {Id: "3", Name: "Bread", Price: 15},
		},
		ListDetails: models.ListDetails{
			TotalPrice:    30,
			SpendingLimit: 50,
		},
	}

	expectedItems := map[string]*models.Item{
		"1": {Id: "1", Name: "Apple", Price: 5},
		"3": {Id: "3", Name: "Bread", Price: 15},
	}

	require.NoError(t, h.DeleteItem(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())
	require.Equal(t, h.Items, expectedItems)
	require.Equal(t, 20, h.ListDetails.TotalPrice)
}

func TestUpdateItemNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/item/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "999"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple"},
			"2": {Id: "2", Name: "Orange"},
			"3": {Id: "3", Name: "Bread"},
		},
	}

	require.NoError(t, h.UpdateItem(c))
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "id (999) not found", rec.Body.String())
}

func TestUpdateItem(t *testing.T) {
	const updateJSON = `{"picked":true}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/item/2", strings.NewReader(updateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "2"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", SeqNum: 0},
			"2": {Id: "2", Name: "Orange", SeqNum: 1},
			"3": {Id: "3", Name: "Bread", SeqNum: 2},
		},
	}

	expectedItems := map[string]*models.Item{
		"1": {Id: "1", Name: "Apple", SeqNum: 0},
		"2": {Id: "2", Name: "Orange", SeqNum: 1, Picked: true},
		"3": {Id: "3", Name: "Bread", SeqNum: 2},
	}

	require.NoError(t, h.UpdateItem(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"id":"2","name":"Orange","seqnum":1,"picked":true,"price":0}`, rec.Body.String())
	require.Equal(t, h.Items, expectedItems)
}

func TestMoveItemUpNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/item/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id/up")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "999"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple"},
			"2": {Id: "2", Name: "Orange"},
			"3": {Id: "3", Name: "Bread"},
		},
	}

	require.NoError(t, h.MoveItemUp(c))
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "id (999) not found", rec.Body.String())
}

func TestMoveItemUpAlreadyTop(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/item/2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id/up")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "1"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", SeqNum: 0},
			"2": {Id: "2", Name: "Orange", SeqNum: 1},
			"3": {Id: "3", Name: "Bread", SeqNum: 2},
		},
	}

	expectedItems := map[string]*models.Item{
		"1": {Id: "1", Name: "Apple", SeqNum: 0},
		"2": {Id: "2", Name: "Orange", SeqNum: 1},
		"3": {Id: "3", Name: "Bread", SeqNum: 2},
	}

	require.NoError(t, h.MoveItemUp(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())
	require.Equal(t, h.Items, expectedItems)
}

func TestMoveItemUp(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/item/2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id/up")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "2"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", SeqNum: 0},
			"2": {Id: "2", Name: "Orange", SeqNum: 1},
			"3": {Id: "3", Name: "Bread", SeqNum: 2},
		},
	}

	expectedItems := map[string]*models.Item{
		"1": {Id: "1", Name: "Apple", SeqNum: 1},
		"2": {Id: "2", Name: "Orange", SeqNum: 0},
		"3": {Id: "3", Name: "Bread", SeqNum: 2},
	}

	require.NoError(t, h.MoveItemUp(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())
	require.Equal(t, h.Items, expectedItems)
}

func TestMoveItemDownNotFound(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/item/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id/down")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "999"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple"},
			"2": {Id: "2", Name: "Orange"},
			"3": {Id: "3", Name: "Bread"},
		},
	}

	require.NoError(t, h.MoveItemDown(c))
	require.Equal(t, http.StatusNotFound, rec.Code)
	require.Equal(t, "id (999) not found", rec.Body.String())
}

func TestMoveItemDownAlreadyBottom(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/item/2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id/down")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "3"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", SeqNum: 0},
			"2": {Id: "2", Name: "Orange", SeqNum: 1},
			"3": {Id: "3", Name: "Bread", SeqNum: 2},
		},
	}

	expectedItems := map[string]*models.Item{
		"1": {Id: "1", Name: "Apple", SeqNum: 0},
		"2": {Id: "2", Name: "Orange", SeqNum: 1},
		"3": {Id: "3", Name: "Bread", SeqNum: 2},
	}

	require.NoError(t, h.MoveItemDown(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())
	require.Equal(t, h.Items, expectedItems)
}

func TestMoveItemDown(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/item/2", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/item/:id/down")
	c.SetPathValues(echo.PathValues{{Name: "id", Value: "2"}})

	h := &handler.Handler{
		ItemMaxID: 3,
		Items: map[string]*models.Item{
			"1": {Id: "1", Name: "Apple", SeqNum: 0},
			"2": {Id: "2", Name: "Orange", SeqNum: 1},
			"3": {Id: "3", Name: "Bread", SeqNum: 2},
		},
	}

	expectedItems := map[string]*models.Item{
		"1": {Id: "1", Name: "Apple", SeqNum: 0},
		"2": {Id: "2", Name: "Orange", SeqNum: 2},
		"3": {Id: "3", Name: "Bread", SeqNum: 1},
	}

	require.NoError(t, h.MoveItemDown(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Empty(t, rec.Body.String())
	require.Equal(t, h.Items, expectedItems)
}

func TestGetListDetails(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/list")

	h := &handler.Handler{
		ListDetails: models.ListDetails{
			TotalPrice:    50,
			SpendingLimit: 200,
		},
	}

	require.NoError(t, h.GetListDetails(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"totalprice":50,"spendingLimit":200}`, rec.Body.String())
}

func TestUpdateListDetails(t *testing.T) {
	const updateJSON = `{"spendingLimit":350}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(updateJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/list")

	h := &handler.Handler{
		ListDetails: models.ListDetails{
			SpendingLimit: 200,
			TotalPrice:    150,
		},
	}

	require.NoError(t, h.UpdateListDetails(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, models.ListDetails{SpendingLimit: 350, TotalPrice: 150}, h.ListDetails)
	require.JSONEq(t, `{"spendingLimit":350,"totalprice":150}`, rec.Body.String())
}

func TestSendListEmailInvalidAddress(t *testing.T) {
	const emailJSON = `{"emailAddress":"invalid"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(emailJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/list/send")

	h := &handler.Handler{
		Items: map[string]*models.Item{
			"1": {Id: "1", SeqNum: 1, Name: "Gin", Price: 30},
			"2": {Id: "2", SeqNum: 0, Name: "Whiskey", Price: 25},
		},
		ListDetails: models.ListDetails{
			TotalPrice: 55,
		},
	}

	require.Error(t, h.SendListEmail(c))
}

const expectedEmail = `Dear Tom,

Please can you pick up the following shopping:

- Whiskey (£25)
- Gin (£30)

The total price of the shop should be £55.

Thanks,

Clare
`

func TestSendListEmail(t *testing.T) {
	const emailJSON = `{"emailAddress":"may.dup@example.com"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(emailJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/list/send")

	mockSender := new(MockEmailSender)
	mockSender.On("Send", "may.dup@example.com", expectedEmail).Return(nil).Once()

	h := &handler.Handler{
		Items: map[string]*models.Item{
			"1": {Id: "1", SeqNum: 1, Name: "Gin", Price: 30},
			"2": {Id: "2", SeqNum: 0, Name: "Whiskey", Price: 25},
		},
		ListDetails: models.ListDetails{
			TotalPrice: 55,
		},
		EmailSender: mockSender,
	}

	require.NoError(t, h.SendListEmail(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "", rec.Body.String())
	mockSender.AssertExpectations(t)
}
