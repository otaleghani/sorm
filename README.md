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
    // Create a new database at given path
    dbPath := "test.db"
    db, err := sorm.CreateDatabase(dbPath)
    if err != nil {
        fatal(err)
    }

    // Create inside of the database a table using struct Product
    err = db.CreateTable(Product{})
    if err != nil {
        fatal(err)
    }

    // Create dummy data by using spg (Simple Placeholder Generator).
    // You can find it in my Github page. In this example is just used
    // to generate random dummy data to fill the database.
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

    // Insert every item in the database, regardless of the type.
    err = db.InsertInto(items)
    if err != nil {
        fatal(err)
    }

    // Here the database is queried to get in return all of the items
    // that have InStock set to true
    selectDest := []Product{}
    err = db.Select(selectDest, "InStock = ?", true)
    if err != nil {
        fatal(err)
    }
    
    // Deletes the database at given path
    err = db.DeleteDatabase(dbPath)
    if err != nil {
        fatal(err)
    }
}
```

## Constraints

`sorm` has some constraints to work. The first is that the primary key
in every table will be the field `Id`. So you'll need to add an `Id`
field inside of your types to add a primary key.

## Flags

You can specify some flags to add constraints to your SQL tables. Every
flag is spciefied with a suffix in the type name. This suffix needs to
be preceded by an underscore, like this: `Name_u string`. There are
several flags that you can use.

- `u` adds the `UNIQUE` constraint;
- `n` adds the `NOT NULL`;
- `nu` adds both `UNIQUE` and `NOT NULL`;
- `id` creates a `FOREIGN KEY` taking the first part of the name as
  the name of the table and linking the `Id` of said table to this
  field.

### Example FOREIGN KEY

``` go
type Category struct {
    Id string
    Name string
}

type Product struct {
    Id string
    Name string
    Quantity int
    Category_id string
}
```

To create a foreign key you need to specify the name of the table that
you want to link and then the suffix `_id`.


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
- [ ] Show the tables inside of a given db
