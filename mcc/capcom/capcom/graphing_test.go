package capcom

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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

var describeInstancesOutput = ec2.DescribeInstancesOutput{
	Reservations: []*ec2.Reservation{
		{
			Instances: []*ec2.Instance{
				{
					State: &ec2.InstanceState{
						Name: aws.String("pending"),
					},
				},
			},
			Groups: []*ec2.GroupIdentifier{
				{
					GroupId: aws.String("sg-12345678"),
				},
			},
		},
	},
}

func (m *mockEC2Client) DescribeInstances(
	in *ec2.DescribeInstancesInput,
) (
	out *ec2.DescribeInstancesOutput,
	err error,
) {
	return &describeInstancesOutput, nil
}

func TestGetInstances(t *testing.T) {
	svc := &mockEC2Client{}
	res := getInstances(svc)
	if *res.Reservations[0].Instances[0].State.Name != "pending" ||
		*res.Reservations[0].Groups[0].GroupId != "sg-12345678" {
		t.Error("Should be equal")
	}
}

func TestGetInstancesStates(t *testing.T) {
	res := getInstancesStates(describeInstancesOutput.Reservations)
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
