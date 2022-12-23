package informer_test

import (
	"context"
	"testing"
	"time"

	//v1beta1 "github.com/kubefirst/kubefirst-watcher-operator/api/v1beta1"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"
)

// TestFakeClient demonstrates how to use a fake client with SharedInformerFactory in tests.
// TODO: Create Test for StartCRDWatcher(clientSet *kubernetes.Clientset, clientCrd *crd.CRDClient, loggerIn *zap.Logger) using this API
// Try full round cycle, create CRD, start pointing to CRD and Add some pods before start

func TestFakeClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create the fake client.
	client := fake.NewSimpleClientset()

	// We will create an informer that writes added pods to a channel.
	pods := make(chan *v1.Pod, 1)
	informers := informers.NewSharedInformerFactory(client, 0)
	podInformer := informers.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(&cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			t.Logf("pod added: %s/%s", pod.Namespace, pod.Name)
			pods <- pod
		},
	})

	// Make sure informers are running.
	informers.Start(ctx.Done())

	// This is not required in tests, but it serves as a proof-of-concept by
	// ensuring that the informer goroutine have warmed up and called List before
	// we send any events to it.
	cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced)

	// The fake client doesn't support resource version. Any writes to the client
	// after the informer's initial LIST and before the informer establishing the
	// watcher will be missed by the informer. Therefore we wait until the watcher
	// starts.
	// Note that the fake client isn't designed to work with informer. It
	// doesn't support resource version. It's encouraged to use a real client
	// in an integration/E2E test if you need to test complex behavior with
	// informer/controllers.

	// Inject an event into the fake client.
	p := &v1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "my-pod"}}
	_, err := client.CoreV1().Pods("test-ns").Create(context.TODO(), p, metav1.CreateOptions{})
	if err != nil {
		t.Fatalf("error injecting pod add: %v", err)
	}

	select {
	case pod := <-pods:
		t.Logf("Got pod from channel: %s/%s", pod.Namespace, pod.Name)
	case <-time.After(wait.ForeverTestTimeout):
		t.Error("Informer did not get the added pod")
	}
}

/*
func TestCRD(t *testing.T) {

	// We will create an informer that writes added pods to a channel.
	sampleWatcher := &v1beta1.Watcher{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sample",
			Namespace: "default",
			Labels: map[string]string{
				"app": "watcher",
			},
		},
	}
	objs := []runtime.Object{sampleWatcher}
	cl := sigsfake.NewFakeClient(objs...)
	opt := client.MatchingLabels(map[string]string{"label-key": "label-value"})
	wactherList := &v1beta1.WatcherList{}
	err := cl.List(context.TODO(), wactherList, opt)
	if err != nil {
		t.Fatalf("list memcached: (%v)", err)
	}

}
*/
