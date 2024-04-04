<p align="center">
    <img src="https://i.imgur.com/w2EFlUe.png" width="128px" height="128px" alt="Ikona bežca"/>
    <h1 align="center">Proboj</h1>
</p>

## Čo je to Proboj?

Proboj je aktivita, ktorej cieľom je naprogramovať vlastného bota to počítačovej hry. Boti následne proti sebe súťažia.

Štandardný Proboj sa skladá z niekoľkých častí:
- **server:** program, ktorý získava od hráčov ich ťahy, posiela im zmeny v hre a pod.
- **observer:** program, ktorý vie prehrávať replay-e hier
- **runner:** program, ktorý zabezpečuje spúšťanie a management hráčskych procesov, ich komunikáciu s ostatnými časťami 
systému
- **bot:** samotný hráč, ktorý hrá hru

Tento repozitár sa primárne zameriava na funkčnosť **runnera** a špecifikáciu komunikácie v Proboji.

## Protokol medzi Runnerom a Serverom

Runner komunikuje so Serverom prostredníctvom stdin/stdout vo forme správ. Každá správa má nasledujúci formát:

```
HEADER (1 line)
PAYLOAD (0+ lines)
.
```

Na začiatku hry Runner spustí proces Servera a pošle mu úvodnú konfiguráciu. Header `CONFIG`, payload sa skladá z riadku
so zoznamom hráčov oddelených medzerou a potom môže nasledovať niekoľko riadkov, ktoré načítavame z game configu `args`.
Zároveň Runner pospúšťa všetkých hráčov.

Od tohto momentu je kontrola nad priebehom hry na strane Serveru, ktorý príkazmi ovláda Runner.

### Runner príkazy

Server posiela Runneru príkazy, pomocou ktorých ovláda hru. Formát príkazov je:

```
COMMAND [args]
DATA
.
```

Odpoveď od Runnera je vo formáte:

```
STATUS
DATA
.
```

Kde status je jedna z hodnôt `OK`, `ERROR`, `DIED`. Dáta obsahujú výsledok operácie, resp. chybu, môžu sa skladať z viacerých riadkov.

- `TO PLAYER player comment...` \
    Príkaz pošle hráčovi `player` celý obsah dát na stdin. Ak je uvedená hodnota `comment`, vypíše sa do hráčovho logu.
- `READ PLAYER player` \
    Prečíta dáta od hráča `player` až po prvý riadok `.` a pošle ich ako dáta odpovede.
- `TO OBSERVER` \
    Zapíše celý obsah dát do observer logu.
- `KILL PLAYER player` \
    Ukončí proces hráča `player` (pošle mu `SIGKILL`).
- `PAUSE PLAYER player` \
    Pozastaví proces hráča `player` (funguje iba na Linuxe/MacOS, pošle `SIGSTOP`).
- `RESUME PLAYER player` \
    Obnoví beh procesu hráča `player` (funguje iba na Linuxe/MacOS, pošle `SIGCONT`).
- `SCORES` \
    Očakáva formát dát ako riadky `player score` (score je `int`). Tieto dáta uloží do súboru `score.json` v priečinku 
    hry vo formáte `{player: score}`.
- `END` \
    Oznámi Runnerovi, že hra skončila. Runner pozabíja zvyšné procesy (vrátane Servera) a poupratuje po sebe.

## Komunikácia medzi hráčom a serverom

Server vie komunikovať pomocou `TO PLAYER` a `READ PLAYER` príkazov. Obsah týchto dát nie je pre Runnera zaujímavý, je
na implementácií hry, aké dáta posiela. Hráč dostáva dáta na stdin a odosiela dáta na stdout.

Hráč obdrží dáta ukončené riadkom `.`, rovnako sa očakáva, že dáta od hráča budú ukončené týmto riadkom.

Ak hráč po zavolaní `READ PLAYER` neukončí svoje dáta `.` do času nastaveného v `timeout` configu, Runner mu zabije
proces. 

Čokoľvek, čo hráč vypíše na stderr sa uloží do jeho log súboru.

## Konfigurácia Runnera

Runner chce dva konfiguračné súbory -- `config.json` a `games.json`.

`config.json` hovorí o nastaveniach runnera. Aktuálne obsahuje tieto voľby:

- `server`: cesta k binárke servera 
- `players`: mapovanie názov hráča → cesta k jeho binárke a jazyk (`command`, `language`)
- `processes_per_player`: počet procesov na hráča (default 1, ak viac, tak sa vyrobia `player_0`, `player_1`...)
- `timeout`: maximálny čas, ktorý môže hráč využiť pri čakaní na výstup v sekundách, pre každý jazyk (mapping jazyk -> timeout)
- `disable_logs`: ak je `true`, deaktivuje ukladanie logov hráčov a servera
- `disable_gzip`: ak je `true`, deaktivuje gzipovanie logov
- `game_root`: cesta, kde sa budú ukladať dáta hier

`games.json` obsahuje definíciu hier:

- `gamefolder`: názov priečinka hry, ktorý sa vyrobí v `game_root`
- `players`: zoznam hráčov, ktorý sa v hre vyskytnú
- `args`: dáta, ktoré dostane server v `CONFIG`u

## Spoločné Go knižnice

Ak sa rozhodneš programovať Server v jazyku Go, v tomto repozitári nájdeš balík `client`, ktorý obsahuje všetko potrebné
na komunikáciu s Runnerom.
