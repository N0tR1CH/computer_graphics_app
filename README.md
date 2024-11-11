## [Ćwiczenie 3 - Formaty plików, format PPM](https://cez.wi.pb.edu.pl/mod/assign/view.php?id=1566)

Zaimplementować w dowolnym języku programowania:

1. Zapis i wczytywanie formatu PBM (1 bit na piksel) w formie tekstowej (P1) i binarnej (P4)
2. Zapis i wczytywanie formatu PGM (8 bitów na piksel) w formie tekstowej (P2) i binarnej (P5)
3. Zapis i wczytywanie formatu PPM (24 bity na piksel) w formie tekstowej (P3) i binarnej (P6)
4. Zaimplementować obsługę komentarzy w formatach tekstowych (zaczynające się od #)

Dobre wyjaśnienie poszczególnych formatów znajdziecie państwo na [wikipedii](https://en.wikipedia.org/wiki/Netpbm#File_formats). Implementacja powinna zostać stworzona bez użycia gotowych bibliotek do wczytywania plików PPM.

Przetestować program na obrazach testowych. **Obrazy testowe powinny zostać prawidłowo wczytane! Obrazy testowe powinny być wczytane w mniej niż 8 sekund!**

Na ocenę maksymalną należy zaimplementować 3 podpunkty z 4 oraz kolejkowanie zdarzeń. Robiąc w grupach należy zaimplementować wszystkie podpunkty.

Asynchroniczne kolejkowanie zdarzeń polega na tym, że:

1. Kod jest wywoływany bez zacinania aplikacji

2. Należy stworzyć klasę która kolejkuje zdarzenia i je sekwencyjnie wywołuje.

3. Można kolejkować komendy do asynchronicznej kolejki i je na koniec wywoływać, tak więc kliknięcie przycisku nie wywołuje bezpośrednio metody tylko dodaje do kolejki komendę która jest przetwarzana jeśli nic innego aktualnie nie było robione. W kodzie powinno to być widoczne tym, że wszystkie wywołania przycisków są tworzone poprzez wzorzec [Command](http://xn--dobry%20pocztek%20architektury%20mvvm-0ee.xn--%20jeli%20chcecie%20si%20pastwo%20w%20tym%20temacie%20rozwija%20to%20proponuj%20z%20klasy%20mainviewmodel%20wycign%20rzeczy%20odpowiedzialne%20za%20rysowanie%20na%20ekranie%20do%20klasy%20drawingservice%20i%20uycie%20jej%20wewntrz%20vm,%20alternatywnie%20mona%20do%20kadego%20trybu%20typu%20triangle%20ellipse%20da%20oddzielny%20viewmodel%20z%20metod%20draw%20ktry%20byby%20wywoywany%20w%20mainviewmodel,%20to%20bardzo%20poprawi%20jako%20kodu-645bn05aca7isf50ayb7mteo46avb237jja46cj25axma901cnbm./) które po uzyskaniu od użytkownika inputu/danych wejściowych (na przykład nazwy pliku) są wrzucane w kolejkę i wykonują się sekwencyjnie.


Ad 1)
Zapis formatu PBM odbywa się przez biblioteke rozszerzającą interfejs standardowej biblioteki natomiast wcztywanie zostało napisane od zera (do binarnych plików została wykonana pomoc w postaci modelu LLM).

Ad 2), Ad 3)
Tu również sprawa ma się jak w Ad 1)

Ad 4)
Podczas zapisywania formaty netpbm mają możliwosć dodawania komenatrzy poprzez interfejs użytkownika, zapisany w formacie plik canvasu w bitmapie będzie posiadał komentarze :)

Operacje asynchroniczne kolejkowania zostały wykonane dzięki wykorzystaniu kanałów w języku Golang. Klasę zastępuje odpowiednia struktura danych gdyż nie jest to język zorientowany obiektowo jest tu raczej zastosowane data oriented approach a przynajmniej moje próby jego stosowania.

Podsumowując wszystkie wymagania zostały zrealizowane. (Kwestia cofania zdarzeń może zostać zaimplementowana poprzez usunięcie zbindowanego elementu usuwającego ostatni kształt - to spowoduje aktualizację interfejsu, reset canvasu i callbacki odwzorywujące stan elementów, które pozostały, gdyż jak wspominałem 1 ćwiczeniu jest to canvas deklaracyjny, więc jeżeli istnieje taka potrzeba to mogę to zaimplementować, to kwestia 1 małej funkcji podpiętej pod przycisk np. cofania)
