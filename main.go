package kapsule_database_interactions

import (
	"fmt"
	"github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"net"
	"net/http"
	"os"
)

func Handle(w http.ResponseWriter, r *http.Request) {

	orgaId := os.Getenv("ORGANIZATION_ID")
	accessKey := os.Getenv("ACCESS_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	regionDB := os.Getenv("DATABASE_REGION")
	regionKapsule := os.Getenv("KAPSULE_REGION")
	instanceId := os.Getenv("DATABASE_INSTANCE_ID")
	clusterId := os.Getenv("KAPSULE_CLUSTER_ID")

	// Create a Scaleway client
	client, err := scw.NewClient(
		scw.WithDefaultOrganizationID(orgaId),
		scw.WithAuth(accessKey, secretKey),
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_, err := fmt.Fprint(w, err.Error())
		if err != nil {
			return
		}
		return
	}

	k8sApi := k8s.NewAPI(client)

	reqListNodes := k8s.ListNodesRequest{
		Region:    scw.Region(regionKapsule),
		ClusterID: clusterId,
	}

	respNode, err := k8sApi.ListNodes(&reqListNodes)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, err.Error())
		if err != nil {
			return
		}
		return
	}
	var IPList []net.IP
	for _, node := range respNode.Nodes {
		IPList = append(IPList, *node.PublicIPV4)
	}

	var rules []*rdb.ACLRuleRequest
	for i, ip := range IPList {
		rule := rdb.ACLRuleRequest{
			IP:          scw.IPNet{IPNet: net.IPNet{IP: ip, Mask: net.IPMask{255, 255, 255, 255}}},
			Description: fmt.Sprintf("IP of the Kapsule node %v", respNode.Nodes[i].Name),
		}
		rules = append(rules, &rule)
	}

	dbApi := rdb.NewAPI(client)
	req1 := &rdb.SetInstanceACLRulesRequest{
		Region:     scw.Region(regionDB),
		InstanceID: instanceId,
		Rules:      rules,
	}
	_, err = dbApi.SetInstanceACLRules(req1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprint(w, err.Error())
		if err != nil {
			return
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string("done"))
}
