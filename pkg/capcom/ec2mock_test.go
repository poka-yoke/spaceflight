package capcom

import (
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
	*ec2.AuthorizeSecurityGroupIngressOutput,
	error,
) {
	return &ec2.AuthorizeSecurityGroupIngressOutput{}, nil
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

func (m *mockEC2Client) RevokeSecurityGroupIngress(
	params *ec2.RevokeSecurityGroupIngressInput,
) (
	*ec2.RevokeSecurityGroupIngressOutput,
	error,
) {
	return &ec2.RevokeSecurityGroupIngressOutput{}, nil
}
