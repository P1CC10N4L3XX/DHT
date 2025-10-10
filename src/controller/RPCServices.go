package controller

import (
	"DHT/src/dao"
	"DHT/src/models"
	pb "DHT/src/proto/stubs"
	"DHT/src/session"
	"DHT/src/utils"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/big"
	"slices"
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
	if len(childs) < 2 {
		return addNodeAsChild(req, childs, childsDao)
	}
	for _, child := range childs {
		if child.ID == req.Next {
			return &pb.JoinResponse{Status: fmt.Sprintf("CONTACT_CHILD:%s:%s", child.Host, child.Port)}, nil
		}
	}
	return nil, nil
}

func (s *DhtServer) LeaveNode(ctx context.Context, req *pb.LeaveRequest) (*pb.LeaveResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	childsDao, err := dao.NewChildsDAO()
	defer childsDao.Close()
	if err != nil {
		return nil, err
	}
	nephewsDao, err := dao.NewNephewsDAO()
	if err != nil {
		return nil, err
	}
	defer nephewsDao.Close()
	resourcesDao, err := dao.NewResourceDAO()
	if err != nil {
		return nil, err
	}
	nodeToLeave := models.Node{ID: req.NodeToLeave.Id, Host: req.NodeToLeave.Host, Port: req.NodeToLeave.Port}
	for _, childReq := range req.Childs {
		child := models.Node{ID: childReq.Id, Host: childReq.Host, Port: childReq.Port}
		if err := nephewsDao.WriteNephew(child); err != nil {
			return nil, err
		}
		if err := contactNodeToChangeParent(child, *session.GetSession().Node); err != nil {
			return nil, err
		}
	}
	for _, nephewReq := range req.Nephews {
		nephew := models.Node{ID: nephewReq.Id, Host: nephewReq.Host, Port: nephewReq.Port}
		if err := nephewsDao.WriteNephew(nephew); err != nil {
			return nil, err
		}
		if err := contactNodeToChangeParent(nephew, *session.GetSession().Node); err != nil {
			return nil, err
		}
	}
	for _, resourceReq := range req.Resources {
		resource := models.Resource{Key: resourceReq.Key, Value: resourceReq.Value}
		if err := resourcesDao.WriteResource(resource); err != nil {
			return nil, err
		}
	}
	if err := childsDao.RemoveChild(nodeToLeave); err != nil {
		return nil, err
	}
	if err := nephewsDao.RemoveNephew(nodeToLeave); err != nil {
		return nil, err
	}

	return &pb.LeaveResponse{Status: "OK"}, nil
}

func (s *DhtServer) ChangeParent(ctx context.Context, req *pb.ChangeParentRequest) (*pb.ChangeParentResponse, error) {
	parentDao, err := dao.NewParentDAO()
	if err != nil {
		return nil, err
	}
	defer parentDao.Close()
	parent := models.Node{ID: req.NewParent.Id, Port: req.NewParent.Port, Host: req.NewParent.Host}
	if err := parentDao.WriteParent(parent); err != nil {
		return nil, err
	}
	log.Printf("Nuovo nodo padre con id:%s, host:%s, port:%s\n", parent.ID, parent.Host, parent.Port)
	return &pb.ChangeParentResponse{Status: "OK"}, nil
}

func (s *DhtServer) PutResource(ctx context.Context, req *pb.PutRequest) (*pb.PutResponse, error) {
	childsDao, err := dao.NewChildsDAO()
	if err != nil {
		return nil, err
	}
	defer childsDao.Close()
	nephewsDao, err := dao.NewNephewsDAO()
	if err != nil {
		return nil, err
	}
	defer nephewsDao.Close()
	childs, err := childsDao.ReadAllChilds()
	if err != nil {
		return nil, err
	}
	for _, child := range childs {
		if child.ID == req.Next {
			return &pb.PutResponse{Status: "CONTACT_CHILD:" + child.Host + ":" + child.Port}, nil
		}
	}
	nephews, err := nephewsDao.ReadAllNephews()
	keyBigInt := new(big.Int)
	keyBigInt.SetString(req.Resource.Key, 16)
	path := utils.BuildPath(keyBigInt)
	for _, nephew := range nephews {
		if slices.Contains(path, nephew.ID) {
			return &pb.PutResponse{Status: "CONTACT_NEPHEW:" + nephew.Host + ":" + nephew.Port + ":" + nephew.ID}, nil
		}
	}

	resourceDao, err := dao.NewResourceDAO()
	if err != nil {
		return nil, err
	}
	defer resourceDao.Close()
	resource := models.Resource{
		Key:   req.Resource.Key,
		Value: req.Resource.Value,
	}
	if err := resourceDao.WriteResource(resource); err != nil {
		return nil, err
	}
	return &pb.PutResponse{Status: "RESOURCE_STORED"}, nil
}

func (s *DhtServer) GetResource(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	resourceDao, err := dao.NewResourceDAO()
	if err != nil {
		return nil, err
	}
	defer resourceDao.Close()
	resource, err := resourceDao.ReadResourceByKey(req.Key)
	if err != nil {
		return nil, err
	}
	if resource != (models.Resource{}) {
		resourceResp := &pb.Resource{Key: resource.Key, Value: resource.Value}
		return &pb.GetResponse{Status: "RESOURCE_DETECTED", Resource: resourceResp}, nil
	}

	childsDao, err := dao.NewChildsDAO()
	if err != nil {
		return nil, err
	}
	defer childsDao.Close()
	childs, err := childsDao.ReadAllChilds()
	if err != nil {
		return nil, err
	}
	for _, child := range childs {
		if child.ID == req.Next {
			return &pb.GetResponse{Status: "CONTACT_CHILD:" + child.Host + ":" + child.Port}, nil
		}
	}
	nephewsDao, err := dao.NewNephewsDAO()
	if err != nil {
		return nil, err
	}
	nephews, err := nephewsDao.ReadAllNephews()
	if err != nil {
		return nil, err
	}
	keyBigInt := new(big.Int)
	keyBigInt.SetString(req.Key, 16)
	path := utils.BuildPath(keyBigInt)
	for _, nephew := range nephews {
		if slices.Contains(path, nephew.ID) {
			return &pb.GetResponse{Status: "CONTACT_NEPHEW:" + nephew.Host + ":" + nephew.Port + ":" + nephew.ID}, nil
		}
	}
	return &pb.GetResponse{Status: "RESOURCE_NOTFOUND"}, nil
}

func addNodeAsChild(req *pb.JoinRequest, childs []models.Node, childsDao *dao.ChildsDAO) (*pb.JoinResponse, error) {
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
	for _, nephew := range nephews {
		if isChild(nephew, nodeToSaveAsChild) {
			if err := nephewsDao.RemoveNephew(nephew); err != nil {
				return nil, err
			}
			if err := contactNodeToChangeParent(nephew, nodeToSaveAsChild); err != nil {
				return nil, err
			}
			ChildsResp = append(ChildsResp, &pb.NodeInfo{Id: nephew.ID, Port: nephew.Port, Host: nephew.Port})
		} else if isNephew(nephew, nodeToSaveAsChild) {
			if err := nephewsDao.RemoveNephew(nephew); err != nil {
				return nil, err
			}
			if err := contactNodeToChangeParent(nephew, nodeToSaveAsChild); err != nil {
				return nil, err
			}
			NephewsResp = append(NephewsResp, &pb.NodeInfo{Id: nephew.ID, Port: nephew.Port, Host: nephew.Host})

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
		Host: parent.Host,
		Port: parent.Port,
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
