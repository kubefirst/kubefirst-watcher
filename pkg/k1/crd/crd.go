package crd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/6za/k1-watcher/pkg/k1/k8s"
	"go.uber.org/zap"
	api "k8s.io/apimachinery/pkg/types"
)

type CRDConfig struct {
	APIVersion   string
	Namespace    string
	InstanceName string
	Resource     string
}

type CRDClient struct {
	CRD    *CRDConfig
	Logger *zap.Logger
}

// UpdateStatus update Status Message on Watcher Status CRD
func (client *CRDClient) UpdateStatus(status string) error {
	client.Logger.Debug(fmt.Sprintf("Watcher Config: #%v ", client.CRD))
	clientSet := k8s.GetK8SConfig()
	myPatch := fmt.Sprintf(`{"status":{"status":"%s"}}`, status)
	client.Logger.Debug(fmt.Sprintf("Watcher Patch: #%v ", myPatch))
	_, err := clientSet.RESTClient().
		Patch(api.MergePatchType).
		AbsPath("/apis/" + client.CRD.APIVersion).
		SubResource("status").
		Namespace(client.CRD.Namespace).
		Resource(client.CRD.Resource).
		Name(client.CRD.InstanceName).
		Body([]byte(myPatch)).
		DoRaw(context.TODO())
	if err != nil {
		client.Logger.Info(fmt.Sprintf("Error updating CRD   #%v ", err))
		return err
	}
	client.Logger.Info(fmt.Sprintf("Update status:  %#v ", client.CRD))
	return nil
}

// GetCRD returns the Watcher CRD Object, used to feed configs
func (client *CRDClient) GetCRD() (Watcher, error) {
	client.Logger.Debug(fmt.Sprintf("Watcher Config: #%v ", client.CRD))
	clientSet := k8s.GetK8SConfig()
	object, err := clientSet.RESTClient().
		Get().
		AbsPath("/apis/" + client.CRD.APIVersion).
		Namespace(client.CRD.Namespace).
		Resource(client.CRD.Resource).
		Name(client.CRD.InstanceName).
		DoRaw(context.TODO())
	client.Logger.Info(fmt.Sprintf("Get CRD   #%v ", string(object)))
	if err != nil {
		client.Logger.Info(fmt.Sprintf("Error loading CRD   #%v ", err))
		return Watcher{}, err
	}
	watcher := Watcher{}
	err = json.Unmarshal(object, &watcher)
	client.Logger.Info(fmt.Sprintf("Watcher   #%v ", watcher))
	if err != nil {
		client.Logger.Info(fmt.Sprintf("Error loading CRD   #%v ", err))
		return Watcher{}, err
	}
	client.Logger.Info(fmt.Sprintf("Update status:  %#v ", client.CRD))
	return watcher, nil
}
