package ga

import (
	"math/rand"
	"sort"
	"sync"
	"time"
)

type IndividualFactory func(rng *rand.Rand) Individual

type Ga struct {
	NewIndividual IndividualFactory
	PopSize       int

	Rnd *rand.Rand

	totalFitness float64
	Generations  int
	NotChanged   int

	Population []Individual
	Best       Individual
}

func (g *Ga) Initialize() {
	// create initial population
	if g.Rnd == nil {
		g.Rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if g.Generations == 0 {
		g.Generations = g.PopSize * 10000
	}

	if g.NotChanged == 0 {
		g.NotChanged = 1000
	}

	//	g.If.Init(g.Rnd)

	g.Population = make([]Individual, 0, g.PopSize)
	for i := 0; i < g.PopSize; i++ {
		g.Population = append(g.Population, g.NewIndividual(g.Rnd))
	}

	g.Best = g.Population[0]
}

func (g *Ga) pick() Individual {
	return g.Population[0]
}

func (g *Ga) evolve() {
	// create new population in place of current
	// crossover
	// mutate

	newPopulation := make([]Individual, 0, g.PopSize)
	var fitnessSum float64

	var wg sync.WaitGroup
	chEvolved := make(chan Individual, g.PopSize)
	for i := 0; i < g.PopSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parent1 := g.pick()
			//			parent2 := g.pick()
			//			child := parent1.Crossover(parent2)
			child := parent1.Clone()
			child.Mutate(g.Rnd)
			chEvolved <- child
		}()
	}

	go func() {
		wg.Wait()
		close(chEvolved)
	}()

	for newIdividual := range chEvolved {
		newPopulation = append(newPopulation, newIdividual)
		fitnessSum += newIdividual.Fitness()
	}

	g.Population = newPopulation
	g.totalFitness = fitnessSum
}

func (g *Ga) findBest() Individual {
	sort.SliceStable(g.Population, func(i, j int) bool {
		return g.Population[i].Fitness() > g.Population[j].Fitness() // it is a "less" function, so we need bigger first
	})

	return g.Population[0]
}

func (g *Ga) Solve() {
	var currBest Individual
	countNotChanged := 0
	for i := 0; i < g.Generations; i++ {
		currBest = g.findBest()
		if g.Best.Fitness() <= currBest.Fitness() {
			countNotChanged++
			g.Best = currBest
		} else {
			countNotChanged = 0
		}

		if countNotChanged > g.NotChanged {
			return
		}
		g.evolve()
	}
}
