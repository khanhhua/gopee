package engine

import (
	stack "github.com/golang-collections/collections/stack"
	"github.com/tealeg/xlsx"
	"github.com/xuri/efp"
)

type NodeType int8

const (
	NodeTypeRoot     NodeType = 1
	NodeTypeLiteral  NodeType = 2
	NodeTypeRef      NodeType = 3
	NodeTypeFunc     NodeType = 4
	NodeTypeOperator NodeType = 5
)

type Node struct {
	name     string
	nodeType NodeType
	value    interface{}
	children []*Node
	parent   *Node
}

type Engine struct {
	root  Node
	stack *stack.Stack
}

func New(parser *efp.Parser) *Engine {
	tokens := parser.Tokens.Items
	root := Node{
		name:     "root",
		nodeType: NodeTypeRoot,
		value:    nil,
		children: nil,
	}
	current := &root
	index := 0
	count := len(tokens)
	var token *efp.Token

	for index < count {
		token = &tokens[index]
		name := token.TValue
		ttype := token.TType
		tsubtype := token.TSubType

		if ttype == efp.TokenTypeFunction && tsubtype == efp.TokenSubTypeStart {
			current = current.makeNode(NodeTypeFunc, name)
		} else if ttype == efp.TokenTypeOperand {
			if tokens[index+1].TType == efp.OperatorsInfix && tokens[index+2].TType == efp.TokenTypeOperand { // Look ahead
				node := current.makeNode(NodeTypeOperator, tokens[index+1].TValue) // Infix-Operators: = + - * /
				node.makeNode(resolveNodeType(ttype, tsubtype), name)
				node.makeNode(resolveNodeType(tokens[index+2].TType, tokens[index+2].TSubType), tokens[index+2].TValue)
				index += 2
				continue
			}

		} else if tsubtype == efp.TokenSubTypeStop {
			current = current.parent
		}

		index++
	}

	engine := Engine{
		root:  root,
		stack: stack.New(),
	}
	return &engine
}

func (parent *Node) makeNode(nodeType NodeType, name string) *Node {
	node := Node{
		name:     name,
		nodeType: nodeType,
		parent:   parent,
		children: nil,
	}
	if parent.children == nil {
		parent.children = []*Node{&node}
	} else {
		parent.children = append(parent.children, &node)
	}

	return &node
}

func resolveNodeType(ttype string, tsubtype string) NodeType {
	if ttype == efp.TokenTypeFunction && tsubtype == efp.TokenSubTypeStart {
		return NodeTypeFunc
	} else {
		return NodeTypeLiteral
	}
}

func (self *Engine) Execute(xlFile *xlsx.File) string {

	return "NaN"
}
