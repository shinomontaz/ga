package ga

import (
	"math/rand"
	"sort"
	"sync"
)

//type IndividualFactory func(rng *rand.Rand) Individual

type IndividualFactory func() Individual

type Ga struct {
	NewIndividual IndividualFactory
	PopSize       int

	totalFitness float64
	Generations  int

	Population []Individual
	Best       Individual
	lock       sync.Mutex
}

func (g *Ga) Initialize() {
	// create initial population
	//	rand.Seed(time.Now().UnixNano())

	if g.Generations == 0 {
		g.Generations = g.PopSize * 10
	}

	g.Population = make([]Individual, 0, g.PopSize)
	for i := 0; i < g.PopSize; i++ {
		g.Population = append(g.Population, g.NewIndividual())
	}

	g.Best = g.Population[0]
}

func (g *Ga) pick() Individual {
	random := rand.Float64()

	// отсортировать популяцию в порядке убывания значений

	if random == 0 { // just get best
		return g.Population[0]
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

	var wg sync.WaitGroup
	chEvolved := make(chan Individual, g.PopSize)
	for i := 0; i < g.PopSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			parent1 := g.pick()
			parent2 := g.pick()
			child := parent1.Crossover(parent2).Mutate()
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

func (g *Ga) Record() Individual {
	sort.Slice(g.Population, func(i, j int) bool {
		return g.Population[i].Fitness() > g.Population[j].Fitness() // it is a "less" function, so we need bigger first
	})

	return g.Population[0]
}
