# Simple Object-Relational Mapping

A simple and easy to use Golang sqlite3 ORM.

## Installation

Add to your imports spg.

``` go
import "github.com/otaleghani/sorm"
```

Afterwards run go get github.com/otaleghani/sorm to get the package.

## Basic usage

Here is a basic example to get you started with sorm.

``` go
package main

import ( 
    "github.com/otaleghani/sorm"
    "github.com/otaleghani/spg" // Simple Placeholder Generator
)

type Product struct {
    Id string
    Name string
    Stock int
    InStock bool
}

func main() {
    dbPath := "test.db"
    err := sorm.CreateDatabase(dbPath)
    if err != nil {
        fatal(err)
    }
    // Create a new database at given path

    err = sorm.CreateTable(Product{})
    if err != nil {
        fatal(err)
    }
    // Create inside of the database a table using struct Product

    gen := spg.New("en-usa") 
    opt := spg.Options{}
    var items []interface{}
    for i := 0; i < 100000; i++ {
        items = append(items, Product{
            Id: gen.Product().UUID(opt),
            Name: gen.Product().ProductName(opt),
            Stock: gen.RandomNumber(200),
            InStock: gen.Boolean(),
        })
    }
    // Create dummy data by using spg (Simple Placeholder Generator). You can find it in my Github page. In this example is just used to generate random dummy data to fill the database.

    err = sorm.InsertInto(items)
    if err != nil {
        fatal(err)
    }
    // Insert every item in the database, regardless of the type.

    selectDest := []Product{}
    err = sorm.Select(selectDest, "InStock = ?", true)
    if err != nil {
        fatal(err)
    }
    // Here the database is queried to get in return all of the items that have InStock set to true
    
    err = DeleteDatabase(dbPath)
    if err != nil {
        fatal(err)
    }
    // Deletes the database at given path
}
```

## Benchmarks

If you batch the inserts you can get around 100.000 writes in just one second. 
As for the read speeds you can get around 500.000 reads in one second. This benchmarks are made with relatively small objects, so results may vary. 

## Disclaimer

This library is made to make developing an MVP faster and more efficient, it is not meant for production so use at you own risk.

## To do

- [x] Update record
- [ ] Image and files management
- [ ] Backups
- [ ] Migration
