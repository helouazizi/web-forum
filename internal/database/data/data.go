package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -lsqlite3
#include "sqlite_helper.c"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	// Open the database
	dbFile := C.CString("test.db")
	defer C.free(unsafe.Pointer(dbFile))

	db := C.open_database(dbFile)
	if db == nil {
		fmt.Println("Failed to open database")
		return
	}
	fmt.Println("Database opened successfully!")

	// Create a table
	sql := C.CString("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT);")
	defer C.free(unsafe.Pointer(sql))
	C.execute_sql(db, sql)

	// Insert a user
	insertSQL := C.CString("INSERT INTO users (name) VALUES ('Alice');")
	defer C.free(unsafe.Pointer(insertSQL))
	C.execute_sql(db, insertSQL)

	fmt.Println("Inserted user: Alice")

	// Fetch users
	fmt.Println("Fetching users from database:")
	C.fetch_users(db)

	// Close the database
	C.close_database(db)
}
