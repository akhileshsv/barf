package barf

import (
	"testing"
)

func TestIntersect(t *testing.T){
	// t.Log("testing point rounding")
	// pt := Pt2d{51.24445, 54.798126}
	// t.Log("input->",pt)
	// pt.SetTol(2)
	// t.Log("rounded to two->",pt)
	t.Log("testing edge intersection")
	var p1, p2, p3, p4, px Pt2d
	var cls string
	p1 = Pt2d{0.5,0.5}
	p2 = Pt2d{1.5,0.5}
	p3 = Pt2d{1,0}
	p4 = Pt2d{1,2}
	cls, px = EdgeInt(p1,p2,p3,p4)
	t.Log("intersection of ",p1, p2, p3, p4, "->",cls,px)
	p1 = Pt2d{0,1}
	p2 = Pt2d{2,3}
	p3 = Pt2d{2,3}
	p4 = Pt2d{0,4}
	cls, px = EdgeInt(p1,p2,p3,p4)
	t.Log("intersection of ",p1, p2, p3, p4, "->",cls,px)

	p1 = Pt2d{1,1}
	p2 = Pt2d{3,3}
	p3 = Pt2d{1,3}
	p4 = Pt2d{3,1}
	cls, px = EdgeInt(p1,p2,p3,p4)
	t.Log("intersection of ",p1, p2, p3, p4, "->",cls,px)

	p1 = Pt2d{0,0}
	p2 = Pt2d{1,1}
	p3 = Pt2d{-0.5,-0.5}
	p4 = Pt2d{1,1}
	cls, px = EdgeInt(p1,p2,p3,p4)
	t.Log("intersection of ",p1, p2, p3, p4, "->",cls,px)

	p1 = Pt2d{0.5,0.5}
	p2 = Pt2d{1.5,0.5}
	p3 = Pt2d{-3,0}
	p4 = Pt2d{-3,-6}
	cls, px = EdgeInt(p1,p2,p3,p4)
	t.Log("intersection of ",p1, p2, p3, p4, "->",cls,px)

}
