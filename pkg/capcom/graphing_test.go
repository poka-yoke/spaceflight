package capcom

import (
	"testing"

	"github.com/poka-yoke/spaceflight/internal/test/mocks"
)

func TestSGInstanceStateGetKeysAndHas(t *testing.T) {
	sg := make(sGInstanceState)
	empty := map[string]int{}

	keyList := []string{"first", "second", "third"}
	for _, v := range keyList {
		sg[v] = empty
	}

	length := len(sg.getKeys())
	if length != 3 {
		t.Errorf("Expected length 3, found %d.\n", length)
	}
	for _, v := range keyList {
		if !sg.has(v) {
			t.Errorf("Expected value \"%s\" not found.", v)
		}
	}
}

func TestGetInstances(t *testing.T) {
	svc := &mocks.EC2Client{}
	res := getInstanceReservations(svc)
	if *res[0].Instances[0].State.Name != "pending" ||
		*res[0].Groups[0].GroupId != "sg-12345678" {
		t.Error("Should be equal")
	}
}

func TestGetInstancesStates(t *testing.T) {
	svc := &mocks.EC2Client{}
	res := getInstancesStates(getInstanceReservations(svc))
	if len(res) != 1 {
		t.Error("Unexpected amount of results")
	}
	if state := res["sg-12345678"]; state != nil {
		if state["pending"] != 1 {
			t.Error("Unexpected values")
		}
	} else {
		t.Error("Expected key missing")
	}
}
