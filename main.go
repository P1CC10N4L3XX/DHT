package main

import (
	"DHT/src/controller"
	"log"
	"os"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <port>")
	}
	join := controller.Join{}
	if join.InitConnection(os.Args[1]) == -1 {
		log.Fatal("Failed to join")
	}

}
