package handler

import (
	"QuotesAPI/internal/model"
	"QuotesAPI/internal/service/idGen"
	"QuotesAPI/pgk/dto"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

var rndGen = rand.New(rand.NewSource(time.Now().UnixNano()))

type QuoteHandler struct {
	quotes    []model.Quote
	quotesMux sync.RWMutex
}

func NewQuoteHandler() *QuoteHandler {
	return &QuoteHandler{
		make([]model.Quote, 0),
		sync.RWMutex{},
	}
}

func (h *QuoteHandler) QuotesPOST(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req dto.RequestQuote

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if r.Body == http.NoBody {
		sendError(w, http.StatusBadRequest, "Bad Request", "пустое тело запроса")
		return
	}

	err := decoder.Decode(&req)
	if err != nil {
		log.Println("ошибка декодирования")
		sendError(w, http.StatusBadRequest, "Bad Request", "неверный формат данных")
		return
	}

	err = quotesValidation(req)
	if err != nil {
		log.Println("ошибка валидации: " + err.Error())
		sendError(w, http.StatusBadRequest, "Validation Error", err.Error())
		return
	}

	newQuote := model.Quote{
		ID:     idGen.NextID(),
		Author: req.Author,
		Quote:  req.Quote,
	}

	h.quotesMux.Lock()
	defer h.quotesMux.Unlock()

	h.quotes = append(h.quotes, newQuote)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(dto.ResponseQuotes{
		ID:     newQuote.ID,
		Quote:  newQuote.Quote,
		Author: newQuote.Author,
	})
}

func (h *QuoteHandler) QuotesGET(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.quotesMux.RLock()
	defer h.quotesMux.RUnlock()

	if len(h.quotes) == 0 {
		log.Println("список цитат пуст!")
		sendError(w, http.StatusInternalServerError, "Server error", "внутреняя ошибка сервера")
		return
	}

	author := r.URL.Query().Get("author")
	var responseQuote []dto.ResponseQuotes

	if author != "" {
		responseQuote = filterQuotesByAuthor(author, h.quotes)

		if len(responseQuote) == 0 {
			sendError(w, http.StatusNotFound, "Not Found", "цитаты с таким автором не найдены")
			return
		}
	} else {
		responseQuote = make([]dto.ResponseQuotes, len(h.quotes))
		for i, q := range h.quotes {
			responseQuote[i] = dto.ResponseQuotes{
				ID:     q.ID,
				Quote:  q.Quote,
				Author: q.Author,
			}
		}
	}

	err := json.NewEncoder(w).Encode(responseQuote)
	if err != nil {
		log.Println("ошибка кодирования ответа" + err.Error())
		http.Error(w, "Server error", http.StatusInternalServerError)
	}
}

func (h *QuoteHandler) QuotesRandom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	h.quotesMux.RLock()
	defer h.quotesMux.RUnlock()

	if len(h.quotes) == 0 {
		log.Println("ошибка нет доступных к выдаче цитат!")
		sendError(w, http.StatusNotFound, "Server error", "нет доступных к выдаче цитат")
		return
	}

	randomIndex := rndGen.Intn(len(h.quotes))
	randomQuote := h.quotes[randomIndex]

	responseQuote := dto.ResponseQuotes{
		ID:     randomQuote.ID,
		Quote:  randomQuote.Quote,
		Author: randomQuote.Author,
	}

	err := json.NewEncoder(w).Encode(responseQuote)
	if err != nil {
		log.Println("ошибка кодирования данных!" + err.Error())
		sendError(w, http.StatusInternalServerError, "Server error", "не удалось загрузить ответ")
	}
}

func (h *QuoteHandler) QuotesDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	if idStr == "" {
		log.Println("пропущен параметр id")
		sendError(w, http.StatusBadRequest, "Bad Request", "ID обязателен для заполнения")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("ошибка конвертации" + err.Error())
		sendError(w, http.StatusBadRequest, "Bad Request", "не удалось преобразовать данные из параметра id")
		return
	}

	if id == 0 {
		log.Println("id должен быть > 0")
		sendError(w, http.StatusBadRequest, "Bad Request", "id должен быть > 0")
		return
	}

	h.quotesMux.Lock()
	defer h.quotesMux.Unlock()

	indexForDelete := -1
	for i, quote := range h.quotes {
		if quote.ID == int64(id) {
			indexForDelete = i
			break
		}
	}

	if indexForDelete == -1 {
		log.Printf("цитата не по id :%d не найдена", id)
		sendError(w, http.StatusNotFound, "Not Found", "цитата не по указанному id не найдена")
		return
	}

	h.quotes = append(h.quotes[:indexForDelete], h.quotes[indexForDelete+1:]...)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode("цитата успешно удаленна")
	if err != nil {
		log.Println("ошибка кодирования данных" + err.Error())
	}
}

func quotesValidation(validQuote dto.RequestQuote) error {
	if validQuote.Author == "" {
		return errors.New("поле Author не может быть пустым")
	}
	if validQuote.Quote == "" {
		return errors.New("поле Quote не может быть пустым")
	}
	return nil
}

func filterQuotesByAuthor(author string, quotes []model.Quote) []dto.ResponseQuotes {
	var filtered []dto.ResponseQuotes

	for _, q := range quotes {
		if strings.EqualFold(q.Author, author) {
			filtered = append(filtered, dto.ResponseQuotes{
				ID:     q.ID,
				Quote:  q.Quote,
				Author: q.Author,
			})
		}
	}
	return filtered
}

func sendError(w http.ResponseWriter, status int, errorType string, msg string) {
	w.WriteHeader(status)
	errResp := dto.ResponseErrors{[]dto.Errors{{
		errorType, msg,
	}},
	}
	err := json.NewEncoder(w).Encode(errResp)
	if err != nil {
		log.Println("ошибка кодирования ответа" + err.Error())
	}
}
