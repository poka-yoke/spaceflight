package mocks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// EC2Client implements The EC2 service insterface
type EC2Client struct {
	ec2iface.EC2API
}

// AuthorizeSecurityGroupIngress mocks the equivalent AWS SDK function
func (m *EC2Client) AuthorizeSecurityGroupIngress(
	params *ec2.AuthorizeSecurityGroupIngressInput,
) (
	*ec2.AuthorizeSecurityGroupIngressOutput,
	error,
) {
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
	}, nil
}

// DescribeSecurityGroups mocks the equivalent AWS SDK function
func (m *EC2Client) DescribeSecurityGroups(
	in *ec2.DescribeSecurityGroupsInput,
) (
	out *ec2.DescribeSecurityGroupsOutput,
	err error,
) {
	out = &ec2.DescribeSecurityGroupsOutput{
		SecurityGroups: []*ec2.SecurityGroup{
			{
				Description: aws.String(""),
				GroupId:     aws.String("sg-1234"),
				GroupName:   aws.String(""),
				IpPermissions: []*ec2.IpPermission{
					{
						IpProtocol: aws.String("tcp"),
						ToPort:     aws.Int64(22),
						IpRanges: []*ec2.IpRange{
							{CidrIp: aws.String("1.2.3.4/32")},
						},
					},
				},
			},
		},
	}
	return
}

// RevokeSecurityGroupIngress mocks the equivalent AWS SDK function
func (m *EC2Client) RevokeSecurityGroupIngress(
	params *ec2.RevokeSecurityGroupIngressInput,
) (
	*ec2.RevokeSecurityGroupIngressOutput,
	error,
) {
	return &ec2.RevokeSecurityGroupIngressOutput{}, nil
}
