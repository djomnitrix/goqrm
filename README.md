# goqrm
## Yet another package to query on database for golang.
if you have used laravel then you must have used the query builder like `User::select("column")->where("condiiton")->get()` it is quite easy to do the query 
if you  do not know much about sql queries. This kind of query builder with editors autocomplete helps a lot.

I am from the php background and from few years working on golang so i have done some experiments to create Laravel like query builder.
if you will start with golang and wanted to do some database operation it will not be that easy to fetch from database and map those records into 
maps or struct. 

This package will give you a good start. it is just begining of package so it doesn't contains all the query builder functions but in future it will.

## Installation
`go get github.com/djomnitrix/goqrm`

for using it with mysql database look the below example.

i woud recomment first connect with database in main function. and create a a seperate model file for each of the table

`users.go`

```go

package main

import "github.com/djomnitrix/goqrm"

type Users struct {
	*goqrm.Model
	Table string
}

func UsersModel() *Users {
	return &Users{Table: "users"}
}


```
and in main function you can query like below 

`main.go`

```go
package main

import (
	"fmt"

	"github.com/djomnitrix/goqrm"
)

func main() {
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		"db_username", "db_password",
		"localhost", "3306", "Demo")

	_, err := goqrm.Connect(connString)
	if err != nil {
		panic(err)
	}

	userModel := goqrm.NewModel(UsersModel()).Where("id", "=", "2").Get()

	fmt.Println(userModel)

}

```
list of methods available 
