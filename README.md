1. Pokrenuti InfluxDB bazu naredbom: docker run --name influxdb -p 8086:8086 influxdb:2.2.0
2. Otvoriti projekt s Goland IDE-om (community edition free verzija)
3. Ispuniti env varijable unutar .env datoteke te kreirat trazene resurse unutar InfluxDB-a (bucket, token)
4. Otvoriti terminal u root projekta i skinuti dependencies s naredbom go mod download
5. Navigirati unutar cmd/apsim-api te pokrenuti main metodu




*NOTE*

- Potrebno je instalirati Go verziju 1.18 (API je zapravo pisan u ranijoj verziji 1.16 pa ne koristi najnovije feature kao sto su genericsi
, ali sam ga prebacio u 1.18 tako da se lakše pokrene) i GCC zbog SQLite drivera.

- Buildana CLI aplikacija mora biti unutar foldera apsim-cli u root projektu te apsim-stage-area direktorij mora postojati, hardkodirano u kodu. Ideja za proširivanje je dodavanje pathova u env. datoteku pa njih koristiti umjesto.

- Ostao je zakomentiran feature unutar YieldService-a koji umjesto value vrijednosti prinosa u influx db zapisuje json s vrijednosti prinosa i listom datuma.

- S obzirom, da je API worlds model otkud se dohvaćaju očitanja dosta nestabilan (nekada ne radi, format podataka se zna promjeniti, neki datumi se izostave, vrijednosti za pojedine datume se izostave itd...) dobra bi ideja bila pokrenuti API pa ručno provjeriti koja su se sva zadnja očitanja dodala u tablicu MicroclimateReadings kako se ne bi "zatrovala" krivim podacima jer će simulacije padati ili zakomentirati liniju pozivanja metode ScheduleBackgroundWorks() u main.go. U trenutku testiranja API-a radi.

- S obzirom na to da sam dodao tablicu naknadno, ako se misli prosirivat funkcionalnost aplikacije i ostaviti identičnima MicroclimateReading i PredictedMicroclimateReading bila bi dobra ideja refaktorirati kod i napraviti zajednicki interface s metodama koje vracaju atribute za MicroclimateReading i PredictedMicroclimateReading tako da se smanji broj funkcija u MicroclimateReadingService-u.

- Trajanje smulacije je ograničeno na 5 god, ali u ovom trenutačnom zadanom vremenu u kontekstu komotno se može izvrtiti i 10 god.

- Batch size-ove dohvaćanja MicroclimateReading i PredictedMicroclimateReading sam uskladio s veličinama kanala, 100, moguće je još veće povećanje kako bi se ubrzao proces.

- InfluxDB provjerava zapisivanje zadnjih 48h, inače se simulacija odvrti odmah, dodao sam to da testiranjem korištenje time series data, ako se hoće od početka provjeravati staviti range(start 0) u YieldRepository-u unutar flux querya.

- Unutar paketa models apsimx.go
- Unutar paketa playground su razne go datoteke unutar kojhi sam testirao funkcionalnosti i učio koncepte Go-a