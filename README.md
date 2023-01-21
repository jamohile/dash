# Dash
A simple experiment manager.

Dash handles running collections of workloads in parallel.

System Architecture:
	Sweep: a single experimental sweep. Runs a number of workers in parallel.
		- generator: generates configs for workers.
		- worker: a function that will be run with a given config
		- results: workers generate results
	Hypersweep:
		- Oh, actually...can this just use a higher-order sweep?
		- Woah.
