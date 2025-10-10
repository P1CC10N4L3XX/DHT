package dao

import (
	"DHT/src/config"
	"DHT/src/models"
	"DHT/src/utils"
	"os"
)

type ParentDAO struct {
	readFile  *os.File
	writeFile *os.File
}

func NewParentDAO() (*ParentDAO, error) {
	wf, err := os.OpenFile(config.Parent_CSV_path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	rf, err := os.OpenFile(config.Parent_CSV_path, os.O_RDONLY, 0644)
	return &ParentDAO{readFile: rf, writeFile: wf}, nil
}

func (dao *ParentDAO) Close() error {
	if err := dao.readFile.Close(); err != nil {
		return err
	}
	if err := dao.writeFile.Close(); err != nil {
		return err
	}
	return nil
}

func (dao *ParentDAO) WriteParent(node models.Node) error {
	if err := utils.WriteNodeToCSV(dao.writeFile, node); err != nil {
		return err
	}
	return nil
}

func (dao *ParentDAO) ReadParent() (models.Node, error) {
	parent, err := utils.ReadAllNodesFromCSV(dao.readFile)
	if err != nil {
		return models.Node{}, err
	}

	return parent[0], nil
}
