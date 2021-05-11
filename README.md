# ![ATC logo](https://raw.githubusercontent.com/speedyhoon/TriumphChallenge/main/winres/32x32.png) All Triumph Challenge
[![go report card](https://goreportcard.com/badge/github.com/speedyhoon/TriumphChallenge)](https://goreportcard.com/report/github.com/speedyhoon/TriumphChallenge)

A simple commandline program to collate event results for the **All Triumph Challenge** motorsport series. Its mission is to streamline result finalisation and make the process easier and more enjoyable for competitors & volunteers.

[![All Triumph Challenge](https://raw.githubusercontent.com/speedyhoon/TriumphChallenge/master/tsoa.png)](https://www.tsoavic.com.au/)

## Sample Results
```
   All Triumph Challenge - Sports Car Track Day
Competitors: 4
Pos     # Driver       Qualify     Secs      Fastest     Secs      Slowest     Secs       Slow Ave    Percentage    Runs   Laps
1st    42 Joe Bloggs   1m4.825s    64.8250   1m4.401s    64.4010   2m40.145s   160.1450   112.48500   57.25296706      3     31
2nd   512 Sam Smith    1m10.226s   70.2260   1m10.226s   70.2260   1m30.972s   90.9720     80.59900   87.13011328      1     16
3rd    47 Jack Black   1m6.113s    66.1130   1m4.931s    64.9310   1m25.592s   85.5920     75.85250   85.60166112      1     14
4th   513 Jac Jones    1m0.347s    60.3470   1m0.347s    60.3470   1m30.707s   90.7070     75.52700   79.90122738      1      6
```

## Instructions
- Locate the event on [Natsoft racing results](http://racing.natsoft.com.au/results/)
- Copy the results. `Ctrl + A` then `Ctrl + C` on Windows or `Command + A` then `Command + C` on Mac.
- Start the TriumphChallenge program
- The program will detect the [Natsoft racing results](http://racing.natsoft.com.au/results/) in your clipboard if present. Otherwise, it will prompt for input.\
  Press `Enter` to continue once the results are copied into your clipboard.
- Type in the list of competitors racing numbers separated by a space.\
  For example: `1 881 4 55 92 5 7 9 13 43`
- Press `Enter`
- Results will be calculated and saved in HTML, Text and XLSX spreadsheet files with a copy of the [Natsoft racing results](http://racing.natsoft.com.au/results/) and competitor list.


## Results Formula
Fastest lap time **÷** ((Slowest lap time **+** Qualifying lap time) **÷** 2) **×** 100

For example:
\
1:04.8932 **÷** ((1:13.7864 **+** 1:08.564) **÷** 2) **×** 100
\
Is converted to seconds:
\
64.8932 **÷** ((73.7864 **+** 68.564) **÷** 2) **×** 100
\
Then is calculated:
\
64.8932 **÷** (142.3504 **÷** 2) **×** 100
\
64.8932 **÷** 71.1752 **×** 100
\
= **91.1739**


## Event Rules
- Laps completed during the practice session don't count towards the total quantity of laps completed.
- The **Fastest lap time** and **Slowest lap time** is calculated for each competitor excluding their practice session.
- The **Qualifying lap time** is calculated for each competitor only during their practice session.
- The highest possible score is **100** and the lowest possible score is **0**.
- Competitors are sorted by:
  - The quantity of runs/sessions they complete in descending order (most sessions completed first)
  - Their Percentage result in descending order (highest number first)
  - Their quantity of laps completed in descending order (most laps wins).
  - If two or more competitors have the same result then both are assigned that position, for example: `1st, =2nd, =2nd, 4th, 5th, etc ...`
- The first lap of each run/session is ignored to allow for competitors to line up on the starting grid.
- Competitors who fail to complete a lap time during qualifying won't be eligible for any placing. 
- Decimal numbers are calculated and sorted using 64 bit precision.
- Percentage results in HTML and text format are displayed with 8 decimal places. Spreadsheet format uses built-in formulas to display decimal numbers (precision varies between software).

## Permissions
The application may need permission granted to allow execution. This is used for:
- Saving files to disk and
- Opening [Natsoft racing results](http://racing.natsoft.com.au/results/) in your default browser.

## Build
Built using [go-winres](https://github.com/tc-hib/go-winres)
```
go-winres make
go build
```
