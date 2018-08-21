package main

import (
	"github.com/rubensayshi/xlspaceship/pkg"
)

func main() {
	s := pkg.NewXLSpaceship()

	pkg.Serve(s)
}
