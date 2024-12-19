package barf

import (
	//"fmt"
	"math"
	"gonum.org/v1/gonum/mat"
)

//MemFrcNp returns member end forces (msqf) and builds pf vec(p - pf = sd)
func MemFrcNp(ms map[int]*MemNp, msp [][]float64, nsc []int, ncjt int, ndof int, vdx bool, pfchn chan []interface{}) {
	//map of loaded members msloaded[0] ex. - 1,1,40.25,0,134.16,0
	//msqf - map of local fixed end forces/end reactions
	//vdx - include effect of shear deformation
	msloaded := make(map[int][][]float64)
	qfchn := make(chan []interface{}, len(msp))
	for _, ldcase := range msp{
		member := int(ldcase[0])
		msloaded[member] = append(msloaded[member], ldcase)
		ms[member].Lds = append(ms[member].Lds, ldcase)
	}
	for member, msloads := range msloaded{
		mem := ms[member]
		go FxdEndFrcNp(member, mem, msloads, ncjt, vdx, qfchn)
	}
	msqf := make(map[int][]float64)
	for member := range ms{
		msqf[member] = make([]float64, 2*ncjt)
	}
	//get member FIXED END forces qf
	for i := 0; i < len(msp); i++{
		r := <-qfchn
		//fmt.Println(ColorBlue,"loaded member->",r,ColorReset)
		member, _ := r[0].(int)
		qf, _ := r[1].([]float64)
		mx, _ := r[2].([]float64)
		//vx, _ := r[3].([]float64)
		ms[member].Mx = mx
		//ms[member].Vx = vx
		for i, f := range qf{
			msqf[member][i] += f
		}
	}
	ffchn := make(chan [][]float64, len(msloaded))
	for member, mem  := range ms{
		switch mem.Frmtyp{
		case "1db": 
			_, ok := msloaded[member]
			if ok{
				qf := msqf[member]
				go ffmatb(mem.Mprp[0], mem.Mprp[1], qf, nsc, ncjt, ndof, ffchn)
			} else {
				continue
			}
		default:
			_, ok := msloaded[member]
			if ok{
				//fmt.Println(ColorBlue,"loaded member>",member,ColorReset)
				qf := msqf[member]
				//fmt.Println(ColorBlue,"qf->",qf,ColorReset)
				go ffmat(mem.Mprp[0], mem.Mprp[1], qf, mem.Tmat, nsc, ncjt, ndof, ffchn)
			} else{
				continue
			}
		}
	}
	//get global fixed end force vector pf
	pf := make([]float64, ndof)
	for i := 0; i < len(msloaded); i++{
		rez := <-ffchn
		for _, r := range rez{
			pf[int(r[0])] += r[1]
		}
	}
	rez := make([]interface{}, 3)
	rez[0] = msqf
	rez[1] = pf
	rez[2] = msloaded
	pfchn <- rez
}

//FxdEndFrcNp calculates member fixed end forces via numerical integration (see hulse section 2.6 )
func FxdEndFrcNp(member int, mem *MemNp, msloads [][]float64, ncjt int, vdx bool, qfchn chan []interface{}){
	vxs := make([]float64, 21)
	mxs := make([]float64, 21)
	l := mem.Lspan
	//memrel := mem.Mprp[4]
	qf := make([]float64, 2*ncjt)
	var fab, fsb, fmb, fae, fse, fme float64
	var di, dj float64
	ei := mem.Em * mem.I0/l
	b := mem.Kc * math.Pow(l,2) * ei
	ai := mem.Ka * math.Pow(l,2) * ei
	aj := mem.Kb * math.Pow(l,2) * ei
	for _, ldcase := range msloads {
		ltyp := int(ldcase[1])
		la := ldcase[4]
		lb := ldcase[5]
		var w, wa, wb, ldcx, ldcvr, W, rl, rr float64
		//simply supported end reactions - rl and rr
		switch ltyp {
		case 1:
			w = ldcase[2]
			ldcx = la
			lb = l - la
			W = w
			rl = w * lb/l
			rr = W - rl
			//rr = w * la/l
			fsb += rl; fse += rr
		case 2:
			w = ldcase[2]
			ldcx = la
			lb = l - la
			rl = -w/l; rr = w/l
			fsb += rl; fse += rr
		case 3:
			w = ldcase[2]
			ldcvr = (l - la - lb)
			ldcx = la + ldcvr/2.0
			W = w*ldcvr
			rl = W * (l-ldcx)/l; rr = W - rl
			fsb += rl; fse += rr
			//fmt.Println("XXXXmember->",member, "W->",W, rl, rr)
		case 4:
			wa = ldcase[2]
			wb = ldcase[3]
			ldcvr = l - la - lb
			if wa >= wb {
				W = wb*ldcvr + (wa-wb)*ldcvr/2.0
				ldcx = la + (ldcvr)*((wa+2.0*wb)/(3.0*(wa+wb)))
			} else {
				W = wa*ldcvr + (wb -wa)*ldcvr/2.0
				ldcx = (ldcvr)*((wb+2.0*wa)/(3.0*(wa+wb)))
				ldcx = la + ldcvr - ldcx
			}
			rl = W * (l - ldcx)/l; rr = W - rl
			fsb += rl; fse += rr
			//fmt.Println(ColorGreen,"trap load","rl",rl,"rr",rr,ColorReset)
		case 5:
			w = ldcase[2]
			ldcx = la
			lb = l - la
			W = w
			fab += W*lb/l; fae += W*la/l 
		case 6:
			w = ldcase[2]
			ldcvr = l - la - lb
			ldcx = la + ldcvr/2.0
			W = w * ldcvr
			fab += w*(l - la -lb)*(l - la + lb)/2.0/l
			fae += w*(l - la -lb)*(l + la - lb)/2.0/l
		case 7:
			w = ldcase[2]
			ldcx = la
			W = w
			//ADD FTB for a 3d frame
		}
		for i, x := range mem.Xs {
			x1 := x - la
			switch ltyp{
				case 1:
				switch {
				case x < la://brez == -1: //before load start
					vxs[i] += rl
					mxs[i] += rl * x
				case x >= la:
					vxs[i] += rl - W
					mxs[i] += rl * x - W*(x1)
					//vxs[i] -= rr
					//mxs[i] -= rr * (l - x)
				}
				case 2:
				switch{
					case x < la:
					mxs[i] += rl * x
					vxs[i] += rl
					case x >= la:
					mxs[i] += rl * x + w
					vxs[i] += rl
				}
				default:
				switch{
					case x < la://before load start
					mxs[i] += rl * x
					vxs[i] += rl
					case x >= la + ldcvr: //after load cover
					mxs[i] += rl*x - W*(x-ldcx)
					vxs[i] += rl - W 
					default: //within load cover
					switch ltyp {
					case 3:
						w1 := w * x1
						ldcx1 := la + x1/2.0
						mxs[i] += rl*x - w1*(x-ldcx1)
						vxs[i] += rl - w1
					case 4:
						if wa >= wb {
							wm := (wa - wb) / ldcvr
							dw := wm * x
							wx := wa - dw
							w1 := (wx * x) + (dw * x / 2.0)
							ldcx1 := (x1)*(wa+2.0*wx)/(3.0*(wa+wx))
							mxs[i] += rl*x - w1*(x-ldcx1)
							vxs[i] += rl - w1
						} else {
							if x == la && wa == 0.0{
								mxs[i] += rl * x
								vxs[i] += rl
							} else {
								wm := (wb - wa) / ldcvr
								dw := wm * x
								wx := wa + dw
								w1 := (wa * x) + (dw * x / 2.0)
								ldcx1 := (x1)*(wx+2.0*wa)/(3.0*(wa+wx))
								mxs[i] += rl * x - w1*ldcx1
								vxs[i] += rl - w1
							}
						}
						//fmt.Println(ColorWhite, "mxs, vxs->", mxs[i], vxs[i],ColorReset)
					}
				}
			}
			
			switch ltyp {
			case 1, 2, 3, 4:
				di += mxs[i] * mem.M4[i]/ei
				dj -= mxs[i] * mem.M5[i]/ei
			}
		}
		switch ltyp {
		case 5,6:
			i5 := MemNpI5(mem, la)
			fab -= W * i5/mem.I4
			fae += W * i5
		default:
			//fmb += (mem.M11 * di + mem.M12* dj) * mem.I0*mem.Em/l
			//fme += (mem.M12 * di + mem.M22 * dj) * mem.I0*mem.Em/l
			fmb += (ai * di + b * dj)
			fme += (aj * dj + b * di)
			fsb += (fmb + fme)/l; fse -= (fmb + fme)/l
		}
		//log.Println(ColorYellow,"mem->",member,"ltyp",ltyp,"shears->",fsb, fse, "end moments",fmb, fme,ColorReset)
	}
	//fmt.Println("FxdEndFrcNp\n","member",member,"frcs",fab, fsb, fmb, fae, fse, fme, "\n slopes",di, dj)
	switch mem.Frmtyp {
	case "1db":
		//beam
		qf = []float64{fsb, fmb, fse, fme}
	case "2df":
		//frame
		qf = []float64{fab, fsb, fmb, fae, fse, fme}
	case "3dg":
		//grid
	case "3df":
		//SPACE.IZ final frontier. Desa da logs of... frame
	}
	rez := make([]interface{},3)
	rez[0] = member
	rez[1] = qf
	//fmt.Println(ColorCyan,"FxdEndFrcNp\n","member",member,"qf->",qf,ColorReset)
	rez[2] = mxs
	qfchn <- rez
}

//MemNpI5 gets the I5 factor for an np member (axial load factor)
func MemNpI5(mem *MemNp, la float64) (i5 float64){
	div := mem.Lspan/20.0
	idx := int(la/div)
	for i := 0; i < 21 - idx; i++{
		switch {
		case i == 0:
			i5 += 1.0/mem.Ax[i + idx]
		case i == 20 - idx:
			i5 += 1.0/mem.Ax[i + idx]
		case i % 2 == 0:
			i5 += 2.0/mem.Ax[i + idx]
		case i % 2 != 0:
			i5 += 4.0/mem.Ax[i + idx]
		}
	}
	return
}

//EndFrcNp gets member end forces for an np member (only for beam/plane frame members so far)
func EndFrcNp(member int, mem *MemNp, msqf map[int][]float64, dglb []float64, nsc []int, mfrtyp,  ncjt , ndof int, rchn chan [][]float64, fchn chan []interface{}) {
	//mp, _ := memprp[0].([]int)
	//bkmat := 
	var rez [][]float64
	var frez []interface{}
	switch mfrtyp{
		case 0:
		//beam
		qf := msqf[member]
		//mem.Qf = msqf[member]
		qfmat := mat.NewDense(len(qf), 1, qf)
		vmat := mat.NewDense(2*ncjt, 1, nil)
		jb := mem.Mprp[0]
		je := mem.Mprp[1]
		var n int
		for i := 1; i <= 2*ncjt; i++ {
			if i <= ncjt {
				n = (jb-1)*ncjt + i
			} else {
				n = (je-1)*ncjt + (i - ncjt)
			}
			vmat.Set(i-1, 0, dglb[n-1])
		}
		fmat := mat.NewDense(2*ncjt, 1, nil)
		fmat.Product(mem.Bkmat, vmat)
		fmat.Add(fmat, qfmat)
		//return index and values of fmat (>ndof) for support reaction r
		rez = [][]float64{}
		fvec := make([]float64,2*ncjt)
		for i := 1; i <= 2*ncjt; i++ {
			fvec[i-1] = fmat.At(i-1,0)
			if i <= ncjt {
				n = (jb-1)*ncjt + i
			} else {
				n = (je-1)*ncjt + (i - ncjt)
			}
			sc := nsc[n-1]
			if sc > ndof {
				rez = append(rez, []float64{float64(sc - ndof - 1), fmat.At(i-1, 0)})
			}
		}
		//WHY NOT MODIFY THE STRUCT HERE *fool
		frez = make([]interface{},2)
		frez[0] = member
		frez[1] = fvec
		case 1:
		//2d frame
		qfmat := mat.NewDense(len(msqf[member]), 1, msqf[member])
		vmat := mat.NewDense(2*ncjt, 1, nil)
		jb := mem.Mprp[0]
		je := mem.Mprp[1]
		var n int
		for i := 1; i <= 2*ncjt; i++ {
			if i <= ncjt {
				n = (jb-1)*ncjt + i
			} else {
				n = (je-1)*ncjt + (i - ncjt)
			}
			//DO DIS 
			//vmat.Set(i-1, 0, dglb[n-1]+mem.Vfs[i-1])
			vmat.Set(i-1, 0, dglb[n-1])
		}
		u := mat.NewDense(2*ncjt, 1, nil)
		u.Mul(mem.Tmat, vmat)
		qmat := mat.NewDense(2*ncjt, 1, nil)
		qmat.Mul(mem.Bkmat,u)
		qmat.Add(qmat,qfmat)
		fmat := mat.NewDense(2*ncjt, 1, nil)
		fmat.Mul(mem.Tmat.T(),qmat)
		//return index and values of fmat (>ndof) for support reaction r
		rez = [][]float64{}
		fvec := make([]float64, 2*ncjt)
		qfvec := make([]float64, 2*ncjt)
		for i := 1; i <= 2*ncjt; i++ {
			fvec[i-1] = fmat.At(i-1,0)
			qfvec[i-1] = qmat.At(i-1,0)
			if i <= ncjt {
				n = (jb-1)*ncjt + i
			} else {
				n = (je-1)*ncjt + (i - ncjt)
			}
			sc := nsc[n-1]
			if sc > ndof {
				rez = append(rez, []float64{float64(sc - ndof - 1), fmat.At(i-1, 0)})
			}
		}
		frez = make([]interface{},3)
		frez[0] = member
		frez[1] = fvec
		frez[2] = qfvec	
	}
	//log.Println("member->",member,"fvec->",fvec)
	//log.Println("member->",member,"rez->",rez)
	rchn <- rez
	fchn <- frez
}
