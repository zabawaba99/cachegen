package main

//go:generate cachegen -key-type=string -value-type=int

//go:generate cachegen -key-type=int -value-type=Customer
type Customer struct {
	ID int
}
