package capcom

import (
	"errors"
	"fmt"
	"net"
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
				GroupId:     aws.String("sg-1234"),
				GroupName:   aws.String(""),
				IpPermissions: []*ec2.IpPermission{
					{
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
	out *ec2.RevokeSecurityGroupIngressOutput,
	err error,
) {
	return
}

var ListSecurityGroupsExpectedOutput = []string{
	fmt.Sprintf("* %10s %20s %s\n", "sg-1234", "", ""),
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

var biptable = []struct {
	origin string
	proto  string
	port   int64
	err    error
}{
	{
		origin: "1.2.3.4/32",
		proto:  "tcp",
		port:   int64(0),
		err:    nil,
	},
	{
		origin: "sg-",
		proto:  "tcp",
		port:   int64(0),
		err:    nil,
	},
	{
		origin: "1.2.3.4/32",
		proto:  "udp",
		port:   int64(0),
		err:    nil,
	},
	{
		origin: "sg-",
		proto:  "icmp",
		port:   int64(0),
		err:    nil,
	},
	{
		origin: "1.2.3./32",
		proto:  "udp",
		port:   int64(0),
		err:    errors.New(""),
	},
}

func TestBuildIPPermission(t *testing.T) {
	for _, tt := range biptable {
		_, err := BuildIPPermission(
			tt.origin,
			tt.proto,
			tt.port,
		)
		if (err != nil && tt.err == nil) ||
			(err == nil && tt.err != nil) {
			t.Error(err)
		}
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

var fsgbntable = []struct {
	name string
	vpc  string
	ret  []string
}{
	{
		name: "",
		vpc:  "",
		ret: []string{
			"sg-1234",
		},
	},
}

func TestFindSGByName(t *testing.T) {
	svc := &mockEC2Client{}
	for _, tt := range fsgbntable {
		ret := FindSGByName(tt.name, tt.vpc, svc)
		for index := range ret {
			if ret[index] != tt.ret[index] {
				t.Error("Unexpected output")
			}
		}
	}
}

var fsgwtable = []struct {
	cidr string
	err  error
	ret  []string
}{
	{
		cidr: "1.2.3.4/32",
		err:  nil,
		ret: []string{
			"sg-1234",
		},
	},
}

func TestFindSecurityGroupsWithRange(t *testing.T) {
	svc := &mockEC2Client{}
	for _, tt := range fsgwtable {
		ret, err := FindSecurityGroupsWithRange(svc, tt.cidr)
		if (err != nil && tt.err == nil) ||
			(err == nil && tt.err != nil) {
			t.Error("Unexpected/mismatched error")
		}
		if len(ret) != len(tt.ret) {
			t.Error("Mismatched results and expectations length")
		}
		for k, v := range ret {
			if v != tt.ret[k] {
				t.Error("Unexpected output")
			}
		}
	}
}

var ncictable = []struct {
	cidr string
	ip   string
	ret  bool
	err  error
}{
	{
		cidr: "0.0.0.0/0",
		ip:   "1.2.3.4/32",
		ret:  true,
		err:  nil,
	},
	{
		cidr: "deadbeef",
		ip:   "1.2.3.4/32",
		ret:  false,
		err:  errors.New(""),
	},
	{
		cidr: "192.168.1.0/24",
		ip:   "192.168.1.1/32",
		ret:  true,
		err:  nil,
	},
	{
		cidr: "192.168.1.0/24",
		ip:   "192.168.3.1/32",
		ret:  false,
		err:  nil,
	},
	{
		cidr: "192.168.1.0/24",
		ip:   "192.168.1.0/24",
		ret:  true,
		err:  nil,
	},
	{
		cidr: "192.168.1.0/24",
		ip:   "192.168.3.0/24",
		ret:  false,
		err:  nil,
	},
	{
		cidr: "192.168.0.0/16",
		ip:   "192.168.10.0/24",
		ret:  true,
		err:  nil,
	},
}

func TestNetworkContainsIPCheck(t *testing.T) {
	for _, tt := range ncictable {
		ip, _, _ := net.ParseCIDR(tt.ip)
		ret, err := NetworkContainsIPCheck(tt.cidr, ip)
		if (err != nil && tt.err == nil) ||
			(err == nil && tt.err != nil) {
			t.Error("Unexpected/mismatched error")
		}
		if ret != tt.ret {
			t.Error("Unexpected mismatch")
		}
	}
}
