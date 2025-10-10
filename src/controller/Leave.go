package controller

import (
	"DHT/src/dao"
	pb "DHT/src/proto/stubs"
	"DHT/src/session"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LeaveController struct{}

func (l *LeaveController) Leave() error {
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
	resourcesDao, err := dao.NewResourceDAO()
	if err != nil {
		return err
	}
	parent, err := parentDao.ReadParent()
	if err != nil {
		return err
	}
	childs, err := childsDao.ReadAllChilds()
	if err != nil {
		return err
	}
	nephews, err := nephewsDao.ReadAllNephews()
	if err != nil {
		return err
	}
	resources, err := resourcesDao.ReadAllResources()
	if err != nil {
		return err
	}
	addr := fmt.Sprintf("%s:%s", parent.Host, parent.Port)

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := pb.NewDHTClient(conn)
	nodeToLeave := &pb.NodeInfo{Id: session.GetSession().Node.ID, Host: session.GetSession().Node.Host, Port: session.GetSession().Node.Port}
	var childsReq []*pb.NodeInfo
	var nephewsReq []*pb.NodeInfo
	var resourcesReq []*pb.Resource
	for _, child := range childs {
		childsReq = append(childsReq, &pb.NodeInfo{Id: child.ID, Host: child.Host, Port: child.Port})
	}
	for _, nephew := range nephews {
		nephewsReq = append(nephewsReq, &pb.NodeInfo{Id: nephew.ID, Host: nephew.Host, Port: nephew.Port})
	}
	for _, resource := range resources {
		resourcesReq = append(resourcesReq, &pb.Resource{Key: resource.Key, Value: resource.Value})
	}
	req := &pb.LeaveRequest{NodeToLeave: nodeToLeave, Childs: childsReq, Nephews: nephewsReq, Resources: resourcesReq}
	resp, err := client.LeaveNode(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.Status == "OK" {
		return nil
	}
	return errors.New(resp.Status)
}
