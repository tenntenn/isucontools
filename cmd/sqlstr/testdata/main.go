package main

import "fmt"

const sql1 = "SELECT * FROM items"

func main() {
	const sql2 = "SELECT * FROM items"
	fmt.Println(sql1, sql2)
}
