package barf

import (
	"html/template"
)

var tindex = template.Must(template.ParseFiles("srvr/templates/index.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

//kass
var tanalyze = template.Must(template.ParseFiles("srvr/templates/analyze.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tcmod = template.Must(template.ParseFiles("srvr/templates/calcmod.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tcnp = template.Must(template.ParseFiles("srvr/templates/calcnp.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tcep = template.Must(template.ParseFiles("srvr/templates/calcep.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

var tcmodrez = template.Must(template.ParseFiles("srvr/templates/modrez.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

//mosh
var trcc = template.Must(template.ParseFiles("srvr/templates/rcc.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trslab = template.Must(template.ParseFiles("srvr/templates/rccslab.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trbeam = template.Must(template.ParseFiles("srvr/templates/rccbeam.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trbeamsec = template.Must(template.ParseFiles("srvr/templates/rccbeamsec.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trssbeam = template.Must(template.ParseFiles("srvr/templates/rccbeamss.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trcsbeam = template.Must(template.ParseFiles("srvr/templates/rccbeamcs.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

var trcol = template.Must(template.ParseFiles("srvr/templates/rcccol.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trftng = template.Must(template.ParseFiles("srvr/templates/rccftng.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trcbeam = template.Must(template.ParseFiles("srvr/templates/rcccbeam.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trcsubfrm = template.Must(template.ParseFiles("srvr/templates/rccsubfrm.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var trcfrm2d = template.Must(template.ParseFiles("srvr/templates/rccframe2d.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

//bash
var tsteel = template.Must(template.ParseFiles("srvr/templates/steel.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tstlbeam = template.Must(template.ParseFiles("srvr/templates/stlbeam.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

var tstlcol = template.Must(template.ParseFiles("srvr/templates/stlcol.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tstlcolfrm = template.Must(template.ParseFiles("srvr/templates/stlcolfrm.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tstlcolstrt = template.Must(template.ParseFiles("srvr/templates/stlcolstrt.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

var tstltrs = template.Must(template.ParseFiles("srvr/templates/stltrs.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tstltrsmodopt = template.Must(template.ParseFiles("srvr/templates/stltrsmodopt.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tstltrsgen = template.Must(template.ParseFiles("srvr/templates/stltrsgen.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tstltrsrez = template.Must(template.ParseFiles("srvr/templates/stltrsrez.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))

//tmbr
var ttimber = template.Must(template.ParseFiles("srvr/templates/timber.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var ttmbrbeam = template.Must(template.ParseFiles("srvr/templates/tmbrbeam.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var ttmbrcol = template.Must(template.ParseFiles("srvr/templates/tmbrcol.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))



//misc/about
var trrez = template.Must(template.ParseFiles("srvr/templates/rccrez.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tabout = template.Must(template.ParseFiles("srvr/templates/about.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var terror = template.Must(template.ParseFiles("srvr/templates/error.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tmsg = template.Must(template.ParseFiles("srvr/templates/message.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))
var tdoc = template.Must(template.ParseFiles("srvr/templates/docs.tmpl","srvr/templates/header.tmpl","srvr/templates/footer.tmpl"))


//kass htmx templatse
var tmod1db = template.Must(template.ParseFiles(
	"srvr/templates/ex/mod1db.html",
))

var tmod2dt = template.Must(template.ParseFiles(
	"srvr/templates/ex/mod2dt.html",
))

var tmod2df = template.Must(template.ParseFiles(
	"srvr/templates/ex/mod2df.html",
))

var tmod3dt = template.Must(template.ParseFiles(
	"srvr/templates/ex/mod3dt.html",
))

var tmod3dg = template.Must(template.ParseFiles(
	"srvr/templates/ex/mod3dg.html",
))

var tmod3df = template.Must(template.ParseFiles(
	"srvr/templates/ex/mod3df.html",
))

//mosh htmx templates
var tslb1w = template.Must(template.ParseFiles(
	"srvr/templates/ex/slb1w.html",
))

var tslb2w = template.Must(template.ParseFiles(
	"srvr/templates/ex/slb2w.html",
))

var tslb2wcs = template.Must(template.ParseFiles(
	"srvr/templates/ex/slb2wcs.html",
))

var tslb1wcs = template.Must(template.ParseFiles(
	"srvr/templates/ex/slb1wcs.html",
))

var tslbclvr = template.Must(template.ParseFiles(
	"srvr/templates/ex/slbclvr.html",
))

var trcbmcsdz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcbmcsdz.html",
))


var trcbmcsopt = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcbmcsopt.html",
))


var trcbmssdz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcbmssdz.html",
))

var trcbmssopt = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcbmssopt.html",
))

var trcsfopt = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcsfopt.html",
))


var trcsfdz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcsfdz.html",
))

var trcf2dopt = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcf2dopt.html",
))


var trcf2ddz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcf2ddz.html",
))

var trcbmdz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcbmdz.html",
))

var trcbmaz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcbmaz.html",
))


var trccoldz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rccoldz.html",
))

var trccolaz = template.Must(template.ParseFiles(
	"srvr/templates/ex/rccolaz.html",
))

var trcftng = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcftng.html",
))

var trcftngpad = template.Must(template.ParseFiles(
	"srvr/templates/ex/rcftngpad.html",
))

//bash htmx templates

var tsmod2dtopt = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlmod2dtopt.html",
))

var tsmod3dtopt = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlmod3dtopt.html",
))

var tstrsgenopt = template.Must(template.ParseFiles(
	"srvr/templates/ex/stltrsgenopt.html",
))

var tstrsgen = template.Must(template.ParseFiles(
	"srvr/templates/ex/stltrsgen.html",
))


var tstlcolfrmdz = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlcolfrmdz.html",
))

var tstlcolfrmchk = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlcolfrmchk.html",
))


var tstlcolstrtdz = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlcolstrtdz.html",
))

var tstlcolstrtchk = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlcolstrtchk.html",
))

var tstlbmprln = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlbmprln.html",
))

var tstlbmsimp = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlbmprln.html",
))

var tstlbmrgd = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlbmprln.html",
))


var tstlbmchk = template.Must(template.ParseFiles(
	"srvr/templates/ex/stlbmchk.html",
))

//tmbr htmx templates
var ttmbmgrp = template.Must(template.ParseFiles(
	"srvr/templates/ex/tmbmgrp.html",
))

var ttmbmprp = template.Must(template.ParseFiles(
	"srvr/templates/ex/tmbmprp.html",
))

var ttmcolgrp = template.Must(template.ParseFiles(
	"srvr/templates/ex/tmcolgrp.html",
))

var ttmcolprp = template.Must(template.ParseFiles(
	"srvr/templates/ex/tmcolprp.html",
))

var tdocs = template.Must(template.ParseGlob("srvr/templates/docs/*tmpl"))
