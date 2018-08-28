package ga

import (
	"math/rand"
	"sort"
	"time"
)

type IndividualFactory func(rng *rand.Rand) Individual

type Ga struct {
	NewIndividual IndividualFactory
	PopSize       int

	totalFitness float64

	Population []Individual
	Best       Individual

	Rnd *rand.Rand

	CreateRate     float64
	KeepRate       float64
	TournamentSize int
}

func (g *Ga) Initialize() {
	// create initial population

	if g.Rnd == nil {
		g.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if g.TournamentSize == 0 {
		g.TournamentSize = 3
	}

	g.Population = make([]Individual, 0, g.PopSize)
	for i := 0; i < g.PopSize; i++ {
		g.Population = append(g.Population, g.NewIndividual(g.Rnd))
	}

	g.Best = g.Population[0]
}

func (g *Ga) pick() Individual {

	return g.tournamentSelection(g.TournamentSize)
	//return g.poolSelection()
}

func (g *Ga) tournamentSelection(n int) Individual {
	var best Individual
	for i := 0; i < n; i++ {
		inst := g.Population[g.Rnd.Intn(len(g.Population))]
		if best == nil || inst.Fitness() > best.Fitness() {
			best = inst
		}
	}
	return best
}

func (g *Ga) poolSelection() Individual {

	random := g.Rnd.Float64()

	if random == 0 {
		return g.Population[0].Clone()
	}

	i := 0
	for ; random > 0; i++ {
		random -= (g.Population[i].Fitness() / g.totalFitness)
	}
	i--
	return g.Population[i]
}

func (g *Ga) Evolve() {
	// create new population in place of current
	// crossover
	// mutate

	newPopulation := make([]Individual, 0, g.PopSize)
	var fitnessSum float64

	if g.CreateRate > 0 {
		num := int(float64(g.PopSize) * g.CreateRate)
		for i := 0; i < num; i++ {
			newPopulation = append(newPopulation, g.NewIndividual(g.Rnd))
		}
	}

	if g.KeepRate > 0 {
		num := int(float64(g.PopSize) * g.KeepRate)
		for i := 0; i < num; i++ {
			newPopulation = append(newPopulation, g.pick())
		}
	}

	for i := len(newPopulation); i < g.PopSize; i++ {
		parent1 := g.pick()
		parent2 := g.pick()
		child := parent1.Crossover(parent2, g.Rnd).Mutate(g.Rnd)
		newPopulation = append(newPopulation, child)
		fitnessSum += child.Fitness()
	}

	g.Population = newPopulation
	g.totalFitness = fitnessSum
}

func (g *Ga) Record() Individual {
	sort.Slice(g.Population, func(i, j int) bool {
		return g.Population[i].Fitness() > g.Population[j].Fitness() // it is a "less" function, so we need bigger first
	})

	return g.Population[0]
}
