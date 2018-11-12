package capcom

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/poka-yoke/spaceflight/internal/test/mocks"
)

func TestListSecurityGroups(t *testing.T) {
	svc := &mocks.EC2Client{}
	svc.SGList = append(
		svc.SGList,
		[]*ec2.SecurityGroup{
			{
				GroupId:     aws.String("sg-1234"),
				GroupName:   aws.String(""),
				Description: aws.String(""),
			},
		}...,
	)

	out := ListSecurityGroups(svc)
	expected := []string{
		fmt.Sprintf("* %10s %20s %s\n", "sg-1234", "", ""),
	}
	for index, line := range out {
		if line != expected[index] {
			t.Error("Unexpected output")
		}
	}
}

func TestCreateSG(t *testing.T) {
	data := []struct {
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

	svc := &mocks.EC2Client{}
	for _, tc := range data {
		t.Run(
			tc.description,
			func(t *testing.T) {
				out := CreateSG(
					tc.name,
					tc.description,
					tc.vpcid,
					svc,
				)
				if out != tc.out {
					t.Error("Unexpected output")
				}
			},
		)
	}
}

func TestFindSGByName(t *testing.T) {
	data := []struct {
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
	svc := &mocks.EC2Client{}
	for _, tc := range data {
		ret := FindSGByName(tc.name, tc.vpc, svc)
		for index := range ret {
			if ret[index] != tc.ret[index] {
				t.Error("Unexpected output")
			}
		}
	}
}

func TestFindSecurityGroupsWithRange(t *testing.T) {
	data := []struct {
		cidr string
		err  error
		ret  []string
	}{
		{
			cidr: "1.2.3.4/32",
			err:  nil,
			ret: []string{
				"sg-1234 22/tcp 1.2.3.4/32",
			},
		},
	}

	svc := &mocks.EC2Client{}
	svc.SGList = append(
		svc.SGList,
		[]*ec2.SecurityGroup{
			{
				GroupId:   aws.String("sg-1234"),
				GroupName: aws.String(""),
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
		}...,
	)

	for _, tc := range data {
		ret, err := FindSecurityGroupsWithRange(svc, tc.cidr)
		if (err != nil && tc.err == nil) ||
			(err == nil && tc.err != nil) {
			t.Error("Unexpected/mismatched error")
		}
		if len(ret) != len(tc.ret) {
			t.Error("Mismatched results and expectations length")
		}
		for k, v := range ret {
			if v != tc.ret[k] {
				t.Errorf(
					"Unexpected output %s != %s",
					v,
					tc.ret[k],
				)
			}
		}
	}
}
