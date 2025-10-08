package UI

import (
	"DHT/src/controller"
	"DHT/src/models"
	"DHT/src/session"
	"DHT/src/utils"
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func StartUI() {
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	showHelp()

	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		args := strings.Split(input, " ")
		cmd := args[0]

		switch cmd {

		case "leave":
			// TODO: implementa leave
			fmt.Println("Leave non ancora implementato.")

		case "put":
			if len(args) < 2 {
				fmt.Println("Uso: put <value>")
				continue
			}
			value := args[1]
			key := utils.Hash(value).Text(16)
			putController := controller.PutController{}
			resource := models.Resource{Value: value, Key: key}
			err := putController.Put(resource)
			if err != nil {
				log.Fatal(err)
			}

		case "get":
			if len(args) < 2 {
				fmt.Println("Uso: get <key>")
				continue
			}
			key := args[1]
			getController := controller.GetController{}
			resource, err := getController.Get(key)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("risorsa trovata Key:" + resource.Key + " Value:" + resource.Value)
			}

		case "show":
			node := session.GetSession().Node
			fmt.Printf("Nodo attivo -> ID=%s, Host=%s, Port=%s\n", node.ID, node.Host, node.Port)
		case "help":
			showHelp()
		default:
			fmt.Println("Comando sconosciuto. Digita help per la lista comandi.")
		}
	}
}

func showHelp() {
	fmt.Println("==== DHT Binary Tree CLI ====")
	fmt.Println("Comandi disponibili:")
	fmt.Println(" leave       -> lascia la rete")
	fmt.Println(" put <file>  -> inserisci risorsa")
	fmt.Println(" get <file>  -> recupera risorsa")
	fmt.Println(" show        -> mostra info del nodo corrente (host, porta, id e risorse in gestione)")
	fmt.Println(" ping <id>   -> ping al nodo con un certo id")
	fmt.Println(" help        -> mostra legenda comandi")
	fmt.Println("==============================")
}
