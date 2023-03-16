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

type ProductInfo struct {
	ID           primitive.ObjectID `bson:"_id"`
	ProductName  string             `bson:"product_name"`
	ProductPrice string             `bson:"product_price"`
	Tags         []string           `bson:"tags"`
	Comments     uint64             `bson:"comments"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

func main() {
	var db database
	var err error

	db.client, err = mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)

	db.ctx = context.Background()
	err = db.client.Connect(db.ctx)

	db.col = db.client.Database("Inventory").Collection("Products")

	mux := http.NewServeMux()
	mux.HandleFunc("/list", db.list)
	mux.HandleFunc("/price", db.price)
	mux.HandleFunc("/create", db.create) // Handler function to handle a user create
	mux.HandleFunc("/read", db.read)     // Handler function to handle a user read
	mux.HandleFunc("/update", db.update) // Handler function to handle a user update
	mux.HandleFunc("/delete", db.delete) // Handler function to handle a user delete
	log.Fatal(http.ListenAndServe(":8000", mux))
}

type database struct {
	mu     sync.RWMutex // Mutex variable
	ctx    context.Context
	client *mongo.Client
	col    *mongo.Collection
}

func (db *database) list(w http.ResponseWriter, req *http.Request) {

	db.mu.RLock()         // Locks so clients can only read
	defer db.mu.RUnlock() // Will unlock at end of function

	var product []ProductInfo

	// filter posts tagged as mongodb
	//filter := bson.M{"product_name": bson.M{"$elemMatch": bson.M{"$eq": "shoes"}}}
	filter := bson.M{"tags": bson.M{"$eq": "mongodb"}}

	item, err := db.col.Find(db.ctx, filter)

	if err == mongo.ErrNoDocuments {
		log.Fatal(err)
	}

	if err = item.All(db.ctx, &product); err == mongo.ErrNoDocuments {
		log.Fatal(err)
	}

	for _, list := range product {
		format := fmt.Sprintf("Item: %s\nPrice: %s\nTags: %s\nComments: %v\n Created at: %s\n Updated at: %s\n",
			list.ProductName, list.ProductPrice, list.Tags, list.Comments, list.CreatedAt, list.UpdatedAt)
		fmt.Fprintf(w, "%s\n", format)
	}

	log.Printf("Items: %+v\n", product)

}

func (db *database) create(w http.ResponseWriter, req *http.Request) { // Create handler function

	db.mu.Lock()         // Locks so only one client can update
	defer db.mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")   // Item is assigned from URL
	price := req.URL.Query().Get("price") // Price is assigned from URL

	// filter posts tagged as mongodb
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true
	var product ProductInfo

	err := db.col.FindOne(db.ctx, filter).Decode(&product)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			checkDatabase = false

		} else {
			checkDatabase = true
		}

	}

	if checkDatabase == true {

		fmt.Fprintf(w, "%s is already created in database\n", item)

	} else {

		_, err := strconv.ParseFloat(price, 32) // Parses price from URL

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Price is invalid\n")

		} else {

			// Insert one
			res, err := db.col.InsertOne(db.ctx, &ProductInfo{
				ID:           primitive.NewObjectID(),
				ProductName:  item,
				Tags:         []string{"mongodb"},
				ProductPrice: price,
				CreatedAt:    time.Now(),
			})

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
			}

			fmt.Fprintf(w, "Entry created -> %s: $%s\n", item, price) // Output confirmation statement

		}

	}

}

func (db *database) read(w http.ResponseWriter, req *http.Request) { // Read handler function

	db.mu.RLock()         // Locks so clients can only read
	defer db.mu.RUnlock() // Will unlock at end of function

	var product []ProductInfo

	// filter posts tagged as mongodb
	filter := bson.M{"tags": bson.M{"$elemMatch": bson.M{"$eq": "mongodb"}}}

	item, err := db.col.Find(db.ctx, filter)

	if err == mongo.ErrNoDocuments {
		log.Fatal(err)
	}

	if err = item.All(db.ctx, &product); err == mongo.ErrNoDocuments {
		log.Fatal(err)
	}

	for _, list := range product {
		format := fmt.Sprintf("Item: %s\nPrice: %s\nTags: %s\nComments: %v\n Created at: %s\n Updated at: %s\n",
			list.ProductName, list.ProductPrice, list.Tags, list.Comments, list.CreatedAt, list.UpdatedAt)
		fmt.Fprintf(w, "%s\n", format)
	}

	log.Printf("Items: %+v\n", product)

}

func (db *database) update(w http.ResponseWriter, req *http.Request) { // Update handler function

	db.mu.Lock()         // Locks so only one client can update
	defer db.mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")   // Item assigned from URL
	price := req.URL.Query().Get("price") // Price assigned from URL

	// filter posts tagged as mongodb
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true
	var product ProductInfo

	err := db.col.FindOne(db.ctx, filter).Decode(&product)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			checkDatabase = false

		} else {
			checkDatabase = true
		}

	}

	if checkDatabase == false {

		fmt.Fprintf(w, "%s is not in the database\n", item)

	} else {

		_, err := strconv.ParseFloat(price, 32) // Parse price value out of URL

		if err != nil {

			fmt.Fprintf(w, "Price is invalid\n")

		} else {

			update := bson.M{"$set": bson.M{"product_price": price}, "$currentDate": bson.M{"updated_at": time.Now()}}

			db.col.UpdateOne(db.ctx, filter, update)

			fmt.Fprintf(w, "Entry Updated -> %s: $%s\n", item, price) // Confirmation statement

		}

	}

}

func (db *database) delete(w http.ResponseWriter, req *http.Request) { // Delete handler function

	db.mu.Lock()         // Locks so only one client can update
	defer db.mu.Unlock() // Will unlock at end of function

	item := req.URL.Query().Get("item") // Get item from URL

	// filter posts tagged as mongodb
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true
	var product ProductInfo

	err := db.col.FindOne(db.ctx, filter).Decode(&product)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			checkDatabase = false

		} else {
			checkDatabase = true
		}

	}

	if checkDatabase == false {

		fmt.Fprintf(w, "%s is not the in database\n", item)

	} else {

		db.col.DeleteOne(db.ctx, filter)
		fmt.Fprintf(w, "%s has been deleted out of the database\n", item)

	}

}

func (db *database) price(w http.ResponseWriter, req *http.Request) {

	db.mu.RLock()         // Locks so clients can only read
	defer db.mu.RUnlock() // Will unlock at end of function

	item := req.URL.Query().Get("item")

	// filter posts tagged as mongodb
	filter := bson.M{"product_name": bson.M{"$eq": item}}

	checkDatabase := true
	var product ProductInfo

	err := db.col.FindOne(db.ctx, filter).Decode(&product)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			checkDatabase = false

		} else {
			checkDatabase = true
		}

	}

	if checkDatabase == false {

		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)

	} else {

		format := fmt.Sprintf("Item: %v at Price: %v\n", product.ProductName, product.ProductPrice)

		fmt.Fprintf(w, "%s\n", format)

	}

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
