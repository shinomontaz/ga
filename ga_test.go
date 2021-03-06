package ga

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

type PathFactory struct {
	N       int
	dMatrix [][]float64
	points  []*LatLng
}

func (f *PathFactory) Create() Individual {
	path := Path{way: rand.Perm(f.N), points: f.points}
	return path
}

func TestTsp(t *testing.T) {
	var N = 10
	rand.Seed(time.Now().UnixNano())

	points := createPoints(N)
	dMatrix := calcDistances(points)

	PFactory := &PathFactory{N: N, dMatrix: dMatrix, points: points}

	var ga = Ga{
		NewIndividual: PFactory.Create,
		PopSize:       N * 20}

	ga.Initialize()
	initial := ga.Population[0].Clone()

	countNotChanged := 0
	bestEver := initial

	for i := 0; i < 10000; i++ {
		currBest := ga.Record()
		if bestEver.Fitness() <= currBest.Fitness() {
			countNotChanged++
			bestEver = currBest
		} else {
			countNotChanged = 0
		}

		if countNotChanged > 10 {
			break
		}
		ga.Evolve()
	}

	//	ga.Solve()
	fmt.Printf("Initial Way: %v, Score: %f\n", initial.(Path).way, initial.Fitness())
	fmt.Printf("Way: %v, Score: %f\n", bestEver.(Path).way, bestEver.(Path).Fitness())
}

type LatLng struct {
	Lat float64
	Lng float64
}

type Path struct {
	way    []int
	points []*LatLng
}

func (p Path) Fitness() float64 {
	distance := 0.0
	for i := 0; i < len(p.way)-1; i++ {
		distance += getDistance(p.points[p.way[i+1]], p.points[p.way[i]])
		//		distance += math.Sqrt(math.Pow(p.points[p.way[i+1]].Lat-p.points[p.way[i]].Lat, 2) + math.Pow(p.points[p.way[i+1]].Lng-p.points[p.way[i]].Lng, 2))
		// getDistance()
	}
	return 1 / (distance + 1)
}

func (p Path) Educate() {

}

// Mutate a Path by applying by permutation mutation and/or splice mutation.
func (p Path) Mutate() Individual {
	dice := rand.Float64()

	if dice > 0.01 {
		return p
	}

	clone := p.Clone().(Path)

	randIndex1, randIndex2 := rand.Intn(len(clone.way)-1), rand.Intn(len(clone.way)-1)
	clone.way[randIndex1], clone.way[randIndex2] = clone.way[randIndex2], clone.way[randIndex1]

	return clone
}

func (p Path) Crossover(q Individual) Individual {
	return p
}

func (p Path) Clone() Individual {
	clone := Path{way: make([]int, len(p.way)), points: p.points}
	copy(clone.way, p.way)
	return clone
}

func createPoints(n int) []*LatLng {
	res := make([]*LatLng, 0)
	for i := 0; i < n; i++ {
		res = append(res, &LatLng{
			Lat: rand.Float64() * 100,
			Lng: rand.Float64() * 100,
		})
	}
	return res
}

func calcDistances(points []*LatLng) [][]float64 {
	res := make([][]float64, 0)
	for _, from := range points {
		row := make([]float64, 0)
		for _, to := range points {
			row = append(row, getDistance(from, to))
		}
		res = append(res, row)
	}

	return res
}

func getDistance(from, to *LatLng) float64 {
	if from == to {
		return 0
	}
	return math.Sqrt(math.Pow(from.Lat-to.Lat, 2) + math.Pow(from.Lng-to.Lng, 2))
}
