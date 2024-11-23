package main

import (
	"fmt"
	"github.com/g3n/engine/math32"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/ungerik/go3d/float64/quaternion"
	"github.com/wcharczuk/go-chart/v2"
	"math"
	"math/rand"
	"os"
	"time"
)

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

func (w Wektor) String() string {
	return fmt.Sprintf("[%v, %v, %v]", w.x, w.y, w.z)
}

func main() {

	rand.NewSource(time.Now().UnixNano())

	const szerokosc int = 1600
	const wysokosc int = 900
	const czas_dzialania_sily_na_pilke = 0.1
	const fps = 200 // fpsy dlatego takie duże żeby uwzglednić jak najmniejnsze zmiany pozycji
	const krok_czasowy = 1.0 / fps
	const g float64 = -9.81
	const stala_wzrostu_malenia = 1
	const grubosc_tablicy = 0.01  // metry
	const szerokosc_tablicy = 1.8 // metry
	const wysokosc_tablicy = 1.05 //metry
	const x_tablicy = 5.0
	const y_tablicy = 3.0
	const z_tablicy = 0.0
	const wspolczynnik_odbicia = 0.8
	const wspolczynnik_odbicia_od_tablicy = 0.6
	const wspolczynnik_odbicia_od_obreczy = 0.5
	const wspolczynnik_momentu = 2.0 / 3.0
	const wspolczynnik_tarcia = 0.5
	const wspolczynnik_tangensa = 10.0
	const srednica_obreczy = 0.5
	const stosunek_promienia_obreczy_do_promienia_przekroju = 0.1
	const max_sily = 50.0

	//100 px to 1 metr przedział

	start_pos_x := 0
	start_pos_y := 1
	start_pos_z := 0
	start_rot_x := 0
	start_rot_y := 0
	start_rot_z := 0

	pilka := Kula{
		promien: 0.12,                 // metry
		posX:    float64(start_pos_x), // metry
		posY:    float64(start_pos_y), // metry
		posZ:    float64(start_pos_z), // metry
		rotX:    float64(start_rot_x), // radiany
		rotY:    float64(start_rot_y), // radiany
		rotZ:    float64(start_rot_z), // radiany
		masa:    0.5,                  // kilogramy
	}

	var punkty_przylozenia_wektorow [3]Punkt_przylozenia
	var wektory_sily [3]Wektor_sily
	var sila_ruszajaca [3]float64
	var wektory_sily_ruszajacej [3]Wektor
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
	var I = wspolczynnik_momentu * pilka.masa * pilka.promien * pilka.promien

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
	for i := 0; i < 3; i++ {

		punkty_przylozenia_wektorow[i] = Punkt_przylozenia{ // wyznaczane przez gracza
			kat_fi:   3 * math.Pi / 2, //radiany
			kat_teta: 0,               // radiany
		}

		wektory_sily[i] = Wektor_sily{
			dlugosc_x:     0, // wyznacza gracz Niutony
			dlugosc_y:     0,
			dlugosc_z:     0,
			przylozenie_x: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Sin(punkty_przylozenia_wektorow[i].kat_fi), //metry
			przylozenie_y: pilka.promien * math.Sin(punkty_przylozenia_wektorow[i].kat_teta),
			przylozenie_z: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Cos(punkty_przylozenia_wektorow[i].kat_fi),
		}

	}

	px := []float64{float64(start_pos_x)}
	py := []float64{float64(start_pos_y)}
	pz := []float64{float64(start_pos_z)}
	ox := []float64{float64(start_rot_x)}
	oy := []float64{float64(start_rot_y)}
	oz := []float64{float64(start_rot_z)}
	vx := []float64{0}
	vy := []float64{0}
	vz := []float64{0}
	wx := []float64{0}
	wy := []float64{0}
	wz := []float64{0}
	czas := []float64{0}

	var liczba_krokow int = 0

	rl.InitWindow(int32(szerokosc), int32(wysokosc), "Koszykówka")

	// Widok kamery na obszar 3D
	camera := rl.Camera{
		Position:   rl.NewVector3(-2.0, 1.0, 0.0),
		Target:     rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	rl.DisableCursor()
	//inicjalizacja obiektow

	//pilka
	sphereMesh := rl.GenMeshSphere(float32(pilka.promien), 32, 32)
	sphereModel := rl.LoadModelFromMesh(sphereMesh)

	texture := rl.LoadTexture("basketball.png")

	materials := sphereModel.GetMaterials()

	rl.SetMaterialTexture(&materials[0], rl.MapDiffuse, texture)

	//punkt przylozenia sily
	punktMesh := rl.GenMeshSphere(float32(pilka.promien/10), 32, 32)
	punktModel := rl.LoadModelFromMesh(punktMesh)

	texture_1 := rl.LoadTexture("bialy.png")

	materials_1 := punktModel.GetMaterials()

	rl.SetMaterialTexture(&materials_1[0], rl.MapDiffuse, texture_1)

	punktPosition := rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ))

	// tablica
	basketMesh := rl.GenMeshCube(grubosc_tablicy, wysokosc_tablicy, szerokosc_tablicy)
	basketModel := rl.LoadModelFromMesh(basketMesh)

	texture_2 := rl.LoadTexture("czarny.png")

	materials_2 := basketModel.GetMaterials()

	rl.SetMaterialTexture(&materials_2[0], rl.MapDiffuse, texture_2)

	basketPosition := rl.NewVector3(x_tablicy, y_tablicy, z_tablicy)

	//obrecz
	hoopMesh := rl.GenMeshTorus(stosunek_promienia_obreczy_do_promienia_przekroju, srednica_obreczy, 16, 32)
	hoopModel := rl.LoadModelFromMesh(hoopMesh)

	texture_3 := rl.LoadTexture("czerwony.jpg")

	materials_3 := hoopModel.GetMaterials()
	rl.SetMaterialTexture(&materials_3[0], rl.MapDiffuse, texture_3)

	hoopPosition := rl.NewVector3(basketPosition.X-0.28, basketPosition.Y-0.3, basketPosition.Z)

	var pozycja_srodka_obreczy = Wektor{
		x: float64(hoopPosition.X),
		y: float64(hoopPosition.Y),
		z: float64(hoopPosition.Z),
	}

	promien_obreczy := srednica_obreczy / 2
	promien_przekroju_obreczy := promien_obreczy * stosunek_promienia_obreczy_do_promienia_przekroju

	rotationAngle := float32(0.0)
	rotationAxis := rl.NewVector3(1.0, 1.0, 0.0)

	rotationAngle_2 := float32(90.0)
	rotationAxis_2 := rl.NewVector3(1.0, 0.0, 0.0)

	//inicjalizacja wartości zmiennych odpowiadających za fazy programu

	licznik := 0
	licz_wartosci := true
	faza_gry := 0
	odbicie := 0
	czy_byla_krawedz := false
	czy_byl_rog := false
	czy_wpadlo := false
	alfa := float32(1)

	rl.SetTargetFPS(fps)

	// Main game loop
	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraFree)

		if faza_gry == 0 {

			if licznik == 15 {
				for j := 0; j < 3; j++ {
					wektory_sily[j].przylozenie_x = pilka.promien * math.Cos(punkty_przylozenia_wektorow[j].kat_teta) * math.Sin(punkty_przylozenia_wektorow[j].kat_fi)
					wektory_sily[j].przylozenie_y = pilka.promien * math.Sin(punkty_przylozenia_wektorow[j].kat_teta)
					wektory_sily[j].przylozenie_z = pilka.promien * math.Cos(punkty_przylozenia_wektorow[j].kat_teta) * math.Cos(punkty_przylozenia_wektorow[j].kat_fi)
				}
				faza_gry += 1

			}
			if licznik == 14 {
				rl.DrawText(fmt.Sprintf("Wartosc 3. wektora sily w kierunku z = %f [N]", wektory_sily[1].dlugosc_z), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[2].dlugosc_z, &licznik, stala_wzrostu_malenia)
				if wektory_sily[2].dlugosc_z > max_sily {
					wektory_sily[2].dlugosc_z = max_sily
				}
				if wektory_sily[2].dlugosc_z < -max_sily {
					wektory_sily[2].dlugosc_z = -max_sily
				}
			}
			if licznik == 13 {
				rl.DrawText(fmt.Sprintf("Wartosc 3. wektora sily w kierunku y = %f [N]", wektory_sily[2].dlugosc_y), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[2].dlugosc_y, &licznik, stala_wzrostu_malenia)
				if wektory_sily[2].dlugosc_y > max_sily {
					wektory_sily[2].dlugosc_y = max_sily
				}
				if wektory_sily[2].dlugosc_y < -max_sily {
					wektory_sily[2].dlugosc_y = -max_sily
				}
			}

			if licznik == 12 {
				rl.DrawText(fmt.Sprintf("Wartosc 3. wektora sily w kierunku x = %f [N]", wektory_sily[1].dlugosc_x), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[2].dlugosc_x, &licznik, stala_wzrostu_malenia)
				if wektory_sily[2].dlugosc_x > max_sily {
					wektory_sily[2].dlugosc_x = max_sily
				}
				if wektory_sily[2].dlugosc_x < -max_sily {
					wektory_sily[2].dlugosc_x = -max_sily
				}
			}
			if licznik == 11 {
				rl.DrawText(fmt.Sprintf("Przylozenie 3. wektora sily kat teta = %f [rad]", punkty_przylozenia_wektorow[2].kat_teta), 10, 10, 20, rl.DarkGray)
				punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[2].kat_teta)*math.Sin(punkty_przylozenia_wektorow[2].kat_fi))
				punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[2].kat_teta))
				punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[2].kat_teta)*math.Cos(punkty_przylozenia_wektorow[2].kat_fi))
				po_wcisnieciu(&punkty_przylozenia_wektorow[2].kat_teta, &licznik, math32.Pi/360)
			}
			if licznik == 10 {
				rl.DrawText(fmt.Sprintf("Przylozenie 3. wektora sily kat fi = %f [rad]", punkty_przylozenia_wektorow[2].kat_fi), 10, 10, 20, rl.DarkGray)
				punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[2].kat_teta)*math.Sin(punkty_przylozenia_wektorow[2].kat_fi))
				punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[2].kat_teta))
				punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[2].kat_teta)*math.Cos(punkty_przylozenia_wektorow[2].kat_fi))
				po_wcisnieciu(&punkty_przylozenia_wektorow[2].kat_fi, &licznik, math32.Pi/360)
			}
			if licznik == 9 {
				rl.DrawText(fmt.Sprintf("Wartosc 2. wektora sily w kierunku z = %f [N]", wektory_sily[1].dlugosc_z), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[1].dlugosc_z, &licznik, stala_wzrostu_malenia)
				if wektory_sily[1].dlugosc_z > max_sily {
					wektory_sily[1].dlugosc_z = max_sily
				}
				if wektory_sily[1].dlugosc_z < -max_sily {
					wektory_sily[1].dlugosc_z = -max_sily
				}
			}
			if licznik == 8 {
				rl.DrawText(fmt.Sprintf("Wartosc 2. wektora sily w kierunku y = %f [N]", wektory_sily[1].dlugosc_y), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[1].dlugosc_y, &licznik, stala_wzrostu_malenia)
				if wektory_sily[1].dlugosc_y > max_sily {
					wektory_sily[1].dlugosc_y = max_sily
				}
				if wektory_sily[1].dlugosc_y < -max_sily {
					wektory_sily[1].dlugosc_y = -max_sily
				}
			}

			if licznik == 7 {
				rl.DrawText(fmt.Sprintf("Wartosc 2. wektora sily w kierunku x = %f [N]", wektory_sily[1].dlugosc_x), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[1].dlugosc_x, &licznik, stala_wzrostu_malenia)
				if wektory_sily[1].dlugosc_x > max_sily {
					wektory_sily[1].dlugosc_x = max_sily
				}
				if wektory_sily[1].dlugosc_x < -max_sily {
					wektory_sily[1].dlugosc_x = -max_sily
				}
			}
			if licznik == 6 {
				rl.DrawText(fmt.Sprintf("Przylozenie 2. wektora sily kat teta = %f [rad]", punkty_przylozenia_wektorow[1].kat_teta), 10, 10, 20, rl.DarkGray)
				punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[1].kat_teta)*math.Sin(punkty_przylozenia_wektorow[1].kat_fi))
				punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[1].kat_teta))
				punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[1].kat_teta)*math.Cos(punkty_przylozenia_wektorow[1].kat_fi))
				po_wcisnieciu(&punkty_przylozenia_wektorow[1].kat_teta, &licznik, math32.Pi/360)
			}
			if licznik == 5 {
				rl.DrawText(fmt.Sprintf("Przylozenie 2. wektora sily kat fi = %f [rad]", punkty_przylozenia_wektorow[1].kat_fi), 10, 10, 20, rl.DarkGray)
				punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[1].kat_teta)*math.Sin(punkty_przylozenia_wektorow[1].kat_fi))
				punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[1].kat_teta))
				punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[1].kat_teta)*math.Cos(punkty_przylozenia_wektorow[1].kat_fi))
				po_wcisnieciu(&punkty_przylozenia_wektorow[1].kat_fi, &licznik, math32.Pi/360)
			}
			if licznik == 4 {
				rl.DrawText(fmt.Sprintf("Wartosc 1. wektora sily w kierunku z = %f [N]", wektory_sily[0].dlugosc_z), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[0].dlugosc_z, &licznik, stala_wzrostu_malenia)
				if wektory_sily[0].dlugosc_z > max_sily {
					wektory_sily[0].dlugosc_z = max_sily
				}
				if wektory_sily[0].dlugosc_z < -max_sily {
					wektory_sily[0].dlugosc_z = -max_sily
				}
			}
			if licznik == 3 {
				rl.DrawText(fmt.Sprintf("Wartosc 1. wektora sily w kierunku y = %f [N]", wektory_sily[0].dlugosc_y), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[0].dlugosc_y, &licznik, stala_wzrostu_malenia)
				if wektory_sily[0].dlugosc_y > max_sily {
					wektory_sily[0].dlugosc_y = max_sily
				}
				if wektory_sily[0].dlugosc_y < -max_sily {
					wektory_sily[0].dlugosc_y = -max_sily
				}
			}

			if licznik == 2 {
				rl.DrawText(fmt.Sprintf("Wartosc 1. wektora sily w kierunku x = %f [N]", wektory_sily[0].dlugosc_x), 10, 10, 20, rl.DarkGray)
				po_wcisnieciu(&wektory_sily[0].dlugosc_x, &licznik, stala_wzrostu_malenia)
				if wektory_sily[0].dlugosc_x > max_sily {
					wektory_sily[0].dlugosc_x = max_sily
				}
				if wektory_sily[0].dlugosc_x < -max_sily {
					wektory_sily[0].dlugosc_x = -max_sily
				}
			}
			if licznik == 1 {
				rl.DrawText(fmt.Sprintf("Przylozenie 1. wektora sily kat teta = %f [rad]", punkty_przylozenia_wektorow[0].kat_teta), 10, 10, 20, rl.DarkGray)
				punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[0].kat_teta)*math.Sin(punkty_przylozenia_wektorow[0].kat_fi))
				punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[0].kat_teta))
				punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[0].kat_teta)*math.Cos(punkty_przylozenia_wektorow[0].kat_fi))
				po_wcisnieciu(&punkty_przylozenia_wektorow[0].kat_teta, &licznik, math32.Pi/360)
			}
			if licznik == 0 {
				rl.DrawText(fmt.Sprintf("Przylozenie 1. wektora sily kat fi = %f [rad]", punkty_przylozenia_wektorow[0].kat_fi), 10, 10, 20, rl.DarkGray)
				punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[0].kat_teta)*math.Sin(punkty_przylozenia_wektorow[0].kat_fi))
				punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[0].kat_teta))
				punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[0].kat_teta)*math.Cos(punkty_przylozenia_wektorow[0].kat_fi))
				po_wcisnieciu(&punkty_przylozenia_wektorow[0].kat_fi, &licznik, math32.Pi/360)
			}
		}

		if faza_gry == 1 {

			if licz_wartosci == true {
				for i := 0; i < 3; i++ {

					sila_ruszajaca[i] = wektory_sily[i].dlugosc_x*wektory_sily[i].przylozenie_x/pilka.promien + wektory_sily[i].dlugosc_y*wektory_sily[i].przylozenie_y/pilka.promien + wektory_sily[i].dlugosc_z*wektory_sily[i].przylozenie_z/pilka.promien

					wektory_sily_ruszajacej[i] = Wektor{
						x: sila_ruszajaca[i] * wektory_sily[i].przylozenie_x / pilka.promien,
						y: sila_ruszajaca[i] * wektory_sily[i].przylozenie_y / pilka.promien,
						z: sila_ruszajaca[i] * wektory_sily[i].przylozenie_z / pilka.promien,
					}
					println("wektor siły ruszjacej")
					println(wektory_sily_ruszajacej[i].String())

					przyspieszenia_z_wektorow[i] = Wektor{
						x: wektory_sily_ruszajacej[i].x / pilka.masa,
						y: wektory_sily_ruszajacej[i].y / pilka.masa,
						z: wektory_sily_ruszajacej[i].z / pilka.masa,
					}
					println("przyspieszenia")
					println(przyspieszenia_z_wektorow[i].String())

					momenty_z_wektorow[i] = Wektor{
						x: wektory_sily[i].przylozenie_y*wektory_sily[i].dlugosc_z - wektory_sily[i].przylozenie_z*wektory_sily[i].dlugosc_y,
						y: wektory_sily[i].przylozenie_z*wektory_sily[i].dlugosc_x - wektory_sily[i].przylozenie_x*wektory_sily[i].dlugosc_z,
						z: wektory_sily[i].przylozenie_x*wektory_sily[i].dlugosc_y - wektory_sily[i].przylozenie_y*wektory_sily[i].dlugosc_x,
					}
					println("momenty wektorow")
					println(momenty_z_wektorow[i].String())

					przyspieszenia_katowe_z_wektorow[i] = Wektor{
						x: momenty_z_wektorow[i].x / I,
						y: momenty_z_wektorow[i].y / I,
						z: momenty_z_wektorow[i].z / I,
					}

					println("przyspieszenia_katowe_z_wektorow")
					println(przyspieszenia_katowe_z_wektorow[i].String())

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
					println("wypadkowa_przyspieszen_katowych")
					println(wypadkowa_przyspieszen_katowych.String())
				}

				predkosci_pilki = Wektor{
					x: wypadkowa_przyspieszen.x * czas_dzialania_sily_na_pilke,
					y: wypadkowa_przyspieszen.y * czas_dzialania_sily_na_pilke,
					z: wypadkowa_przyspieszen.z * czas_dzialania_sily_na_pilke,
				}
				predkosci_pilki = Wektor{
					x: rownaj_do_zera(predkosci_pilki.x),
					y: rownaj_do_zera(predkosci_pilki.y),
					z: rownaj_do_zera(predkosci_pilki.z),
				}
				println("predkosci_pilki")
				println(predkosci_pilki.String())
				predkosci_katowe_pilki = Wektor{
					x: wypadkowa_przyspieszen_katowych.x * czas_dzialania_sily_na_pilke,
					y: wypadkowa_przyspieszen_katowych.y * czas_dzialania_sily_na_pilke,
					z: wypadkowa_przyspieszen_katowych.z * czas_dzialania_sily_na_pilke,
				}
				predkosci_katowe_pilki = Wektor{
					x: rownaj_do_zera(predkosci_katowe_pilki.x),
					y: rownaj_do_zera(predkosci_katowe_pilki.y),
					z: rownaj_do_zera(predkosci_katowe_pilki.z),
				}

				println("predkosci_katowe_pilki")
				println(predkosci_katowe_pilki.String())

				licz_wartosci = false
			}

			//odbicie od ziemii

			if pilka.posY-pilka.promien < 0 {

				var odleglosc_punktu_od_pilki = Wektor{
					x: 0,
					y: -pilka.promien,
					z: 0,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				// zapewnienie wyjscia
				for pilka.posY-pilka.promien < 0 {
					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)
				}
				odbicie += 1

			}

			//powierzchnia tablicy przod

			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posX+pilka.promien > x_tablicy-grubosc_tablicy/2 && pilka.posX < x_tablicy {

				println("Uderzona tablica")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: 0,
					z: 0,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia_od_tablicy, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				for pilka.posY < y_tablicy+wysokosc_tablicy/2 &&
					pilka.posY > y_tablicy-wysokosc_tablicy/2 &&
					pilka.posZ > z_tablicy-szerokosc_tablicy/2 &&
					pilka.posZ < z_tablicy+szerokosc_tablicy/2 &&
					pilka.posX+pilka.promien > x_tablicy-grubosc_tablicy/2 &&
					pilka.posX < x_tablicy {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)
					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}

			}

			// powerzchnia tablicy tyl

			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posX+pilka.promien < x_tablicy+grubosc_tablicy/2 && pilka.posX > x_tablicy {

				println("Uderzona tablica od tylu")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: 0,
					z: 0,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia_od_tablicy, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				for pilka.posY < y_tablicy+wysokosc_tablicy/2 &&
					pilka.posY > y_tablicy-wysokosc_tablicy/2 &&
					pilka.posZ > z_tablicy-szerokosc_tablicy/2 &&
					pilka.posZ < z_tablicy+szerokosc_tablicy/2 &&
					pilka.posX+pilka.promien < x_tablicy+grubosc_tablicy/2 &&
					pilka.posX > x_tablicy {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)
					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}

			}

			//krawedz prawa

			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posZ > z_tablicy+szerokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) && !czy_byl_rog {
				println("krawedz prawa")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: 0,
					z: z_tablicy + szerokosc_tablicy/2 - pilka.posZ,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia_od_tablicy, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z obszaru odbicia
				for pilka.posY < y_tablicy+wysokosc_tablicy/2 &&
					pilka.posY > y_tablicy-wysokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posZ > z_tablicy+szerokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)
					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}
				czy_byla_krawedz = true
			}

			//krawedz lewa

			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posZ < z_tablicy-szerokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy+szerokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) && !czy_byl_rog {
				println("krawedz lewa")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: 0,
					z: z_tablicy - szerokosc_tablicy/2 - pilka.posZ,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia_od_tablicy, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z obszaru odbicia
				for pilka.posY < y_tablicy+wysokosc_tablicy/2 &&
					pilka.posY > y_tablicy-wysokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posZ < z_tablicy-szerokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy+szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)
					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}
				czy_byla_krawedz = true
			}

			//krawedz gora

			if pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posY > y_tablicy+wysokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) && !czy_byl_rog {
				println("krawedz gora")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: y_tablicy + wysokosc_tablicy/2 - pilka.posY,
					z: 0,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia_od_tablicy, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z obszaru odbicia
				for pilka.posZ < z_tablicy+szerokosc_tablicy/2 &&
					pilka.posZ > z_tablicy-szerokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posY > y_tablicy+wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}
				czy_byla_krawedz = true
			}

			//krawedz dol

			if pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posY < y_tablicy-wysokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) && !czy_byl_rog {
				println("krawedz dol")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: y_tablicy - wysokosc_tablicy/2 - pilka.posY,
					z: 0,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia_od_tablicy, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z obszaru odbicia
				for pilka.posZ < z_tablicy+szerokosc_tablicy/2 &&
					pilka.posZ > z_tablicy-szerokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posY < y_tablicy-wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}
				czy_byla_krawedz = true
			}

			// prawy gorny rog
			if pilka.posZ > z_tablicy+szerokosc_tablicy/2 && pilka.posY > y_tablicy+wysokosc_tablicy/2 && math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) && !czy_byla_krawedz {
				println("prawy gorny rog")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: y_tablicy + wysokosc_tablicy/2 - pilka.posY,
					z: z_tablicy + szerokosc_tablicy/2 - pilka.posZ,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z miejsca odbicia
				for pilka.posZ > z_tablicy+szerokosc_tablicy/2 &&
					pilka.posY > y_tablicy+wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}
				czy_byl_rog = true
			}

			// lewy gorny rog
			if pilka.posZ < z_tablicy-szerokosc_tablicy/2 && pilka.posY > y_tablicy+wysokosc_tablicy/2 && math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy+szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) && !czy_byla_krawedz {
				println("lewy gorny rog")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: y_tablicy + wysokosc_tablicy/2 - pilka.posY,
					z: z_tablicy - szerokosc_tablicy/2 - pilka.posZ,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z miejsca odbicia
				for pilka.posZ < z_tablicy-szerokosc_tablicy/2 &&
					pilka.posY > y_tablicy+wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}
				czy_byl_rog = true
			}

			// lewy dolny rog
			if pilka.posZ < z_tablicy-szerokosc_tablicy/2 && pilka.posY < y_tablicy-wysokosc_tablicy/2 && math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy+szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) && !czy_byla_krawedz {
				println("lewy dolny rog")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: y_tablicy - wysokosc_tablicy/2 - pilka.posY,
					z: z_tablicy - szerokosc_tablicy/2 - pilka.posZ,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z miejsca odbicia
				for pilka.posZ < z_tablicy-szerokosc_tablicy/2 &&
					pilka.posY < y_tablicy-wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy+szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}
				czy_byl_rog = true
			}

			// prawy dolny rog
			if pilka.posZ > z_tablicy+szerokosc_tablicy/2 && pilka.posY < y_tablicy-wysokosc_tablicy/2 && math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) && !czy_byla_krawedz {
				println("prawy dolny rog")

				var odleglosc_punktu_od_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: y_tablicy - wysokosc_tablicy/2 - pilka.posY,
					z: z_tablicy + szerokosc_tablicy/2 - pilka.posZ,
				}

				odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia, wspolczynnik_tarcia, wspolczynnik_tangensa, I, odleglosc_punktu_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

				//zapewnienie wyjścia piłki z miejsca odbicia
				for pilka.posZ > z_tablicy+szerokosc_tablicy/2 &&
					pilka.posY < y_tablicy-wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

					zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

					rl.DrawModelEx(
						sphereModel,
						rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
						rotationAxis,                 // Rotation axis
						rotationAngle,                // Rotation angle
						rl.NewVector3(1.0, 1.0, 1.0), // Scale
						rl.White,                     // Tint
					)
				}

				czy_byl_rog = true
			}

			//obrecz

			if math.Pow(pilka.posX-pozycja_srodka_obreczy.x, 2)+math.Pow(pilka.posZ-pozycja_srodka_obreczy.z, 2) < math.Pow(promien_obreczy+promien_przekroju_obreczy+pilka.promien, 2) && math.Pow(pilka.posX-pozycja_srodka_obreczy.x, 2)+math.Pow(pilka.posZ-pozycja_srodka_obreczy.z, 2) > math.Pow(promien_obreczy-promien_przekroju_obreczy-pilka.promien, 2) && pilka.posY < pozycja_srodka_obreczy.y+promien_przekroju_obreczy+pilka.promien && pilka.posY > pozycja_srodka_obreczy.y-promien_przekroju_obreczy-pilka.promien {

				// wektor z ktorego wyznaczamy kierunek normlanej na ktorej lezy punkt w srodku obreczy najblizszy srodkowi pilki

				var wektor_odleglosci_pilki_od_centrum_obreczy = Wektor{
					x: pilka.posX - pozycja_srodka_obreczy.x,
					y: 0,
					z: pilka.posZ - pozycja_srodka_obreczy.z,
				}

				dlugosc_wektora := math.Sqrt(math.Pow(wektor_odleglosci_pilki_od_centrum_obreczy.x, 2) + math.Pow(wektor_odleglosci_pilki_od_centrum_obreczy.z, 2))

				var normalna = Wektor{
					x: wektor_odleglosci_pilki_od_centrum_obreczy.x / dlugosc_wektora,
					y: 0,
					z: wektor_odleglosci_pilki_od_centrum_obreczy.z / dlugosc_wektora,
				}

				var punkt_w_obreczy_najblizej_pilki = Wektor{
					x: pozycja_srodka_obreczy.x + normalna.x*promien_obreczy,
					y: pozycja_srodka_obreczy.y,
					z: pozycja_srodka_obreczy.z + normalna.z*promien_obreczy,
				}

				var wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki = Wektor{
					x: punkt_w_obreczy_najblizej_pilki.x - pilka.posX,
					y: punkt_w_obreczy_najblizej_pilki.y - pilka.posY,
					z: punkt_w_obreczy_najblizej_pilki.z - pilka.posZ,
				}

				odleglosc_pilki_od_punktu_w_obreczy := math.Sqrt(math.Pow(wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki.x, 2) + math.Pow(wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki.y, 2) + math.Pow(wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki.z, 2))

				if odleglosc_pilki_od_punktu_w_obreczy < pilka.promien+promien_przekroju_obreczy {
					println("obrecz")

					odbij(pilka.promien, pilka.masa, krok_czasowy, wspolczynnik_odbicia_od_obreczy, wspolczynnik_tarcia, wspolczynnik_tangensa, I, wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki, &predkosci_pilki, &predkosci_katowe_pilki)

					//zapewnienie wyjścia piłki z miejsca odbicia
					for odleglosc_pilki_od_punktu_w_obreczy < pilka.promien+promien_przekroju_obreczy {

						wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki = Wektor{
							x: punkt_w_obreczy_najblizej_pilki.x - pilka.posX,
							y: punkt_w_obreczy_najblizej_pilki.y - pilka.posY,
							z: punkt_w_obreczy_najblizej_pilki.z - pilka.posZ,
						}

						odleglosc_pilki_od_punktu_w_obreczy = math.Sqrt(math.Pow(wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki.x, 2) + math.Pow(wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki.y, 2) + math.Pow(wektor_odleglosci_punktu_w_srodku_obreczy_od_pilki.z, 2))

						zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

						rl.DrawModelEx(
							sphereModel,
							rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
							rotationAxis,                 // Rotation axis
							rotationAngle,                // Rotation angle
							rl.NewVector3(1.0, 1.0, 1.0), // Scale
							rl.White,                     // Tint
						)
					}
				}
			}

			//wykrycie trafienia
			if predkosci_pilki.y < 0 && pilka.posY < pozycja_srodka_obreczy.y && pilka.posY > pozycja_srodka_obreczy.y-pilka.promien && math.Pow(pilka.posX-pozycja_srodka_obreczy.x, 2)+math.Pow(pilka.posZ-pozycja_srodka_obreczy.z, 2) < math.Pow(promien_obreczy-pilka.promien, 2) {
				czy_wpadlo = true
			}

			if czy_wpadlo == true {
				rl.DrawRectangle(0, 0, int32(szerokosc), int32(wysokosc), rl.Fade(rl.Green, alfa))
				alfa -= 2 * krok_czasowy
				if alfa < 0 {
					alfa = 1
					czy_wpadlo = false
				}
			}

			// obliczanie w powietrzu
			zmiana_parametrow_w_czasie(&pilka, &predkosci_pilki, &predkosci_katowe_pilki, g, krok_czasowy, &liczba_krokow, &px, &py, &pz, &vx, &vy, &vz, &ox, &oy, &oz, &wx, &wy, &wz, &czas)

			// powrot do rzutu gdy pilka sie 5 razy odbije od ziemii albo za daleko poleci
			if odbicie == 5 || pilka.posX > 10 || pilka.posX < -10 || pilka.posZ < -10 || pilka.posZ > 10 {
				pilka.posX = float64(start_pos_x)
				pilka.posY = float64(start_pos_y)
				pilka.posZ = float64(start_pos_z)
				pilka.rotX = float64(start_rot_x)
				pilka.rotY = float64(start_rot_y)
				pilka.rotZ = float64(start_rot_z)
				licznik = 0
				licz_wartosci = true
				faza_gry = 0
				odbicie = 0
				for i := 0; i < 3; i++ {

					sila_ruszajaca[i] = 0
					wektory_sily_ruszajacej[i] = zeruj_wektor(wektory_sily_ruszajacej[i])
					przyspieszenia_z_wektorow[i] = zeruj_wektor(przyspieszenia_z_wektorow[i])
					momenty_z_wektorow[i] = zeruj_wektor(momenty_z_wektorow[i])
					przyspieszenia_katowe_z_wektorow[i] = zeruj_wektor(przyspieszenia_katowe_z_wektorow[i])

				}

				wypadkowa_przyspieszen = zeruj_wektor(wypadkowa_przyspieszen)
				wypadkowa_przyspieszen_katowych = zeruj_wektor(wypadkowa_przyspieszen_katowych)
				predkosci_pilki = zeruj_wektor(predkosci_pilki)
				predkosci_katowe_pilki = zeruj_wektor(predkosci_katowe_pilki)

			}
			qX := quaternion.FromXAxisAngle(pilka.rotX)
			qY := quaternion.FromYAxisAngle(pilka.rotY)
			qZ := quaternion.FromZAxisAngle(pilka.rotZ)

			kwaternion_glowny := quaternion.Mul3(&qX, &qY, &qZ)

			oska, kat := kwaternion_glowny.AxisAngle()

			rotationAxis = rl.NewVector3(float32(oska[0]), float32(oska[1]), float32(oska[2]))
			rotationAngle = float32(kat * 180 / math.Pi)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		rl.DrawModelEx(
			sphereModel,
			rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
			rotationAxis,                 // Rotation axis
			rotationAngle,                // Rotation angle
			rl.NewVector3(1.0, 1.0, 1.0), // Scale
			rl.White,                     // Tint
		)

		if faza_gry == 0 {
			if licznik < 5 {
				rl.DrawModel(
					punktModel,
					punktPosition,
					1,
					rl.DarkBlue,
				)
			} else if licznik < 10 {
				rl.DrawModel(
					punktModel,
					punktPosition,
					1,
					rl.DarkGreen,
				)
			} else {
				rl.DrawModel(
					punktModel,
					punktPosition,
					1,
					rl.Red,
				)
			}
		}
		rl.DrawModel(
			basketModel,
			basketPosition, // Position
			1,
			rl.White, // Tint
		)

		rl.DrawModelEx(
			hoopModel,
			hoopPosition,
			rotationAxis_2,
			rotationAngle_2,
			rl.NewVector3(1.0, 1.0, 1.0),
			rl.White,
		)

		rl.DrawGrid(100, 1.0)

		rl.EndMode3D()

		rl.EndDrawing()
	}
	// De-initialization

	rl.UnloadTexture(texture)
	rl.UnloadModel(sphereModel)

	rl.CloseWindow()

	zapisz_obraz(czas, px, "pozycja_x.png", "Czas [s]", "Pozycja wzdluz osi X [m]")
	zapisz_obraz(czas, py, "pozycja_y.png", "Czas [s]", "Pozycja wzdluz osi Y [m]")
	zapisz_obraz(czas, pz, "pozycja_z.png", "Czas [s]", "Pozycja wzdluz osi Z [m]")
	zapisz_obraz(czas, ox, "obrot_x.png", "Czas [s]", "Obrot wokol osi X [rad]")
	zapisz_obraz(czas, oy, "obrot_y.png", "Czas [s]", "Obrot wokol osi Y [rad]")
	zapisz_obraz(czas, oz, "obrot_z.png", "Czas [s]", "Obrot wokol osi Z [rad]")
	zapisz_obraz(czas, vx, "predkosc_liniowa_x.png", "Czas [s]", "Predkosc wzdluz osi X [m/s]")
	zapisz_obraz(czas, vy, "predkosc_liniowa_y.png", "Czas [s]", "Predkosc wzdluz osi Y [m/s]")
	zapisz_obraz(czas, vz, "predkosc_liniowa_z.png", "Czas [s]", "Predkosc wzdluz osi Z [m/s]")
	zapisz_obraz(czas, wx, "predkosc_katowa_x.png", "Czas [s]", "Predkosc obrotowa wokol osi X [rad/s]")
	zapisz_obraz(czas, wy, "predkosc_katowa_y.png", "Czas [s]", "Predkosc obrotowa wokol osi Y [rad/s]")
	zapisz_obraz(czas, wz, "predkosc_katowa_z.png", "Czas [s]", "Predkosc obrotowa wokol osi Z [rad/s]")
}

func po_wcisnieciu(zmienna *float64, licznik *int, zmiana float64) {
	if rl.IsKeyPressed(rl.KeyI) {
		*zmienna += 10 * zmiana
	}
	if rl.IsKeyPressed(rl.KeyU) {
		*zmienna += zmiana
	}
	if rl.IsKeyPressed(rl.KeyY) {
		*zmienna -= zmiana
	}
	if rl.IsKeyPressed(rl.KeyT) {
		*zmienna -= 10 * zmiana
	}
	if rl.IsKeyPressed(rl.KeyEnter) {
		*licznik++
	}
}
func iloczyn_wektorowy(wektor_1 Wektor, wektor_2 Wektor) Wektor {
	return Wektor{
		x: wektor_1.y*wektor_2.z - wektor_1.z*wektor_2.y,
		y: wektor_1.z*wektor_2.x - wektor_1.x*wektor_2.z,
		z: wektor_1.x*wektor_2.y - wektor_1.y*wektor_2.x,
	}
}
func zapisz_obraz(tab_x []float64, tab_y []float64, nazwa string, nazwa_osi_x string, nazwa_osi_y string) {

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name: nazwa_osi_x,
		},
		YAxis: chart.YAxis{
			Name: nazwa_osi_y,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: tab_x,
				YValues: tab_y,
			},
		},
	}

	dlugosc := len(tab_y)
	czy_same_zera := true

	for i := 0; i < dlugosc; i++ {
		if tab_y[i] != 0 {
			czy_same_zera = false
		}
	}
	if czy_same_zera {

		graph = chart.Chart{
			XAxis: chart.XAxis{
				Name: nazwa_osi_x,
			},
			YAxis: chart.YAxis{
				Name: nazwa_osi_y,

				Range: &chart.ContinuousRange{
					Min: -1,
					Max: 1,
				},
			},
			Series: []chart.Series{
				chart.ContinuousSeries{
					XValues: tab_x,
					YValues: tab_y,
				},
			},
		}
	}
	// Tworzenie pliku
	f, err := os.Create(nazwa)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	// Renderowanie wykresu do pliku PNG
	err = graph.Render(chart.PNG, f)
	if err != nil {
		panic(err)
	}
}
func rownaj_do_zera(zmienna float64) float64 {
	if zmienna < 0.000001 && zmienna > -0.000001 {
		zmienna = 0
	}
	return zmienna
}
func zeruj_wektor(wektor Wektor) Wektor {
	wektor = Wektor{
		x: 0,
		y: 0,
		z: 0,
	}
	return wektor
}
func odbij(promien float64, masa float64, krok float64, wspolczynnik_odbicia float64, wspolczynnik_tarcia float64, wspolczynnik_tanh float64, I float64, odleglosc_punktu_od_pilki Wektor, predkosci_pilki *Wektor, predkosci_katowe_pilki *Wektor) {

	var wersor_normalnej = Wektor{
		x: odleglosc_punktu_od_pilki.x / promien,
		y: odleglosc_punktu_od_pilki.y / promien,
		z: odleglosc_punktu_od_pilki.z / promien,
	}

	var iloczyn_skalarny_predkosci_z_wersorem = predkosci_pilki.x*wersor_normalnej.x + predkosci_pilki.y*wersor_normalnej.y + predkosci_pilki.z*wersor_normalnej.z
	var iloczyn_skalarny_wersora_z_wersorem = wersor_normalnej.x*wersor_normalnej.x + wersor_normalnej.y*wersor_normalnej.y + wersor_normalnej.z*wersor_normalnej.z

	var predkosc_normalna = Wektor{
		x: iloczyn_skalarny_predkosci_z_wersorem / iloczyn_skalarny_wersora_z_wersorem * wersor_normalnej.x,
		y: iloczyn_skalarny_predkosci_z_wersorem / iloczyn_skalarny_wersora_z_wersorem * wersor_normalnej.y,
		z: iloczyn_skalarny_predkosci_z_wersorem / iloczyn_skalarny_wersora_z_wersorem * wersor_normalnej.z,
	}

	var wartosc_predkosci_normalnej = math.Sqrt(math.Pow(predkosc_normalna.x, 2) + math.Pow(predkosc_normalna.y, 2) + math.Pow(predkosc_normalna.z, 2))

	var omega_wektorowo_z_promieniem = iloczyn_wektorowy(*predkosci_katowe_pilki, odleglosc_punktu_od_pilki)

	var predkosc_styczna = Wektor{
		x: omega_wektorowo_z_promieniem.x + predkosci_pilki.x - predkosc_normalna.x,
		y: omega_wektorowo_z_promieniem.y + predkosci_pilki.y - predkosc_normalna.y,
		z: omega_wektorowo_z_promieniem.z + predkosci_pilki.z - predkosc_normalna.z,
	}

	var wartosc_predkosci_stycznej = math.Sqrt(math.Pow(predkosc_styczna.x, 2) + math.Pow(predkosc_styczna.y, 2) + math.Pow(predkosc_styczna.z, 2))

	if wartosc_predkosci_stycznej != 0 {
		var wersor_stycznej = Wektor{
			x: predkosc_styczna.x / wartosc_predkosci_stycznej,
			y: predkosc_styczna.y / wartosc_predkosci_stycznej,
			z: predkosc_styczna.z / wartosc_predkosci_stycznej,
		}

		var wartosc_sily_nacisku = masa * wartosc_predkosci_normalnej * (1 + wspolczynnik_odbicia) / krok

		var sila_tarcia = Wektor{
			x: -wspolczynnik_tarcia * wartosc_sily_nacisku * math.Tanh(wartosc_predkosci_stycznej/wspolczynnik_tanh) * wersor_stycznej.x,
			y: -wspolczynnik_tarcia * wartosc_sily_nacisku * math.Tanh(wartosc_predkosci_stycznej/wspolczynnik_tanh) * wersor_stycznej.y,
			z: -wspolczynnik_tarcia * wartosc_sily_nacisku * math.Tanh(wartosc_predkosci_stycznej/wspolczynnik_tanh) * wersor_stycznej.z,
		}

		var moment_sily = iloczyn_wektorowy(odleglosc_punktu_od_pilki, sila_tarcia)

		var przyspieszenie_katowe = Wektor{
			x: moment_sily.x / I,
			y: moment_sily.y / I,
			z: moment_sily.z / I,
		}

		nowe_predkosci_pilki := Wektor{
			x: predkosci_pilki.x - (1+wspolczynnik_odbicia)*predkosc_normalna.x + sila_tarcia.x*krok/masa,
			y: predkosci_pilki.y - (1+wspolczynnik_odbicia)*predkosc_normalna.y + sila_tarcia.y*krok/masa,
			z: predkosci_pilki.z - (1+wspolczynnik_odbicia)*predkosc_normalna.z + sila_tarcia.z*krok/masa,
		}

		nowe_predkosci_katowe_pilki := Wektor{
			x: predkosci_katowe_pilki.x + przyspieszenie_katowe.x*krok,
			y: predkosci_katowe_pilki.y + przyspieszenie_katowe.y*krok,
			z: predkosci_katowe_pilki.z + przyspieszenie_katowe.z*krok,
		}

		*predkosci_pilki = nowe_predkosci_pilki
		*predkosci_katowe_pilki = nowe_predkosci_katowe_pilki

	}
	if wartosc_predkosci_stycznej == 0 {
		nowe_predkosci_pilki := Wektor{
			x: predkosci_pilki.x - (1+wspolczynnik_odbicia)*predkosc_normalna.x,
			y: predkosci_pilki.y - (1+wspolczynnik_odbicia)*predkosc_normalna.y,
			z: predkosci_pilki.z - (1+wspolczynnik_odbicia)*predkosc_normalna.z,
		}

		*predkosci_pilki = nowe_predkosci_pilki
		// a predkosci obrotowe sie nie zmieniaja
	}
}
func zmiana_parametrow_w_czasie(pilka *Kula, predkosci_pilki *Wektor, predkosci_katowe_pilki *Wektor, g float64, krok_czasowy float64, liczba_krokow *int, px *[]float64, py *[]float64, pz *[]float64, vx *[]float64, vy *[]float64, vz *[]float64, ox *[]float64, oy *[]float64, oz *[]float64, wx *[]float64, wy *[]float64, wz *[]float64, czas *[]float64) {

	pilka.posX = pilka.posX + predkosci_pilki.x*krok_czasowy
	pilka.posY = pilka.posY + predkosci_pilki.y*krok_czasowy + g*krok_czasowy*krok_czasowy/2
	predkosci_pilki.y = predkosci_pilki.y + g*krok_czasowy
	pilka.posZ = pilka.posZ + predkosci_pilki.z*krok_czasowy

	pilka.rotX = pilka.rotX + predkosci_katowe_pilki.x*krok_czasowy
	pilka.rotY = pilka.rotY + predkosci_katowe_pilki.y*krok_czasowy
	pilka.rotZ = pilka.rotZ + predkosci_katowe_pilki.z*krok_czasowy

	*liczba_krokow += 1
	*px = append(*px, pilka.posX)
	*py = append(*py, pilka.posY)
	*pz = append(*pz, pilka.posZ)
	*ox = append(*ox, pilka.rotX)
	*oy = append(*oy, pilka.rotY)
	*oz = append(*oz, pilka.rotZ)
	*vx = append(*vx, predkosci_pilki.x)
	*vy = append(*vy, predkosci_pilki.y)
	*vz = append(*vz, predkosci_pilki.z)
	*wx = append(*wx, predkosci_katowe_pilki.x)
	*wy = append(*wy, predkosci_katowe_pilki.y)
	*wz = append(*wz, predkosci_katowe_pilki.z)
	*czas = append(*czas, float64(*liczba_krokow)*krok_czasowy)

}
