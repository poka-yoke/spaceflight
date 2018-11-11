package capcom

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Permission represents a Permission for a Security Group
type Permission struct {
	sgid, cidr, protocol string
	port                 int64
}

// NewPermission returns a pointer to a new Permission object
func NewPermission(origin, protocol string, port int64) (*Permission, error) {
	perm := Permission{protocol: protocol, port: port}
	switch {
	case strings.HasPrefix(origin, "sg-"):
		// It's a security group
		perm.sgid = origin

	case isCIDR(origin):
		// It's a valid CIDR
		perm.cidr = origin
	default:
		err := fmt.Errorf(
			"%s is neither sgid nor IP range in CIDR notation",
			origin,
		)
		return nil, err
	}
	return &perm, nil
}

// buildIPPermission provides an IpPermission object fully populated
func (p *Permission) buildIPPermission() (
	perm *ec2.IpPermission,
) {
	perm = &ec2.IpPermission{
		FromPort:   &p.port,
		ToPort:     &p.port,
		IpProtocol: &p.protocol,
	}
	if p.sgid != "" {
		perm.UserIdGroupPairs = []*ec2.UserIdGroupPair{{GroupId: &p.sgid}}
	}
	if p.cidr != "" {
		perm.IpRanges = []*ec2.IpRange{{CidrIp: &p.cidr}}
	}
	return perm
}

// AddToSG adds the permission to the specified Security Group ID
// using the service
func (p *Permission) AddToSG(svc ec2iface.EC2API, sgid string) bool {
	return authorizeAccessToSecurityGroup(svc, p, sgid)
}

// RemoveToSG removes the permission to the specified Security Group
// ID using the service
func (p *Permission) RemoveToSG(svc ec2iface.EC2API, sgid string) bool {
	return revokeAccessToSecurityGroup(svc, p, sgid)
}

// AuthorizeAccessToSecurityGroup adds the specified permissions to the Ingress
// list of the destination security group on protocol and port
func authorizeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *Permission,
	destination string,
) bool {
	out, error := svc.AuthorizeSecurityGroupIngress(
		&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       &destination,
			IpPermissions: []*ec2.IpPermission{perm.buildIPPermission()},
		})
	if error != nil {
		log.Panic(error)
	}
	if !(out != nil) {
		return false
	}
	return true
}

// RevokeAccessToSecurityGroup adds the specified permissions to the Ingress
// list of the destination security group on protocol and port
func revokeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *Permission,
	destination string,
) bool {
	out, error := svc.RevokeSecurityGroupIngress(
		&ec2.RevokeSecurityGroupIngressInput{
			GroupId:       &destination,
			IpPermissions: []*ec2.IpPermission{perm.buildIPPermission()},
		})
	if error != nil {
		log.Panic(error)
	}
	if !(out != nil) {
		return false
	}
	return true
}
