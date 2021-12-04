# go-ogpgen

# Overview

go-ogpgen can generate ogp image with text or image.

![overview](https://user-images.githubusercontent.com/7589567/144695117-61ef81e7-04ce-4f4d-b5f8-77bc2596f787.png)

([The Gopher](https://blog.golang.org/gopher) on the base iamge was designed by [Renée French.](http://reneefrench.blogspot.com/))

# Usage

## code example

a simple example composed background image and text.

```
func main() {
	g := ogpgen.NewRandomBackGround()
	if err := g.AttachText(&ogpgen.TextCompositionParams{
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

The code for the more complex case of combining images and text, as shown in the overview, is available in [/demo](./demo)
