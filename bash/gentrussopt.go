package barf

// //OptMod is the entry func for truss (axial force) model based optimization funcs
// func OptTrsMod(mod kass.Model)(mrez kass.Model, err error){
// 	switch mod.Frmstr{
// 		case "2dt","3dt":
// 		switch mod.Opt{
// 			case 1,11,12,13:
// 			//g.a, adapt. ga, nruns ga, nruns adapt ga 
// 			return trsoptga(mod)
// 			case 2,21,22,23:
// 			//pso, pso w improvement criteria, nruns pso, nruns pso w/impr
// 			return trsoptpso(mod)
// 		}
// 		default:
// 		err = fmt.Errorf("%s model opt is not written yet",mod.Frmstr)
// 	}
// 	return
// }

