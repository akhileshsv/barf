#!"C:\Users\Akhilesh SV\Desktop\venvz\zero\Scripts\python.exe"
import matplotlib.pyplot as plt
from matplotlib.collections import PolyCollection
from matplotlib import patheffects
import vapeplot
import json
#import mplcyberpunk
#plt.style.use("cyberpunk")
vapeplot.set_palette('jazzcup')
import sys
import random
from datetime import datetime

random.seed(datetime.now().timestamp())

def wire_plot(nodecords,supports, cols, beams, slabnodes):
    with plt.xkcd():
        plt.figure(figsize=(11,8))
        #plt.style.use("grayscale")
        #plt.rcParams['path.effects'] = [patheffects.withStroke(linewidth=0.1)]
        #plt.rcParams['figure.facecolor'] = 'black'
        plt.rcParams.update({'font.size': 6})
        ax = plt.axes(projection='3d')
        ax.xaxis.pane.fill = False
        ax.yaxis.pane.fill = False
        ax.zaxis.pane.fill = False
        
        # Now set color to white (or whatever is "invisible")
        #ax.xaxis.pane.set_edgecolor('b')
        #ax.yaxis.pane.set_edgecolor('b')
        #ax.zaxis.pane.set_edgecolor('b')
        for node, cords in nodecords.items():
            ax.scatter(cords[0],cords[1],cords[2],marker='o')
            ax.text(cords[0],cords[1],cords[2],f'  {node}',fontsize='medium')
        for col in cols:
            xs,ys,zs = zip(nodecords[col[0]],nodecords[col[1]])
            ax.plot(xs, ys, zs,linewidth=3)
        for beam in beamxs:
            xs,ys,zs = zip(nodecords[beam[0]],nodecords[beam[1]])
            ax.plot(xs, ys, zs,linestyle='dashed',linewidth=2)
        for beam in beamys:
            xs,ys,zs = zip(nodecords[beam[0]],nodecords[beam[1]])
            ax.plot(xs, ys, zs,linestyle='dashdot',linewidth=2)
        for sup in supports:
            if sum(supports[sup]) == -6:
                ax.scatter(nodecords[sup][0],nodecords[sup][1],nodecords[sup][2],marker='x')
        verts = []
        for slab in slabnodes:
            xs = [nodecords[s][0] for s in slab]
            ys = [nodecords[s][1] for s in slab]
            zs = [nodecords[s][2] for s in slab]
            ax.plot_trisurf(xs,ys,zs,alpha=0.3)
        
        

        fname = f'{folder}/frame.png'
        plt.tight_layout()
        #plt.savefig(fname,facecolor='black',edgecolor='black',dpi=300)
        plt.savefig(fname,dpi=300)
        plt.show()
    #print(fname)


with open(sys.argv[1], 'r') as file:
    data = file.read()
fdat = json.loads(data)
folder = sys.argv[2]
nodecords = fdat["Nodecords"]
supports = fdat["Supports"]
cols = fdat["Cols"]
beamxs = fdat["Beamxs"]
beamys = fdat["Beamys"]
slabnodes = fdat["Slabnodes"]
nodecords = {int(i):x for i,x in nodecords.items()}
supports = {int(i):x for i,x in supports.items()}
xf = fdat["Xf"]
yf = fdat["Yf"]
wire_plot(nodecords,supports, cols, beams, slabnodes)
frame_plot(xf, nodecords)
       
'''
#cols = {int(i):x for i,x in nodecords.items}

#print(type(cols))
#print(cols)


rez = data.split("|")
nodecords = ast.literal_eval(rez[1])
supports = ast.literal_eval(rez[2])
cols = ast.literal_eval(rez[3])
beamxs = ast.literal_eval(rez[4])
beamys = ast.literal_eval(rez[5])
slabnodes = ast.literal_eval(rez[6])


'''

