package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
	"time"
)

func rysuj_linie_na_pilce(posX float64, posY float64, radius float64, alpha1 float64, alpha2 float64) {
	rl.DrawLine(int32(posX-radius+(radius*math.Sqrt(2-2*math.Cos(alpha1*math.Pi/180)))*math.Cos(((180-alpha1)/2)*math.Pi/180)), int32(posY-(radius*math.Sqrt(2-2*math.Cos(alpha1*math.Pi/180)))*math.Sin(((180-alpha1)/2)*math.Pi/180)), int32(posX-radius+(radius*math.Sqrt(2-2*math.Cos(alpha2*math.Pi/180)))*math.Cos(((180-alpha2)/2)*math.Pi/180)), int32(posY-(radius*math.Sqrt(2-2*math.Cos(alpha2*math.Pi/180)))*math.Sin(((180-alpha2)/2)*math.Pi/180)), rl.Black)
}

var ball struct {
	radius float64
	posX   float64
	posY   float64
}

func main() {

	rand.NewSource(time.Now().UnixNano())

	const szerokosc int32 = 1600
	const wysokosc int32 = 900

	rl.InitWindow(szerokosc, wysokosc, "Koszykówka")
	defer rl.CloseWindow()

	ball.radius = 30
	ball.posX = float64(rand.Intn(int(szerokosc * 3 / 4)))
	ball.posY = float64(wysokosc - 200)

	var zmiana_kata float64 = 0
	alpha := [8]float64{0, 60, 90, 120, 180, 240, 270, 300}

	// Główna pętla gry
	for !rl.WindowShouldClose() {

		// changing angle and position
		//ball.posX += 0.02
		zmiana_kata = 0.02

		//nadpisanie kata
		for i := 0; i < 8; i++ {
			alpha[i] += zmiana_kata
		}

		for i := 0; i < 8; i++ {
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
		rl.DrawCircle(int32(ball.posX), int32(ball.posY), float32(ball.radius), rl.Orange)

		rysuj_linie_na_pilce(ball.posX, ball.posY, ball.radius, alpha[0], alpha[4])
		rysuj_linie_na_pilce(ball.posX, ball.posY, ball.radius, alpha[1], alpha[7])
		rysuj_linie_na_pilce(ball.posX, ball.posY, ball.radius, alpha[2], alpha[6])
		rysuj_linie_na_pilce(ball.posX, ball.posY, ball.radius, alpha[3], alpha[5])

		rl.EndDrawing()
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
