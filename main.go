package main

import (
	"DHT/src/UI"
	"DHT/src/controller"
	"DHT/src/models"
	pb "DHT/src/proto/stubs"
	"DHT/src/session"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func startGRPCServer(node models.Node) {
	lis, err := net.Listen("tcp", ":"+node.Port)
	if err != nil {
		log.Fatalf("Errore nell'apertura della porta: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterDHTServer(grpcServer, &controller.DhtServer{})

	log.Printf("Nodo con id %s e hostname %s in ascolto sulla porta %s...\n", node.ID, node.Host, node.Port)
	if err := grpcServer.Serve(lis); err != nil {

		log.Fatalf("Errore nell'avvio del server gRPC: %v", err)
	}
}

func closeTerminal() {
	fmt.Println("Chiudo il terminale...")
	cmd := "exit"
	_ = syscall.Exec("/bin/sh", []string{"sh", "-c", cmd}, os.Environ())
}

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// Goroutine per gestire il segnale
	go func() {
		<-sigs
		fmt.Println("\n🛑 Chiusura richiesta...")

		closeTerminal()
	}()

	join := controller.JoinController{}
	if len(os.Args) > 1 && os.Args[1] == "-entry" {
		if err := join.InitConnectionAsEntry(); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := join.InitConnection(); err != nil {
			log.Fatal(err)
		}
	}

	go startGRPCServer(*session.GetSession().Node)

	fmt.Printf("\n\n%s:%s --> Effettuata la Join all'interno della rete con id %s\n\n", session.GetSession().Node.Host, session.GetSession().Node.Port, session.GetSession().Node.ID)
	UI.StartUI()
}
