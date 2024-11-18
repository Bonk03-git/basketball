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
	const fps = 200
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
	for i := 0; i < 1; i++ { //todo 1 na 3

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
		Position:   rl.NewVector3(1.0, 1.0, 0.0),
		Target:     rl.NewVector3(0, 0, 0),
		Up:         rl.NewVector3(0.0, 1.0, 0.0),
		Fovy:       45.0,
		Projection: rl.CameraPerspective,
	}

	rl.DisableCursor()
	//inicjalizacja obiektow

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

	//inicjalizacja wartości zmiennych odpowiadających za fazy programu

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
					punkty_przylozenia_wektorow[i].kat_teta = po_wcisnieciu(punkty_przylozenia_wektorow[i].kat_teta, 2*math32.Pi/360)
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}
				if licznik == 3 {
					rl.DrawText(fmt.Sprintf("Przylozenie sily kat fi = %f", punkty_przylozenia_wektorow[i].kat_fi), 10, 10, 20, rl.DarkGray)
					punktPosition.X = float32(pilka.posX + pilka.promien*math.Cos(punkty_przylozenia_wektorow[i].kat_teta)*math.Sin(punkty_przylozenia_wektorow[i].kat_fi))
					punktPosition.Y = float32(pilka.posY + pilka.promien*math.Sin(punkty_przylozenia_wektorow[i].kat_teta))
					punktPosition.Z = float32(pilka.posZ + pilka.promien*math.Cos(punkty_przylozenia_wektorow[i].kat_teta)*math.Cos(punkty_przylozenia_wektorow[i].kat_fi))
					punkty_przylozenia_wektorow[i].kat_fi = po_wcisnieciu(punkty_przylozenia_wektorow[i].kat_fi, 2*math32.Pi/360)
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}
				if licznik == 2 {
					rl.DrawText(fmt.Sprintf("Wartosc wektora w kierunku z= %f", wektory_sily[i].dlugosc_z), 10, 10, 20, rl.DarkGray)
					wektory_sily[i].dlugosc_z = po_wcisnieciu(wektory_sily[i].dlugosc_z, stala_wzrostu_malenia)
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}
				if licznik == 1 {
					rl.DrawText(fmt.Sprintf("Wartosc wektora w kierunku y= %f", wektory_sily[i].dlugosc_y), 10, 10, 20, rl.DarkGray)
					wektory_sily[i].dlugosc_y = po_wcisnieciu(wektory_sily[i].dlugosc_y, stala_wzrostu_malenia)
					if rl.IsKeyPressed(rl.KeyEnter) {
						licznik += 1
					}
				}

				if licznik == 0 {
					rl.DrawText(fmt.Sprintf("Wartosc wektora w kierunku x= %f", wektory_sily[i].dlugosc_x), 10, 10, 20, rl.DarkGray)
					wektory_sily[i].dlugosc_x = po_wcisnieciu(wektory_sily[i].dlugosc_x, stala_wzrostu_malenia)
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
				pilka.posY = pilka.promien
				bufor_x := predkosci_pilki.x
				bufor_z := predkosci_pilki.z
				predkosci_pilki.y = -predkosci_pilki.y * wspolczynnik_odbicia
				predkosci_pilki.x = wspolczynnik_odbicia * (bufor_x - wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.z)
				predkosci_pilki.z = wspolczynnik_odbicia * (bufor_z + wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.x)
				predkosci_katowe_pilki.x = wspolczynnik_odbicia * (predkosci_katowe_pilki.x + wspolczynnik_momentu*bufor_z/pilka.promien)
				predkosci_katowe_pilki.z = wspolczynnik_odbicia * (predkosci_katowe_pilki.z - wspolczynnik_momentu*bufor_x/pilka.promien)

				// zapewnienie wyjścia
				for pilka.posY-pilka.promien < 0 {

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

					liczba_krokow += 1
					px = append(px, pilka.posX)
					py = append(py, pilka.posY)
					pz = append(pz, pilka.posZ)
					ox = append(ox, pilka.rotX)
					oy = append(oy, pilka.rotY)
					oz = append(oz, pilka.rotZ)
					vx = append(vx, predkosci_pilki.x)
					vy = append(vy, predkosci_pilki.y)
					vz = append(vz, predkosci_pilki.z)
					wx = append(wx, predkosci_katowe_pilki.x)
					wy = append(wy, predkosci_katowe_pilki.y)
					wz = append(wz, predkosci_katowe_pilki.z)
					czas = append(czas, float64(liczba_krokow)*krok_czasowy)
				}
				odbicie += 1

			}

			//odbicia od tablicy

			//powierzchnia
			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posX+pilka.promien > x_tablicy-grubosc_tablicy/2 && pilka.posX < x_tablicy+grubosc_tablicy/2 {
				//wyjscie pilki z obszru odbicia
				pilka.posX = x_tablicy - grubosc_tablicy/2 - pilka.promien
				//predkosci
				bufor_y := predkosci_pilki.y
				bufor_z := predkosci_pilki.z
				predkosci_pilki.x = -predkosci_pilki.x * wspolczynnik_odbicia
				predkosci_pilki.y = wspolczynnik_odbicia * (bufor_y - wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.z)
				predkosci_pilki.z = wspolczynnik_odbicia * (bufor_z + wspolczynnik_momentu*pilka.promien*predkosci_katowe_pilki.y)
				predkosci_katowe_pilki.y = wspolczynnik_odbicia * (predkosci_katowe_pilki.y + wspolczynnik_momentu*bufor_z/pilka.promien)
				predkosci_katowe_pilki.z = wspolczynnik_odbicia * (predkosci_katowe_pilki.z - wspolczynnik_momentu*bufor_y/pilka.promien)

				//zapewnienie wyjścia piłki z obszaru odbicia
				for pilka.posY < y_tablicy+wysokosc_tablicy/2 &&
					pilka.posY > y_tablicy-wysokosc_tablicy/2 &&
					pilka.posZ > z_tablicy-szerokosc_tablicy/2 &&
					pilka.posZ < z_tablicy+szerokosc_tablicy/2 &&
					pilka.posX+pilka.promien > x_tablicy-grubosc_tablicy/2 &&
					pilka.posX < x_tablicy+grubosc_tablicy/2 {

					liczba_krokow += 1

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

					px = append(px, pilka.posX)
					py = append(py, pilka.posY)
					pz = append(pz, pilka.posZ)
					ox = append(ox, pilka.rotX)
					oy = append(oy, pilka.rotY)
					oz = append(oz, pilka.rotZ)
					vx = append(vx, predkosci_pilki.x)
					vy = append(vy, predkosci_pilki.y)
					vz = append(vz, predkosci_pilki.z)
					wx = append(wx, predkosci_katowe_pilki.x)
					wy = append(wy, predkosci_katowe_pilki.y)
					wz = append(wz, predkosci_katowe_pilki.z)
					czas = append(czas, float64(liczba_krokow)*krok_czasowy)
				}
				println("Uderzona tablica")
			}

			//krawedz prawa
			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posZ > z_tablicy+szerokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) {
				println("krawedz prawa")
				var predkosci_pocz = Wektor{
					x: predkosci_pilki.x,
					y: predkosci_pilki.y,
					z: predkosci_pilki.z,
				}

				//wektor jednostkowy normalnej

				var normalna = Wektor{
					x: pilka.posX - x_tablicy,
					y: 0.0,
					z: pilka.posZ - z_tablicy - szerokosc_tablicy/2,
				}
				var dlugosc_wektora = math.Sqrt(math.Pow(normalna.x, 2) + math.Pow(normalna.y, 2) + math.Pow(normalna.z, 2))
				normalna.x = normalna.x / dlugosc_wektora
				normalna.y = normalna.y / dlugosc_wektora
				normalna.z = normalna.z / dlugosc_wektora

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
				for pilka.posY < y_tablicy+wysokosc_tablicy/2 &&
					pilka.posY > y_tablicy-wysokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posZ > z_tablicy+szerokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

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

					liczba_krokow += 1
					px = append(px, pilka.posX)
					py = append(py, pilka.posY)
					pz = append(pz, pilka.posZ)
					ox = append(ox, pilka.rotX)
					oy = append(oy, pilka.rotY)
					oz = append(oz, pilka.rotZ)
					vx = append(vx, predkosci_pilki.x)
					vy = append(vy, predkosci_pilki.y)
					vz = append(vz, predkosci_pilki.z)
					wx = append(wx, predkosci_katowe_pilki.x)
					wy = append(wy, predkosci_katowe_pilki.y)
					wz = append(wz, predkosci_katowe_pilki.z)
					czas = append(czas, float64(liczba_krokow)*krok_czasowy)
				}

			}

			//krawedz lewa
			if pilka.posY < y_tablicy+wysokosc_tablicy/2 && pilka.posY > y_tablicy-wysokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posZ < z_tablicy-szerokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy+szerokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) {

				println("krawedz lewa")
				var predkosci_pocz = Wektor{
					x: predkosci_pilki.x,
					y: predkosci_pilki.y,
					z: predkosci_pilki.z,
				}
				//wektor jednostkowy normalnej

				var normalna = Wektor{
					x: pilka.posX - x_tablicy,
					y: 0.0,
					z: pilka.posZ - z_tablicy + szerokosc_tablicy/2,
				}
				var dlugosc_wektora = math.Sqrt(math.Pow(normalna.x, 2) + math.Pow(normalna.y, 2) + math.Pow(normalna.z, 2))
				normalna.x = normalna.x / dlugosc_wektora
				normalna.y = normalna.y / dlugosc_wektora
				normalna.z = normalna.z / dlugosc_wektora

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
				for pilka.posY < y_tablicy+wysokosc_tablicy/2 &&
					pilka.posY > y_tablicy-wysokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posZ < z_tablicy-szerokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posZ-z_tablicy+szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

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

					liczba_krokow += 1
					px = append(px, pilka.posX)
					py = append(py, pilka.posY)
					pz = append(pz, pilka.posZ)
					ox = append(ox, pilka.rotX)
					oy = append(oy, pilka.rotY)
					oz = append(oz, pilka.rotZ)
					vx = append(vx, predkosci_pilki.x)
					vy = append(vy, predkosci_pilki.y)
					vz = append(vz, predkosci_pilki.z)
					wx = append(wx, predkosci_katowe_pilki.x)
					wy = append(wy, predkosci_katowe_pilki.y)
					wz = append(wz, predkosci_katowe_pilki.z)
					czas = append(czas, float64(liczba_krokow)*krok_czasowy)
				}
			}

			//krawedz gora
			if pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posY > y_tablicy+wysokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) {

				println("krawedz gora")
				var predkosci_pocz = Wektor{
					x: predkosci_pilki.x,
					y: predkosci_pilki.y,
					z: predkosci_pilki.z,
				}
				//wektor jednostkowy normalnej

				var normalna = Wektor{
					x: pilka.posX - x_tablicy,
					y: pilka.posY - y_tablicy - wysokosc_tablicy/2,
					z: 0.0,
				}

				var dlugosc_wektora = math.Sqrt(math.Pow(normalna.x, 2) + math.Pow(normalna.y, 2) + math.Pow(normalna.z, 2))
				normalna.x = normalna.x / dlugosc_wektora
				normalna.y = normalna.y / dlugosc_wektora
				normalna.z = normalna.z / dlugosc_wektora

				//r to wektor od środka piłki do punktu styku z krawędzią
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
				for pilka.posZ < z_tablicy+szerokosc_tablicy/2 &&
					pilka.posZ > z_tablicy-szerokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posY > y_tablicy+wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

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

					liczba_krokow += 1
					px = append(px, pilka.posX)
					py = append(py, pilka.posY)
					pz = append(pz, pilka.posZ)
					ox = append(ox, pilka.rotX)
					oy = append(oy, pilka.rotY)
					oz = append(oz, pilka.rotZ)
					vx = append(vx, predkosci_pilki.x)
					vy = append(vy, predkosci_pilki.y)
					vz = append(vz, predkosci_pilki.z)
					wx = append(wx, predkosci_katowe_pilki.x)
					wy = append(wy, predkosci_katowe_pilki.y)
					wz = append(wz, predkosci_katowe_pilki.z)
					czas = append(czas, float64(liczba_krokow)*krok_czasowy)
				}
			}

			//krawedz dol
			if pilka.posZ < z_tablicy+szerokosc_tablicy/2 && pilka.posZ > z_tablicy-szerokosc_tablicy/2 && pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien && pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien && pilka.posY < y_tablicy-wysokosc_tablicy/2 && (math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2)) < math.Pow(pilka.promien, 2) {

				println("krawedz dol")
				var predkosci_pocz = Wektor{
					x: predkosci_pilki.x,
					y: predkosci_pilki.y,
					z: predkosci_pilki.z,
				}
				//wektor jednostkowy normalnej

				var normalna = Wektor{
					x: pilka.posX - x_tablicy,
					y: pilka.posY - y_tablicy + wysokosc_tablicy/2,
					z: 0.0,
				}

				var dlugosc_wektora = math.Sqrt(math.Pow(normalna.x, 2) + math.Pow(normalna.y, 2) + math.Pow(normalna.z, 2))
				normalna.x = normalna.x / dlugosc_wektora
				normalna.y = normalna.y / dlugosc_wektora
				normalna.z = normalna.z / dlugosc_wektora

				//r to wektor od środka piłki do punktu styku z krawędzią
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
				for pilka.posZ < z_tablicy+szerokosc_tablicy/2 &&
					pilka.posZ > z_tablicy-szerokosc_tablicy/2 &&
					pilka.posX > x_tablicy-grubosc_tablicy/2-pilka.promien &&
					pilka.posX < x_tablicy+grubosc_tablicy/2+pilka.promien &&
					pilka.posY < y_tablicy-wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy+wysokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

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

					liczba_krokow += 1
					px = append(px, pilka.posX)
					py = append(py, pilka.posY)
					pz = append(pz, pilka.posZ)
					ox = append(ox, pilka.rotX)
					oy = append(oy, pilka.rotY)
					oz = append(oz, pilka.rotZ)
					vx = append(vx, predkosci_pilki.x)
					vy = append(vy, predkosci_pilki.y)
					vz = append(vz, predkosci_pilki.z)
					wx = append(wx, predkosci_katowe_pilki.x)
					wy = append(wy, predkosci_katowe_pilki.y)
					wz = append(wz, predkosci_katowe_pilki.z)
					czas = append(czas, float64(liczba_krokow)*krok_czasowy)
				}
			}

			// prawy gorny rog
			if pilka.posZ > z_tablicy+szerokosc_tablicy/2 && pilka.posY > y_tablicy+wysokosc_tablicy/2 && math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

				println("prawy gorny rog")

				var punkt_odbicia_wzgl_pilki = Wektor{
					x: x_tablicy - pilka.posX,
					y: y_tablicy + wysokosc_tablicy/2 - pilka.posY,
					z: z_tablicy + szerokosc_tablicy/2 - pilka.posZ,
				}

				dlugosc_wektora_puktu_odbicia_od_pilki := math.Sqrt(math.Pow(punkt_odbicia_wzgl_pilki.x, 2) + math.Pow(punkt_odbicia_wzgl_pilki.y, 2) + math.Pow(punkt_odbicia_wzgl_pilki.z, 2))

				var predkosc_katowa_z_punktem_odbicia = iloczyn_wektorowy(predkosci_katowe_pilki, punkt_odbicia_wzgl_pilki)

				var predkosc_punktu_styku_przed_odbiciem = dodaj_wektory(predkosci_pilki, predkosc_katowa_z_punktem_odbicia)

				//wektor normalny

				var normalna = Wektor{
					x: punkt_odbicia_wzgl_pilki.x / dlugosc_wektora_puktu_odbicia_od_pilki,
					y: punkt_odbicia_wzgl_pilki.y / dlugosc_wektora_puktu_odbicia_od_pilki,
					z: punkt_odbicia_wzgl_pilki.z / dlugosc_wektora_puktu_odbicia_od_pilki,
				}

				predkosc_normalna := predkosci_pilki.x*normalna.x + predkosci_pilki.y*normalna.y + predkosci_pilki.z*normalna.z

				var predkosc_styczna = Wektor{
					x: predkosc_punktu_styku_przed_odbiciem.x - predkosc_normalna*normalna.x,
					y: predkosc_punktu_styku_przed_odbiciem.y - predkosc_normalna*normalna.y,
					z: predkosc_punktu_styku_przed_odbiciem.z - predkosc_normalna*normalna.z,
				}

				predkosc_normalna = -wspolczynnik_odbicia * predkosc_normalna

				var predkosc_punktu_styku_po_odbiciu = Wektor{
					x: predkosc_normalna*normalna.x + predkosc_styczna.x,
					y: predkosc_normalna*normalna.y + predkosc_styczna.y,
					z: predkosc_normalna*normalna.z + predkosc_styczna.z,
				}

				var impuls = Wektor{
					x: -pilka.masa * (predkosc_punktu_styku_po_odbiciu.x - predkosc_punktu_styku_przed_odbiciem.x),
					y: -pilka.masa * (predkosc_punktu_styku_po_odbiciu.y - predkosc_punktu_styku_przed_odbiciem.y),
					z: -pilka.masa * (predkosc_punktu_styku_po_odbiciu.z - predkosc_punktu_styku_przed_odbiciem.z),
				}
				predkosci_pilki = Wektor{
					x: predkosci_pilki.x + impuls.x/pilka.masa,
					y: predkosci_pilki.y + impuls.y/pilka.masa,
					z: predkosci_pilki.z + impuls.z/pilka.masa,
				}
				var punkt_odbicia_wzgl_pilki_z_impulsem = iloczyn_wektorowy(punkt_odbicia_wzgl_pilki, impuls)

				predkosci_katowe_pilki = Wektor{
					x: predkosci_katowe_pilki.x + punkt_odbicia_wzgl_pilki_z_impulsem.x/I,
					y: predkosci_katowe_pilki.y + punkt_odbicia_wzgl_pilki_z_impulsem.y/I,
					z: predkosci_katowe_pilki.z + punkt_odbicia_wzgl_pilki_z_impulsem.z/I,
				}

				//zapewnienie wyjścia piłki z miejsca odbicia
				for pilka.posZ > z_tablicy+szerokosc_tablicy/2 &&
					pilka.posY > y_tablicy+wysokosc_tablicy/2 &&
					math.Pow(pilka.posX-x_tablicy, 2)+math.Pow(pilka.posY-y_tablicy-wysokosc_tablicy/2, 2)+math.Pow(pilka.posZ-z_tablicy-szerokosc_tablicy/2, 2) < math.Pow(pilka.promien, 2) {

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
				}
			}

			pilka.posX = pilka.posX + predkosci_pilki.x*krok_czasowy
			pilka.posY = pilka.posY + predkosci_pilki.y*krok_czasowy + g*math.Pow(krok_czasowy, 2)/2
			predkosci_pilki.y = predkosci_pilki.y + g*krok_czasowy
			pilka.posZ = pilka.posZ + predkosci_pilki.z*krok_czasowy

			pilka.rotX = pilka.rotX + predkosci_katowe_pilki.x*krok_czasowy
			pilka.rotY = pilka.rotY + predkosci_katowe_pilki.y*krok_czasowy
			pilka.rotZ = pilka.rotZ + predkosci_katowe_pilki.z*krok_czasowy

			pilka.rotX = pilnowanie_zakresu_kata(pilka.rotX)
			pilka.rotY = pilnowanie_zakresu_kata(pilka.rotY)
			pilka.rotZ = pilnowanie_zakresu_kata(pilka.rotZ)

			liczba_krokow += 1
			px = append(px, pilka.posX)
			py = append(py, pilka.posY)
			pz = append(pz, pilka.posZ)
			ox = append(ox, pilka.rotX)
			oy = append(oy, pilka.rotY)
			oz = append(oz, pilka.rotZ)
			vx = append(vx, predkosci_pilki.x)
			vy = append(vy, predkosci_pilki.y)
			vz = append(vz, predkosci_pilki.z)
			wx = append(wx, predkosci_katowe_pilki.x)
			wy = append(wy, predkosci_katowe_pilki.y)
			wz = append(wz, predkosci_katowe_pilki.z)
			czas = append(czas, float64(liczba_krokow)*krok_czasowy)

			qX := quaternion.FromXAxisAngle(pilka.rotX)
			qY := quaternion.FromYAxisAngle(pilka.rotY)
			qZ := quaternion.FromZAxisAngle(pilka.rotZ)

			kwaternion_glowny := quaternion.Mul3(&qX, &qY, &qZ)

			oska, kat := kwaternion_glowny.AxisAngle()

			rotationAxis = rl.NewVector3(float32(oska[0]), float32(oska[1]), float32(oska[2]))
			rotationAngle = float32(kat * 180 / math.Pi)
		}

		// powrot do rzutu gdy pilka sie 5 razy odbije od ziemii, za wysoko poleci albo za dalko
		if odbicie == 5 || pilka.posY > 30 || pilka.posX > 10 || pilka.posX < -10 || pilka.posZ < -10 || pilka.posZ > 10 {
			pilka.posX = float64(start_pos_x)
			pilka.posY = float64(start_pos_y)
			pilka.posZ = float64(start_pos_z)
			pilka.rotX = float64(start_rot_x)
			pilka.rotY = float64(start_rot_y)
			pilka.rotZ = float64(start_rot_z)
			licznik = 0
			licz_wartosci = true
			faza_gry = 2
			odbicie = 0
			for i := 0; i < 1; i++ { // todo 1 na 3

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

		if faza_gry == 2 {

			zapisz_obraz(czas, px, "pozycja_x.png")
			zapisz_obraz(czas, py, "pozycja_y.png")
			zapisz_obraz(czas, pz, "pozycja_z.png")
			zapisz_obraz(czas, ox, "obrot_x.png")
			zapisz_obraz(czas, oy, "obrot_y.png")
			zapisz_obraz(czas, oz, "obrot_z.png")
			zapisz_obraz(czas, vx, "predkosc_liniowa_x.png")
			zapisz_obraz(czas, vy, "predkosc_liniowa_y.png")
			zapisz_obraz(czas, vz, "predkosc_liniowa_z.png")
			zapisz_obraz(czas, wx, "predkosc_katowa_x.png")
			zapisz_obraz(czas, wy, "predkosc_katowa_y.png")
			zapisz_obraz(czas, wz, "predkosc_katowa_z.png")

			faza_gry = 0
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

		rl.DrawGrid(100, 1.0)

		rl.EndMode3D()

		rl.EndDrawing()
	}
	// De-initialization

	rl.UnloadTexture(texture)
	rl.UnloadModel(sphereModel)

	rl.CloseWindow()
}

func po_wcisnieciu(zmienna float64, zmiana float64) float64 {
	if rl.IsKeyPressed(rl.KeyI) {
		zmienna += 10 * zmiana
	}
	if rl.IsKeyPressed(rl.KeyU) {
		zmienna += zmiana
	}
	if rl.IsKeyPressed(rl.KeyY) {
		zmienna -= zmiana
	}
	if rl.IsKeyPressed(rl.KeyT) {
		zmienna -= 10 * zmiana
	}
	return zmienna
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
}

func iloczyn_wektorowy(wektor_1 Wektor, wektor_2 Wektor) Wektor {
	return Wektor{
		x: wektor_1.y*wektor_2.z - wektor_1.z*wektor_2.y,
		y: wektor_1.z*wektor_2.x - wektor_1.x - wektor_2.z,
		z: wektor_1.x*wektor_2.y - wektor_1.y - wektor_2.x,
	}
}
func dodaj_wektory(wektor_1 Wektor, wektor_2 Wektor) Wektor {
	return Wektor{
		x: wektor_1.x + wektor_2.x,
		y: wektor_1.y + wektor_2.y,
		z: wektor_1.z + wektor_2.z,
	}
}

func zapisz_obraz(tab_x []float64, tab_y []float64, nazwa string) {
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: tab_x,
				YValues: tab_y,
			},
		},
	}
	f, _ := os.Create(nazwa)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	err := graph.Render(chart.PNG, f)
	if err != nil {
		return
	}
	return
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
