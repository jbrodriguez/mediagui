package storage

import (
	"database/sql"
	"fmt"
)

func commit(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		fmt.Printf("w: COMMIT: %s\n", err)
	}
}
