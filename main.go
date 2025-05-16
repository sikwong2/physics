package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/sdl"
)

const (
	WIDTH  = 1280
	HEIGHT = 960

	// Physics constants
	GRAVITY  = 0.05 // Reduced gravity for slower falling
	FRICTION = 0.99

	// Frame rate control
	TARGET_FPS = 60
	N          = 15
)

type Ball struct {
	x      float32
	y      float32
	dx     float32 // x velocity component
	dy     float32 // y velocity component
	mass   float32
	radius float32
	color  struct{ r, g, b, a uint8 }
}

func NewBall(mass float32, radius float32) *Ball {
	x := radius + rand.Float32()*(WIDTH-2*radius)
	y := radius + rand.Float32()*(HEIGHT-2*radius)

	// Initial random velocity - slower velocities
	dx := (rand.Float32()*2 - 1) * 2
	dy := (rand.Float32()*2 - 1) * 1.5

	// fmt.Printf("Created ball at (%f, %f) with velocity (%f, %f)\n", x, y, dx, dy)

	ball := &Ball{
		x:      x,
		y:      y,
		dx:     dx,
		dy:     dy,
		mass:   mass,
		radius: radius,
		color: struct{ r, g, b, a uint8 }{
			r: uint8(rand.Intn(200) + 55),
			g: uint8(rand.Intn(200) + 55),
			b: uint8(rand.Intn(200) + 55),
			a: 255,
		},
	}
	return ball
}

func (b *Ball) Update() {
	// Apply gravity
	b.dy += GRAVITY

	// Update position based on velocity
	b.x += b.dx
	b.y += b.dy

	// Handle collisions with walls

	// Right wall
	if b.x+b.radius > WIDTH {
		b.x = WIDTH - b.radius
		b.dx = -b.dx * FRICTION
	}

	// Left wall
	if b.x-b.radius < 0 {
		b.x = b.radius
		b.dx = -b.dx * FRICTION
	}

	// Bottom wall
	if b.y+b.radius > HEIGHT {
		b.y = HEIGHT - b.radius
		b.dy = -b.dy * FRICTION

		// Apply some horizontal friction when hitting the floor
		b.dx *= FRICTION
	}

	// Top wall
	if b.y-b.radius < 0 {
		b.y = b.radius
		b.dy = -b.dy * FRICTION
	}
}

func DrawCircle(renderer *sdl.Renderer, centreX, centreY, radius float32, color struct{ r, g, b, a uint8 }) {
	// Set the color for drawing
	renderer.SetDrawColor(color.r, color.g, color.b, color.a)

	diameter := radius * 2
	x := radius - 1
	y := float32(0)
	tx := float32(1)
	ty := float32(1)
	error := tx - diameter

	for x >= y {
		// Draw 8 points for the circle using the 8-way symmetry
		renderer.RenderPoint(centreX+x, centreY-y)
		renderer.RenderPoint(centreX+x, centreY+y)
		renderer.RenderPoint(centreX-x, centreY-y)
		renderer.RenderPoint(centreX-x, centreY+y)
		renderer.RenderPoint(centreX+y, centreY-x)
		renderer.RenderPoint(centreX+y, centreY+x)
		renderer.RenderPoint(centreX-y, centreY-x)
		renderer.RenderPoint(centreX-y, centreY+x)

		if error <= 0 {
			y++
			error += ty
			ty += 2
		}

		if error > 0 {
			x--
			tx += 2
			error += tx - diameter
		}
	}
}

func main() {
	// Seed the random number generator

	defer binsdl.Load().Unload() // sdl.LoadLibrary(sdl.Path())
	defer sdl.Quit()

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	window, renderer, err := sdl.CreateWindowAndRenderer("Bouncing Ball Simulation", WIDTH, HEIGHT, 0)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	defer window.Destroy()

	// Create several bouncing balls
	balls := []*Ball{
		NewBall(10.0, 20.0),
		NewBall(5.0, 15.0),
		NewBall(15.0, 25.0),
		NewBall(8.0, 18.0),
		NewBall(12.0, 22.0),
	}

	frameCount := 0
	lastTime := time.Now()
	frameDelay := time.Second / TARGET_FPS

	sdl.RunLoop(func() error {
		var event sdl.Event
		for sdl.PollEvent(&event) {
			if event.Type == sdl.EVENT_QUIT {
				return sdl.EndLoop
			}
		}

		// Clear screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// Display FPS counter and info
		frameCount++
		currentTime := time.Now()
		elapsed := currentTime.Sub(lastTime)
		fps := 1.0 / elapsed.Seconds()

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DebugText(10, 10, fmt.Sprintf("Frame: %d | FPS: %.1f", frameCount, fps))
		renderer.DebugText(10, 30, fmt.Sprintf("Balls: %d | Gravity: %.3f", len(balls), GRAVITY))

		// Update and draw each ball
		for _, ball := range balls {
			ball.Update()
			DrawCircle(renderer, ball.x, ball.y, ball.radius, ball.color)
		}

		// Add a new ball every 180 frames (about 3 seconds at 60fps)
		if frameCount%180 == 0 && len(balls) < N {
			radius := float32(rand.Intn(20) + 10)
			balls = append(balls, NewBall(radius/2, radius))
		}

		renderer.Present()

		// Control frame rate
		frameTime := time.Since(currentTime)
		if frameDelay > frameTime {
			time.Sleep(frameDelay - frameTime)
		}

		lastTime = time.Now()
		return nil
	})
}
