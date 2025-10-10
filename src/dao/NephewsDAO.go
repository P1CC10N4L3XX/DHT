package dao

import (
	"DHT/src/config"
	"DHT/src/models"
	"DHT/src/utils"
	"os"
)

type NephewsDAO struct {
	file *os.File
}

func NewNephewsDAO() (*NephewsDAO, error) {
	f, err := os.OpenFile(config.Nephews_CSV_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &NephewsDAO{f}, nil
}

func (dao *NephewsDAO) Close() error {
	return dao.file.Close()
}

func (dao *NephewsDAO) WriteNephew(node models.Node) error {
	if err := utils.WriteNodeToCSV(dao.file, node); err != nil {
		return err
	}
	return nil
}

func (dao *NephewsDAO) ReadAllNephews() ([]models.Node, error) {
	err, nephews := utils.ReadAllNodesFromCSV(dao.file)
	if err != nil {
		return err, nil
	}
	return nil, nephews
}

func (dao *NephewsDAO) WriteNephews(nodes []models.Node) error {
	for _, node := range nodes {
		if err := dao.WriteNephew(node); err != nil {
			return err
		}
	}
	return nil
}

func (dao *NephewsDAO) RemoveNephew(node models.Node) error {
	if err := utils.RemoveNodeFromCSV(dao.file, node); err != nil {
		return err
	}
	return nil
}
