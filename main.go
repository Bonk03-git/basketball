package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

func main() {

	rl.InitWindow(800, 600, "Koszykówka")
	defer rl.CloseWindow()

	var radius float32 = 30
	var posX float32 = 400
	var posY float32 = 300

	// Główna pętla gry
	for !rl.WindowShouldClose() {

		// Rysowanie
		posX += 0.02
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite) // Czyszczenie tła

		// Pilka
		rl.DrawCircle(int32(posX), int32(posY), radius, rl.Orange)
		rl.DrawLine(int32(posX-radius), int32(posY), int32(posX+radius), int32(posY), rl.Black)
		rl.DrawLine(int32(posX), int32(posY-radius), int32(posX), int32(posY+radius), rl.Black)
		rl.DrawLine(int32(posX-radius/2), int32(float64(posY)-float64(radius)*math.Sqrt(3)/2), int32(posX-radius/2), int32(float64(posY)+float64(radius)*math.Sqrt(3)/2), rl.Black)
		rl.DrawLine(int32(posX+radius/2), int32(float64(posY)-float64(radius)*math.Sqrt(3)/2), int32(posX+radius/2), int32(float64(posY)+float64(radius)*math.Sqrt(3)/2), rl.Black)
		rl.EndDrawing()
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
