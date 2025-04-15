package main

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/mux"
)

var pointsMap map[string]int

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	Description string `json:"shortDescription"`
	Price       string `json:"price"`
}

func process(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt

	json.NewDecoder(r.Body).Decode(&receipt)

	id := removeNonAlphanumeric(receipt.Retailer) + receipt.PurchaseDate + receipt.PurchaseTime

	points, err := calculatePoints(&receipt)
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	pointsMap[id] = points

	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculatePoints(r *Receipt) (int, error) {
	var points int

	// One point for every alphanumeric character in the name
	points += countAlphanumeric(r.Retailer)

	// 50 points if total is an even dollar amount
	// 25 points if total is a multiple of 0.25
	total_raw, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return 0, errors.New("error converting total")
	}
	total := int(total_raw * 100)
	if total%100 == 0 {
		points += 50
	}
	if total%25 == 0 {
		points += 25
	}

	// 5 points for every 2 items on the receipt
	points += 5 * (len(r.Items) / 2)

	// If item desc - whitespace is a multiple of 3, points += ceil(price *.02)
	for _, item := range r.Items {
		desc := strings.TrimSpace(item.Description)
		if len(desc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				return 0, errors.New("error converting price")
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	// If day is odd, add 6
	purchaseDate, err := time.Parse(time.DateOnly, r.PurchaseDate)
	if err != nil {
		return 0, errors.New("error parsing purchase date")
	}
	if purchaseDate.Day()%2 == 1 {
		points += 6
	}

	// If time is between 2:00pm and 4:00 pm add 10, assuming strictly between
	purchaseTime, err := time.Parse("15:04", r.PurchaseTime)
	if err != nil {
		return 0, errors.New("error parsing purchase time")
	}
	if ((purchaseTime.Hour() >= 14 && purchaseTime.Minute() > 0) || purchaseTime.Hour() > 14) && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points, nil
}

func countAlphanumeric(s string) int {
	count := 0
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			count++
		}
	}
	return count
}

func removeNonAlphanumeric(s string) string {
	var cleaned strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			cleaned.WriteRune(r)
		}
	}
	return cleaned.String()
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	points, present := pointsMap[id]
	if !present {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}
	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {

	pointsMap = make(map[string]int)

	router := mux.NewRouter()
	router.HandleFunc("/receipts/process", process).Methods(http.MethodPost)
	router.HandleFunc("/receipts/{id}/points", getPoints).Methods(http.MethodGet)

	http.ListenAndServe(":8080", router)
}
