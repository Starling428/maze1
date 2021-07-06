package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"time"
)

func main() {
	const height, width = 75, 75
	type block struct {
		//x, y int
		visited, iswall bool
	}
	maze := [height][width]block{}
	var linecolor, bgcolor color.Gray16
	linecolor = color.Black
	bgcolor = color.White

	type Point struct{ X, Y int }
	var tempBlock Point
	var stack []Point
	var stackTemp []Point

	//stack = append(stack, Point{1,7}) // Push
	//stack = append(stack, Point{2,8})
	//stack = append(stack, Point{3,9}) // Push
	//stack = append(stack, Point{4,10})
	//stack = append(stack, Point{5,11}) // Push
	//stack = append(stack, Point{6,12})

	A := (height - 1) * (height - 1) / 4
	// img creating
	img := image.NewRGBA(image.Rect(0, 0, height, height))
	for y := 1; y <= height; y++ {
		for x := 1; x <= height; x++ {
			img.Set(x, y, bgcolor)
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
				img.Set(i, j, linecolor)
			}
		}
	}

	//text debug
	//for i:=0; i<height; i++ {
	//	for j:=0; j<width; j++ {
	//		if maze[i][j].visited {print(1)} else {print (0)}
	//	}
	//	println()
	//}
	//for i:=1; i < height; i+=2 {
	//	for j:=1; j < width; j+=2 {
	//			maze[i][j].visited=true
	//			A--
	//			img.Set(i, j, color.RGBA{R:255, A:255})
	//	}
	//}

	//for i:=1; i < height; i+=2 {
	//	for j:=1; j < width; j+=2 {
	//		maze[i][j].iswall=true
	//		A--
	//		img.Set(i, j, color.RGBA{R:255, A:255})
	//	}
	//}

	x0, y0 := 5, 5
	maze[x0][y0].visited = true
	A = (height - 1) * (height - 1) / 4
	rand.Seed(time.Now().UTC().UnixNano())
	for A > 0 {
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
			fmt.Println(len(stackTemp), rand.Intn(len(stackTemp)))
			tempBlock = stackTemp[rand.Intn(len(stackTemp))]
			stackTemp = stackTemp[:0]
			maze[(x0+tempBlock.X)/2][(y0+tempBlock.Y)/2].iswall = false
			img.Set((x0+tempBlock.X)/2, (y0+tempBlock.Y)/2, color.White)
			x0 = tempBlock.X
			y0 = tempBlock.Y
			maze[x0][y0].visited = true
			A--
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
