## Test bounding tree performance
printf "Starting nbt eval:"
printf "\n======================\n"

printf "\n\nug 100x500\n"
go run osmgw.go -m=2 -x=100 -y=500 -f=planet-coastlines.pbf
printf "\n\nug 100x500 NBT\n"
go run osmgw.go -m=2 -x=100 -y=500 -f=planet-coastlines.pbf -nbt
printf "\n\nug 1000x1000\n"
go run osmgw.go -m=2 -x=1000 -y=1000 -f=planet-coastlines.pbf
printf "\n\nug 1000x1000 NBT"
go run osmgw.go -m=2 -x=1000 -y=1000 -f=planet-coastlines.pbf -nbt

printf "\n======================\n"
printf "Done."