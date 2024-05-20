package orm

import (
	"fmt"
	"testing"
)

type Item struct {
	Id          int
	Description string
	Price       float64
	InStock     bool
}

func Test_Database(t *testing.T) {
	if err := New("test.db"); err != nil {
		t.Fatal(err)
	}
	if err := CreateTable(Item{}); err != nil {
		t.Fatal(err)
	}

	item := Item{Id: 1, Description: "A sample item", Price: 9.99, InStock: true}
	if err := InsertInto(item); err != nil {
		t.Fatal(err)
	}

	var selectResult []Item
	if err := Select(&selectResult); err != nil {
		t.Fatal(err)
	}

	for _, item := range selectResult {
		fmt.Printf("Item: %d - %s - %f - %t\n", item.Id, item.Description, item.Price, item.InStock)
	}

}
