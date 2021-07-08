package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func drawBlock(imag *image.RGBA, x0 int, y0 int, scale int, col color.Color) {
	if scale == 1 {
		imag.Set(x0, y0, col)
	} else {
		xt := x0 * scale
		yt := y0 * scale
		for i := 0; i < scale; i++ {
			for j := 0; j < scale; j++ {
				imag.Set(xt+i, yt+j, col)
			}
		}
	}
}

func main() {

	size := flag.Int("size", 201, "Length of the rib of the maze (the maze is a square")
	scale := flag.Int("scale", 3, "Pixels in one block")
	startx := flag.Int("startx", 101, "Start point x coordinate")
	starty := flag.Int("starty", 101, "Start point y coordinate")
	inversion := flag.Bool("inversion", false, "Invert colors")

	flag.Parse()

	if *scale < 1 {
		*scale = 2
	}

	if *size < 5 {
		*size = 63
	} else {
		if *size%2 == 0 {
			*size++
		}
	}

	if *startx < 1 || *startx > *size-2 {
		*startx = 1
	}
	if *starty < 1 || *starty > *size-2 {
		*starty = 1
	}

	if *startx%2 == 0 {
		if *startx >= 2 {
			*startx--
		} else {
			*startx++
		}
	}
	if *starty%2 == 0 {
		if *starty >= 2 {
			*starty--
		} else {
			*starty++
		}
	}

	var linecolor, bgcolor color.Gray16

	if *inversion {
		linecolor = color.White
		bgcolor = color.Black
	} else {
		linecolor = color.Black
		bgcolor = color.White
	}

	var height, width = *size, *size
	type block struct {
		//x, y int
		visited, iswall bool
	}
	maze := [][]block{}

	for i := 0; i < height; i++ {
		maze = append(maze, []block{})
		for j := 0; j < width; j++ {
			maze[i] = append(maze[i], block{})
		}
	}

	type Point struct{ X, Y int }
	var tempBlock Point
	var stack []Point
	var stackTemp []Point

	Unvisited := (height - 1) * (height - 1) / 4
	// img creating
	img := image.NewRGBA(image.Rect(0, 0, *scale*height, *scale*height))
	for y := 1; y <= height; y++ {
		for x := 1; x <= height; x++ {
			drawBlock(img, x, y, *scale, bgcolor)
			//img.Set(x, y, bgcolor)
		}
	}

	// img painting
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if (i%2 != 0 && j%2 != 0) && (i < height-1 && j < width-1) {
				maze[i][j].iswall = false
				maze[i][j].visited = false

			} else {
				maze[i][j].iswall = true
				maze[i][j].visited = true
				//img.Set(i, j, linecolor)
				drawBlock(img, i, j, *scale, linecolor)
			}
		}
	}

	x0, y0 := *startx, *starty
	maze[x0][y0].visited = true
	Unvisited = (height - 1) * (height - 1) / 4
	rand.Seed(time.Now().UTC().UnixNano())
	for Unvisited > 0 {
		if x0 > 1 {
			if !maze[x0-2][y0].visited {
				stackTemp = append(stackTemp, Point{x0 - 2, y0})
			}
		}
		if x0 < height-2 {
			if !maze[x0+2][y0].visited {
				stackTemp = append(stackTemp, Point{x0 + 2, y0})
			}
		}
		if y0 > 1 {
			if !maze[x0][y0-2].visited {
				stackTemp = append(stackTemp, Point{x0, y0 - 2})
			}
		}
		if y0 < height-2 {
			if !maze[x0][y0+2].visited {
				stackTemp = append(stackTemp, Point{x0, y0 + 2})
			}
		}
		if len(stackTemp) != 0 {
			stack = append(stack, Point{x0, y0})
			tempBlock = stackTemp[rand.Intn(len(stackTemp))]
			stackTemp = stackTemp[:0]
			maze[(x0+tempBlock.X)/2][(y0+tempBlock.Y)/2].iswall = false
			//img.Set((x0+tempBlock.X)/2, (y0+tempBlock.Y)/2, color.White)
			drawBlock(img, (x0+tempBlock.X)/2, (y0+tempBlock.Y)/2, *scale, bgcolor)
			x0 = tempBlock.X
			y0 = tempBlock.Y
			maze[x0][y0].visited = true
			Unvisited--
		} else if len(stack) != 0 {
			tempBlock = stack[len(stack)-1]
			x0 = tempBlock.X
			y0 = tempBlock.Y
			stack = stack[:len(stack)-1]
		} else {
			break
		}
	}

	// file creating
	f, err := os.OpenFile("out.png", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println("Can't create a file")
		os.Exit(10)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			fmt.Println("Can't close a file")
			os.Exit(20)
		}
	}(f)
	err = png.Encode(f, img)
	if err != nil {
		fmt.Println("Can't encode an image")
		os.Exit(30)
	}
}
