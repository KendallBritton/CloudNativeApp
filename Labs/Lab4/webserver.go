package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create) // Add comment
	mux.HandleFunc("/read", db.read)     // Add comment
	mux.HandleFunc("/update", db.update) // Add comment
	mux.HandleFunc("/delete", db.delete) // Add comment
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) { // Add comment
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	_, checkDatabase := db[item]

	if checkDatabase == true {

		fmt.Fprintf(w, "%s is already created in database", item)

	} else {

		tempParse, _ := strconv.ParseFloat(price, 32)

		tempFloat32 := float32(tempParse)

		convertedPrice := dollars(tempFloat32)

		db[item] = convertedPrice

		fmt.Fprintf(w, "Entry created: %s: $%s\n", item, price)

	}

}

func (db database) read(w http.ResponseWriter, req *http.Request) { // Add comment
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) { // Add comment
	//item := req.URL.Query().Get("item")
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) { // Add comment
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
