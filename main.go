package main

import (
	"context"
	"hypermedia/internal/component"
	"os"
)

func main() {
	c := component.Hello("John")
	c.Render(context.Background(), os.Stdout)
}
