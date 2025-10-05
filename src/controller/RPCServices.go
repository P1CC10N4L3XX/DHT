package controller

import (
	"DHT/src/dao"
	"DHT/src/models"
	pb "DHT/src/proto/stubs"
	"DHT/src/session"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
	"sync"
)

type DhtServer struct {
	pb.UnimplementedDHTServer
	mu sync.Mutex
}

func (s *DhtServer) JoinNode(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	childsDao, err := dao.NewChildsDAO()
	defer childsDao.Close()
	if err != nil {
		return nil, err
	}
	defer childsDao.Close()
	childs, err := childsDao.ReadAllChilds()
	if err != nil {
		return nil, err
	}
	fmt.Println(len(childs))
	if len(childs) < 2 {
		return addNodeAsChild(req, childs, childsDao)
	}
	for _, child := range childs {
		if child.ID == fmt.Sprintf("%d", req.Next) {
			return &pb.JoinResponse{Status: fmt.Sprintf("CONTACT_CHILD:%s:%s", child.Host, child.Port)}, nil
		}
	}
	return nil, nil
}

func (s *DhtServer) ChangeParent(ctx context.Context, req *pb.ChangeParentRequest) (*pb.ChangeParentResponse, error) {
	parentDao, err := dao.NewParentDAO()
	if err != nil {
		return nil, err
	}
	defer parentDao.Close()
	parent := models.Node{ID: req.NewParent.Id, Port: req.NewParent.Host, Host: req.NewParent.Port}
	if err := parentDao.WriteParent(parent); err != nil {
		return nil, err
	}
	log.Printf("Nuovo nodo padre con id:%s, host:%s, port:%s\n", parent.ID, parent.Host, parent.Port)
	return &pb.ChangeParentResponse{Status: "OK"}, nil
}

func addNodeAsChild(req *pb.JoinRequest, childs []models.Node, childsDao *dao.ChildsDAO) (*pb.JoinResponse, error) {
	fmt.Println("entrata in addNodeAsChild")
	nephewsDao, err := dao.NewNephewsDAO()
	if err != nil {
		return nil, err
	}
	defer nephewsDao.Close()
	myIntegerId, err := strconv.Atoi(session.GetSession().Node.ID)
	if err != nil {
		return nil, err
	}
	var id int
	switch len(childs) {
	case 0:
		id = 2*myIntegerId + 1
		break
	case 1:
		fmt.Println("entrata in case 1")
		if childs[0].ID == fmt.Sprintf("%d", 2*myIntegerId+1) {
			id = 2*myIntegerId + 2
		} else {
			id = 2*myIntegerId + 1
		}
		break
	}
	nodeToSaveAsChild := models.Node{
		Host: req.Host,
		Port: req.Port,
		ID:   fmt.Sprintf("%d", id),
	}
	if err := childsDao.WriteChild(nodeToSaveAsChild); err != nil {
		return nil, err
	}
	nephews, err := nephewsDao.ReadAllNephews()
	if err != nil {
		return nil, err
	}
	var ChildsResp []*pb.NodeInfo
	var NephewsResp []*pb.NodeInfo
	for nephew := range nephews {
		if isChild(nephews[nephew], nodeToSaveAsChild) {
			if err := contactNodeToChangeParent(nephews[nephew], nodeToSaveAsChild); err != nil {
				return nil, err
			}
			ChildsResp = append(ChildsResp, &pb.NodeInfo{Id: nephews[nephew].ID, Port: nephews[nephew].Port, Host: nephews[nephew].Port})
		} else if isNephew(nephews[nephew], nodeToSaveAsChild) {
			if err := contactNodeToChangeParent(nephews[nephew], nodeToSaveAsChild); err != nil {
				return nil, err
			}
			NephewsResp = append(NephewsResp, &pb.NodeInfo{Id: nephews[nephew].ID, Port: nephews[nephew].Port, Host: nephews[nephew].Host})

		}
	}

	return &pb.JoinResponse{Status: fmt.Sprintf("NEED_CHILD:%d", id), Childs: ChildsResp, Nephews: NephewsResp}, nil
}

func contactNodeToChangeParent(node models.Node, parent models.Node) error {
	addr := fmt.Sprintf("%s:%s", node.Host, node.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewDHTClient(conn)
	newParent := &pb.NodeInfo{
		Id:   parent.ID,
		Host: node.Host,
		Port: node.Port,
	}
	req := &pb.ChangeParentRequest{NewParent: newParent}

	resp, err := client.ChangeParent(context.Background(), req)

	if err != nil {
		return err
	}
	if resp.Status != "OK" {
		return errors.New(resp.Status)
	}
	return nil
}

func isChild(node models.Node, parent models.Node) bool {
	integerNodeID, err := strconv.Atoi(node.ID)
	if err != nil {
		log.Fatal(err)
	}
	integerParentID, err := strconv.Atoi(parent.ID)
	if err != nil {
		log.Fatal(err)
	}
	return (integerNodeID-1)/2 == integerParentID
}

func isNephew(nephew models.Node, node models.Node) bool {
	integerNodeID, err := strconv.Atoi(node.ID)
	if err != nil {
		log.Fatal(err)
	}
	integerNephewID, err := strconv.Atoi(nephew.ID)
	if err != nil {
		log.Fatal(err)
	}
	for integerNodeID > 0 {
		parent := (integerNodeID - 1) / 2
		if parent == integerNephewID {
			return true
		}
		integerNodeID = parent
	}
	return false
}
