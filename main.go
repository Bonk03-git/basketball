package main

// PAMIĘTAĆ ZAMIENIAĆ WSPOLRZEDNE X Z Y NA KONIEC PROJEKTU!!!
import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"math/rand"
	"time"
)

// 50px -> 1m

type Kula struct {
	promien float64
	posX    float64
	posY    float64
	posZ    float64
	rotX    float64
	rotY    float64
	rotZ    float64
	masa    float64
}

type Punkt_przylozenia struct {
	kat_fi   float64 // wspolrzedne sferyczne punktu przylozenia 0 ; 2pi
	kat_teta float64 // -pi/2 ; pi/2
}

type Wektor_sily struct {
	dlugosc_x     float64 // pixele
	dlugosc_y     float64 // pixele
	dlugosc_z     float64 // pixele
	przylozenie_x float64 // pixele przylozenia sa wzgledem srodka pilki
	przylozenie_y float64 // pixele
	przylozenie_z float64 // pixele
}

type Wektor struct {
	x float64
	y float64
	z float64
}

func main() {

	rand.NewSource(time.Now().UnixNano())

	const szerokosc int = 1600
	const wysokosc int = 900
	const czas_dzialania_sily_na_pilke = 0.5
	const krok_czasowy = 0.01
	const g = -981

	pilka := Kula{
		promien: 12,                            // pixele
		posX:    float64(rand.Intn(wysokosc)),  // pixele
		posY:    float64(rand.Intn(szerokosc)), // pixele
		posZ:    200,                           // pixele
		rotX:    0,                             // radiany
		rotY:    0,                             // radiany
		rotZ:    0,                             // radiany
		masa:    0.5,                           // kilogramy
	}

	var punkty_przylozenia_wektorow [3]Punkt_przylozenia
	var wektory_sily [3]Wektor_sily
	var przyspieszenia_z_wektorow [3]Wektor
	var momenty_z_wektorow [3]Wektor
	var przyspieszenia_katowe_z_wektorow [3]Wektor
	var wypadkowa_przyspieszen = Wektor{
		x: 0,
		y: 0,
		z: 0,
	}
	var wypadkowa_przyspieszen_katowych = Wektor{
		x: 0,
		y: 0,
		z: 0,
	}
	var I = 2 / 5 * (pilka.masa) * pilka.promien * pilka.promien

	for i := 0; i < 3; i++ {

		punkty_przylozenia_wektorow[i] = Punkt_przylozenia{ // wyznaczane przez gracza
			kat_fi:   math.Pi / 2,
			kat_teta: math.Pi / 4,
		}

		wektory_sily[i] = Wektor_sily{
			dlugosc_x:     100, // wyznacza gracz
			dlugosc_y:     100, // wyznacza gracz
			dlugosc_z:     100, // wyznacza gracz
			przylozenie_x: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Cos(punkty_przylozenia_wektorow[i].kat_fi),
			przylozenie_y: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Sin(punkty_przylozenia_wektorow[i].kat_fi),
			przylozenie_z: pilka.promien * math.Sin(punkty_przylozenia_wektorow[i].kat_teta),
		}

		przyspieszenia_z_wektorow[i] = Wektor{
			x: wektory_sily[i].dlugosc_x / pilka.masa,
			y: wektory_sily[i].dlugosc_y / pilka.masa,
			z: wektory_sily[i].dlugosc_z / pilka.masa,
		}

		momenty_z_wektorow[i] = Wektor{
			x: wektory_sily[i].przylozenie_y*wektory_sily[i].dlugosc_z - wektory_sily[i].przylozenie_z*wektory_sily[i].dlugosc_y,
			y: -wektory_sily[i].przylozenie_x*wektory_sily[i].dlugosc_z + wektory_sily[i].przylozenie_z*wektory_sily[i].dlugosc_x,
			z: wektory_sily[i].przylozenie_x*wektory_sily[i].dlugosc_y - wektory_sily[i].przylozenie_y*wektory_sily[i].dlugosc_x,
		}

		przyspieszenia_katowe_z_wektorow[i] = Wektor{
			x: momenty_z_wektorow[i].x / I,
			y: momenty_z_wektorow[i].y / I,
			z: momenty_z_wektorow[i].z / I,
		}

		wypadkowa_przyspieszen = Wektor{
			x: wypadkowa_przyspieszen.x + przyspieszenia_z_wektorow[i].x,
			y: wypadkowa_przyspieszen.y + przyspieszenia_z_wektorow[i].y,
			z: wypadkowa_przyspieszen.z + przyspieszenia_z_wektorow[i].z,
		}

		wypadkowa_przyspieszen_katowych = Wektor{
			x: wypadkowa_przyspieszen_katowych.x + przyspieszenia_katowe_z_wektorow[i].x,
			y: wypadkowa_przyspieszen_katowych.y + przyspieszenia_katowe_z_wektorow[i].y,
			z: wypadkowa_przyspieszen_katowych.z + przyspieszenia_katowe_z_wektorow[i].z,
		}
	}

	var predkosci_pilki = Wektor{
		x: wypadkowa_przyspieszen.x * czas_dzialania_sily_na_pilke,
		y: wypadkowa_przyspieszen.y * czas_dzialania_sily_na_pilke,
		z: wypadkowa_przyspieszen.z * czas_dzialania_sily_na_pilke,
	}

	var predkosci_katowe_pilki = Wektor{
		x: wypadkowa_przyspieszen_katowych.x * czas_dzialania_sily_na_pilke,
		y: wypadkowa_przyspieszen_katowych.y * czas_dzialania_sily_na_pilke,
		z: wypadkowa_przyspieszen_katowych.z * czas_dzialania_sily_na_pilke,
	}

	rl.InitWindow(int32(szerokosc), int32(wysokosc), "Koszykówka")
	defer rl.CloseWindow()
	// Glowna petla gry
	for !rl.WindowShouldClose() {

		pilka.posX = pilka.posX + predkosci_pilki.x*krok_czasowy
		pilka.posY = pilka.posY + predkosci_pilki.y*krok_czasowy
		pilka.posZ = pilka.posZ + predkosci_pilki.z*krok_czasowy + g*krok_czasowy*krok_czasowy/2
		pilka.rotX = pilka.rotX + predkosci_katowe_pilki.x*krok_czasowy // pamietac skasowac radiany powyzej 2pi
		pilka.rotY = pilka.rotY + predkosci_katowe_pilki.y*krok_czasowy // to co wyzej
		pilka.rotZ = pilka.rotZ + predkosci_katowe_pilki.z*krok_czasowy // to co wyzej
		predkosci_pilki.z = predkosci_pilki.z + g*krok_czasowy

		rl.BeginDrawing()
		rl.ClearBackground(rl.Color{R: 255, G: 79, B: 45, A: 255}) // Czyszczenie tła

		// os wspolrzednych
		rl.DrawCircleLines(int32(szerokosc/50), int32(szerokosc/50), float32(szerokosc/200), rl.White)
		rl.DrawCircle(int32(szerokosc/50), int32(szerokosc/50), float32(szerokosc/400), rl.White)
		rl.DrawLine(int32(szerokosc/40), int32(szerokosc/50), int32(szerokosc/15), int32(szerokosc/50), rl.White)
		rl.DrawLine(int32(szerokosc/50), int32(szerokosc/40), int32(szerokosc/50), int32(szerokosc/15), rl.White)
		wektor_1_osi_x := rl.Vector2{float32(szerokosc / 50), float32(szerokosc / 15)}
		wektor_2_osi_x := rl.Vector2{float32(szerokosc / 40), float32(szerokosc / 20)}
		wektor_3_osi_x := rl.Vector2{float32(szerokosc * 3 / 200), float32(szerokosc / 20)}

		rl.DrawTriangle(wektor_1_osi_x, wektor_2_osi_x, wektor_3_osi_x, rl.White)
		// boisko
		rl.DrawRectangle(int32(szerokosc*15/16), int32(wysokosc/12), int32(szerokosc/320), int32(wysokosc*10/12), rl.White)

		rl.EndDrawing()
	}
}
