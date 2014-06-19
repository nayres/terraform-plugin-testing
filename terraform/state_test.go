package terraform

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/config"
)

func TestResourceState_MergeDiff(t *testing.T) {
	rs := ResourceState{
		ID: "foo",
		Attributes: map[string]string{
			"foo": "bar",
		},
	}

	diff := &ResourceDiff{
		Attributes: map[string]*ResourceAttrDiff{
			"foo": &ResourceAttrDiff{
				Old: "bar",
				New: "baz",
			},
			"bar": &ResourceAttrDiff{
				Old: "",
				New: "foo",
			},
			"baz": &ResourceAttrDiff{
				Old:         "",
				New:         "foo",
				NewComputed: true,
			},
		},
	}

	rs2 := rs.MergeDiff(diff)

	expected := map[string]string{
		"foo": "baz",
		"bar": "foo",
		"baz": config.UnknownVariableValue,
	}

	if !reflect.DeepEqual(expected, rs2.Attributes) {
		t.Fatalf("bad: %#v", rs2.Attributes)
	}
}

func TestResourceState_MergeDiff_nil(t *testing.T) {
	var rs *ResourceState = nil

	diff := &ResourceDiff{
		Attributes: map[string]*ResourceAttrDiff{
			"foo": &ResourceAttrDiff{
				Old: "",
				New: "baz",
			},
		},
	}

	rs2 := rs.MergeDiff(diff)

	expected := map[string]string{
		"foo": "baz",
	}

	if !reflect.DeepEqual(expected, rs2.Attributes) {
		t.Fatalf("bad: %#v", rs2.Attributes)
	}
}

func TestReadWriteState(t *testing.T) {
	state := &State{
		Resources: map[string]*ResourceState{
			"foo": &ResourceState{
				ID: "bar",
			},
		},
	}

	buf := new(bytes.Buffer)
	if err := WriteState(state, buf); err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err := ReadState(buf)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(actual, state) {
		t.Fatalf("bad: %#v", actual)
	}
}
