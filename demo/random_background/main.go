package main

import "ogpgen"

func main() {
	g := ogpgen.NewRandomBackground()
	g.AttachText(&ogpgen.TextCompositionParams{
		Text: "Hello World by go-ogpgen",
	})
	g.Save("output.jpeg")
}
