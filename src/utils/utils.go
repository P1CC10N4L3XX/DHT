package utils

import (
	"DHT/src/config"
	"DHT/src/models"
	"crypto/sha1"
	"encoding/csv"
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
