package cleanup

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)


func Run(db *sqlx.DB) {
	fmt.Println("Cleanup started")
	for {
		select {
		case <-time.After(time.Minute * 10):
			db.MustExec("DELETE FROM codes")
			fmt.Printf("codes cleaned up")

		}
	}
}
