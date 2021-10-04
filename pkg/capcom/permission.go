package capcom

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Permission represents a Permission for a Security Group
type Permission struct {
	sgid, cidr, protocol string
	port                 int64
	errs                 []error
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

// AddToSG adds the permission to the specified Security Group ID
// using the service
func (p *Permission) AddToSG(svc ec2iface.EC2API, sgid string) bool {
	_, err := svc.AuthorizeSecurityGroupIngress(
		&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       &sgid,
			IpPermissions: []*ec2.IpPermission{p.buildIPPermission()},
		},
	)
	if err != nil {
		p.errs = append(
			p.errs,
			NewPermissionError("adding", sgid, err),
		)
		return false
	}
	return true
}

// RemoveToSG removes the permission to the specified Security Group
// ID using the service
func (p *Permission) RemoveToSG(svc ec2iface.EC2API, sgid string) bool {
	_, err := svc.RevokeSecurityGroupIngress(
		&ec2.RevokeSecurityGroupIngressInput{
			GroupId:       &sgid,
			IpPermissions: []*ec2.IpPermission{p.buildIPPermission()},
		},
	)
	if err != nil {
		p.errs = append(
			p.errs,
			NewPermissionError("revoking", sgid, err),
		)
		return false
	}
	return true
}

// Err returns any error that may have occurred during the execution
// and a bool to signal if there are more errors pending to be checked
func (p *Permission) Err() (more bool, out error) {
	switch len(p.errs) {
	case 0:
		return false, nil
	case 1:
		out, p.errs = p.errs[0], nil
		return false, out
	default:
		out, p.errs = p.errs[0], p.errs[1:]
		return true, out
	}

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

// PermissionError implements error interface to encapsulate errors on
// Permission operations
type PermissionError struct {
	err       error
	sgid      string
	operation string
}

// NewPermissionError creates a new PermissionError
func NewPermissionError(operation, sgid string, err error) PermissionError {
	return PermissionError{err: err, sgid: sgid, operation: operation}
}

// Error is required to satisfy the error interface
func (e PermissionError) Error() string {
	return fmt.Sprintf(
		"Error while %s on %s: %s",
		e.operation,
		e.sgid,
		e.err.Error(),
	)
}
