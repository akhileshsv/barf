package barf

var (
	//main_menus = []string{icon_kass,icon_mosh,icon_bash,icon_tmbr,icon_exit}
	main_menus = []string{"analysis","rcc design","steel design","timber design","exit"}
	mosh_icons = []string{icon_subframe, icon_frame2d, icon_frame3d, icon_col, icon_beam, icon_slab, icon_ftng}
	bash_icons = []string{icon_col, icon_beam, icon_bolt, icon_weld}
	kass_icons = []string{icon_beam, icon_truss, icon_frame2d, icon_truss3d, icon_grid, icon_frame3d}
	//add calcep
	/*
	kass_menus = []string{
		"1d beam analysis",
		"2d truss analysis",
		"2d frame analysis",
		"3d truss analysis",
		"3d grid analysis",
		"3d frame analysis",
		"1d non uniform beam analysis",
		"2d non uniform frame analysis",
		"bolt group analysis",
		"weld group analysis",
		"exit",
	}
	*/
	input_menus = []string{
		"read json text",
		"read json file",
	}
	kass_menus = []string{
		"beam",
		"2d truss",
		"2d frame",
		"3d truss",
		"grid",
		"3d frame",
		"connections",
		"exit",
	}
	mosh_menus = []string{
		"slab",
		"beam",
		"column",
		"footing",
		"continuous beam",
		"sub frame",
		"2d frame",
		"3d frame",
		"exit",
	}
	bash_menus = []string{
		"beam",
		"column",
		"2d truss",
		"exit",
	}
	tmbr_menus = []string{
		"beam",
		"column",
		"2d truss",
		"exit",
	}
	flay_menus = []string{
		"craft",
		"squarify",
		"exit",
	}
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m" //this is actually magenta
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"

)

const (
	icon_barf =`
 __                ___ 
|  |--.---.-.----.'  _|
|  _  |  _  |   _|   _|
|_____|___._|__| |__|  

   bar force calculator                      
`
	
	icon_ftng = `
  ___             __   __             
.'  _.-----.-----|  |_|__.-----.-----.
|   _|  _  |  _  |   _|  |     |  _  |
|__| |_____|_____|____|__|__|__|___  |
                               |_____|
        
                      footing design                                                          
`
	icon_mosh= `
                      __    
.--------.-----.-----|  |--.
|        |  _  |__ --|     |
|__|__|__|_____|_____|__|__|
                            
              rcc design                                                          
`
	icon_col = `
             __                       
.----.-----.|  |.--.--.--------.-----.
|  __|  _  ||  ||  |  |        |     |
|____|_____||__||_____|__|__|__|__|__|
                                                                                       
                      column design                         
`
	icon_slab = `
        __         __    
.-----.|  |.---.-.|  |--.
|__ --||  ||  _  ||  _  |
|_____||__||___._||_____|
                         
            slab design                         
`
	icon_bash = `
 __                 __    
|  |--.---.-.-----.|  |--.
|  _  |  _  |__ --||     |
|_____|___._|_____||__|__|
                          
             steel design
`
	icon_kass = `
 __                      
|  |--.---.-.-----.-----.
|    <|  _  |__ --|__ --|
|__|__|___._|_____|_____|
                                                
               analysis
`
	icon_beam = `
 __                         
|  |--.-----.---.-.--------.
|  _  |  -__|  _  |        |
|_____|_____|___._|__|__|__|
                            
          ...cool beams?
`
	icon_frame2d = `

  ___              __ ___     _ __  
 | __|_ _  _ __   / /|_  ) __| |\ \ 
 | _|| '_|| '  \ | |  / / / _  | | |
 |_| |_|  |_|_|_|| | /___|\__,_| | |
                  \_\           /_/ 

                   2d frame design
`
	icon_frame3d = `

  ___              __ ____    _ __  
 | __|_ _  _ __   / /|__ / __| |\ \ 
 | _|| '_|| '  \ | |  |_ \/ _  | | |
 |_| |_|  |_|_|_|| | |___/\__,_| | |
                  \_\           /_/ 

                   space frame design
`
	icon_subframe = `

          _     __           
  ____  _| |__ / _|_ _ _ __  
 (_-< || | '_ \  _| '_| '  \ 
 /__/\_,_|_.__/_| |_| |_|_|_|
                             
            sub frame design
`
	icon_weld = `
                 __     __ 
.--.--.--.-----.|  |.--|  |
|  |  |  |  -__||  ||  _  |
|________|_____||__||_____|

        weld group analysis
`
	icon_bolt = `
 __           __ __   
|  |--.-----.|  |  |_ 
|  _  |  _  ||  |   _|
|_____|_____||__|____|
                      
    bolt group analysis
`
	icon_truss = `
 __                          
|  |_.----.--.--.-----.-----.
|   _|   _|  |  |__ --|__ --|
|____|__| |_____|_____|_____|

         plane truss analysis                             
`
	icon_truss3d = `

 __                           ______     __ 
|  |_.----.--.--.-----.-----.|__    |.--|  |
|   _|   _|  |  |__ --|__ --||__    ||  _  |
|____|__| |_____|_____|_____||______||_____|
                                            
                       space truss analysis
`
	icon_grid = `
             __     __ 
.-----.----.|__|.--|  |
|  _  |   _||  ||  _  |
|___  |__|  |__||_____|
|_____|                
       3d grid analysis
`
	icon_tmbr = `
 __            __         
|  |_.--------|  |--.----.
|   _|        |  _  |   _|
|____|__|__|__|_____|__|  

           timber design            
`
	icon_flay = `
  ___  __               
.'  _||  |.---.-..--.--.
|   _||  ||  _  ||  |  |
|__|  |__||___._||___  |
                 |_____|

          facility layout             
`
	icon_conn = `
               __          
.-----..--.--.|  |_ .-----.
|     ||  |  ||   _||__ --|
|__|__||_____||____||_____|
                           
             bolts n welds                      
`
	icon_warning = `
 __           __ __   __ 
|  |--.---.-.|  |  |_|  |
|     |  _  ||  |   _|__|
|__|__|___._||__|____|__|
                            
           this is broken
                  |(*_*)-
        and does not work
               -(*_*)|
`

	icon_exit = `
              __ __   
.-----.--.--.|__|  |_ 
|  -__|_   _||  |   _|
|_____|__.__||__|____|

             exit/close                      
`
)

