package controller

import (
	"DHT/src/models"
	pb "DHT/src/proto/stubs"
	"DHT/src/utils"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/big"
	"os"
	"strings"
)

type GetController struct{}

func (getController *GetController) Get(key string) (models.Resource, error) {
	bigIntKey := new(big.Int)
	_, ok := bigIntKey.SetString(key, 16)
	if !ok {
		return models.Resource{}, errors.New("invalid key")
	}
	path := utils.BuildPath(bigIntKey)
	currentAddr := fmt.Sprintf("%s:%s", os.Getenv("ENTRY_HOST"), os.Getenv("ENTRY_PORT"))
	for i := 1; i < len(path); i++ {
		next := path[i]
		conn, err := grpc.Dial(currentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return models.Resource{}, err
		}
		defer conn.Close()
		client := pb.NewDHTClient(conn)
		req := &pb.GetRequest{Key: key, Next: next}
		resp, err := client.GetResource(context.Background(), req)
		if err != nil {
			return models.Resource{}, err
		}

		if resp.Status == "RESOURCE_NOT_FOUND" {
			return models.Resource{}, errors.New("resource not found")
		} else if resp.Status == "RESOURCE_DETECTED" {
			return models.Resource{Key: resp.Resource.Key, Value: resp.Resource.Value}, nil
		} else if strings.Split(resp.Status, ":")[0] == "CONTACT_CHILD" {
			host := strings.Split(resp.Status, ":")[1]
			port := strings.Split(resp.Status, ":")[2]
			currentAddr = fmt.Sprintf("%s:%s", host, port)
			println(currentAddr)
		}

	}
	return models.Resource{}, errors.New("resource not found")
}
