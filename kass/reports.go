package barf

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"github.com/olekukonko/tablewriter"
)

//Trs2dSrvTable returns a table of 2d truss superposed service load results
func Trs2dSrvTable(mod *Model) (string){
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorYellow)
	hdr := fmt.Sprintf("\n2d plane %s truss \nnumber of nodes %v\nnumber of members %v\n",BaseMatz[mod.Mtyp],len(mod.Js), len(mod.Ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)

	// table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member", "grp id","len","area","cmax","grp","tmax","grp","fc","ft","units"})
	table.SetCaption(true, "ult. limit state values")
	for i:=0; i < len(mod.Ms); i++ {
		member := mod.Ms[i+1]
		cp := member.Mprp[3]
		cgrp := mod.Grez[cp][0]
		tgrp := mod.Grez[cp][1]
		row := fmt.Sprintf("%v, %v, %.3f, %.3f, %.3f, %.3f,%.3f, %.3f,%.3f, %.3f,%s", member.Id, cp,
			member.Geoms[0],member.Geoms[2],member.Cmax,cgrp,member.Tmax,tgrp,member.Cmax/member.Geoms[2],member.Tmax/member.Geoms[2], mod.Units)
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//FrcTable prints a report of applied model forces (nodal and member)
func FrcTable(mod *Model) (string){
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorCyan+"\n\n")
	table := tablewriter.NewWriter(rezstr)
	if len(mod.Jloads) > 0{
		table.SetCaption(true,"nodal forces")
		sort.Slice(mod.Jloads, func(i, j int) bool{
			return mod.Jloads[i][0] < mod.Jloads[j][0]
		})
		switch mod.Frmstr{
			case "1db":
			table.SetHeader([]string{"Node","Fy","Mz","Vec"})
			for _, vec := range mod.Jloads{
				row := fmt.Sprintf("%.f, %.3f, %.3f, %.3f", vec[0],vec[1],vec[2],vec)
				table.Append(strings.Split(row,","))
			}
			case "2dt":
			table.SetHeader([]string{"Node","Fx","Fy","Vec"})
			for _, vec := range mod.Jloads{
				row := fmt.Sprintf("%.f, %.3f, %.3f, %.3f", vec[0],vec[1],vec[2],vec)
				table.Append(strings.Split(row,","))
			}
			case "2df":
			table.SetHeader([]string{"Node","Fx","Fy","Mz","Vec"})
			for _, vec := range mod.Jloads{
				row := fmt.Sprintf("%.f, %.3f, %.3f, %.3f", vec[0],vec[1],vec[2],vec)
				table.Append(strings.Split(row,","))
			}
			case "3dt":
			table.SetHeader([]string{"Node","Fx","Fy","Fz","Vec"})
			for _, vec := range mod.Jloads{
				row := fmt.Sprintf("%.f, %.3f, %.3f, %.3f", vec[0],vec[1],vec[2],vec)
				table.Append(strings.Split(row,","))
			}
			case "3dg":
			table.SetHeader([]string{"Node","Fx","Fy","Fz","Vec"})
			for _, vec := range mod.Jloads{
				row := fmt.Sprintf("%.f, %.3f, %.3f, %.3f", vec[0],vec[1],vec[2],vec)
				table.Append(strings.Split(row,","))
			}
			case "3df":
			table.SetHeader([]string{"Node","Fx","Fy","Fz","Vec"})
			for _, vec := range mod.Jloads{
				row := fmt.Sprintf("%.f, %.3f, %.3f, %.3f", vec[0],vec[1],vec[2],vec)
				table.Append(strings.Split(row,","))
			}
		}	
	}
	table.Render()
	table = tablewriter.NewWriter(rezstr)
	rezstr.WriteString(ColorPurple+"\n\n")
	if len(mod.Msloads) > 0{
		table.SetCaption(true,"member forces")
		sort.Slice(mod.Msloads, func(i, j int) bool{
			return mod.Msloads[i][0] < mod.Msloads[j][0]
		})
		table.SetHeader([]string{"Member", "Ltyp","Wa","Wb","La","Lb","Vec"})
		for _, vec := range mod.Msloads{
			row := fmt.Sprintf("%.f, %.f, %.3f, %.3f,%.3f,%.3f,%.3f", vec[0],vec[1],vec[2],vec[3],vec[4],vec[5],vec)
			table.Append(strings.Split(row,","))
		}	
	}
	table.Render()
	rezstr.WriteString(ColorReset+"\n\n")
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//Frm3dTable prints calculated resutls of a 3d frame
//ADD REZ REPORTS
func Frm3dTable (js map[int]*Node,ms map[int]*Mem, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string) {
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("3d space frame analysis \n number of nodes %v \n number of members %v\n",len(js), len(ms))
	//hdr += fmt.Sprintf("model units (force-length) %s \n",mod.Units)
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	//rezstr.WriteString("\n\nNODE TABLE\n\n")
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Dsp x Dsp y Dsp z","Rot x Rot y Rot Z","FAx FSy FSz","Mx(T) My Mz"})
	table.SetCaption(true,"node table")
	//var data [][]string
	var row string
	for i:=0; i < len(js); i++{
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %.2f, %.2f, %.2f, %.2f, ", node.Id, node.Coords, node.Supports,node.Displ[:3],node.Displ[3:],node.React[:3],node.React[3:])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)
	//rezstr.WriteString("\n\nMEMBER TABLE\n\n")
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Member", "Prp","Geoms","FAxb FSyb FSzb","Mxb Myb Mzb"})
	table.SetCaption(true, "member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]
		row = fmt.Sprintf("%v, %v, %.1f, %.2f, %.2f", member.Id, member.Mprp, member.Geoms[:3],member.Qf[:3],member.Qf[3:7])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//Frm2dTable prints calc results of a 2d frame
func Frm2dTable (js map[int]*Node,ms map[int]*Mem, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string) {
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorYellow)
	hdr := fmt.Sprintf("\n2d frame analysis\n number of nodes %v \n number of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	//rezstr.WriteString("\n\nNODE TABLE\n\n")
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp x","Disp y", "Rot Z","Rx","Ry","Mz"})
	table.SetCaption(true,"node table")
	//var data [][]string
	var row string
	for i:=0; i < len(js); i++ {
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %f, %f, %f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.Displ[2],node.React[0], node.React[1],node.React[2])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)
	table.SetCaption(true, "member table")
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member", "prp","length","em","area","iz","fab","fsb","fmb","fae","fse","fme"})
	
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]
		row = fmt.Sprintf("%v, %v, %.3f,%.3f,%.3f,%.3f, %.3f,%.3f,%.3f,%.3f,%.3f,%.3f", member.Id, member.Mprp, member.Geoms[0],member.Geoms[1],member.Geoms[2],member.Geoms[3],member.Qf[0],member.Qf[1],member.Qf[2],member.Qf[3],member.Qf[4],member.Qf[5])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//FrmEpTable prints (elastic plastic) calc results of a 2d frame
func Frm2dEpTable (hinges [][]int, mlfs, mrez [][]float64, plm, prez []float64,js map[int]*Node,ms map[int]*Mem) (string) {
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorYellow)

	hdr := fmt.Sprintf("\n2d frame elastic-plastic analysis\n number of nodes %v \n number of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp x","Disp y", "Rot Z","Rx","Ry","Mz"})	
	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++ {
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %f, %f, %f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.Displ[2],node.React[0], node.React[1],node.React[2])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)


	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member", "prp","length","em","area","iz","mp","axial load","mb","me"})
	table.SetCaption(true, "member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]
		row = fmt.Sprintf("%v, %v, %.3f,%.3f,%.3f,%.3f, %.3f,%.3f,%.3f,%.3f", member.Id, member.Mprp, member.Geoms[0],member.Geoms[1],member.Geoms[2],member.Geoms[3],plm[i],prez[i],mrez[i][0],mrez[i][1])
		table.Append(strings.Split(row,","))
	}
	table.Render()

	
	rezstr.WriteString(ColorGreen)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member", "hinge at node","load factor","c.l.f"})
	table.SetCaption(true, "sequence of hinge formation")
	for i, h := range hinges{
		row = fmt.Sprintf("%v,%v,%.4f,%.4f",h[0],h[1],mlfs[i][0],mlfs[i][1])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}


//Grd3dTable prints calc results of a 3d grillage
func Grd3dTable (js map[int]*Node,ms map[int]*Mem, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string) {
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("3d grid analysis\n number of nodes %v \n number of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp y","Rot X", "Rot Z","Ry","Mx","Mz"})
	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++ {
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.Displ[2],node.React[0], node.React[1],node.React[2])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)

	table = tablewriter.NewWriter(rezstr)
	//{l, e, gu, iz, jv, cx, cy, cz}
	table.SetHeader([]string{"member", "prp","len","em","g\n(shear mod.)","iz","j\n(polar i)","fsb","ftb","fmb","fse","fte","fme"})
	table.SetCaption(true, "member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]

		row = fmt.Sprintf("%v, %v, %.3f, %.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f", member.Id, member.Mprp, member.Geoms[0],member.Geoms[1],member.Geoms[2],member.Geoms[3],member.Geoms[4],member.Qf[0],member.Qf[1],member.Qf[2],member.Qf[3],member.Qf[4],member.Qf[5])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//Bm1dTable prints calc results of a beam model
func Bm1dTable (js map[int]*Node,ms map[int]*Mem, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string) {
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("1d beam analysis\n number of nodes %v \n number of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)

	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp y","Rot Z", "Ry","Mz"})

	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++{
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %.4f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.React[0], node.React[1])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member", "prp","len","em","iz","ar","fsb", "fmb","fse","fme"})
	table.SetCaption(true, "member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]

		row = fmt.Sprintf("%v, %v, %.3f, %.3f,%.3f,%.3f,%.3f,%.3f,%.3f,%.3f", member.Id, member.Mprp, member.Geoms[0],member.Geoms[1],member.Geoms[2],member.Geoms[3],member.Qf[0],member.Qf[1],member.Qf[2],member.Qf[3])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}


//BmEpTable prints (elastic plastic) calc results of a beam model
func BmEpTable (hinges [][]int,mlfs, mrez [][]float64, plm, prez []float64, js map[int]*Node,ms map[int]*Mem) (string) {
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorYellow)
	hdr := fmt.Sprintf("\n1d beam elastic-plastic analysis\n number of nodes %v \n number of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)

	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp y","Rot Z", "Ry","Mz"})
	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++{
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %.4f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.React[0], node.React[1])
		table.Append(strings.Split(row,","))
	}

	
	table.Render()
	rezstr.WriteString(ColorPurple)



	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"mem", "prp","length","em","iz","mp","axial load","mb","me"})
	table.SetCaption(true, "member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]
		row = fmt.Sprintf("%v, %v, %.3f,%.3f,%.3f, %.3f,%.3f,%.3f,%.3f", member.Id, member.Mprp, member.Geoms[0],member.Geoms[1],member.Geoms[2],plm[i],prez[i],mrez[i][0],mrez[i][1])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	
	rezstr.WriteString(ColorGreen)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member", "hinge at node","load factor","c.l.f"})
	table.SetCaption(true, "sequence of hinge formation")
	for i, h := range hinges{
		row = fmt.Sprintf("%v,%v,%.4f,%.4f",h[0],h[1],mlfs[i][0],mlfs[i][1])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}


//Trs3dTable prints calc results of a 3d truss
func Trs3dTable (js map[int]*Node,ms map[int]*Mem, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string) {
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorYellow)
	
	hdr := fmt.Sprintf("\n3d space truss analysis\nnumber of nodes %v\nnumber of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp x","Disp y", "Disp z", "Rx","Ry","Rz"})
	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++ {
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.Displ[2],node.React[0], node.React[1], node.React[2])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	
	rezstr.WriteString(ColorPurple)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member","prp","len","em","area","fab", "fae","T/C"})
	table.SetCaption(true,"member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]
		frctyp := "T"
		if member.Qf[0] > 0 {
			frctyp = "C"
		}
		row = fmt.Sprintf("%v, %v, %.3f, %.3f,%.3f,%.3f,%.3f,%s", member.Id, member.Mprp, member.Geoms[0],member.Geoms[1],member.Geoms[2],member.Qf[0], member.Qf[1],frctyp)
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//Trs2dTable prints calc results of a 2d truss model
func Trs2dTable(js map[int]*Node,ms map[int]*Mem, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string){
	rezstr := new(strings.Builder)
	rezstr.WriteString(ColorYellow)
	hdr := fmt.Sprintf("\n2d plane truss analysis\nnumber of nodes %v\nnumber of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","SD-x", "SD-y","Disp x","Disp y", "Rx","Ry"})
	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++ {
		node := js[i+1]
		row = fmt.Sprintf("%v, %v, %v, %.4f, %.4f, %.4f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Sdj[0], node.Sdj[1],node.Displ[0], node.Displ[1], node.React[0], node.React[1])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)


	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"member", "prp","len","em","area","fab","fae","T/C","stress"})
	table.SetCaption(true, "member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]
		frctyp := "T"
		if member.Qf[0] > 0 {
			frctyp = "C"
		}
		row = fmt.Sprintf("%v, %v, %.3f, %.3f, %.3f, %.3f, %.3f,%s,%.3f", member.Id, member.Mprp, member.Geoms[0],member.Geoms[1],member.Geoms[2],member.Qf[0],member.Qf[2],frctyp, math.Abs(member.Qf[0]/member.Geoms[1]))
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//BmNpTable prints calc results of an np beam model
func BmNpTable (js map[int]*Node,ms map[int]*MemNp, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string) {
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("\n1d non-uniform beam analysis\nnumber of nodes %v\nnumber of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp y","Rot Z", "Ry","Mz"})
	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++ {
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %.4f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.React[0], node.React[1])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Member","Prp","Type","Length","Depth","Width","Dims","Styp","Qf (local)", "Gf"})
	table.SetCaption(true, "member table")
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]

		row = fmt.Sprintf("%v, %v, %v, %.2f, %.2f, %.2f, %.2f, %v, %.2f, %.2f", member.Id, member.Mprp, member.Ts,member.Ls,member.Ds,member.Bs, member.Dims,member.Styp, member.Qf, member.Gf)
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}

//F2dNpTable prints calc results of a 2d np frame model
func F2dNpTable (js map[int]*Node,ms map[int]*MemNp, dglb,rnode []float64,nsc []int,ndof,ncjt int) (string) {
	rezstr := new(strings.Builder)
	hdr := fmt.Sprintf("\n2d non-uniform frame analysis\n number of nodes %v \n number of members %v\n",len(js), len(ms))
	rezstr.WriteString(hdr)
	rezstr.WriteString(ColorCyan)
	table := tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Node", "Coords","Sup","Disp x","Disp y", "Rot Z","Rx","Ry","Mz"})
	table.SetCaption(true,"node table")
	var row string
	for i:=0; i < len(js); i++ {
		node := js[i+1]
		row = fmt.Sprintf("%v, %.2f, %v, %f, %f, %f, %.4f, %.4f, %.4f", node.Id, node.Coords, node.Supports,node.Displ[0], node.Displ[1],node.Displ[2],node.React[0], node.React[1],node.React[2])
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorPurple)

	table.SetCaption(true,"member table")
	table = tablewriter.NewWriter(rezstr)
	table.SetHeader([]string{"Member","Prp","Type","Length","Depth","Width","Dims","Styp","Qf (local)", "Gf"})
	
	for i:=0; i < len(ms); i++ {
		member := ms[i+1]

		row = fmt.Sprintf("%v, %v, %v, %.2f, %.2f, %.2f, %.2f, %v, %.2f, %.2f", member.Id, member.Mprp, member.Ts,member.Ls,member.Ds,member.Bs, member.Dims,member.Styp, member.Qf, member.Gf)
		table.Append(strings.Split(row,","))
	}
	table.Render()
	rezstr.WriteString(ColorReset)
	rez := fmt.Sprintf("%s",rezstr)
	return rez
}
