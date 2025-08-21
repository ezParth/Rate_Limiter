package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Status string
}

var status = []string{"Pending", "Shipped", "Delivered"}

func main() {
	var wg sync.WaitGroup

	orders := createOrder(20)

	wg.Add(3)

	go func() {
		defer wg.Done()
		processOrders(orders)
	}()

	go func() {
		defer wg.Done()
		updateOrders(orders)
	}()

	go func() {
		defer wg.Done()
		reportOrderStatus(orders)
	}()

	fmt.Println("All orders completed, Exiting!")
}

func updateOrders(orders []*Order) {
	for _, order := range orders {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		order.Status = status[rand.Intn(3)]
		fmt.Printf("Updated Orders %d, status %s \n", order.ID, order.Status)
	}
}

func reportOrderStatus(orders []*Order) {
	for i := 0; i < 5; i++ {
		time.Sleep(1000 * time.Millisecond)
		fmt.Println("\n --- Orders Status Report --- ")

		for _, order := range orders {
			fmt.Printf("Order %d: %s\n", order.ID, order.Status)
		}
	}

	fmt.Print("-------------------------------\n")
}

func processOrders(orders []*Order) {
	for _, order := range orders {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		fmt.Println("Porcessing Order: ", order.ID)
	}
}

func createOrder(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{ID: i + 1, Status: "Pending"}
	}
	return orders
}
