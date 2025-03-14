package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 800
	screenHeight = 450
)

/*
Author: Alistair Chambers


Create the classic game of Breakout!  The player’s goal is to clear a grid of blocks using a paddle and a bouncing ball. Check class recording to see the game in motion. Start from scratch with this one, don’t start with my code!

The ball should begin locked to the player paddle. After the player presses space, the ball should launch up, and move left or right depending on where the paddle is moving at time of launch. The ball should travel faster to the right the further to the right of the paddle it hits, and vice versa for the left side. See diagram 1 for info.

The ball should bounce off the left, right, and top of the screen.
Hitting the top should set the ball to move down
Hitting the left wall should set the ball to move right
Hitting the right wall should set the ball to move left

Leaving the bottom of the screen is a failure state, and should reset the game.

If the ball hits a block, the block should disappear.
Hitting block from bottom should set ball to move down
Hitting block from top should set ball to move up
Hitting block from left should set ball to move left
Hitting block from right should set ball to move right

Additionally, the ball should speed up a little each time it destroys a block. Find a reasonable speed on your own.Destroying all the blocks should reset the game.

*/

type Ball struct {
	Position rl.Vector2
	Radius   float32
	Speed    rl.Vector2
	Color    rl.Color
}

type Paddle struct {
	Position rl.Vector2
	Size     rl.Vector2
	Speed    float32
	Color    rl.Color
}

type Block struct {
	Position rl.Vector2
	Size     rl.Vector2
	Color    rl.Color
}

func main() {
	rl.InitWindow(screenWidth, screenHeight, "Breakout")
	rl.SetTargetFPS(60)

	ball := initializeBall()
	paddle := initializePaddle()
	blocks := initializeBlocks()

	gameOver := false

	for !rl.WindowShouldClose() {
		if !gameOver {
			updatePaddle(&paddle)
			updateBall(&ball, &paddle, &blocks, &gameOver)

			if len(blocks) == 0 {
				ball = initializeBall()
				paddle = initializePaddle()
				blocks = initializeBlocks()
				gameOver = false
			}
		} else if gameOver{
			// Reset game when 'R' is pressed
			ball = initializeBall()
			paddle = initializePaddle()
			blocks = initializeBlocks()
			gameOver = false
		}

		drawGame(ball, paddle, blocks)
	}

	rl.CloseWindow()
}
//initializeBall creates a new Ball object with the specified position, radius, speed, and color.
func initializeBall() Ball {
	return Ball{
		Position: rl.NewVector2(screenWidth/2, screenHeight-50),
		Radius:   10,
		Speed:    rl.NewVector2(0, 0),
		Color:    rl.White,
	}
}
//initializePaddle creates a new Paddle object with the specified position, size, speed, and color.
func initializePaddle() Paddle {
	return Paddle{
		Position: rl.NewVector2(screenWidth/2, screenHeight-20),
		Size:     rl.NewVector2(100, 10),
		Speed:    5,
		Color:    rl.Gray,
	}
}
//initializeBlocks creates a slice of Block objects with the specified size, spacing, and offset.
func initializeBlocks() []Block {
	 
	blockSize := rl.NewVector2(60, 20)
	blockSpacing := rl.NewVector2(12, 10)
	blockOffset := rl.NewVector2(10, 10)
	blocks := make([]Block, 0)
	//tested with 2 blocks to check for win condition and reset, it works
	for y := 0; y < 5; y++ {
		for x := 0; x < 11; x++ { //based on offset and spacing, 11 blocks fit in the screen perfectly
			block := Block{
				Position: rl.NewVector2(blockOffset.X+float32(x)*(blockSize.X+blockSpacing.X),
					blockOffset.Y+float32(y)*(blockSize.Y+blockSpacing.Y)),
				Size:  blockSize,
				Color: rl.SkyBlue,
			}
			blocks = append(blocks, block)
		}
	}
	return blocks
}
//updatePaddle updates the paddle's position based on user input.
func updatePaddle(paddle *Paddle) {
	if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA) {
		paddle.Position.X -= paddle.Speed
		if paddle.Position.X < 0 {
			paddle.Position.X = 0
		}
	}
	if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
		paddle.Position.X += paddle.Speed
		if paddle.Position.X+paddle.Size.X > screenWidth {
			paddle.Position.X = screenWidth - paddle.Size.X
		}
	}
}
//updateBall updates the ball's position and speed based on collisions with the paddle, walls, and blocks.
func updateBall(ball *Ball, paddle *Paddle, blocks *[]Block, gameOver *bool) {
	if ball.Speed.Y == 0 {
		// Ball is locked to the paddle until the player launches it
		ball.Position.X = paddle.Position.X + paddle.Size.X/2
		ball.Position.Y = paddle.Position.Y - ball.Radius
	}
	// Launch the ball when space is pressed, if the paddle is not moving then the ball will move straight up , if the paddle is moving left or right the ball will move in that direction but at a 45 degree angle
	if rl.IsKeyPressed(rl.KeySpace) && ball.Speed.Y == 0 {
		ball.Speed = rl.NewVector2(0, -5) // Default up movement
		
		if rl.IsKeyDown(rl.KeyLeft) || rl.IsKeyDown(rl.KeyA){
			ball.Speed = rl.NewVector2(-5, -5) // Move diagonally left at 45 degree angle
		} else if rl.IsKeyDown(rl.KeyRight) || rl.IsKeyDown(rl.KeyD) {
			ball.Speed = rl.NewVector2(5, -5) // Move diagonally right at 45 degree angle
		}
	}

	ball.Position.X += ball.Speed.X
	ball.Position.Y += ball.Speed.Y

	if ball.Position.X < ball.Radius || ball.Position.X > screenWidth-ball.Radius {
		ball.Speed.X *= -1
	}
	if ball.Position.Y < ball.Radius {
		ball.Speed.Y *= -1
	}
	if ball.Position.Y > screenHeight {
		*gameOver = true
	}
	if ball.Position.X+ball.Radius > paddle.Position.X && ball.Position.X-ball.Radius < paddle.Position.X+paddle.Size.X &&
		ball.Position.Y+ball.Radius > paddle.Position.Y && ball.Position.Y-ball.Radius < paddle.Position.Y+paddle.Size.Y {

		// Reverse Y direction
		ball.Speed.Y *= -1

		// Adjust X speed based on where it hit the paddle
		ball.Speed.X = (ball.Position.X - (paddle.Position.X + paddle.Size.X/2)) / 5
	}
	//check for collision with blocks
	for i := 0; i < len(*blocks); i++ {
		block := (*blocks)[i]

		if ball.Position.X+ball.Radius > block.Position.X && ball.Position.X-ball.Radius < block.Position.X+block.Size.X &&
			ball.Position.Y+ball.Radius > block.Position.Y && ball.Position.Y-ball.Radius < block.Position.Y+block.Size.Y {

			// Determine which side the ball is hitting from
			ballCenterX := ball.Position.X
			ballCenterY := ball.Position.Y
			blockCenterX := block.Position.X + block.Size.X/2
			blockCenterY := block.Position.Y + block.Size.Y/2

			overlapX := (block.Size.X/2 + ball.Radius) - float32(math.Abs(float64(ballCenterX-blockCenterX)))
			overlapY := (block.Size.Y/2 + ball.Radius) - float32(math.Abs(float64(ballCenterY-blockCenterY)))

			if overlapX < overlapY {
				ball.Speed.X *= -1 // Ball hits left or right side
			} else {
				ball.Speed.Y *= -1 // Ball hits top or bottom
			}

			// Remove the block
			*blocks = append((*blocks)[:i], (*blocks)[i+1:]...)
			break
		}
	}
}
//drawGame draws the game elements on the screen, including the ball, paddle, blocks
func drawGame(ball Ball, paddle Paddle, blocks []Block) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)

	rl.DrawCircleV(ball.Position, ball.Radius, ball.Color)
	rl.DrawRectangleV(paddle.Position, paddle.Size, paddle.Color)

	for i := 0; i < len(blocks); i++ {
		rl.DrawRectangleV(blocks[i].Position, blocks[i].Size, blocks[i].Color)
	}


	rl.EndDrawing()
}
