![Logo](https://user-images.githubusercontent.com/17712692/213901290-2d6da4e4-4497-4af6-8eaf-aaaae7440a3e.png)
# Dash (WIP)

A simple experiment manager for collecting data across parallel workloads.

## Motivation

I recently started working with a codebase for research. It involves an extremely long-running simulation, and a requirement to collect data across dozens of runs, and dozens of configurations of those runs.

I'd originally modified its python runner to be highly parallel, but ran into some limitations when I wanted to test more abstract cases. Instead of hacking it in, it seemed more future-friendly to solve this problem in the general case. Enter, Dash.

The original codebase is in Python. While, to be honest, it would have made sense to stick with it for compatibility...Dash uses Go for the following reasons.

1. Python parallelism and multi-threading is less than ideal. Conversely, this is one of Go's selling points.
2. Python external program execution is a bit of messy.
3. I recently started using Go, and wanted an excuse to build something in it.

System Architecture:
Sweep: a single experimental sweep. Runs a number of workers in parallel. - generator: generates configs for workers. - worker: a function that will be run with a given config - results: workers generate results
