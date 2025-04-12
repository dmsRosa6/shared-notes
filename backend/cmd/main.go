package main

import (
	"fmt"

	"github.com/dmsrosa/shared-notes/internal/handlers"
)

func main(){
	handler, _ := handlers.New()

	fmt.Print("Starting Server...")
	handler.Start()
}