package controller

import (
	"DHT/src/config"
	"DHT/src/dao"
	"DHT/src/models"
	"DHT/src/session"
	"DHT/src/utils"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"os"
	"strings"
	"time"

	pb "DHT/src/proto/stubs"
)

type JoinController struct{}

func (j *JoinController) InitConnectionAsEntry() error {
	hostname, _ := os.Hostname()
	node := models.Node{
		ID:   "0",
		Port: config.Port,
		Host: hostname,
	}

	session.GetSession().Node = &node
	return nil
}

func (j *JoinController) InitConnection() error {
	hostname, _ := os.Hostname()
	ts := time.Now().Unix()
	meta := fmt.Sprintf("%s-%d", hostname, ts)
	id := utils.Hash(meta)
	fmt.Println(id)
	path := utils.BuildPath(id)
	currentAddr := os.Getenv("ENTRY_HOST") + ":" + os.Getenv("ENTRY_PORT")
	for i := 1; i < len(path); i++ {
		target := path[i]
		conn, err := grpc.Dial(currentAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		client := pb.NewDHTClient(conn)
		req := &pb.JoinRequest{Host: hostname, Port: config.Port, Next: target}

		resp, err := client.JoinNode(context.Background(), req)
		if err != nil {
			return err
		}
		if strings.Split(resp.Status, ":")[0] == "NEED_CHILD" {
			parent := models.Node{
				ID:   fmt.Sprintf("%s", path[i-1]),
				Port: strings.Split(currentAddr, ":")[1],
				Host: strings.Split(currentAddr, ":")[0],
			}
			node := models.Node{
				ID:   strings.Split(resp.Status, ":")[1],
				Port: config.Port,
				Host: hostname,
			}
			if err := initState(resp, node, parent); err != nil {
				return err
			}
			return nil
		} else if strings.Split(resp.Status, ":")[0] == "CONTACT_CHILD" {
			currentAddr = strings.Split(resp.Status, ":")[1] + ":" + strings.Split(resp.Status, ":")[2]
		}
		conn.Close()
	}
	return errors.New("JOIN_FAILED")
}

func initState(resp *pb.JoinResponse, node models.Node, parent models.Node) error {
	parentDao, err := dao.NewParentDAO()
	if err != nil {
		return err
	}
	defer parentDao.Close()
	childsDao, err := dao.NewChildsDAO()
	if err != nil {
		return err
	}
	defer childsDao.Close()
	nephewsDao, err := dao.NewNephewsDAO()
	if err != nil {
		return err
	}
	defer nephewsDao.Close()

	if err := parentDao.WriteParent(parent); err != nil {
		return err
	}
	session.GetSession().Node.ID = node.ID
	session.GetSession().Node.Host = node.Host
	session.GetSession().Node.Port = node.Port

	for i := 0; i < len(resp.Childs); i++ {
		child := models.Node{
			ID:   resp.Childs[i].Id,
			Host: resp.Childs[i].Host,
			Port: resp.Childs[i].Port,
		}
		if err := childsDao.WriteChild(child); err != nil {
			return err
		}
	}
	for i := 0; i < len(resp.Nephews); i++ {
		nephew := models.Node{
			ID:   resp.Nephews[i].Id,
			Host: resp.Nephews[i].Host,
			Port: resp.Nephews[i].Port,
		}
		if err := nephewsDao.WriteNephew(nephew); err != nil {
			return err
		}
	}

	return nil
}
