package barf

//mostly scratch funcs, just wanted to figure out graphs

import (
	"fmt"
)

//Tupal is a float tuple
type Tupal struct {
	I,J float64
}

//Stk is an int stack
type Stk struct {
	Items []int
}

//Push pushes a val onto the stack
func (s *Stk) Push (val int) {
	s.Items = append(s.Items, val)
}

//Pop pops the stack
func (s *Stk) Pop() {
	s.Items = s.Items[:len(s.Items)-1]
}

//Stack is a string stack
type Stack struct {
	Items []string
}

//Push pushes val onto a stack
func (s *Stack) Push(val string) {
	s.Items = append(s.Items, val)
}

//Node is a multipurpose graph node struct
type Node struct{
	Id      int
	Title   string
	State   int
	Parent  *Node
	Child   []*Node
}

//Pop pops the stack
func (s *Stack) Pop() {
	s.Items = s.Items[:len(s.Items)-1]
}

//Graph is here a graph struct of nodes, edges; all string vals
type Graph struct{
	Nodes []string
	Edges map[string][]string
	Weights map[string][]float64
}

//NewGraph inits a new graph
func NewGraph() (g Graph){
	edges := make(map[string][]string)
	nodes := make([]string,0)
	g = Graph{Nodes: nodes, Edges: edges}
	return g
}

//AddNode adds a new node to a graph
func (g *Graph) AddNode(node string) {
	g.Nodes = append(g.Nodes,node)
	g.Edges[node] = []string{}
}

//AddEdge adds an edge from n1, n2 with weight w1
//HAVE THIS RETURN AN ERROR
func (g *Graph) AddEdge(n1,n2 string) {
	t1 := false
	t2 := false
	for _, n := range g.Nodes {
		if n == n1 {
			t1 = true
		}
		if n == n2 {
			t2 = true
		}
	}
	if t1 && t2 {
		g.Edges[n1] = append(g.Edges[n1],n2)
		g.Edges[n2] = append(g.Edges[n2],n1)
	} else {
		fmt.Println("Node not found")
	}
}

//InVec checks if string val is in string vec
func InVec(vec []string, val string) (bool){
	for _, i := range vec {
		if i == val {
			return true
		}
	}
	return false
}

//Getshortpath returns the shortest path by number of visits i guess? what trollery
func Getshortpath(allpaths [][]string) ([]string){
	var shortest, val int
	val = len(allpaths[0])
	for idx, path := range allpaths {
		if len(path) < val {
			val = len(path)
			shortest = idx
		}
	}
	return allpaths[shortest]
}

//Dfs performs a depth first search on a graph for paths from src to dest 
func Dfs(g Graph, src, dest string, path []string, allpaths [][]string) ([][]string, bool){
	haspath := false
	s := Stack{}
	s.Push(src)
	path = append(path, src)
	//fmt.Printf("at node %s \n", src)
	if src == dest {
		allpaths = append(allpaths, path)
		haspath = true
		path = []string{}
		return allpaths, haspath
	}
	for _, node := range g.Edges[src] {
		return Dfs(g, node, dest, path, allpaths)
	}
	path = []string{}
	return allpaths, haspath
}

//Dfsgrid performs depth first search on a grid
func Dfsgrid(grid [][]int,i,j,check int, visited map[Tupil]bool, ncells int) (map[Tupil]bool,[]Tupil){
	height := len(grid)
	width := len(grid[0])
	stk := []Tupil{}
	stk = append(stk,Tupil{i, j})
	//visited := make(map[Tupil]bool)
	cells := []Tupil{}
	var current Tupil
	var iter int
	for len(stk) > 0 && iter == 0{
		if ncells != -1 && len(cells) == ncells {
			iter = -1
			break			
		}		
		current, stk = stk[len(stk)-1], stk[:len(stk)-1]
		i = current.I; j = current.J
		if grid[i][j] != check{continue}
		visited[current] = true
		if ncells == -1 {
			cells = append(cells,current)
		} else if len(cells) < ncells {
			cells= append(cells,current)
		}
		if len(cells) == ncells{
			iter = -1
			break
		}
		for _, t := range []Tupil{
			{i-1,j},
			{i+1,j},
			{i,j+1},
			{i,j-1},
		} {
			if (t.I >=0 && t.I < height) && (t.J >=0 && t.J< width) && (grid[t.I][t.J]==check){
				if _, ok := visited[t]; !ok {
					stk = append(stk, t)
					
				}
			}
		}
	}
	return visited, cells
}

//PrintPath prints a path
func PrintPath(path []string) {
	rez := ""
	for idx, node := range path {
		if idx == len(path)-1 {
			rez += fmt.Sprintf("%s", node)
		}else {
			rez += fmt.Sprintf("%s-->",node)
		}
	}
	fmt.Println(rez)
	
}

//HasPath checks if a path exists from src to dest
//ALL THESE ARE FROM ALWYN'S DP VIDEO (i think)
func HasPath(g Graph, src, dest string, path []string, allpaths [][]string, visited map[string]bool) ([][]string, bool){

	if visited[src] {
		path = []string{}
		return allpaths, false
	}
	
	fmt.Println("at node", src)
	path = append(path,src)
	PrintPath(path)
	visited[src] = true
	
	if src == dest {
		allpaths = append(allpaths, path)
		visited[src] = false
		path = []string{}
		return allpaths, true
	}
	
	for _, node := range g.Edges[src] {
		if allpaths, b := HasPath(g, node, dest, path, allpaths,visited); b {
			return allpaths, true
		}
	}
	
	return allpaths, false
}

//Bfsprint prints the bfs path
func Bfsprint(g Graph,src string) {
	q := []string{}
	q = append(q, src)
	for len(q) > 0 {
		current := q[0]
		q = q[1:]
		fmt.Println("at ",current)
		for _, node := range g.Edges[current] {q = append(q,node)}
	}
}

//celldiv flips rooms rsmol and rlarge and rewrites a grid
func celldiv(nsmol,rsmol,rlarge int, grid[][]int, src Tupil, smolcs, largecs []*Cell) ([][]int,[]Tupil,[]Tupil) {
	gridedit := make([][]int, len(grid))
	lcs := []Tupil{}
	for i, row := range grid{
		gridedit[i] = make([]int, len(row))
		copy(gridedit[i], grid[i])
	}
	
	
	for _, cell := range smolcs{
		gridedit[cell.Row][cell.Col] = rlarge
		lcs = append(lcs, Tupil{cell.Row, cell.Col})
	}
	for _, cell := range largecs{
		gridedit[cell.Row][cell.Col] = 666
	}
	visited := make(map[Tupil]bool)
	//Dfsgrid(grid [][]int,i,j,check int, visited map[Tupil]bool, ncells int)
	visited, scs := Dfsgrid(gridedit,src.I,src.J,666,visited,nsmol)
	for _, t := range scs{
		gridedit[t.I][t.J] = rsmol
	}
	for _, cell := range largecs{
		if gridedit[cell.Row][cell.Col] == 666{
			gridedit[cell.Row][cell.Col] = rlarge
			lcs = append(lcs,Tupil{cell.Row, cell.Col})
		}
	}

	return gridedit, scs, lcs
}

//Bfsgrid performs a breadth first search on a grid 
func Bfsgrid(grid [][]int,i,j,check int, visited map[Tupil]bool, ncells int) (map[Tupil]bool,[]Tupil){
	height := len(grid)
	width := len(grid[0])
	q := []Tupil{}
	q = append(q,Tupil{i, j})
	//visited = make(map[Tupil]bool)
	cells := []Tupil{}
	for len(q) > 0 {
		if ncells != -1 && len(cells) == ncells {
			break			
		}		
		current := q[0]
		q = q[1:]
		i = current.I; j = current.J
		if _, ok := visited[current]; !ok{
			visited[current] = true
			cells = append(cells, Tupil{i,j})
		}
		for _, t := range []Tupil{
			{i-1,j},
			{i+1,j},
			{i,j-1},
			{i,j+1},
		} {
			if (t.I >=0 && t.I < height) && (t.J >=0 && t.J< width) && (grid[t.I][t.J]==check){
				if _, ok := visited[t]; !ok {
					q = append(q, t)
				}
			}
		}
	}
	return visited, cells
}

//Nbcomponents checks for the number of connected components in a grid
//in 4 directions (l, r, t, b)
func Nbcomponents(grid [][]int,check int) ([][]Tupil, []Tupil, int){
	height := len(grid)
	width := len(grid[0])
	visited := make(map[Tupil]bool)
	rezcells := [][]Tupil{}
	cells := []Tupil{}
	ncomps := 0
	for i := 0; i < height; i++ {
		for j:=0; j < width; j++ {
			current := Tupil{i,j}
			if (grid[i][j] == check) {
				if _, ok := visited[current]; !ok{
					ncomps++
					visited, cells = Bfsgrid(grid,i,j,check,visited,-1)
					rezcells = append(rezcells, cells)
				}
			} 
		}
	}
	return rezcells, cells, ncomps
}

