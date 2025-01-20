package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"sync"
	"unicode"
	"strconv"
	"errors"
	)

type Receipt struct {
	Retailer		string	`json:"retailer"`
	PurchaseDate	string	`json:"purchaseDate"`
	PurchaseTime	string	`json:"purchaseTime"`
	Items			[]item	`json:"items"`
	Total			string	`json:"total"`
}

type Item struct {
	Description		string	`json:"shortDescription"`
	price			string	`json:"price"`
}

func process(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receipt)
	if err != nil || r.Method != http.MethodPost {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	} 

	id := receipt.Retailer + receipt.PurchaseDate + receipt.purchaseTime

	points, err := calculatePoints(&receipt)
	if err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	} 

	// Add points to map

	http.HandleFunc("/receipts/" + p.id + "/points", getPoints)	
}

func calculatePoints(r *Receipt) (points int, err error) {
	var points int

	// One point for every alphanumeric character in the name
	points += countAlphanumeric(r.Retailer)

	// 50 points if total is an even dollar amount
	// 25 points if total is a multiple of 0.25
	total_raw, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		return 0, errors.New("Error converting total")
	}
	total := int(total_raw * 100)
	if total % 100 == 0 {
		points += 50
	}
	if total % 25 == 0 {
		points += 25
	}

	// 5 points for every 2 items on the receipt
	points += 5 * (len(r.Items) / 2) 

	// If item desc - whitespace is a multiple of 3, points += ceil(price *.02)
	

	// Nice try

	// If day is odd, add 6

	// If time is between 2:00pm and 4:00 pm add 10, assuming strictly between
	
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

func processTotal(p int) int 

func getPoints(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "points: 10\n")
}

func main() {
	http.HandleFunc("/receipts/process", process)
}
