package main

import (
	"fmt"
	"log"
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
	initMap()
	go keyPress()
	go movePlayer()
	go render()
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

var clear = map[string]string{
	"linux":   "clear",
	"darwin":  "clear",
	"windows": "cls",
}

var controller = make(chan Position, 1)

func clearScreen() {
	value, ok := clear[runtime.GOOS]
	if ok {
		cmd := exec.Command(value)
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatalln(fmt.Errorf("error on run clear: %w", err))
		}
	} else {
		log.Fatalln("your platform is unsupported! I can't clear terminal screen :(")
	}
}

func keyPress() {
	// disable input buffering
	err := exec.Command("stty", "-f", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		log.Fatalln(fmt.Errorf("error on run disable input buffering: %w", err))
	}
	// do not display entered characters on the screen
	err = exec.Command("stty", "-f", "/dev/tty", "-echo").Run()
	if err != nil {
		log.Fatalln(fmt.Errorf("error on run command: %w", err))
	}
	// restore the echoing state when exiting
	defer func() {
		err := exec.Command("stty", "-f", "/dev/tty", "echo").Run()
		if err != nil {
			log.Fatalln(fmt.Errorf("error on run reset: %w", err))
		}
	}()

	var b = make([]byte, 1)
	for {
		_, err = os.Stdin.Read(b)
		if err != nil {
			log.Fatalln(fmt.Errorf("error on read standard input: %w", err))
		}
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

func movePlayer() {
	for {
		k := <-controller
		if !mpp[player.X+k.X][player.Y+k.Y] {
			player.X += k.X
			player.Y += k.Y
		}
	}
}

func render() {
	for {
		clearScreen()
		for y := 0; y < HEIGHT; y++ {
			for x := 0; x < WIDTH; x++ {
				if mpp[x][y] {
					val := 0

					if x > 0 && mpp[x-1][y] {
						val += LEFT
					}

					if x < WIDTH-1 && mpp[x+1][y] {
						val += RIGHT
					}

					if y > 0 && mpp[x][y-1] {
						val += UP
					}

					if y < HEIGHT-1 && mpp[x][y+1] {
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
		time.Sleep(time.Millisecond * 100)
	}
}

func initMap() {
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
