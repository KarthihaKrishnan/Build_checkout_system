package config

import (
	"database/sql"
	"fmt"
)



// Handling error
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
