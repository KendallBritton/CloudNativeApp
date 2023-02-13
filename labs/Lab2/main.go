package main

import (
	"fmt"
	"labs/Lab2/lru"
)

func main() {

	testlru := lru.NewCache(3)

	testkey := "key"
	testval := "val"
	testlru.Put(testkey, testval)

	fmt.Println(testlru.Get(testkey))

	testlru = lru.NewCache(3)

	testlru.Put("key1", "val1")
	testlru.Put("key2", "val2")
	testlru.Put("key3", "val3")
	testlru.Put("key4", "val4")

	testlru.Get("key1")

	fmt.Println(testlru.Get("key1"))

}
