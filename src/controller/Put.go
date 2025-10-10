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

type PutController struct{}

func (*PutController) Put(resource models.Resource) error {
	nodeIdToStore := new(big.Int)
	if _, ok := nodeIdToStore.SetString(resource.Key, 16); !ok {
		return errors.New("conversione fallita")
	}

	path := utils.BuildPath(nodeIdToStore)

	currentAddr := os.Getenv("ENTRY_HOST") + ":" + os.Getenv("ENTRY_PORT")

	for i := 1; i < len(path); i++ {
		target := path[i]
		conn, err := grpc.Dial(currentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		defer conn.Close()
		client := pb.NewDHTClient(conn)
		ResourceReq := &pb.Resource{Key: resource.Key, Value: resource.Value}
		req := &pb.PutRequest{Resource: ResourceReq, Next: target}

		resp, err := client.PutResource(context.Background(), req)

		if err != nil {
			return err
		}

		if resp.Status == "RESOURCE_STORED" {
			fmt.Println("Risorsa VALUE=" + resource.Value + " KEY=" + resource.Key + " assegnata al nodo con id " + path[i-1])
			return nil
		} else if strings.Split(resp.Status, ":")[0] == "CONTACT_CHILD" {
			currentAddr = strings.Split(resp.Status, ":")[1] + ":" + strings.Split(resp.Status, ":")[2]
		} else if strings.Split(resp.Status, ":")[0] == "CONTACT_NEPHEW" {
			currentAddr = strings.Split(resp.Status, ":")[1] + ":" + strings.Split(resp.Status, ":")[2]
			i = utils.IndexOf(path, strings.Split(resp.Status, ":")[3])
		}
	}

	return nil
}
