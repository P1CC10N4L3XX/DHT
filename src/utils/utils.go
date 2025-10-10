package utils

import (
	"DHT/src/config"
	"DHT/src/models"
	"crypto/sha1"
	"encoding/csv"
	"fmt"
	"math/big"
	"os"
)

func Hash(input string) *big.Int {
	h := sha1.New()
	h.Write([]byte(input))
	bs := h.Sum(nil)

	num := new(big.Int).SetBytes(bs)
	mod := new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(config.M)), nil)
	num.Mod(num, mod)

	return num
}

func BuildPath(id *big.Int) []string {
	var path []string
	zero := big.NewInt(0)
	one := big.NewInt(1)
	two := big.NewInt(2)
	// Cloniamo id per non modificarlo
	current := new(big.Int).Set(id)
	for current.Cmp(zero) > 0 {
		// Aggiungiamo la stringa corrispondente
		path = append([]string{current.String()}, path...)

		// id = (id - 1) / 2
		temp := new(big.Int).Sub(current, one)
		current.Div(temp, two)
	}
	// Aggiungiamo la radice "0" all'inizio
	path = append([]string{"0"}, path...)
	return path
}

func WriteNodeToCSV(file *os.File, node models.Node) error {
	writer := csv.NewWriter(file)
	defer writer.Flush()
	record := []string{
		node.ID,
		node.Host,
		node.Port,
	}
	if err := writer.Write(record); err != nil {
		return err
	}

	return nil
}

func ReadAllNodesFromCSV(file *os.File) ([]models.Node, error) {
	reader := csv.NewReader(file)

	// Legge tutte le righe
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var nodes []models.Node
	for _, record := range records {
		if len(record) < 3 {
			continue
		}
		nodes = append(nodes, models.Node{
			ID:   record[0],
			Host: record[1],
			Port: record[2],
		})
	}

	return nodes, nil
}

func WriteResourceToCSV(file *os.File, resource models.Resource) error {
	writer := csv.NewWriter(file)
	defer writer.Flush()
	record := []string{
		resource.Key,
		resource.Value,
	}
	if err := writer.Write(record); err != nil {
		return err
	}
	return nil
}

func ReadAllResourcesFromCSV(file *os.File) ([]models.Resource, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	var resources []models.Resource
	for _, record := range records {
		if len(record) < 2 {
			continue
		}
		resources = append(resources, models.Resource{
			Key:   record[0],
			Value: record[1],
		})
	}
	return resources, nil
}

func RemoveNodeFromCSV(file *os.File, node models.Node) error {
	// Riporta il cursore all'inizio del file
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	newRecords := [][]string{}

	for _, record := range records {
		if len(record) < 3 {
			newRecords = append(newRecords, record)
			continue
		}

		if record[0] == node.ID && record[1] == node.Host && record[2] == node.Port {
			continue
		}

		newRecords = append(newRecords, record)
	}

	// Svuota e riscrive il file
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}

	writer := csv.NewWriter(file)
	if err := writer.WriteAll(newRecords); err != nil {
		return err
	}
	writer.Flush()

	return writer.Error()
}

func IndexOf(slice []string, target string) int {
	fmt.Println("target: " + target)
	for i, v := range slice {
		if v == target {
			return i
		}
	}
	return -1
}

func PrintNodesTable(nodes []models.Node) {
	fmt.Printf("%-10s | %-15s | %-6s\n", "ID", "Host", "Port")
	fmt.Println("----------------------------------------")
	for _, node := range nodes {
		fmt.Printf("%-10s | %-15s | %-6s\n", node.ID, node.Host, node.Port)
	}
}
func PrintResourcesTable(resources []models.Resource) {
	fmt.Printf("%-10s | %-20s\n", "Key", "Value")
	fmt.Println("----------------------------------")
	for _, r := range resources {
		fmt.Printf("%-10s | %-20s\n", r.Key, r.Value)
	}
}
