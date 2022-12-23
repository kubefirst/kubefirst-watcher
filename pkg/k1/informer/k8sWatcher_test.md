package informer_test

import (
	"context"
	"testing"

	"github.com/kubefirst/kubefirst-watcher/pkg/k1/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

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
	cl := fake.NewFakeClient(objs...)
	opt := client.MatchingLabels(map[string]string{"label-key": "label-value"})
	wactherList := &v1beta1.WatcherList{}
	err := cl.List(context.TODO(), wactherList, opt)
	if err != nil {
		t.Fatalf("list memcached: (%v)", err)
	}

}
