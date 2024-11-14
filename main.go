package main

import (
	"fmt"
	"github.com/g3n/engine/math32"
	"github.com/gen2brain/raylib-go/raylib"
	"github.com/ungerik/go3d/float64/quaternion"
	"math"
	"math/rand"
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
	const fps = 100
	const krok_czasowy = 1.0 / fps
	const g float64 = -9.81
	const wspolczynnik_odbicia = 0.7
	const stala_wzrostu_malenia = 1
	const grubosc_tablicy = 0.01  // metry
	const szerokosc_tablicy = 1.8 // metry
	const wysokosc_tablicy = 1.05 //metry
	const x_tablicy = 5
	const y_tablicy = 3
	const z_tablicy = 0
	const wspolczynnik_momentu = 2.0 / 3.0

	//100 px to 1 metr przedział

	start_pos_x := 0
	start_pos_y := 1
	start_pos_z := 5
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

	for i := 0; i < 1; i++ { //todo 1 na 3

		punkty_przylozenia_wektorow[i] = Punkt_przylozenia{ // wyznaczane przez gracza
			kat_fi:   3 * math.Pi / 2, //radiany
			kat_teta: 0,               // radiany
		}

		wektory_sily[i] = Wektor_sily{
			dlugosc_x:     25,                                                                                                                  // wyznacza gracz Niutony
			dlugosc_y:     40,                                                                                                                  // wyznacza gracz
			dlugosc_z:     -25,                                                                                                                 // wyznacza gracz
			przylozenie_x: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Sin(punkty_przylozenia_wektorow[i].kat_fi), //metry
			przylozenie_y: pilka.promien * math.Sin(punkty_przylozenia_wektorow[i].kat_teta),
			przylozenie_z: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Cos(punkty_przylozenia_wektorow[i].kat_fi),
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

	// Widok kamery na obszar 3D
	camera := rl.Camera{
		Position:   rl.NewVector3(1.0, 1.0, 0.0), // pozycja kamery
		Target:     rl.NewVector3(0, 0, 0),       // punkt patrzenia
		Up:         rl.NewVector3(0.0, 1.0, 0.0), // Camera up vector (rotation towards target)
		Fovy:       45.0,                         // Camera field-of-view Y
		Projection: rl.CameraPerspective,         // Camera mode type
	}

	rl.DisableCursor()
	//pilka
	sphereMesh := rl.GenMeshSphere(float32(pilka.promien), 32, 32)
	sphereModel := rl.LoadModelFromMesh(sphereMesh)

	texture := rl.LoadTexture("basketball.png") // Ensure this file exists

	materials := sphereModel.GetMaterials()

	rl.SetMaterialTexture(&materials[0], rl.MapDiffuse, texture)
	//punkt przylozenia sily
	punktMesh := rl.GenMeshSphere(float32(pilka.promien/10), 32, 32)
	punktModel := rl.LoadModelFromMesh(punktMesh)

	texture_1 := rl.LoadTexture("niebieski.png") // Ensure this file exists

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
	hoopMesh := rl.GenMeshTorus(0.1, 0.5, 16, 32)
	hoopModel := rl.LoadModelFromMesh(hoopMesh)

	texture_3 := rl.LoadTexture("czerwony.jpg")

	materials_3 := hoopModel.GetMaterials()
	rl.SetMaterialTexture(&materials_3[0], rl.MapDiffuse, texture_3)

	hoopPosition := rl.NewVector3(basketPosition.X-0.28, basketPosition.Y-0.3, basketPosition.Z)

	rotationAngle := float32(0.0)
	rotationAxis := rl.NewVector3(1.0, 1.0, 0.0)

	rotationAngle_2 := float32(90.0)
	rotationAxis_2 := rl.NewVector3(1.0, 0.0, 0.0)

	licznik := 0
	licz_wartosci := true
	faza_gry := 0
	odbicie := 0

	rl.SetTargetFPS(fps)

	// Main game loop
	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraFree)

		if rl.IsKeyPressed(rl.KeyZ) {
			camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
		}
		if faza_gry == 0 {
			for i := 0; i < 1; i++ { // todo 1 na 3
				if licznik == 5 {
					for j := 0; j < 1; j++ {
						wektory_sily[j].przylozenie_x = pilka.promien * math.Cos(punkty_przylozenia_wektorow[j].kat_teta) * math.Sin(punkty_przylozenia_wektorow[j].kat_fi)
						wektory_sily[j].przylozenie_y = pilka.promien * math.Sin(punkty_przylozenia_wektorow[j].kat_teta)
						wektory_sily[j].przylozenie_z = pilka.promien * math.Cos(punkty_przylozenia_wektorow[j].kat_teta) * math.Cos(punkty_przylozenia_wektorow[j].kat_fi)
					}
					faza_gry += 1

				}

				if licznik == 4 {
					rl.DrawText(fmt.Sprintf("Przylozenie sily kat teta = %f", punkty_przylozenia_wektorow[i].kat_teta), 10, 10, 20, rl.DarkGray)
					punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[i].kat_teta)*math.Sin(punkty_przylozenia_wektorow[i].kat_fi))
					punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[i].kat_teta))
					punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[i].kat_teta)*math.Cos(punkty_przylozenia_wektorow[i].kat_fi))
					if rl.IsKeyPressed(rl.KeyU) {
						punkty_przylozenia_wektorow[i].kat_teta += 2 * math32.Pi / 360
					}
					if rl.IsKeyPressed(rl.KeyY) {
						punkty_przylozenia_wektorow[i].kat_teta -= 2 * math32.Pi / 360
					}
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}
				if licznik == 3 {
					rl.DrawText(fmt.Sprintf("Przylozenie sily kat fi = %f", punkty_przylozenia_wektorow[i].kat_fi), 10, 10, 20, rl.DarkGray)
					punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[i].kat_teta)*math.Sin(punkty_przylozenia_wektorow[i].kat_fi))
					punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[i].kat_teta))
					punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[i].kat_teta)*math.Cos(punkty_przylozenia_wektorow[i].kat_fi))
					if rl.IsKeyPressed(rl.KeyU) {
						punkty_przylozenia_wektorow[i].kat_fi += 2 * math32.Pi / 360
					}
					if rl.IsKeyPressed(rl.KeyY) {
						punkty_przylozenia_wektorow[i].kat_fi -= 2 * math32.Pi / 360
					}
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}
				if licznik == 2 {
					rl.DrawText(fmt.Sprintf("Wartosc wektora w kierunku z= %f", wektory_sily[i].dlugosc_z), 10, 10, 20, rl.DarkGray)
					if rl.IsKeyPressed(rl.KeyU) {
						wektory_sily[i].dlugosc_z += stala_wzrostu_malenia
					}
					if rl.IsKeyPressed(rl.KeyY) {
						wektory_sily[i].dlugosc_z -= stala_wzrostu_malenia
					}
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}
				if licznik == 1 {
					rl.DrawText(fmt.Sprintf("Wartosc wektora w kierunku y= %f", wektory_sily[i].dlugosc_y), 10, 10, 20, rl.DarkGray)
					if rl.IsKeyPressed(rl.KeyU) {
						wektory_sily[i].dlugosc_y += stala_wzrostu_malenia
					}
					if rl.IsKeyPressed(rl.KeyY) {
						wektory_sily[i].dlugosc_y -= stala_wzrostu_malenia
					}
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}

				if licznik == 0 {
					rl.DrawText(fmt.Sprintf("Wartosc wektora w kierunku x= %f", wektory_sily[i].dlugosc_x), 10, 10, 20, rl.DarkGray)
					if rl.IsKeyPressed(rl.KeyU) {
						wektory_sily[i].dlugosc_x += stala_wzrostu_malenia
					}
					if rl.IsKeyPressed(rl.KeyY) {
						wektory_sily[i].dlugosc_x -= stala_wzrostu_malenia
					}
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}

			}
		}

		if faza_gry == 1 {

			if licz_wartosci == true {
				for i := 0; i < 1; i++ { // todo: zmiana 1 na 3

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
				println("predkosci_pilki")
				println(predkosci_pilki.String())
				predkosci_katowe_pilki = Wektor{
					x: wypadkowa_przyspieszen_katowych.x * czas_dzialania_sily_na_pilke,
					y: wypadkowa_przyspieszen_katowych.y * czas_dzialania_sily_na_pilke,
					z: wypadkowa_przyspieszen_katowych.z * czas_dzialania_sily_na_pilke,
				}
				println("predkosci_katowe_pilki")
				println(predkosci_katowe_pilki.String())

				licz_wartosci = false
			}

			//odbicie od ziemii
			if pilka.posY-pilka.promien < 0 {
				pilka.posY = pilka.promien
				bufor_x := predkosci_pilki.x
				bufor_z := predkosci_pilki.z
				predkosci_pilki.y = -predkosci_pilki.y * wspolczynnik_odbicia
				predkosci_pilki.x = wspolczynnik_odbicia * (bufor_x - wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.z)
				predkosci_pilki.z = wspolczynnik_odbicia * (bufor_z + wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.x)
				predkosci_katowe_pilki.x = wspolczynnik_odbicia * (predkosci_katowe_pilki.x + wspolczynnik_momentu*bufor_z/pilka.promien)
				predkosci_katowe_pilki.z = wspolczynnik_odbicia * (predkosci_katowe_pilki.z - wspolczynnik_momentu*bufor_x/pilka.promien)

				odbicie += 1

			}

			//odbicia od tablicy todo sprawdzac wartosc pozycji o klatke wczesniej!
			//powierzchnia

			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posX+pilka.promien > x_tablicy-grubosc_tablicy/2 && pilka.posX+pilka.promien < x_tablicy+grubosc_tablicy/2 {
				pilka.posX = x_tablicy - grubosc_tablicy/2 - pilka.promien
				bufor_y := predkosci_pilki.y
				bufor_z := predkosci_pilki.z
				predkosci_pilki.x = -predkosci_pilki.x * wspolczynnik_odbicia
				predkosci_pilki.y = wspolczynnik_odbicia * (bufor_y - wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.z)
				predkosci_pilki.z = wspolczynnik_odbicia * (bufor_z + wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.y)
				predkosci_katowe_pilki.y = wspolczynnik_odbicia * (predkosci_katowe_pilki.y + wspolczynnik_momentu*bufor_z/pilka.promien)
				predkosci_katowe_pilki.z = wspolczynnik_odbicia * (predkosci_katowe_pilki.z - wspolczynnik_momentu*bufor_y/pilka.promien)

				println("Uderzona tablica")
			}

			//krawedz prawa

			if pilka.posY <= y_tablicy+wysokosc_tablicy/2 && pilka.posY >= y_tablicy-wysokosc_tablicy/2 && pilka.posX >= x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX <= x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posZ >= z_tablicy+szerokosc_tablicy/2+pilka.promien && ((pilka.posX-x_tablicy+grubosc_tablicy/2)*(pilka.posX-x_tablicy+grubosc_tablicy/2)+(pilka.posZ-z_tablicy-szerokosc_tablicy/2)*(pilka.posZ-z_tablicy-szerokosc_tablicy/2)) < pilka.promien*pilka.promien {
				println("krawedz")
				var predkosci_pocz = Wektor{
					x: predkosci_pilki.x,
					y: predkosci_pilki.y,
					z: predkosci_pilki.z,
				}

				//wektor jednostkowy normalnej

				var normalna = Wektor{
					x: pilka.posX - x_tablicy + grubosc_tablicy/2,
					y: 0.0,
					z: pilka.posZ - z_tablicy - szerokosc_tablicy/2,
				}
				normalna.x = normalna.x / math.Sqrt(normalna.x*normalna.x+normalna.z*normalna.z)
				normalna.z = normalna.z / math.Sqrt(normalna.x*normalna.x+normalna.z*normalna.z)

				//wektor od środka piłki do punktu styku z krawędzią
				var r_przez_I = Wektor{
					x: -pilka.promien * normalna.x / I,
					y: -pilka.promien * normalna.y / I,
					z: -pilka.promien * normalna.z / I,
				}

				var delta_liniowa = Wektor{
					x: -(1 + wspolczynnik_odbicia) * (predkosci_pocz.x*normalna.x + predkosci_pocz.y*normalna.y + predkosci_pocz.z*normalna.z) * normalna.x,
					y: -(1 + wspolczynnik_odbicia) * (predkosci_pocz.x*normalna.x + predkosci_pocz.y*normalna.y + predkosci_pocz.z*normalna.z) * normalna.y,
					z: -(1 + wspolczynnik_odbicia) * (predkosci_pocz.x*normalna.x + predkosci_pocz.y*normalna.y + predkosci_pocz.z*normalna.z) * normalna.z,
				}

				var delta_katowa = iloczyn_wektorowy(r_przez_I, delta_liniowa)

				//predkosci liniowe po odbiciu
				predkosci_pilki.x = predkosci_pocz.x + delta_liniowa.x
				predkosci_pilki.y = predkosci_pocz.y + delta_liniowa.y
				predkosci_pilki.z = predkosci_pocz.z + delta_liniowa.z

				//predkosci katowe po odbiciu
				predkosci_katowe_pilki.x = predkosci_katowe_pilki.x + delta_katowa.x
				predkosci_katowe_pilki.y = predkosci_katowe_pilki.y + delta_katowa.y
				predkosci_katowe_pilki.z = predkosci_katowe_pilki.z + delta_katowa.z
				//zapewnienie wyjścia piłki z obszaru odbicia

				pilka.posX = pilka.posX + predkosci_pilki.x*krok_czasowy
				pilka.posY = pilka.posY + predkosci_pilki.y*krok_czasowy + g*krok_czasowy*krok_czasowy/2
				predkosci_pilki.y = predkosci_pilki.y + g*krok_czasowy
				pilka.posZ = pilka.posZ + predkosci_pilki.z*krok_czasowy

				pilka.rotX = pilka.rotX + predkosci_katowe_pilki.x*krok_czasowy
				pilka.rotY = pilka.rotY + predkosci_katowe_pilki.y*krok_czasowy
				pilka.rotZ = pilka.rotZ + predkosci_katowe_pilki.z*krok_czasowy
			}

			// powrot do rzutu gdy pilka sie 3 razy odbije lub za wysoko poleci albo za dalko
			if odbicie == 5 || pilka.posY > 30 || pilka.posX > 10 || pilka.posX < -10 {
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
				for i := 0; i < 1; i++ { // todo 1 na 3
					sila_ruszajaca[i] = 0

					wektory_sily_ruszajacej[i] = Wektor{
						x: 0,
						y: 0,
						z: 0,
					}

					przyspieszenia_z_wektorow[i] = Wektor{
						x: 0,
						y: 0,
						z: 0,
					}

					momenty_z_wektorow[i] = Wektor{
						x: 0,
						y: 0,
						z: 0,
					}

					przyspieszenia_katowe_z_wektorow[i] = Wektor{
						x: 0,
						y: 0,
						z: 0,
					}
				}
				wypadkowa_przyspieszen = Wektor{
					x: 0,
					y: 0,
					z: 0,
				}

				wypadkowa_przyspieszen_katowych = Wektor{
					x: 0,
					y: 0,
					z: 0,
				}

				predkosci_pilki = Wektor{
					x: 0,
					y: 0,
					z: 0,
				}

				predkosci_katowe_pilki = Wektor{
					x: 0,
					y: 0,
					z: 0,
				}

			}

			pilka.posX = pilka.posX + predkosci_pilki.x*krok_czasowy
			pilka.posY = pilka.posY + predkosci_pilki.y*krok_czasowy + g*krok_czasowy*krok_czasowy/2
			predkosci_pilki.y = predkosci_pilki.y + g*krok_czasowy
			pilka.posZ = pilka.posZ + predkosci_pilki.z*krok_czasowy

			pilka.rotX = pilka.rotX + predkosci_katowe_pilki.x*krok_czasowy
			pilka.rotY = pilka.rotY + predkosci_katowe_pilki.y*krok_czasowy
			pilka.rotZ = pilka.rotZ + predkosci_katowe_pilki.z*krok_czasowy

			pilka.rotX = pilnowanie_zakresu_kata(pilka.rotX)
			pilka.rotY = pilnowanie_zakresu_kata(pilka.rotY)
			pilka.rotZ = pilnowanie_zakresu_kata(pilka.rotZ)

			qX := quaternion.FromXAxisAngle(pilka.rotX)
			qY := quaternion.FromYAxisAngle(pilka.rotY)
			qZ := quaternion.FromZAxisAngle(pilka.rotZ)

			kwaternion_glowny := quaternion.Mul3(&qX, &qY, &qZ)

			os, kat := kwaternion_glowny.AxisAngle()

			//println(os.String(), " ", kat, " ", pilka.posY)

			rotationAxis = rl.NewVector3(float32(os[0]), float32(os[1]), float32(os[2]))
			rotationAngle = float32(kat * 180 / math.Pi)
		}
		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		// Draw the model with rotation
		rl.DrawModelEx(
			sphereModel,
			rl.NewVector3(float32(pilka.posX), float32(pilka.posY), float32(pilka.posZ)), // Position
			rotationAxis,                 // Rotation axis
			rotationAngle,                // Rotation angle
			rl.NewVector3(1.0, 1.0, 1.0), // Scale
			rl.White,                     // Tint
		)
		if faza_gry == 0 && licznik > 2 {
			rl.DrawModel(
				punktModel,
				punktPosition,
				1,
				rl.White,
			)
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

		// Optionally, draw a grid for reference
		rl.DrawGrid(100, 1.0)

		rl.EndMode3D()

		rl.EndDrawing()
	}

	// De-initialization
	rl.UnloadTexture(texture)   // Unload texture
	rl.UnloadModel(sphereModel) // Unload model

	rl.CloseWindow() // Close window and OpenGL context
}

func pilnowanie_zakresu_kata(rotacja float64) float64 {
	for i := 0; i < 1; i++ {

		if rotacja > 2*math.Pi {
			rotacja = rotacja - 2*math.Pi
			i--
		}
		if rotacja < 0 {
			rotacja = rotacja + 2*math.Pi
			i--
		}
	}
	return rotacja
} /*
func Sign(x float64) float64 {
	if x > 0.001 {
		return 1
	} else if x < -0.001 {
		return -1
	}
	return 0
}*/
func iloczyn_wektorowy(wektor_1 Wektor, wektor_2 Wektor) Wektor {
	return Wektor{
		x: wektor_1.y*wektor_2.z - wektor_1.z*wektor_2.y,
		y: wektor_1.z*wektor_2.x - wektor_1.x - wektor_2.z,
		z: wektor_1.x*wektor_2.y - wektor_1.y - wektor_2.x,
	}
}
