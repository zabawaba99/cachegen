package main

import (
	"log"
	"time"
)

//go:generate cachegen -key-type=string -value-type=int

//go:generate cachegen -key-type=int -value-type=Customer
type Customer struct {
	ID int
}

func main() {
	cc := NewCustomerCache(time.Second, 10*time.Second)

	cust := Customer{ID: 42}
	cc.Add(cust.ID, cust)

	time.Sleep(time.Second)

	_, ok := cc.Get(cust.ID)
	log.Println("Customer is expired:", !ok)
}
