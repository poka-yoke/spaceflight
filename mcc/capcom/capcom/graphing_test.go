package capcom

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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
