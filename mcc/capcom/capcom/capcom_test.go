package capcom

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

type mockEC2Client struct {
	ec2iface.EC2API
}

func (m *mockEC2Client) AuthorizeSecurityGroupIngress(
	params *ec2.AuthorizeSecurityGroupIngressInput,
) (
	out *ec2.AuthorizeSecurityGroupIngressOutput,
	err error,
) {
	return
}

func (m *mockEC2Client) CreateSecurityGroup(
	params *ec2.CreateSecurityGroupInput,
) (
	out *ec2.CreateSecurityGroupOutput,
	err error,
) {
	err = params.Validate()
	out = &ec2.CreateSecurityGroupOutput{
		GroupId: aws.String("sg-12345678"),
	}
	return
}

func (m *mockEC2Client) DescribeSecurityGroups(
	in *ec2.DescribeSecurityGroupsInput,
) (
	out *ec2.DescribeSecurityGroupsOutput,
	err error,
) {
	out = &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []*ec2.SecurityGroup{
			{
				Description: aws.String(""),
				GroupId:     aws.String(""),
				GroupName:   aws.String(""),
			},
		},
	}
	return
}

func (m *mockEC2Client) RevokeSecurityGroupIngress(
	params *ec2.RevokeSecurityGroupIngressInput,
) (
	out *ec2.RevokeSecurityGroupIngressOutput,
	err error,
) {
	return
}

var ListSecurityGroupsExpectedOutput = []string{
	fmt.Sprintf("* %10s %20s %s\n", "", "", ""),
}

func TestListSecurityGroups(t *testing.T) {
	svc := &mockEC2Client{}
	out := ListSecurityGroups(svc)
	expected := ListSecurityGroupsExpectedOutput
	for index, line := range out {
		if line != expected[index] {
			t.Error("Unexpected output")
		}
	}
}

func TestAuthorizeAccessToSecurityGroup(t *testing.T) {
	svc := &mockEC2Client{}
	origin := "/32"
	proto := ""
	port := int64(0)
	destination := "sg-"
	_, err := AuthorizeAccessToSecurityGroup(svc, origin, proto, port, destination)
	if err != nil {
		t.Error(err)
	}

}
func TestRevokeAccessToSecurityGroup(t *testing.T) {
	svc := &mockEC2Client{}
	origin := "/32"
	proto := ""
	port := int64(0)
	destination := "sg-"
	_, err := RevokeAccessToSecurityGroup(svc, origin, proto, port, destination)
	if err != nil {
		t.Error(err)
	}

}

var csgtable = []struct {
	name        string
	description string
	vpcid       string
	out         string
}{
	{
		name:        "",
		description: "Non-VPC success",
		vpcid:       "",
		out:         "sg-12345678",
	},
	{
		name:        "",
		description: "VPC success",
		vpcid:       "vpc-12345678",
		out:         "sg-12345678",
	},
}

func TestCreateSG(t *testing.T) {
	svc := &mockEC2Client{}
	for _, tt := range csgtable {
		t.Run(
			tt.description,
			func(t *testing.T) {
				out := CreateSG(
					tt.name,
					tt.description,
					tt.vpcid,
					svc,
				)
				if out != tt.out {
					t.Error("Unexpected output")
				}
			},
		)
	}
}
