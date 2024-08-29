package Api

import (
	"encoding/json"
	"fmt"
	"go_currency_converter/Application"
	"go_currency_converter/Application/Clients"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	apiKey string
}

func NewServer(apiKey string) *Server {
	return &Server{
		apiKey: apiKey,
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) Start() {
	client := Clients.NewHttpClient(s.apiKey)
	baseService := Application.NewCurrencyService(client)
	cachingService := Application.NewCachingService(baseService, 10*time.Minute)

	mux := http.NewServeMux()

	mux.HandleFunc("/rates", func(w http.ResponseWriter, r *http.Request) {
		currency := r.URL.Query().Get("currency")

		rates, err := cachingService.GetRates(currency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rates)
	})

	mux.HandleFunc("/convert", func(w http.ResponseWriter, r *http.Request) {
		fromCurrency := r.URL.Query().Get("from")
		toCurrency := r.URL.Query().Get("to")
		rawAmount := r.URL.Query().Get("amount")

		amount, err := strconv.ParseFloat(rawAmount, 64)
		if err != nil {
			http.Error(w, "invalid amount", http.StatusBadRequest)
			return
		}

		convertedAmount, err := baseService.Convert(amount, fromCurrency, toCurrency)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]float64{
			"amount": convertedAmount,
		})
	})

	handler := enableCORS(mux)

	fmt.Println("running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
