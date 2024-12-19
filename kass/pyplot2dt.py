#!"C:\Users\Admin\AppData\Local\Programs\Python\Python38\python.exe"
import matplotlib.pyplot as plt
import vapeplot
from shapely.geometry import Point, Polygon, LineString

def xkcdtruss(b):
    plt.xkcd()
    ldict = {}
    tdict = {1:"chord",2:"rafter",3:"tie/strut",4:"webs",5:"purlins",6:"steel"}
    cmap = vapeplot.cmap('macplus')
    norm = Normalize(vmin=1, vmax=7)
    plt.rcParams["font.family"] = "Humor Sans"
    ax = plt.axes()
    for idx, pt in enumerate (b["Coords"]):
        plt.scatter(pt[0],pt[1])
        plt.text(pt[0]+25.0,pt[1]+25.0,idx+1)
    for idx, mem in enumerate(b["Mprp"]):
        
        xs = b["Coords"][mem[0]-1][0],b["Coords"][mem[1]-1][0]
        ys = b["Coords"][mem[0]-1][1],b["Coords"][mem[1]-1][1]
        pb = b["Coords"][mem[0]-1]
        pe = b["Coords"][mem[1]-1]
        cp = mem[3]
        
        lgnd = "_"
        if cp in ldict:
            lgnd = "_"
        else:
            dims = b["Dims"][cp-1]
            lgnd = f'{tdict[cp]}\n{dims}'
            ldict[cp] = True
        col = cmap(norm(cp))
        plt.plot(xs, ys, linestyle="--", color=col)
        poly = b["polys"][idx]
        plt.fill(*poly.exterior.xy, alpha=0.5, fc=col, ec="black",label=lgnd)
    jp = []
    skl = b["skl"]
    try:
        jp = b["Jsloads"]
    except Exception as e:
        print(e)
    else:  
        for p in jp:
            x, y = b['points'][int(p[0])-1][0], b['points'][int(p[0])-1][1]
            if p[1]:
                xa = x + abs(p[1]) * skl
                ya = y
                astyle = "<-"
                if p[1] > 0:
                    astyle = "->"
                plt.annotate(text='', xy=(xa,ya), xytext=(x,y),
                             arrowprops=dict(arrowstyle=astyle),annotation_clip=False)
                plt.text((x + xa+50)/2.0, (y + ya+50)/2.0,
                         f'{str(p[2])}N', fontsize="xx-small")
            if p[2]:
                xa = x
                ya = y + abs(p[2]) * skl
                astyle = "<-"
                if p[2] > 0:
                    astyle = "->"
                plt.annotate(text='', xy=(xa,ya), xytext=(x,y),
                             arrowprops=dict(arrowstyle=astyle),annotation_clip=False)
                plt.text((x + xa+50)/2.0, (y + ya+50)/2.0,
                         f'{str(p[2])}N', fontsize="xx-small")
    plt.xlabel("mm")
    plt.ylabel("mm")
    plt.axis('scaled')
    plt.title(f"{b['title']} view", fontsize="small", loc = "left")
    #plt.autoscale()
    #plt.grid(alpha=0.1)
    plt.legend(bbox_to_anchor=(1.1, 0.2),loc="lower right", frameon = False, fontsize = "xx-small")
    tblock = f"dukka\n{b['title']}\nakhilesh@dukka.in\n{date.today()}"
    plt.figtext(0.95, 0.05,
                tblock, 
                horizontalalignment ="center",
                wrap = True, fontsize = "xx-small")
    ax.spines[['right', 'top']].set_visible(False)
    plt.tight_layout()
    plt.savefig(os.path.join(b["filedir"],"plan.png"),format="png",orientation="landscape",dpi=300)
    plt.show()
