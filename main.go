package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"time"
)

const (
	WIDTH  = 32
	HEIGHT = 24
	// UDLR
	RIGHT = 1
	LEFT  = 2
	DOWN  = 4
	UP    = 8
)

var (
	player Player
	mpp    [WIDTH][HEIGHT]bool
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	player = Player{
		Position: Position{
			X: 1,
			Y: 1,
		},
	}
	InitMap()
	go KeyPress()
	go MovePlayer()
	go Render()
	<-c
	fmt.Println("Done!")
}

type Position struct {
	X int
	Y int
}

type Player struct {
	Position
}

var clear map[string]func()
var controller = make(chan Position, 1)

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cls") //Windows example it is untested, but I think its working
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearScreen() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func KeyPress() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()

	var b []byte = make([]byte, 1)
	for {
		os.Stdin.Read(b)
		switch string(b) {
		case "a", "A":
			controller <- Position{X: -1, Y: 0}
		case "s", "S":
			controller <- Position{X: 0, Y: 1}
		case "w", "W":
			controller <- Position{X: 0, Y: -1}
		case "d", "D":
			controller <- Position{X: 1, Y: 0}
		}
	}
}

func MovePlayer() {
	for {
		k := <-controller
		if !mpp[player.X+k.X][player.Y+k.Y] {
			player.X += k.X
			player.Y += k.Y
		}
	}
}

func Render() {
	for {
		ClearScreen()
		for y := 0; y < HEIGHT; y++ {
			for x := 0; x < WIDTH; x++ {
				if mpp[x][y] {
					val := 0

					if x > 0 && mpp[x - 1][y] {
						val += LEFT
					}

					if x < WIDTH - 1 && mpp[x + 1][y] {
						val += RIGHT
					}

					if y > 0 && mpp[x][y - 1] {
						val += UP
					}

					if y < HEIGHT - 1 && mpp[x][y + 1] {
						val += DOWN
					}

					switch val {
					case UP, DOWN, UP + DOWN:
						fmt.Print("║")
					case LEFT, RIGHT, LEFT + RIGHT:
						fmt.Print("═")
					case RIGHT + DOWN:
						fmt.Print("╔")
					case RIGHT + UP:
						fmt.Print("╚")
					case LEFT + DOWN:
						fmt.Print("╗")
					case LEFT + UP:
						fmt.Print("╝")
					case UP + DOWN + RIGHT:
						fmt.Print("╠")
					case UP + DOWN + LEFT:
						fmt.Print("╣")
					case UP + LEFT + RIGHT:
						fmt.Print("╩")
					case DOWN + LEFT + RIGHT:
						fmt.Print("╦")
					case UP + DOWN + LEFT + RIGHT:
						fmt.Print("╬")
					default:
						fmt.Print("═")
					}
				} else if player.X == x && player.Y == y {
					fmt.Print("☺")
				} else {
					fmt.Print(string(32))
				}
			}
			fmt.Print("\n")
		}
		time.Sleep(time.Millisecond*100)
	}
}

func InitMap() {
	rand.Seed(int64(time.Now().Nanosecond()))
	for y := HEIGHT - 1; y >= 0; y-- {
		for x := WIDTH - 1; x >= 0; x-- {
			if y == 0 || y == HEIGHT-1 || x == 0 || x == WIDTH-1 {
				mpp[x][y] = true
			} else {
				mpp[x][y] = rand.Int()%2 == 0
				if !mpp[x][y] {
					player.X = x
					player.Y = y
				}
			}
		}
	}
}
