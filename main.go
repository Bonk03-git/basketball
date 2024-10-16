package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

func main() {

	var szerokosc int32 = 1600
	var wysokosc int32 = 900

	rl.InitWindow(szerokosc, wysokosc, "Koszykówka")
	defer rl.CloseWindow()

	var radius float64 = 30
	var posX float64 = 200
	var posY float64 = float64(wysokosc - 200)
	var zmiana_kata float64 = 0
	alpha := [5]float64{0, 60, 120, 240, 300}

	// Główna pętla gry
	for !rl.WindowShouldClose() {

		// Rysowanie
		posX += 0.02
		zmiana_kata = 0.02

		for i := 0; i < 5; i++ {
			alpha[i] += zmiana_kata
		}

		for i := 0; i < 5; i++ {
			if alpha[i] > 360 {
				for j := 0; j < 1; j++ {
					if alpha[i] > 360 {
						alpha[i] -= 360
						j--
					}
				}
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite) // Czyszczenie tła

		// Pilka
		rl.DrawCircle(int32(posX), int32(posY), float32(radius), rl.Orange)

		// radius * math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180))

		rl.DrawLine(int32(posX-radius+(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Cos(((180-alpha[0])/2)*math.Pi/180)), int32(posY-(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Sin(((180-alpha[0])/2)*math.Pi/180)), int32(posX+radius-(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Cos(((180-alpha[0])/2)*math.Pi/180)), int32(posY+(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Sin(((180-alpha[0])/2)*math.Pi/180)), rl.Black)
		rl.DrawLine(int32(posX+(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Sin(((180-alpha[0])/2)*math.Pi/180)), int32(posY-radius+(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Cos(((180-alpha[0])/2)*math.Pi/180)), int32(posX-(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Sin(((180-alpha[0])/2)*math.Pi/180)), int32(posY+radius-(radius*math.Sqrt(2-2*math.Cos(alpha[0]*math.Pi/180)))*math.Cos(((180-alpha[0])/2)*math.Pi/180)), rl.Black)
		rl.DrawLine(int32(posX-radius+(radius*math.Sqrt(2-2*math.Cos(alpha[1]*math.Pi/180)))*math.Cos(((180-alpha[1])/2)*math.Pi/180)), int32(posY-(radius*math.Sqrt(2-2*math.Cos(alpha[1]*math.Pi/180)))*math.Sin(((180-alpha[1])/2)*math.Pi/180)), int32(posX-radius+(radius*math.Sqrt(2-2*math.Cos(alpha[4]*math.Pi/180)))*math.Cos(((180-alpha[4])/2)*math.Pi/180)), int32(posY-(radius*math.Sqrt(2-2*math.Cos(alpha[4]*math.Pi/180)))*math.Sin(((180-alpha[4])/2)*math.Pi/180)), rl.Black)
		rl.DrawLine(int32(posX-radius+(radius*math.Sqrt(2-2*math.Cos(alpha[2]*math.Pi/180)))*math.Cos(((180-alpha[2])/2)*math.Pi/180)), int32(posY-(radius*math.Sqrt(2-2*math.Cos(alpha[2]*math.Pi/180)))*math.Sin(((180-alpha[2])/2)*math.Pi/180)), int32(posX-radius+(radius*math.Sqrt(2-2*math.Cos(alpha[3]*math.Pi/180)))*math.Cos(((180-alpha[3])/2)*math.Pi/180)), int32(posY-(radius*math.Sqrt(2-2*math.Cos(alpha[3]*math.Pi/180)))*math.Sin(((180-alpha[3])/2)*math.Pi/180)), rl.Black)
		rl.EndDrawing()
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
