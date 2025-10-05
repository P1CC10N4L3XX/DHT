package UI

import (
	"DHT/src/session"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func StartUI() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("==== DHT Binary Tree CLI ====")
	fmt.Println("Comandi disponibili:")
	fmt.Println(" leave       -> lascia la rete")
	fmt.Println(" put <file>  -> inserisci risorsa")
	fmt.Println(" get <file>  -> recupera risorsa")
	fmt.Println(" show        -> mostra info del nodo corrente (host, porta, id e risorse in gestione)")
	fmt.Println(" ping <id>    -> ping al nodo con un certo id")
	fmt.Println("==============================")

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
				fmt.Println("Uso: put <file>")
				continue
			}
			file := args[1]
			// TODO: chiama controller.Put(file)
			fmt.Printf("Inserimento risorsa: %s\n", file)

		case "get":
			if len(args) < 2 {
				fmt.Println("Uso: get <file>")
				continue
			}
			file := args[1]
			// TODO: chiama controller.Get(file)
			fmt.Printf("Recupero risorsa: %s\n", file)

		case "show":
			node := session.GetSession().Node
			fmt.Printf("Nodo attivo -> ID=%s, Host=%s, Port=%s\n", node.ID, node.Host, node.Port)

		default:
			fmt.Println("Comando sconosciuto. Digita help per la lista comandi.")
		}
	}
}
