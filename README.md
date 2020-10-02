# FaPra OSM 2020
Repository for 'Fachpraktikum OpenStreetMap Daten' 2020 UNI Stuttgart by Jonas Vogelsang and Christian Krauter.

Branch of Christian Krauter. </br>
Tested with go version go1.14.2 windows/amd64

# Instructions
In the main directory run 'go run osmgw.go' or osmgw.exe</br>
Open a browser at 'lohalhost:8081'
## CLI options
```
-m int
      Select Mode:
        0: Start server
        1: Evaluate grid creation
        2: Evaluate wayfinding
        3: Evaluate reading pbf
        4: Evaluate ug neighbours
        5: Test routes and neighbours
        6: Add canals to grid
-x int
      Grid size in x direction. (default 1000)
-y int
      Grid size in y direction. (default 1000)
-r int
      Number of runs for wayfinding evaluation. (default 1000)
-f string
      Name of the pbf file inside data/ (default "antarctica-latest.osm.pbf")
-bg
      Create a basic (non-uniform) grid.
-nbt
      Do not use a tree structure for the bounding boxes. (default true)
-bpip
      Use the basic 2D point in polygon test.
-coastline
      Create coastline geoJSON.
-lm
      Use memory efficient method to read unpruned pbf files.
-n string
      Additional note for evaluations.
```