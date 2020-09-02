## Run wayfinding eval for all algorithms
R=${1:-100}
printf "Starting eval of wayfinding with $R runs:"
printf "\n===========================\n"

printf "\n\nUniform Grids:"
printf "\n\n100x500 ug"
printf "\n---------------------------\n"
go run osmgw.go -m=3 -a=0 -r=$R -x=100 -y=500
printf "\n\n"
go run osmgw.go -m=3 -a=1 -r=$R -x=100 -y=500
printf "\n\n"
go run osmgw.go -m=3 -a=2 -r=$R -x=100 -y=500
printf "\n\n"
go run osmgw.go -m=3 -a=3 -r=$R -x=100 -y=500
printf "\n"

printf "\n\n1000x500 ug"
printf "\n---------------------------\n"
go run osmgw.go -m=3 -a=0 -r=$R -x=1000 -y=500
printf "\n\n"
go run osmgw.go -m=3 -a=1 -r=$R -x=1000 -y=500
printf "\n\n"
go run osmgw.go -m=3 -a=2 -r=$R -x=1000 -y=500
printf "\n\n"
go run osmgw.go -m=3 -a=3 -r=$R -x=1000 -y=500
printf "\n"

printf "\n\n1000x1000 ug"
printf "\n---------------------------\n"
go run osmgw.go -m=3 -a=0 -r=$R -x=1000 -y=1000
printf "\n\n"
go run osmgw.go -m=3 -a=1 -r=$R -x=1000 -y=1000
printf "\n\n"
go run osmgw.go -m=3 -a=2 -r=$R -x=1000 -y=1000
printf "\n\n"
go run osmgw.go -m=3 -a=3 -r=$R -x=1000 -y=1000
printf "\n"


printf "\n\nBasic Grids:"
printf "\n\n1000x500 ug"
printf "\n---------------------------\n"
go run osmgw.go -m=3 -a=0 -r=$R -x=1000 -y=500 -bg
printf "\n\n"
go run osmgw.go -m=3 -a=1 -r=$R -x=1000 -y=500 -bg
printf "\n\n"
go run osmgw.go -m=3 -a=2 -r=$R -x=1000 -y=500 -bg
printf "\n\n"
go run osmgw.go -m=3 -a=3 -r=$R -x=1000 -y=500 -bg
printf "\n"

printf "\n\n1000x1000 ug"
printf "\n---------------------------\n"
go run osmgw.go -m=3 -a=0 -r=$R -x=1000 -y=1000 -bg
printf "\n\n"
go run osmgw.go -m=3 -a=1 -r=$R -x=1000 -y=1000 -bg
printf "\n\n"
go run osmgw.go -m=3 -a=2 -r=$R -x=1000 -y=1000 -bg
printf "\n\n"
go run osmgw.go -m=3 -a=3 -r=$R -x=1000 -y=1000 -bg
printf "\n"


printf "\n===========================\n"
printf "Done."
