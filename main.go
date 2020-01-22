package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var bin = []byte("bin")

func setupRouter(db *bolt.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		var val []byte
		err := db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket(bin)
			if bucket == nil {
				fmt.Println(
					"msg", "can't find the bucket bin",
				)
				return status.Error(codes.Internal, "Can't find the bucket")
			}
			val = bucket.Get([]byte("status"))

			return nil
		})
		if err != nil {
			fmt.Println("err", err, "msg", "can't fetch data")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while fetching data"})
		}
		c.JSON(http.StatusOK, gin.H{"status": string(val)})
	})

	return r
}

func main() {
	db, err := bolt.Open("itrash.db", 0644, nil)
	if err != nil {
		fmt.Println(
			"err", err,
			"msg", "can't open the db",
		)
		os.Exit(2)
	}
	defer db.Close()

	r := setupRouter(db)
	r.Run(":8080")
}
