package routing

import (
	"strings"
)

const (
	//
	charDelimiter string = `/`
	charPlus      string = `+`
	charTail      string = `#`
	//

	// типы узлов
	nodeTypeString NodeType = iota
	nodeTypePrefix          // /str+/
	nodeTypeSuffix          // /+str/
	nodeTypePlus            // /+/
	nodeTypeTail            // /#
)

type (
	NodeType int

	INode interface {
		Type() NodeType
		Validate(p string) bool
	}
	//
	nodeString struct {
		v string
	}

	nodePrefix struct {
		v string
	}

	nodeSuffix struct {
		v string
	}

	nodePlus struct{}

	nodeTail struct{}
)

func (n *nodeString) Validate(s string) bool {
	return n.v == s
}

func (n *nodeString) Type() NodeType {
	return nodeTypeString
}

func (n *nodePrefix) Validate(s string) bool {
	return strings.HasPrefix(s, n.v)
}

func (n *nodePrefix) Type() NodeType {
	return nodeTypePrefix
}

func (n *nodeSuffix) Validate(s string) bool {
	return strings.HasSuffix(s, n.v)
}

func (n *nodeSuffix) Type() NodeType {
	return nodeTypeSuffix
}

func (n *nodePlus) Validate(s string) bool {
	return len(s) > 0
}

func (n *nodePlus) Type() NodeType {
	return nodeTypePlus
}

func (n *nodeTail) Validate(s string) bool {
	return len(s) > 0
}

func (n *nodeTail) Type() NodeType {
	return nodeTypeTail
}

//
// createNodeValue - создает узел нужного типа
//
func createNodeValue(s string) (INode, error) {
	countPlus := strings.Count(s, charPlus)
	if len(s) > 0 && countPlus < 2 {

		if s == charTail {
			return &nodeTail{}, nil
		}

		if s == charPlus {
			return &nodePlus{}, nil
		}
		if countPlus == 0 {
			return &nodeString{v: s}, nil
		}
		if strings.HasPrefix(s, charPlus) {
			return &nodeSuffix{v: strings.TrimPrefix(s, charPlus)}, nil
		}
		return &nodePrefix{v: strings.TrimSuffix(s, charPlus)}, nil
	}
	return nil, errInvalidNodeValue
}
