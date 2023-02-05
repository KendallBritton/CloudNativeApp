package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create) // Handler function to handle a user create
	mux.HandleFunc("/read", db.read)     // Handler function to handle a user read
	mux.HandleFunc("/update", db.update) // Handler function to handle a user update
	mux.HandleFunc("/delete", db.delete) // Handler function to handle a user delete
	log.Fatal(http.ListenAndServe("localhost:8000", mux))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {

	var mu sync.RWMutex // Mutex variable

	mu.RLock()         // Locks so clients can only read
	defer mu.RUnlock() // Will unlock at end of function

	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}
}

func (db database) create(w http.ResponseWriter, req *http.Request) { // Create handler function

	var mu sync.RWMutex // Mutex variable

	mu.Lock()         // Locks so only one client can update
	defer mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")   // Item is assigned from URL
	price := req.URL.Query().Get("price") // Price is assigned from URL

	_, checkDatabase := db[item] // Checks database for item

	if checkDatabase == true {

		fmt.Fprintf(w, "%s is already created in database\n", item)

	} else {

		tempParse, err := strconv.ParseFloat(price, 32) // Parses price from URL

		if err != nil {

			fmt.Fprintf(w, "Price is invalid\n")

		} else {

			tempFloat32 := float32(tempParse) // Converts type from float64 to float32

			db[item] = dollars(tempFloat32) // Assigns price to map

			fmt.Fprintf(w, "Entry created -> %s: $%s\n", item, price) // Output confirmation statement

		}

	}

}

func (db database) read(w http.ResponseWriter, req *http.Request) { // Read handler function

	var mu sync.RWMutex // Mutex variable

	mu.RLock()         // Locks so clients can only read
	defer mu.RUnlock() // Will unlock at end of function

	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price) // Prints map values
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) { // Update handler function

	var mu sync.RWMutex // Mutex variable

	mu.Lock()         // Locks so only one client can update
	defer mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")   // Item assigned from URL
	price := req.URL.Query().Get("price") // Price assigned from URL

	_, checkDatabase := db[item] // Checks database for item

	if checkDatabase == false {

		fmt.Fprintf(w, "%s is not in the database\n", item)

	} else {

		tempParse, err := strconv.ParseFloat(price, 32) // Parse price value out of URL

		if err != nil {

			fmt.Fprintf(w, "Price is invalid\n")

		} else {

			tempFloat32 := float32(tempParse) // Convert float64 to float32

			db[item] = dollars(tempFloat32) // Assign value to temp

			fmt.Fprintf(w, "Entry Updated -> %s: $%s\n", item, price) // Confirmation statement

		}

	}

}

func (db database) delete(w http.ResponseWriter, req *http.Request) { // Delete handler function

	var mu sync.RWMutex // Mutex variable

	mu.Lock()         // Locks so only one client can update
	defer mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item") // Get item from URL

	_, checkDatabase := db[item] // Check database for item

	if checkDatabase == false {

		fmt.Fprintf(w, "%s is not the in database\n", item)

	} else {

		delete(db, item)
		fmt.Fprintf(w, "%s has been deleted out of the database\n", item)

	}

}

func (db database) price(w http.ResponseWriter, req *http.Request) {

	var mu sync.RWMutex // Mutex variable

	mu.RLock()         // Locks so clients can only read
	defer mu.RUnlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
