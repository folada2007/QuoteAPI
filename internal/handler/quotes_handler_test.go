package handler

import (
	"QuotesAPI/internal/model"
	"QuotesAPI/pgk/dto"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func TestQuotePOST(t *testing.T) {
	h := &QuoteHandler{
		quotes:    []model.Quote{},
		quotesMux: sync.RWMutex{},
	}

	body := `{"author":"Confucius","quote":"quoteTest"}`
	req := httptest.NewRequest(http.MethodPost, "/quotes", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.QuotesPOST(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Ожидаемый статус 201, Полученный статус %d", res.StatusCode)
	}

	var resp dto.ResponseQuotes
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}

	if resp.Author != "Confucius" || resp.Quote == "" {
		t.Errorf("Ответ с неверными данными")
	}

	h.quotesMux.RLock()
	defer h.quotesMux.RUnlock()
	if len(h.quotes) != 1 {
		t.Errorf("Цитата не добавлена в срез")
	}
}

func TestQuoteGET(t *testing.T) {
	h := &QuoteHandler{
		quotes: []model.Quote{
			{ID: 1, Author: "Confucius", Quote: "quoteTest"},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes", nil)
	w := httptest.NewRecorder()

	h.QuotesGET(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый статус 200, Полученный статус %d", res.StatusCode)
	}

	var quotes []dto.ResponseQuotes
	err := json.NewDecoder(res.Body).Decode(&quotes)
	if err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}

	if len(quotes) == 0 {
		t.Errorf("Список цитат пуст")
	}
}

func TestQuotesRandom(t *testing.T) {
	h := &QuoteHandler{
		quotes: []model.Quote{
			{ID: 1, Author: "Confucius", Quote: "quoteTest"},
			{ID: 2, Author: "Socrates", Quote: "quoteTest"},
		},
		quotesMux: sync.RWMutex{},
	}

	req := httptest.NewRequest(http.MethodGet, "/quotes/random", nil)
	w := httptest.NewRecorder()

	h.QuotesRandom(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, Полученный статус %d", res.StatusCode)
	}

	var quote dto.ResponseQuotes
	err := json.NewDecoder(res.Body).Decode(&quote)
	if err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}

	if quote.ID != 1 && quote.ID != 2 {
		t.Errorf("Получен неожиданный ID: %d", quote.ID)
	}
	if quote.Quote == "" || quote.Author == "" {
		t.Errorf("Получена пустая цитата: %+v", quote)
	}
}

func TestQuotesDelete(t *testing.T) {
	h := &QuoteHandler{
		quotes: []model.Quote{
			{ID: 1, Author: "Confucius", Quote: "quoteTest"},
			{ID: 2, Author: "Socrates", Quote: "quoteTest"},
		},
		quotesMux: sync.RWMutex{},
	}

	req := httptest.NewRequest(http.MethodDelete, "/quotes/1", nil)

	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	h.QuotesDelete(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый статус 200, Полученный статус %d", res.StatusCode)
	}

	var resp map[string]string
	err := json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Fatalf("Ошибка декодирования: %v", err)
	}

	expectedMsg := "цитата успешно удалена"
	if resp["message"] != expectedMsg {
		t.Errorf("Ожидаемое сообщение %q, Полученное сообщение %q", expectedMsg, resp)
	}

	req2 := httptest.NewRequest(http.MethodDelete, "/quotes/999", nil)
	req2 = mux.SetURLVars(req2, map[string]string{"id": "999"})
	w2 := httptest.NewRecorder()

	h.QuotesDelete(w2, req2)

	res2 := w2.Result()
	if res2.StatusCode != http.StatusNotFound {
		t.Errorf("Ожидаемый статус 404, Полученный статус %d", res2.StatusCode)
	}
}
