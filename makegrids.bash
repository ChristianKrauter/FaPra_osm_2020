## Make all grids
printf "Starting grid creation:"
printf "\n======================\n"

printf "\n\nug 100x500\n"
go run osmgw.go -m=2 -x=100 -y=500 -f=planet-coastlines.pbf -nbt
printf "\n\nug 360x360\n"
go run osmgw.go -m=2 -x=360 -y=360 -f=planet-coastlines.pbf -nbt
printf "\n\nug 360x360 bpip\n"
go run osmgw.go -m=2 -x=360 -y=360 -f=planet-coastlines.pbf -bpip -nbt
printf "\n\nug 1000x500"
go run osmgw.go -m=2 -x=1000 -y=500 -f=planet-coastlines.pbf -nbt
printf "\n\nug 1000x1000"
go run osmgw.go -m=2 -x=1000 -y=1000 -f=planet-coastlines.pbf -nbt

printf "\n\nbg 360x360"
go run osmgw.go -m=2 -bg -x=360 -y=360 -f=planet-coastlines.pbf -nbt
printf "\n\nbg 360x360 bpip"
go run osmgw.go -m=2 -bg -x=360 -y=360 -f=planet-coastlines.pbf -bpip -nbt
printf "\n\nbg 1000x500"
go run osmgw.go -m=2 -bg -x=1000 -y=500 -f=planet-coastlines.pbf -nbt
printf "\n\nbg 1000x1000"
go run osmgw.go -m=2 -bg -x=1000 -y=1000 -f=planet-coastlines.pbf -nbt


printf "\n======================\n"
printf "Done."
