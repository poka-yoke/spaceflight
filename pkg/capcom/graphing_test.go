package capcom

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

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

func TestGetInstanceReservations(t *testing.T) {
	data := []struct {
		state, id *string
	}{
		{
			state: aws.String("pending"),
			id:    aws.String("sg-12345678"),
		},
		{
			state: aws.String("running"),
			id:    aws.String("sg-23456789"),
		},
	}

	svc := &mocks.EC2Client{}
	for _, tc := range data {
		svc.ReservationList = []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						State: &ec2.InstanceState{
							Name: tc.state,
						},
					},
				},
				Groups: []*ec2.GroupIdentifier{
					{
						GroupId: tc.id,
					},
				},
			},
		}
		t.Run(
			fmt.Sprintf("%s %s", *tc.state, *tc.id),
			func(t *testing.T) {
				res := getInstanceReservations(svc)
				if *res[0].Instances[0].State.Name != *tc.state {
					t.Errorf(
						"Expected state to be %s, but got: %s",
						*tc.state,
						*res[0].Instances[0].State.Name,
					)
				}
				if *res[0].Groups[0].GroupId != *tc.id {
					t.Errorf(
						"Expected ID to be %s, but got: %s",
						*tc.id,
						*res[0].Groups[0].GroupId,
					)
				}
			},
		)
	}
}

func TestGetInstancesStates(t *testing.T) {
	data := []struct {
		state, id *string
	}{
		{
			state: aws.String("pending"),
			id:    aws.String("sg-12345678"),
		},
		{
			state: aws.String("running"),
			id:    aws.String("sg-23456789"),
		},
	}

	for _, tc := range data {
		reservationList := []*ec2.Reservation{
			{
				Instances: []*ec2.Instance{
					{
						State: &ec2.InstanceState{
							Name: tc.state,
						},
					},
				},
				Groups: []*ec2.GroupIdentifier{
					{
						GroupId: tc.id,
					},
				},
			},
		}
		t.Run(
			fmt.Sprintf("%s %s", *tc.state, *tc.id),
			func(t *testing.T) {
				res := getInstancesStates(reservationList)
				if len(res) != 1 {
					t.Errorf(
						"Unexpected amount of results. Expected %d but got %d",
						1,
						len(res),
					)
				}
				state, ok := res[*tc.id]
				if !ok {
					t.Errorf(
						"Expected key %s not found",
						*tc.id,
					)
				}
				if state[*tc.state] != 1 {
					t.Errorf(
						"Expected state %s to equal %d",
						*tc.state,
						1,
					)
				}
			},
		)
	}
}
