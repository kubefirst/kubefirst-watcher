package informer

import (
	"reflect"
	"testing"

	"github.com/thoas/go-funk"
)

/*
func Test_Subset(t *testing.T) {
	funk.Contains([]string{"foo", "bar"}, "bar")
	labelAll := map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"}
	labelSubset := map[string]string{"label2": "value2", "label3": "value3"}

	isSubset := funk.Contains(labelAll, labelSubset)
	//funk.Subset(labelSubset, labelAll)
	if !isSubset {
		t.Errorf("Not isSubset, got: %v ", isSubset)
	}

}
*/

func Test_Intersect(t *testing.T) {
	labelAll := map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"}
	labelSubset := map[string]string{"label2": "value2", "label3": "value3"}

	match, _ := IsMapPresent(labelAll, labelSubset)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_NotIntersect(t *testing.T) {
	labelAll := map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"}
	labelSubset := map[string]string{"label2": "value2", "label4": "value4"}

	keysAll := funk.Keys(labelAll)
	keysSubset := funk.Keys(labelSubset)

	intersect := funk.Intersect(keysAll, keysSubset)
	if reflect.DeepEqual(intersect, keysSubset) {
		t.Errorf("Not intersect, got: %v ", intersect)
	}
}

func Test_IntersectDifferentOrder(t *testing.T) {
	labelAll := map[string]string{"label1": "value1", "label3": "value3", "label2": "value2"}
	labelSubset := map[string]string{"label2": "value2", "label3": "value3"}

	match, _ := IsMapPresent(labelAll, labelSubset)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_IntersectNotMatcht(t *testing.T) {
	labelAll := map[string]string{"label1": "value1", "label": "value", "label2": "value2"}
	labelSubset := map[string]string{"label2": "value2", "label3": "value3"}

	match, _ := IsMapPresent(labelAll, labelSubset)
	if match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_IntersectNilSubset(t *testing.T) {
	labelAll := map[string]string{"label1": "value1", "label3": "value3", "label2": "value2"}
	labelSubset := map[string]string{}

	match, _ := IsMapPresent(labelAll, labelSubset)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}

func Test_IntersectNilNil(t *testing.T) {
	labelAll := map[string]string{}
	labelSubset := map[string]string{}

	match, _ := IsMapPresent(labelAll, labelSubset)
	if !match {
		t.Errorf("Not match, got: %v ", match)
	}
}
