package main

type Booking struct {
	Hotel       string `json:"hotel"`
	Star        int    `json:"star"`
	Price       int    `json:"price"`
	TaxesPrices int    `json:"taxes_price"`
	CheckIn     string `json:"check_in"`
	CheckOut    string `json:"check_out"`
	Guests      int    `json:"quests"`
}
