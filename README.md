# FaPra OSM 2020
Repository for 'Fachpraktikum OpenStreetMap Daten' 2020 UNI Stuttgart by Jonas Vogelsang and Christian Krauter

# Instructions
## Schein (Tasks 1-6)
### Tasks 1-4
    1. Put planet-coastline.pbf into '/data'
        - For other pbf files please change readPBF.go:206
    2. Run 'go run readPBF.go' or readPBF.exe
    3. You can find the resulting in '/schein/temp'

### Tasks 5 & 6
    1. Run 'go run server.go' or server.exe
    2. Visit 'localhost:8081'
    
## Final Project (Task 7)
    1. Run go 'run osmgw.go' or osmgw.exe (tested with version: go1.14.2 windows/amd64)
    2. Visit 'localhost:8081'
