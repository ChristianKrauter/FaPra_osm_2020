echo "--------------------"
echo "Making grids"
echo $(date)
echo "--------------------"

printf "ug 100x500"
go run osmgw.go -m=2 -x=100 -y=500 -f=planet-coastlines.pbf
printf "ug 1000x500"
go run osmgw.go -m=2 -x=1000 -y=500 -f=planet-coastlines.pbf
printf "ug 1000x1000"
go run osmgw.go -m=2 -x=1000 -y=1000 -f=planet-coastlines.pbf

printf "bg 360x360"
go run osmgw.go -m=2 -bg -x=360 -y=360 -f=planet-coastlines.pbf
printf "bg 1000x500"
go run osmgw.go -m=2 -bg -x=1000 -y=500 -f=planet-coastlines.pbf
printf "bg 1000x1000"
go run osmgw.go -m=2 -bg -x=1000 -y=1000 -f=planet-coastlines.pbf


echo "--------------------"
printf "Done making grids"
echo $(date)
echo "--------------------"
