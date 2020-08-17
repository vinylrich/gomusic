package main

import (
	"log"

	"github.com/ajtwoddltka/gomusic/backend/src/rest"
)

func main() {
	log.Println("Main log.....")
	rest.RunAPI("127.0.0.1:8000")
}
