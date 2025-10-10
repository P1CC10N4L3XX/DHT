package dao

import (
	"DHT/src/config"
	"DHT/src/models"
	"DHT/src/utils"
	"os"
)

type ResourceDAO struct {
	file *os.File
}

func NewResourceDAO() (*ResourceDAO, error) {
	f, err := os.OpenFile(config.Resource_CSV_path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &ResourceDAO{f}, nil
}

func (dao *ResourceDAO) Close() error {
	return dao.file.Close()
}

func (dao *ResourceDAO) WriteResource(resource models.Resource) error {
	return utils.WriteResourceToCSV(dao.file, resource)
}

func (dao *ResourceDAO) ReadResourceByKey(key string) (models.Resource, error) {
	resources, err := utils.ReadAllResourcesFromCSV(dao.file)
	if err != nil {
		return models.Resource{}, err
	}
	for _, resource := range resources {
		if resource.Key == key {
			return resource, nil
		}
	}
	return models.Resource{}, nil
}

func (dao *ResourceDAO) ReadAllResources() ([]models.Resource, error) {
	resources, err := utils.ReadAllResourcesFromCSV(dao.file)
	if err != nil {
		return []models.Resource{}, err
	}
	return resources, nil
}
