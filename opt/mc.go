package barf

import (
	"fmt"
	"math/rand"
)

type Node struct{
	//if typ = 1, binary tree node tis
	//else general graph node
	Val int
	Typ int
	Left *Node
	Right *Node
	Adj []*Node	
}

type Tree struct{
	Left *Tree
	Val int
	Right *Tree
	Adj []*Tree
	Root bool
}

func NewTree(typ, k int) (t *Tree){
	switch typ{
		case 0:
		case 1:
		for _, v := range rand.Perm(10) {
			t = Insert(t, (1+v)*k)
		}
	}
	return t
}

func Insert(t *Tree, v int) (*Tree){
	if t == nil {
		return &Tree{Left:nil, Val:v, Right:nil}
	}
	if v < t.Val {
		t.Left = Insert(t.Left, v)
	} else {
		t.Right = Insert(t.Right, v)
	}
	return t
}

func (t *Tree) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Val)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}
