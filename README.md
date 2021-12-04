# go-ogp-generator

# Overview

go-ogp-generator can generate ogp image with text or image.

![overview](https://user-images.githubusercontent.com/7589567/144695117-61ef81e7-04ce-4f4d-b5f8-77bc2596f787.png)

([The Gopher](https://blog.golang.org/gopher) on the base iamge was designed by [Renée French.](http://reneefrench.blogspot.com/))

# Usage

## code example

```
func main() {
	baseImg, err := readFile("base_image.png")
	if err != nil {
		panic(err)
	}
	g, err := ogp.NewGenerator(baseImg)
	if err != nil {
		panic(err)
	}

	if err := g.AttachText(&ogp.TextCompositionParams{
		Text: "Hello, World!",
	}); err != nil {
		panic(err)
	}

    if err := g.Save("output.jpg"); err != nil {
		panic(err)
	}
```

```
open ./output.jpg
```
