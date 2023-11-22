package storage

import (
	"database/sql"
	"fmt"
)

func rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		fmt.Printf("w: ROLLBACK: %s\n", err)
	}
}
