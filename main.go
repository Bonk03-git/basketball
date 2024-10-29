package main

// PAMIĘTAĆ ZAMIENIAĆ WSPOLRZEDNE X Z Y NA KONIEC PROJEKTU!!!
import (
	"github.com/g3n/engine/math32"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/mathgl/mgl32"
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
	const g = -9.81

	pilka := Kula{
		promien: 0.12, // metry
		posX:    0,    // pixele float64(rand.Intn(10))
		posY:    0,    // pixele float64(rand.Intn(10))
		posZ:    10,   // pixele
		rotX:    0,    // radiany
		rotY:    0,    // radiany
		rotZ:    0,    // radiany
		masa:    0.5,  // kilogramy
	}

	var punkty_przylozenia_wektorow [3]Punkt_przylozenia
	var wektory_sily [3]Wektor_sily
	var sila_ruszajaca [3]float64
	var wektory_sily_ruszajacej [3]Wektor
	var wektory_sily_obrotowej [3]Wektor
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

	var I = 2.0 / 5.0 * pilka.masa * pilka.promien * pilka.promien

	for i := 0; i < 3; i++ {

		punkty_przylozenia_wektorow[i] = Punkt_przylozenia{ // wyznaczane przez gracza
			kat_fi:   math.Pi / 4,
			kat_teta: 0,
		}

		wektory_sily[i] = Wektor_sily{
			dlugosc_x:     -0.1, // wyznacza gracz
			dlugosc_y:     0,    // wyznacza gracz
			dlugosc_z:     0,    // wyznacza gracz
			przylozenie_x: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Cos(punkty_przylozenia_wektorow[i].kat_fi),
			przylozenie_y: pilka.promien * math.Cos(punkty_przylozenia_wektorow[i].kat_teta) * math.Sin(punkty_przylozenia_wektorow[i].kat_fi),
			przylozenie_z: pilka.promien * math.Sin(punkty_przylozenia_wektorow[i].kat_teta),
		}

		sila_ruszajaca[i] = wektory_sily[i].dlugosc_x*wektory_sily[i].przylozenie_x/pilka.promien + wektory_sily[i].dlugosc_y*wektory_sily[i].przylozenie_y/pilka.promien + wektory_sily[i].dlugosc_z*wektory_sily[i].przylozenie_z/pilka.promien

		wektory_sily_ruszajacej[i] = Wektor{
			x: sila_ruszajaca[i] * wektory_sily[i].przylozenie_x / pilka.promien,
			y: sila_ruszajaca[i] * wektory_sily[i].przylozenie_y / pilka.promien,
			z: sila_ruszajaca[i] * wektory_sily[i].przylozenie_z / pilka.promien,
		}

		wektory_sily_obrotowej[i] = Wektor{
			x: wektory_sily[i].dlugosc_x - wektory_sily_ruszajacej[i].x,
			y: wektory_sily[i].dlugosc_y - wektory_sily_ruszajacej[i].y,
			z: wektory_sily[i].dlugosc_z - wektory_sily_ruszajacej[i].z,
		}

		przyspieszenia_z_wektorow[i] = Wektor{
			x: wektory_sily_ruszajacej[i].x / pilka.masa,
			y: wektory_sily_ruszajacej[i].y / pilka.masa,
			z: wektory_sily_ruszajacej[i].z / pilka.masa,
		}

		momenty_z_wektorow[i] = Wektor{
			x: wektory_sily[i].przylozenie_y*wektory_sily_obrotowej[i].z - wektory_sily[i].przylozenie_z*wektory_sily_obrotowej[i].y,
			y: -wektory_sily[i].przylozenie_x*wektory_sily_obrotowej[i].z + wektory_sily[i].przylozenie_z*wektory_sily_obrotowej[i].x,
			z: wektory_sily[i].przylozenie_x*wektory_sily_obrotowej[i].y - wektory_sily[i].przylozenie_y*wektory_sily_obrotowej[i].x,
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

	// Widok kamery na obszar 3D
	camera := rl.Camera{
		Position:   rl.NewVector3(1.0, 1.0, 0.0), // pozycja kamery
		Target:     rl.NewVector3(0, 0, 0),       // punkt patrzenia
		Up:         rl.NewVector3(0.0, 1.0, 0.0), // Camera up vector (rotation towards target)
		Fovy:       45.0,                         // Camera field-of-view Y
		Projection: rl.CameraPerspective,         // Camera mode type
	}

	rl.DisableCursor()

	// Create a sphere mesh and model
	sphereMesh := rl.GenMeshSphere(0.12, 32, 32)
	sphereModel := rl.LoadModelFromMesh(sphereMesh)

	// Load the basketball texture
	texture := rl.LoadTexture("basketball.png") // Ensure this file exists

	// Access the materials using GetMaterials()
	materials := sphereModel.GetMaterials()

	// Use SetMaterialTexture to set the texture
	rl.SetMaterialTexture(&materials[0], rl.MapDiffuse, texture)

	basketMesh := rl.GenMeshCube(0.1, 1.05, 1.8)
	basketmodel := rl.LoadModelFromMesh(basketMesh)

	texture_2 := rl.LoadTexture("kosz.jpg")

	materials = basketmodel.GetMaterials()
	// Use SetMaterialTexture to set the texture
	rl.SetMaterialTexture(&materials[0], rl.MapDiffuse, texture_2)

	// Set initial rotation angle
	rotationAngle := float32(0.0)
	rotationAxis := rl.NewVector3(1.0, 1.0, 0.0)
	rotationAngle_2 := float32(0.0)
	rotationAxis_2 := rl.NewVector3(0.0, 1.0, 0.0)

	rl.SetTargetFPS(60) // Set our game to run at 60 frames-per-second

	// Main game loop
	for !rl.WindowShouldClose() {

		rl.UpdateCamera(&camera, rl.CameraFree)

		if rl.IsKeyPressed(rl.KeyZ) {
			camera.Target = rl.NewVector3(0.0, 0.0, 0.0)
		}

		pilka.posX = pilka.posX + predkosci_pilki.x*krok_czasowy
		pilka.posY = pilka.posY + predkosci_pilki.y*krok_czasowy
		pilka.posZ = pilka.posZ + predkosci_pilki.z*krok_czasowy + g*krok_czasowy*krok_czasowy/2
		predkosci_pilki.z = predkosci_pilki.z + g*krok_czasowy
		pilka.rotX = pilka.rotX + predkosci_katowe_pilki.x*krok_czasowy
		pilka.rotY = pilka.rotY + predkosci_katowe_pilki.y*krok_czasowy
		pilka.rotZ = pilka.rotZ + predkosci_katowe_pilki.z*krok_czasowy

		pilka.rotX = pilnowanie_zakresu_kata(pilka.rotX)
		pilka.rotY = pilnowanie_zakresu_kata(pilka.rotY)
		pilka.rotZ = pilnowanie_zakresu_kata(pilka.rotZ)

		//kwaterniony

		q_x := mgl32.QuatRotate(float32(pilka.rotX), mgl32.Vec3{1, 0, 0})
		q_y := mgl32.QuatRotate(float32(pilka.rotY), mgl32.Vec3{0, 1, 0})
		q_z := mgl32.QuatRotate(float32(pilka.rotZ), mgl32.Vec3{0, 0, 1})
		q_wynik := q_z.Mul(q_y).Mul(q_x)

		kat_obrotu := 2 * math.Acos(float64(q_wynik.W))
		var os_obrotu math32.Vector3

		if kat_obrotu > 0 {
			// Normalizacja wektora osi obrotu
			s := math.Sqrt(float64(1 - q_wynik.W*q_wynik.W))
			os_obrotu = *math32.NewVector3(q_wynik.X()/float32(s), q_wynik.Y()/float32(s), q_wynik.Z()/float32(s))
		} else {
			// Oś obrotu jest nieokreślona, możemy przyjąć dowolną oś
			os_obrotu = *math32.NewVector3(1, 0, 0) // Przykładowa oś
		}

		rotationAxis = rl.NewVector3(os_obrotu.X, os_obrotu.Y, os_obrotu.Z)
		rotationAngle = float32(kat_obrotu)
		println(rotationAngle)

		// Draw
		rl.BeginDrawing()

		rl.ClearBackground(rl.RayWhite)

		rl.BeginMode3D(camera)

		// Draw the model with rotation
		rl.DrawModelEx(
			sphereModel,
			rl.NewVector3(1.0, 0.0, 0.0), // Position
			rotationAxis,                 // Rotation axis
			rotationAngle,                // Rotation angle
			rl.NewVector3(1.0, 1.0, 1.0), // Scale
			rl.White,                     // Tint
		)

		rl.DrawModelEx(
			basketmodel,
			rl.NewVector3(10.0, 0.0, 0.0), // Position
			rotationAxis_2,                // Rotation axis
			rotationAngle_2,               // Rotation angle
			rl.NewVector3(1.0, 1.0, 1.0),  // Scale
			rl.White,                      // Tint
		)

		// Optionally, draw a grid for reference
		rl.DrawGrid(100, 1.0)

		rl.EndMode3D()

		rl.DrawText("Rotating Basketball", 10, 10, 20, rl.DarkGray)

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
}
