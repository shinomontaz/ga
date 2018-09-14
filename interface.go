package ga

type Individual interface {
	Fitness() float64
	Mutate() Individual
	Crossover(parner Individual) (child Individual)
	Clone() Individual
}
