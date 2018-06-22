package ga

import "math/rand"

type Individual interface {
	Fitness() float64
	Mutate(rng *rand.Rand)
	Crossover(parner Individual, rng *rand.Rand) (child Individual)
	Clone() Individual
}
