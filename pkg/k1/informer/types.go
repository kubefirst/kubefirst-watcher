package informer

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type SelectorName string

const (
	POD       SelectorName = "Pod"
	JOB       SelectorName = "Job"
	CONFIGMAP SelectorName = "ConfigMap"
	SECRET    SelectorName = "Secret"
	UNKOWN    SelectorName = ""
)

func (sn SelectorName) GetGVK() *schema.GroupVersionKind {
	switch sn {
	case POD:
		return &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}
	case JOB:
		return &schema.GroupVersionKind{Group: "", Version: "batch/v1", Kind: "Job"}

	case CONFIGMAP:
		return &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}

	case SECRET:
		return &schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}

	default:
		return nil
	}
}

func GetSelector(gvk schema.GroupVersionKind) SelectorName {
	switch gvk {
	case schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}:
		return POD

	case schema.GroupVersionKind{Group: "", Version: "batch/v1", Kind: "Job"}:
		return JOB

	case schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}:
		return CONFIGMAP

	case schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}:
		return SECRET

	default:
		return UNKOWN
	}
}
