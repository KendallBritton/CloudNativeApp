package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://172.17.0.2:27017" // Find this from the Mongo container
)

// Structure to hold the database entries
type ProductInfo struct {
	ID           primitive.ObjectID `bson:"_id"`
	ProductName  string             `bson:"product_name"`
	ProductPrice string             `bson:"product_price"`
	Tags         []string           `bson:"tags"`
	Comments     uint64             `bson:"comments"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

// Main function
func main() {
	var db database
	var err error

	db.client, err = mongo.NewClient( // Client connects to mongo database
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)

	db.ctx = context.Background() // Provides context of database
	err = db.client.Connect(db.ctx)

	db.col = db.client.Database("Inventory").Collection("Products") // Establishes collection within database

	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)     // Handler function to handle a user list
	mux.HandleFunc("/price", db.price)   // Handler function to handle a user price
	mux.HandleFunc("/create", db.create) // Handler function to handle a user create
	mux.HandleFunc("/read", db.read)     // Handler function to handle a user read
	mux.HandleFunc("/update", db.update) // Handler function to handle a user update
	mux.HandleFunc("/delete", db.delete) // Handler function to handle a user delete
	log.Fatal(http.ListenAndServe(":8000", mux))
}

// Database structure to hold all client contents
type database struct {
	mu     sync.RWMutex // Mutex variable
	ctx    context.Context
	client *mongo.Client
	col    *mongo.Collection
}

// List function operation
func (db *database) list(w http.ResponseWriter, req *http.Request) {

	db.mu.RLock()         // Locks so clients can only read
	defer db.mu.RUnlock() // Will unlock at end of function

	var product []ProductInfo // Variable to hold results from search

	// filter posts tagged as mongodb
	filter := bson.M{"tags": bson.M{"$eq": "mongodb"}}

	item, err := db.col.Find(db.ctx, filter) // Performs search of database

	if err == mongo.ErrNoDocuments { // Ends if no documents
		log.Fatal(err)
	}

	if err = item.All(db.ctx, &product); err == mongo.ErrNoDocuments { // Ends if no documents in results
		log.Fatal(err)
	}

	for _, list := range product {
		format := fmt.Sprintf("Item: %s\nPrice: %s\nTags: %s\nComments: %v\n Created at: %s\n Updated at: %s\n", // Formats result to print form
			list.ProductName, list.ProductPrice, list.Tags, list.Comments, list.CreatedAt, list.UpdatedAt)
		fmt.Fprintf(w, "%s\n", format)
	}

	log.Printf("Items: %+v\n", product) // Outputs print statement

}

// Create function operation
func (db *database) create(w http.ResponseWriter, req *http.Request) {

	db.mu.Lock()         // Locks so only one client can update
	defer db.mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")   // Item is assigned from URL
	price := req.URL.Query().Get("price") // Price is assigned from URL

	// filters for item given by user
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true   // Assume item is in database
	var product ProductInfo // Variable to hold product info

	err := db.col.FindOne(db.ctx, filter).Decode(&product) // Performs search of database

	if err != nil {

		if err == mongo.ErrNoDocuments { // Item not in database
			checkDatabase = false

		} else {
			checkDatabase = true // Item in database
		}

	}

	if checkDatabase == true {

		fmt.Fprintf(w, "%s is already created in database\n", item) // Item already in database

	} else {

		_, err := strconv.ParseFloat(price, 32) // Parses price from URL (checks if price is valid)

		if err != nil {
			w.WriteHeader(http.StatusNotFound) // Stops insertion if price is invalid
			fmt.Fprintf(w, "Price is invalid\n")

		} else {

			// Insert one product into database
			res, err := db.col.InsertOne(db.ctx, &ProductInfo{
				ID:           primitive.NewObjectID(),
				ProductName:  item,
				Tags:         []string{"mongodb"},
				ProductPrice: price,
				CreatedAt:    time.Now(),
			})

			if err != nil { // Ends if there is an error
				log.Fatal(err)
			} else {
				fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex()) // Confirmation of insertion with ID
			}

			fmt.Fprintf(w, "Entry created -> %s: $%s\n", item, price) // Output confirmation statement

		}

	}

}

// Read function operation
func (db *database) read(w http.ResponseWriter, req *http.Request) {

	db.mu.RLock()         // Locks so clients can only read
	defer db.mu.RUnlock() // Will unlock at end of function

	var product []ProductInfo

	// filter posts tagged as mongodb
	filter := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "mongodb"}}}

	item, err := db.col.Find(db.ctx, filter) // Performs search of database

	if err == mongo.ErrNoDocuments { // Ends if there is match in database
		log.Fatal(err)
	}

	if err = item.All(db.ctx, &product); err == mongo.ErrNoDocuments { // Ends if there is no match in copying
		log.Fatal(err)
	}

	for _, list := range product {
		format := fmt.Sprintf("Item: %s\nPrice: %s\nTags: %s\nComments: %v\n Created at: %s\n Updated at: %s\n", // Formats output for print
			list.ProductName, list.ProductPrice, list.Tags, list.Comments, list.CreatedAt, list.UpdatedAt)
		fmt.Fprintf(w, "%s\n", format)
	}

	log.Printf("Items: %+v\n", product) // Prints output

}

// Update function operation
func (db *database) update(w http.ResponseWriter, req *http.Request) {

	db.mu.Lock()         // Locks so only one client can update
	defer db.mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")   // Item assigned from URL
	price := req.URL.Query().Get("price") // Price assigned from URL

	// filter item given by user
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true   // Variable to check if item is in database
	var product ProductInfo // Variable to hold product info

	err := db.col.FindOne(db.ctx, filter).Decode(&product) // Performs search of database

	if err != nil {

		if err == mongo.ErrNoDocuments { // Item is not in database
			checkDatabase = false

		} else {
			checkDatabase = true // Item is in database
		}

	}

	if checkDatabase == false {

		fmt.Fprintf(w, "%s is not in the database\n", item) // Prints item not in database

	} else {

		_, err := strconv.ParseFloat(price, 32) // Parse price value out of URL (checks if price is valid)

		if err != nil {

			fmt.Fprintf(w, "Price is invalid\n") // Prints price is invalid

		} else {

			update := bson.M{"$set": bson.M{"product_price": price}, "$currentDate": bson.M{"updated_at": time.Now()}} // Update info

			db.col.UpdateOne(db.ctx, filter, update) // Performs update in database

			fmt.Fprintf(w, "Entry Updated -> %s: $%s\n", item, price) // Confirmation statement

		}

	}

}

// Delete function operation
func (db *database) delete(w http.ResponseWriter, req *http.Request) {

	db.mu.Lock()         // Locks so only one client can update
	defer db.mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item") // Get item from URL

	// filter for item given by user
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true   // Variable to check item in database
	var product ProductInfo // Variable to hold product info

	err := db.col.FindOne(db.ctx, filter).Decode(&product) // Performs search in database

	if err != nil {

		if err == mongo.ErrNoDocuments { // Item not in database
			checkDatabase = false

		} else {
			checkDatabase = true // Item in database
		}

	}

	if checkDatabase == false {

		fmt.Fprintf(w, "%s is not the in database\n", item) // Prints item not in database

	} else {

		db.col.DeleteOne(db.ctx, filter) // Deletes item out of database
		fmt.Fprintf(w, "%s has been deleted out of the database\n", item)

	}

}

// Price function operation
func (db *database) price(w http.ResponseWriter, req *http.Request) {

	db.mu.RLock()         // Locks so clients can only read
	defer db.mu.RUnlock() // Will unlock at end of function

	item := req.URL.Query().Get("item") // Gets item from URL

	// filters for item given by user
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true   // Variable to check for item in database
	var product ProductInfo // Variable to hold product info

	err := db.col.FindOne(db.ctx, filter).Decode(&product) // Performs search on database

	if err != nil {

		if err == mongo.ErrNoDocuments { // Item not in database
			checkDatabase = false

		} else {
			checkDatabase = true // Item in database
		}

	}

	if checkDatabase == false {

		w.WriteHeader(http.StatusNotFound)         // 404
		fmt.Fprintf(w, "no such item: %q\n", item) // Indicates no item in database, so no price

	} else {

		format := fmt.Sprintf("Item: %v at Price: %v\n", product.ProductName, product.ProductPrice) // Formats output to be printed

		fmt.Fprintf(w, "%s\n", format) // Outputs price of item

	}

}

// Error check function
func checkError(err error) {
	if err != nil {
		log.Fatal(err) // Ends if error
	}
}
