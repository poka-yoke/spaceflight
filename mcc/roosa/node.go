package roosa

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

// Node type represents the reference between data.
type Node struct {
	parent   *Node
	children []*Node
	content  *route53.ResourceRecordSet
	indent   int
}

// IsRoot returns true if n is a root node.
func (n *Node) IsRoot() bool {
	return n.parent == nil
}

// String method allows printing of nodes and its children.
func (n *Node) String() (output string) {
	indents := ""
	for i := 0; i < n.indent; i++ {
		indents += "\t"
	}
	extra := ""
	l := len(n.content.ResourceRecords)
	for _, record := range n.content.ResourceRecords[:l-1] {
		extra += fmt.Sprintf("%v, ", *record.Value)
	}
	extra += *n.content.ResourceRecords[l-1].Value
	output = fmt.Sprintf(
		"%v%v %v %v\n",
		indents,
		*n.content.Name,
		*n.content.Type,
		extra,
	)
	for _, child := range n.children {
		child.indent = n.indent + 1
		output += child.String()
	}
	if n.IsRoot() {
		output = strings.TrimSuffix(output, "\n")
	}
	return
}
