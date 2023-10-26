import matplotlib.pyplot as plt
from matplotlib.collections import PolyCollection
from matplotlib import patheffects
import vapeplot
#import mplcyberpunk
#plt.style.use("cyberpunk")
vapeplot.set_palette('jazzcup')
import sys
import ast
import random
from datetime import datetime


with open(sys.argv[1], 'r') as datafile:
    data = datafile.read()
rez = data.split("|")
print(rez[1])
print(rez[2])
coords = ast.literal_eval(rez[1])
ms = ast.literal_eval(rez[2])
folder = sys.argv[2]

with plt.xkcd():
    plt.figure(figsize=(11,8))
    plt.rcParams.update({'font.size': 6})
    #ax = plt.axes()
    #ax.xaxis.pane.fill = False
    #ax.yaxis.pane.fill = False
    #ax.zaxis.pane.fill = False
    for c in coords:
        plt.scatter(c[0],c[1],marker='o')
    for mem in ms:
        xs,ys = zip(coords[mem[0]-1],coords[mem[1]-1])
        plt.plot(xs, ys, linewidth=3)
    fname = f'{folder}/subframe_{random.random()}.png'
    plt.tight_layout()
    #plt.savefig(fname,facecolor='black',edgecolor='black',dpi=300)
    plt.savefig(fname,dpi=300)
    plt.show()
    print(fname)
    
