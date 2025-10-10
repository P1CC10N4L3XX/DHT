package dao

import (
	"DHT/src/config"
	"DHT/src/models"
	"DHT/src/utils"
	"os"
)

type ChildsDAO struct {
	file *os.File
}

func NewChildsDAO() (*ChildsDAO, error) {
	f, err := os.OpenFile(config.Childs_CSV_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &ChildsDAO{f}, nil
}

func (dao *ChildsDAO) Close() error {
	return dao.file.Close()
}

func (dao *ChildsDAO) WriteChild(node models.Node) error {
	if err := utils.WriteNodeToCSV(dao.file, node); err != nil {
		return err
	}
	return nil
}

func (dao *ChildsDAO) WriteChilds(nodes []models.Node) error {
	for _, node := range nodes {
		if err := utils.WriteNodeToCSV(dao.file, node); err != nil {
			return err
		}
	}
	return nil
}

func (dao *ChildsDAO) RemoveChild(node models.Node) error {
	if err := utils.RemoveNodeFromCSV(dao.file, node); err != nil {
		return err
	}
	return nil
}

func (dao *ChildsDAO) ReadAllChilds() ([]models.Node, error) {
	err, childs := utils.ReadAllNodesFromCSV(dao.file)
	if err != nil {
		return err, nil
	}
	return nil, childs
}
