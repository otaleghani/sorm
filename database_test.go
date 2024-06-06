package sorm

import (
	"fmt"
	"github.com/otaleghani/spg"
	"strconv"
	"testing"
)

type Item struct {
	Id      string
	Desc    string
	Price   float64
	InStock bool
	CatId   string
}

type Cat struct {
	Id      string
	CatName string
	CatDesc string
}

type ItemCat struct {
	Id      string
	Desc    string
	Price   float64
	InStock bool
	IdCat   string
	IdCat2  string
	CatName string
	CatDesc string
}

var path = "test.db"

func Test_Delete(t *testing.T) {
	if err := DeleteDatabase(path); err != nil {
		t.Fatal(err)
	}
}

func Test_Create(t *testing.T) {
	fmt.Println("im here")
	if err := CreateDatabase("test.db"); err != nil {
		t.Fatal(err)
	}
}

func Test_CreateTable(t *testing.T) {
	if err := CreateTable(Item{}); err != nil {
		t.Fatal(err)
	}
	if err := CreateTable(Cat{}); err != nil {
		t.Fatal(err)
	}
}

func Test_Insert(t *testing.T) {
	generator := spg.New("en-usa")
	var opt = spg.Options{}

	numItems := 10
	numCat := 10
	items := make([]interface{}, 0, numItems+numCat)
	for i := 0; i < numItems; i++ {
		items = append(items, Item{
			Id:      generator.Product().UUID(opt),
			Desc:    generator.Person().FirstName(opt),
			Price:   0.99,
			InStock: generator.Boolean(),
			CatId:   strconv.Itoa(i),
		})
	}
	for i := 0; i < numCat; i++ {
		items = append(items, Cat{
			Id:      strconv.Itoa(i),
			CatName: generator.Product().Technology(opt),
			CatDesc: generator.Product().ProductDescription(opt),
		})
	}

	if err := InsertInto(items...); err != nil {
		t.Fatal(err)
	}
}

func Test_Select(t *testing.T) {
	var selectResult []Item
	for i := 0; i < 1; i++ {
		if err := Select(&selectResult, ""); err != nil {
			t.Fatal(err)
		}
	}
}

func Test_Delete_Items(t *testing.T) {
	err := Delete(Item{}, "InStock = ?", false)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Select_After_Delete(t *testing.T) {
	var selectResult []Item
	if err := Select(&selectResult, ""); err != nil {
		t.Fatal(err)
	}
	for _, row := range selectResult {
		fmt.Println(row)
	}
}

func Test_Join(t *testing.T) {
	var joinResult []ItemCat
	if err := Join(&joinResult, Item{}, Cat{}, "Item.CatId = Cat.Id"); err != nil {
		t.Fatal(err)
	}
	for _, val := range joinResult {
		fmt.Println(val)
	}
}

func Test_Update(t *testing.T) {
	items := make([]interface{}, 0, 1)
	items = append(items, Item{
		Id:      "sandro",
		Desc:    "something",
		Price:   0.99,
		InStock: false,
		CatId:   "1",
	})

	if err := InsertInto(items...); err != nil {
		t.Fatal(err)
	}

	updatedItem := Item{
		Id:   "sandro",
		Desc: "something else",
		// Price:        0.99,
		// InStock:      false,
		// CatId:        "1",
	}
	err := Update(updatedItem, "id = ?", "sandro")
	if err != nil {
		t.Fatal(err)
	}

	var selectResult []Item
	if err := Select(&selectResult, "id = ?", "sandro"); err != nil {
		t.Fatal(err)
	}
	for _, row := range selectResult {
		fmt.Println(row)
	}
}
