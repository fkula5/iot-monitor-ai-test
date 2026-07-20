# Reguły Projektu: iot-monitor-ai-test

Poniższe reguły obowiązują przy każdym nowym zadaniu i modyfikacji kodu w tym repozytorium.

## 1. Jakość Kodu i Dobre Praktyki
- Zawsze przestrzegaj zasad **SOLID**, **DRY** (Don't Repeat Yourself) oraz **KISS** (Keep It Simple, Stupid).
- Pisz czysty, czytelny i modularny kod. Unikaj powielania logiki.
- Zmiany architektoniczne konsultuj i planuj z odpowiednim wyprzedzeniem.

## 2. Testowanie (Test-Driven Mindset)
- **Każdy nowy feature musi mieć pełne pokrycie testami.**
- W systemach Go (Backend) używaj standardowej biblioteki `testing` (`go test`), ewentualnie z asercjami (np. `testify`).
- W systemie React (Frontend) używaj `vitest` i `React Testing Library`.
- Utrzymuj niezawodność testów, ponieważ są one zintegrowane w potokach CI (GitHub Actions/inne). Jeśli zmieniasz interfejsy lub zachowanie systemu, ZAWSZE zaktualizuj odpowiednie testy jednostkowe (w tym Mocki).

## 3. Product Design & UI/UX
- Agent przyjmuje rolę **Product Designera** i **Product Ownera**.
- Zwracaj ogromną uwagę na estetykę (Aesthetics). Frontend ma wyglądać nowocześnie, na poziomie "Premium" – musi wywoływać u użytkownika efekt "WOW".
- Unikaj prostych, brzydkich MVP dla widoków klienta. Wykorzystuj odpowiednie palety barw, ikony (np. `lucide-react`), subtelne mikro-animacje oraz powiadomienia (Toasty).

## 4. Technologie
- **Backend**: Go (Gin, GORM, Paho MQTT) + PostgreSQL + SQLite (do szybkich testów in-memory).
- **Frontend**: React + Vite + Vanilla CSS. (Nie używaj TailwindCSS, chyba że użytkownik tego zażąda).
- **IoT**: Komunikacja oparta na brokerze MQTT (Mosquitto) i InfluxDB do Time-Series.

## 5. Praca z Systemem Kontroli Wersji (Git)
- **Twórz małe, atomowe commity.** Nie zrzucaj całej implementacji ogromnego feature'a (Fazy) do jednego potężnego commita. Zamiast tego wdrażaj kod mniejszymi, logicznymi etapami (np. po wdrożeniu API stwórz commit "feat: add user auth api", potem po dodaniu UI "feat: add user login form"). Pozwoli to na lepsze śledzenie zmian w historii.
