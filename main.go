// Resik
// Author: Dedy Fajar Setyawan
// See: https://github.com/dedyf5/resik

// Package main serves as the entry point for the Resik application,
// delegating execution flow to the root command handler in the cmd package.
package main

import "github.com/dedyf5/resik/cmd"

func main() {
	cmd.Execute()
}
