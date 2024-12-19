package barf

import (
	"fmt"
	"math"
	"errors"
)

//beam design funcs (shah/is456)

//MinOf returns the minimum of given vals
func MinOf(vars ...float64) float64 {
	//use interface pls
	min := vars[0]
	for _, i := range vars {
		if min > i && i != 0.0{
			min = i
		}
	}
	return min
}

//BmSecAzIs performs (is code) beam section analysis - returns ultimate moment of resistance mur,
//astmax - max steel for a singly reinforced balanced section
//given bm dims, b.Fck, fy, ast, asc
func BmSecAzIs(b *RccBm) (mur, astmax float64){
	psfs := 1.15; k1 := 0.36; k2 := 0.416
	kumax, _, _, _ := BalSecKis(b.Fck, b.Fy, b.DM, psfs, k1, k2)
	effd := b.Dused - b.Cvrt
	xumax := kumax * effd
	fyd := b.Fy/psfs
	var fsc, xu float64
	if b.Bf == 0.0 {
		b.Bf = b.Tyb * (6.0*b.Df + (b.L0/6.0)) + b.Bw
	}
	bfmbw := b.Bf - b.Bw
	if b.Df/effd > 0.2 {
		yf := 0.15 * kumax + 0.65*b.Df
		if yf < b.Df {b.Df = yf} 
	}
	astmax = (0.361 * b.Fck * b.Bw * xumax + 0.446* b.Fck * b.Df * bfmbw)/fyd
	if b. Ast == 0.0 {
		b.Ast = astmax; xu = xumax
	} else {xu = fyd * b.Ast/(0.361*b.Fck*b.Bf)}
	var yf, mur1, mur2 float64
	if b.Asc == 0.0 {
		//singly reinforced
		switch {
		case b.Tyb == 0.0 || xu < b.Df:
			//rectangular beam/ NA within flange
			if b.Tyb > 0.0 {b.Bw = b.Bf}
			bfmbw = 0.0
		case xu > b.Df:
			//NA within web
			//rect stress block inside flange
			xu = (fyd * b.Ast - 0.2899* b.Fck * b.Df * bfmbw)/(b.Fck * (0.361*b.Bw + 0.0669 * bfmbw))
			if xu < 7.0 * b.Df/3.0 {
				yf = 0.15*xu + 0.65 * b.Df
				if yf < b.Df {b.Df = yf}
				
			} else {
				//rect stress block outside flange
				xu = (fyd * b.Ast - 0.446 * b.Fck * b.Df * bfmbw)/(0.361*b.Fck*b.Bw)
				
			}
		}
		if xu > xumax {xu = xumax}
		mur1 = 0.361* b.Fck * b.Bw * xu * (effd - 0.416* xu) * 1e-6
		mur2 = 0.446 * b.Fck * b.Df * bfmbw * (effd - b.Df/2.0) * 1e-6
		mur = mur1 + mur2
	} else {
		//doubly reinforced
		var xua, bw1, df1, df2, M, N, xur, mur1, mur2, mur3 float64
		switch b.Fy {
		case 250.0:
			//mild steel bars Fe 250, linear rel1
			A := 0.361 * b.Fck * b.Bw
			B := 700.0 * b.Asc - fyd*b.Ast + 0.445 * b.Fck* bfmbw * b.Df
			C := -700.0*b.Cvrc*b.Asc
			xu = SolveQuad(A,B,C)
			if xu > 7.0*b.Df/3.0 {
				fsc = 700.0 * (1.0 - b.Cvrc/xu)
				xu = (fyd * (b.Ast - b.Asc) - 0.446 * b.Fck * bfmbw* b.Df)/(0.361 * b.Fck* b.Bw)
			} else {
				xu = (fyd * (b.Ast - b.Asc) - 0.2899 * b.Fck * bfmbw * b.Df)/(0.361* b.Fck * b.Bw + 0.0669 * b.Fck*bfmbw)
			}
			xua = xu
			if xua > xumax {xua = xumax}
		default:
			//HYSD bars
			xua = xumax
			esc := 0.0035 * (1.0 - b.Cvrc/xua)
			fsc = RbrFrcIs(b.Fy, esc)
			if b.Tyb == 0.0 || xua < b.Df {
				b.Bw = b.Bf
				bfmbw = 0.0
			} else if xua < 7.0 * b.Df/3.0 {
				yf = 0.15 * xua + 0.65 * b.Df
				if yf > b.Df {df1 = b.Df} else {df1 = yf}
			} else {
				df1 = b.Df
			}
			M = 0.361 * b.Fck * b.Bw * xua + fsc * b.Asc + 0.446 * b.Fck * bfmbw * df1
			N = fyd * b.Ast
			if M < N {
				switch {
				//seems okay, but shite way.
				case b.Tyb == 0.0 || xua < b.Df:
					bw1 = b.Bf
				case b.Tyb > 0.0 && xua < 7.0 * b.Df/3.0:
					yf = 0.15 * xua + 0.65 * b.Df
					if yf > b.Df {df2 = b.Df} else {df2 = yf}

				case b.Tyb > 0.0 && xua > 7.0 * b.Df/3.0:
					df2 = b.Df
					bw1 = b.Bw
					
				}
				mur1 = 0.361 * b.Fck * bw1 * xua * (effd - 0.416*xua)* 1e-6
				mur2 = fsc * b.Asc * (effd - b.Cvrc) * 1e-6
				mur3 = 0.446 * b.Fck * (b.Bf - bw1) * df2 * (effd - df2/2.0)* 1e-6
				mur = mur1 + mur2 + mur3
				return
			}
			xudiff := math.Abs(xua - xur)
			for xudiff > 1.0 {
				switch {
				case b.Tyb == 0.0:
					bw1 = b.Bw
				case b.Tyb > 0.0 && xua < b.Df:
					bw1 = b.Bf
					df2 = b.Df
				case b.Tyb > 0.0 && 3.0*xua/7.0 < b.Df:
					df2 = 0.15 * xua + 0.65 * b.Df
					bw1 = b.Bw
				case b.Tyb > 0.0 && 3.0 * xua/7.0 > b.Df:
					df2 = b.Df
					bw1 = b.Bw
				}
				xur = (fyd * b.Ast - 0.446 * b.Fck * (b.Bf - bw1)* df2 - fsc * b.Asc)/(0.361*b.Fck*b.Bw)
				xudiff = math.Abs(xur - xua)
				switch {
				case xur < xua:
					switch {
					case xudiff > 20.0:
						xur = xur + xudiff/2.0
					case xudiff > 10.0:
						xur = xur + 5.0
					case xudiff > 5.0:
						xur = xur + 2.0
					case xudiff > 1.0:
						xur = xur + 1.0
					}
				case xur > xua:
					switch {
					case xudiff > 20.0:
						xur = xur - xudiff/2.0
					case xudiff > 10.0:
						xur = xur - 5.0
					case xudiff > 5.0:
						xur = xur - 2.0
					case xudiff > 1.0:
						xur = xur - 0.9
						
					}
				}
				xua = xur
				esc = 0.0035 * (1.0 - b.Cvrc/xua)
				fsc = RbrFrcIs(b.Fy, esc)
			}
		}
		if b.Tyb == 0.0 || xua < b.Df {
			bw1 = b.Bf
		} else if 3.0 * xua/7.0 < b.Df {
			yf = 0.15 * xua + 0.65 * b.Df
			if yf > b.Df {df2 = b.Df} else {df2 = yf}
		} else {
			df2 = b.Df
		}
		mur1 = 0.361 * b.Fck * bw1 * xua * (effd - 0.416*xua)* 1e-6
		mur2 = fsc * b.Asc * (effd - b.Cvrc) * 1e-6
		mur3 = 0.446 * b.Fck * (b.Bf - bw1) * df2 * (effd - df2/2.0)* 1e-6
		mur = mur1 + mur2 + mur3
	}
	return 
} 


//getBf calcs the breadth of flange as per is456
func getBf(b *RccBm) {
	switch b.Endc{
		case 0, 1:
		b.L0 = b.Lspan
		case 2:
		b.L0 = b.Lspan * 0.7
	}
	b.Bf = math.Round(b.Tyb*(6.0*b.Df + b.L0/6.0) + b.Bw)
}

//getBfBs calcs the breadth of flange as per bs8110
func getBfBs(b *RccBm) {
	b.Bf = math.Round(b.Bw + (0.7 * b.Lspan * b.Tyb)/5.0)
}

//BmDIs calcs the area of steel for a beam given a design moment mdu
//as seen in shah section 5.5 
func BmDIs(b *RccBm, mdu float64) (err error, astmax float64){
	//get areas of steel for a (max) bending moment
	psfs := 1.15; k1 := 0.36; k2 := 0.416
	kumax, _, rumax, _ := BalSecKis(b.Fck, b.Fy, b.DM, psfs, k1, k2)
	if b.Bf == 0.0 {getBf(b)}
	dr := math.Sqrt(mdu * 1e6/(rumax*b.Bw)) + b.Cvrt 
	dt := (mdu * 1e6 + 0.15* b.Fck * b.Bf * math.Pow(b.Df,2))/(0.361*b.Fck*b.Bf*b.Df)
	var dldb float64
	if b.Lbd != 0.0 {dldb = b.Lspan/b.Lbd} else {dldb = b.Lspan/12.0}//WARN HERE
	if b.Dused == 0.0 {
		b.Dused = MinOf(dr,dt,dldb)
	}
	err, astmax = BmAstlIs(b, mdu, kumax, rumax)
	//ADD DETAILING FUNCS GODDAMN CHECK FOR NROWS, MIN.AREA
	return
}

//BmAstlIs designs a beam given mur, fck, fy, bw, bf (yes bf), l/d, effcvr, tyb
//returns/updates ast, asv
func BmAstlIs(b *RccBm, mdu, kumax, rumax float64) (err error, astmax float64){
	//depth calc by ?	
	//fmt.Println("effective depth",effd)
	var iter, kiter int
	for iter != -1{
		kiter++
		effd := b.Dused - b.Cvrt
		xumax := kumax * effd
		fyd := b.Fy/psfs
		fy := b.Fy
		var mu1, mu2, murmx1, murmx2, murmax, yfmax, df, yf, ast, xu float64
		bfmbw := b.Bf - b.Bw
		df = b.Df
		if b.Tyb > 0.0 && b.Df/effd > 0.2 {
			yfmax = 0.15 * xumax + 0.65*b.Df
			//if yfmax < b.Df {b.Df = yfmax}
			if yfmax < df {df = yfmax}
		}
		murmx1 = 0.361 * b.Fck * b.Bw * xumax * (effd -0.416*xumax) * 1e-6
		murmx2 = 0.446 * b.Fck * b. Df * bfmbw * (effd - df/2.0) * 1e-6
		murmax = murmx1 + murmx2
		astmax = (0.361 * b.Fck * b.Bw * xumax + 0.446 * b.Fck * bfmbw * df)/fyd
		if mdu < murmax {
			//singly reinforced section
			//mu1 for xu =df
			mu1 = 0.361 * b.Fck * b.Bf* df * (effd - 0.416* df)* 1e-6
			if mdu < mu1 || b.Tyb == 0.0 {
				//rectangular, xu < df
				
				ast = 0.5 * b.Fck * b.Bf * effd * (1.0 - math.Sqrt(1.0- (4.6*mdu*1e6/(b.Fck*b.Bf*math.Pow(effd,2)))))/fy
				//fmt.Println(ast)
				b.Ast = ast
				b.Asc = 0.0
				if b.Flip{
					b.Asc, b.Ast = b.Ast, b.Asc
				}
				//return
			} else {
				//flanged
				mu2 = 0.8423 * b.Fck* b.Bw* df * (effd - 0.98*df)* 1e-6 + 0.446*b.Fck * df * bfmbw * (effd - df/2.0)* 1e-6
				if mdu < mu2 {
					//3xu/7 < df
					A := -b.Fck * (0.1451585 * b.Bw + 0.0050175* b.Bf)
					B := b.Fck * (effd * (0.2941 * b.Bw + 0.0669* b.Bf) - 0.043485 * df * bfmbw)
					C := b.Fck * df * bfmbw * (0.2899* effd - 0.0942175*df) - mdu* 1e6
					xu = SolveQuad(A,B,C)
				} else {
					A := -0.150176 * b.Fck* b.Bw
					B := 0.361 * b.Fck * b.Bw * effd
					C := 0.446 * b.Fck * df * bfmbw * (effd - df/2.0) - mdu* 1e6
					xu = SolveQuad(A,B,C)
				}
				if xu > xumax {xu = xumax}
				yf = 0.15 * xu + 0.65 * df
				if yf > df {yf = b. Df}
				ast = (0.361 * b.Fck * b.Bw * xu + 0.446 * b.Fck * yf* bfmbw)/fyd
				b.Ast = ast
				b.Asc = 0.0
				if b.Flip{
					b.Asc, b.Ast = b.Ast, b.Asc
				}
			}
		} else {	
			//doubly reinforced section
			b.Csteel = true
			xu = xumax
			var bw1, df1, ast1, ast2, mu float64
			if b.Tyb == 0.0 || xu < b. Df {
				//rectangular section
				bw1 = b.Bf
			} else if xu < 7.0* df/3.0 {
				//phlanj
				yf = 0.15 * xu + 0.65 * df
				df1 = yf
				bw1 = b.Bw
			} else {
				df1 = df
				bw1 = b.Bw
			}
			ast1 = (0.361 * b.Fck * bw1* xu + 0.446 * b.Fck* df1 * (b.Bf - bw1))/fyd
			mu = mdu - murmax
			if b.Cvrc == 0.0 {b.Cvrc = b.Cvrt}
			ast2 = mu* 1e6/(fyd*(effd - b.Cvrc))
			ast = ast1 + ast2
			esc := 0.0035* (1.0 - b.Cvrc/xu)
			fsc := RbrFrcIs(fy, esc)
			b.Asc = fyd * ast2/(fsc - 0.446 * b.Fck)
			b.Ast = ast
			if b.Flip{
				b.Asc, b.Ast = b.Ast, b.Asc
			}
		}
		if b.Rslb {
			err = b.RBarGen()
			return
		}
		//fmt.Println("ast, asc->", b.Ast, b.Asc)
		e, mincvr := b.BarGen()
		if e != nil{
			//fmt.Println("ast, asc->", b.Ast, b.Asc)
			//fmt.Println(kiter,e)
			b.Cvrc = mincvr
			b.Cvrt = mincvr
			b.Ast = 0.0
			b.Asc = 0.0
		} else {
			//fmt.Println("ast, asc, rbrc and rbrt->", b.Ast, b.Asc, b.Rbrc, b.Rbrt)
			iter = -1
			err = b.BarLay()
		}
		if kiter > 6{
			iter = -1
			err = errors.New(fmt.Sprint(e," and iteration error"))
			return
		}
	}
	//fmt.Println("df stl->",b.Mid,b.Id,b.Df)
	return
}

//BmShrDIs designs a beam section for shear (shah section 5.6)
//not really used, see beamshear.go
func BmShrDIs(b *RccBm, bsup float64, xs, vs []float64, ibent bool) ([]int, []int, []float64, error){
	/*
	   shah shear design func
	*/
	var idxs, nlegs []int
	var spacing []float64
	fcktucs := map[float64]float64{15.0:2.5,20.0:2.8,25.0:3.1,30.0:3.5,35.0:3.7,40.0:4.0,45.0:4.3,50.0:4.6,55.0:4.8}
	tucmax := fcktucs[b.Fck]
	effd := b.Dused - b.Cvrt
	vcrx := (b.Dused - b.Cvrt + bsup/2.0)/1e3
	vcr := lerp(xs, vs, vcrx)
	//fmt.Println("max allowable shear",tucmax * b.Bw * effd/1000.0," kn")
	if vcr > tucmax * b.Bw * effd/1000.0 {
		return idxs, nlegs, spacing, errors.New("insufficient depth for shear")
	}
	var vus, vusv, vusb, vuc, asb, ptup, beta, vusvmin, vurmin, fy1 float64
	var dnsr int
	//compute shear resisted by concrete vuc
	if ibent {
		//FORGET THIS- IS ALSO WRENG
		asb = RbrArea(b.Diabent) * b.Nbent
		fy1 = b.Fy
		if fy1 > 415.0 {fy1 = 415.0}
		vusb = 0.707 * 0.87 * fy1 * asb/1000.0
	}
	if vusb > 0.5 * vus {
		vusv = 0.5 * vus 
	} else {
		vusv =  vus - vusb
	}
	ptup = 100.0 * (b.Ast + asb)/(b.Bw * effd)
	beta = 0.8 * b.Fck /(6.89 * ptup)
	if beta < 1.0 {beta = 1.0}
	tuc := 0.8499999* math.Sqrt(0.8 * b.Fck) * (math.Sqrt(1.0 + 5.0 * beta)-1)/(6.0 * beta)
	vuc = tuc * b.Bw * effd/1e3
	vusvmin = 0.4 * b.Bw * effd/1000.0
	
	vurmin = vuc + vusvmin
	vusv = vcr - vuc
	dnsr = 1
	if vcr < vurmin {
		vusv = vusvmin
		dnsr = 0
	}
	svmin := 100.0
	minidx := 0
	maxidx := 2
	mxleg := 2
	idxs, nlegs, spacing = RbrShearLink(fy, b.Bw, effd, vusv, svmin, minidx, maxidx, mxleg, dnsr) 
	return idxs, nlegs, spacing, nil
}

