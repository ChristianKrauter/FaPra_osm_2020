# Ocean Wayfinding on OSM Data (FaPra OSM 2020)
Repository for 'Fachpraktikum OpenStreetMap Daten' 2020 UNI Stuttgart by Jonas Vogelsang and Christian Krauter.

Branch of Christian Krauter. </br>
Tested with go version go1.14.2 windows/amd64

# Instructions
In the main directory run 'go run osmgw.go' or osmgw.exe</br>
Open a browser at 'lohalhost:8081'

## CLI Options
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

# To-Do

1. Improve routing through North Pole
2. Fix routing around South Pole
3. Improve Performance by ordering JPS neighbour directions (counter) clockwise
4. Implement [Improved JPS](https://www.researchgate.net/publication/287338108_Improving_jump_point_search)
5. ...

# Screenshots

## Canals
<img alt="A*" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/canals.png?raw=true"/>
<img alt="A*" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/canals_2.png?raw=true"/>

## Algorithms (The images show the expanded nodes.)

### Dijkstra
<img alt="Dijkstra" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/alg_dij.png?raw=true"/>

### Bidirectional Dijkstra
<img alt="Bidirectional Dijkstra" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/alg_bi_dij.png?raw=true"/>

### A-Star
<img alt="A-Star" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/alg_astar.png?raw=true"/>

## Bidirectional A-Star
<img alt="Bidirectional A-Star" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/alg_bi_astar.png?raw=true"/>

## A-Star Jump-Point-Search
<img alt="A-Star Jump-Point-Search" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/alg_astar_jps.png?raw=true"/>

Comparison of expanded nodes vs.

<img alt="jps pq" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/jps_pq.png?raw=true"/>

"viewed" nodes

<img alt="jps viewed" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/jps_viewed.png?raw=true"/>

For a jump point search explanation please see
[Online graph pruning for pathfinding on grid maps](https://dl.acm.org/doi/10.5555/2900423.2900600).

# Problems

## Bidirectional

### Bidirectional algorithms don't always find the same route as dijkstra and A*.

<img alt="bi problem 1" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/bi_prob.png?raw=true"/>

### For some routes even dijkstra/A* find different routes depending on the direction.

<img alt="bi problem 2" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/bi_prob_2.png?raw=true"/>
<img alt="bi problem 3" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/bi_prob_3.png?raw=true"/>

## Jump-Point-Search

### Routing through the North Pole

<img alt="jps problem 1" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/jps_prob.png?raw=true"/>
<img alt="jps problem 2" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/jps_prob_2.png?raw=true"/>

### Routing around the South Pole

<img alt="jps problem 3" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/jps_prob_3.png?raw=true"/>
<img alt="jps problem 4" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/jps_prob_4.png?raw=true"/>

# Evaluation

## Small Grid (129600 nodes)

### Speed

<img alt="small speed" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/ug_129600_speed.jpg?raw=true"/>

### Priority Queue Pops

<img alt="small pq" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/ug_129600_pqpops.jpg?raw=true"/>

### Route Equality

<img alt="small optimality" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/ug_129600_equality.jpg?raw=true"/>

## Big Grid (1000000 nodes)

### Speed

<img alt="big speed" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/ug_1000000_speed.jpg?raw=true"/>

### Priority Queue Pops

<img alt="big pq" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/ug_1000000_pqpops.jpg?raw=true"/>

### Route Equality

<img alt="big optimality" src="https://github.com/ChristianKrauter/FaPra_osm_2020/blob/chris/images/ug_1000000_equality.jpg?raw=true"/>
