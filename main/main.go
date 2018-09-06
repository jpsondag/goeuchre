package main

import (
	"fmt"

	"github.com/jpsondag/goeuchre/euchre"
)

func main() {
	g := euchre.New()
	g.Run();
	fmt.Println(g)
}
