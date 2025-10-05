package utils

import (
	"DHT/src/config"
	"DHT/src/models"
	"crypto/sha1"
	"encoding/binary"
	"encoding/csv"
	"math"
	"os"
)

func Hash(input string) int64 {
	h := sha1.New()
	h.Write([]byte(input))
	bs := h.Sum(nil)
	return int64((binary.BigEndian.Uint64(bs)) % uint64(math.Pow(2, config.M)))
}

func BuildPath(id int64) []int64 {
	var path []int64
	for id > 0 {
		path = append([]int64{int64(id)}, path...)
		id = (id - 1) / 2
	}
	return append([]int64{0}, path...)
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
