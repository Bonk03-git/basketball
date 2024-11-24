PROJEKT RZUTU PIŁKI DO KOSZA 

Bartłomiej Bąk 193634

Podczas wykonywania projektu został wykorzystany Chat GPT oraz kilka przykładowych rozwiązań graficznych z oficjalnej strony raylib.com


Zasada działania programu:

	Za pomocą klawiszy na klawiaturze można poruszać się w przestrzeni:

		- W,S,A,D - poruszanie się
		- Q,E - obrót kamery

	Opis przestrzeni:
		
		- oś X - oś prostopadła do tablicy kosza (ruch przód, tył)
		- oś Y - oś pionowa (ruch góra, dół)
		- oś Z - oś równoległa do tablicy kosza (ruch prawo, lewo)

	Za pomocą klawiszy można zmieniać wartości kątów i wartości wektorów:

		- T - dla wartości wektorów spadek wartości o 10 [N], dla kątów spadek wartości o 0,087266 [rad]
		- Y - dla wartości wektorów spadek wartości o 1 [N], dla kątów spadek wartości o 0,0087266 [rad]
		- U - dla wartości wektorów wzrost wartości o 1 [N], dla kątów wzrost wartości o 0,0087266 [rad]
		- I - dla wartości wektorów wzrost wartości o 10 [N], dla kątów wzrost wartości o 0,087266 [rad]
	
	Aby piłka nabierała nowych wartości, wprowadzone liczby należy zatwierdzać ENTEREM.

	Wektory posiadają ograniczenia wartości (-50, 50[N]).

	Piłka po zatwierdzeniu wszystkich wartości dla trzech wektorów zaczyna się poruszać.

	W momencie trafienia do obręczy program sygnalizuje o trafieniu kolorując tło na zielono.

	Piłka kończy ruch i wraca do pozycji domowej w momencie gdy:

		- Poleci za daleko w kierunku osi X lub osi Z
		- Odbije się 5 razy od podłoża

Po zamknięciu programu przyciskiem ESCAPE program zapisuje wykresy pozycji piłki w przestrzeni, jej obrotu wokół danych osi, prędkości liniowej i prędkości obrotowej w zależności od czasu (od początku pierwszego rzutu do zamknięcia okna), które można później otworzyć poza programem. Obrazki się nie utworzą jeżeli piłka nie zostanie wprawiona w ruch.
	
Wszystkie zmienne w programie są liczone w jednostkach układu SI.

Możliwe jest wystąpienie błędu polegającego na ominięciu przez piłkę obiektów jeżeli piłka pokryje się z obiektem w małym stopniu i posiada dużą prędkość liniową, wynika to z kroku czsowego zastosowanego w programie.
