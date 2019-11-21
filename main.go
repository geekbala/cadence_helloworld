package main

import (
	"cadencedemo/cadence"
)

func main() {
	cadence.StartCadenceWorker()
	cadence.StartWorkflow()
}
