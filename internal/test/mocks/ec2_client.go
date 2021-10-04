package mocks

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// EC2Client implements The EC2 service insterface
type EC2Client struct {
	ec2iface.EC2API
	SGList          []*ec2.SecurityGroup
	ReservationList []*ec2.Reservation
	FailAuthorizeSG bool // Forces AuthorizeSecurityGroupIngress to fail
	FailRevokeSG    bool // Forces RevokeSecurityGroupIngress to fail
}

// AuthorizeSecurityGroupIngress mocks the equivalent AWS SDK function
func (m *EC2Client) AuthorizeSecurityGroupIngress(
	params *ec2.AuthorizeSecurityGroupIngressInput,
) (
	*ec2.AuthorizeSecurityGroupIngressOutput,
	error,
) {
	if m.FailAuthorizeSG {
		return nil, fmt.Errorf("it had to fail")
	}
	return &ec2.AuthorizeSecurityGroupIngressOutput{}, nil
}

// CreateSecurityGroup mocks the equivalent AWS SDK function
func (m *EC2Client) CreateSecurityGroup(
	params *ec2.CreateSecurityGroupInput,
) (
	out *ec2.CreateSecurityGroupOutput,
	err error,
) {
	err = params.Validate()
	if err != nil {
		return nil, err
	}
	return &ec2.CreateSecurityGroupOutput{
		GroupId: aws.String("sg-12345678"),
	}, nil
}

// DescribeInstances mocks the equivalent AWS SDK function
func (m *EC2Client) DescribeInstances(
	in *ec2.DescribeInstancesInput,
) (
	out *ec2.DescribeInstancesOutput,
	err error,
) {
	return &ec2.DescribeInstancesOutput{
		Reservations: m.ReservationList,
	}, nil
}

// DescribeSecurityGroups mocks the equivalent AWS SDK function
func (m *EC2Client) DescribeSecurityGroups(
	in *ec2.DescribeSecurityGroupsInput,
) (
	out *ec2.DescribeSecurityGroupsOutput,
	err error,
) {
	return &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: m.SGList,
	}, nil
}

// RevokeSecurityGroupIngress mocks the equivalent AWS SDK function
func (m *EC2Client) RevokeSecurityGroupIngress(
	params *ec2.RevokeSecurityGroupIngressInput,
) (
	*ec2.RevokeSecurityGroupIngressOutput,
	error,
) {
	if m.FailRevokeSG {
		return nil, fmt.Errorf("it had to fail")
	}

	return &ec2.RevokeSecurityGroupIngressOutput{}, nil
}
