## Test all different options
printf "Starting tests:"
printf "\n===========================================================\n"

printf "\n\nServer tests (you have to terminate each server manually):"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go
printf "\n"
go run osmgw.go -bpip
printf "\n"
go run osmgw.go -bg
printf "\n"
go run osmgw.go -bg -bpip
printf "\n"

printf "\nDataprocessing tests:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=1
printf "\n\n"
go run osmgw.go -m=1 -lm
printf "\n\n"
go run osmgw.go -m=1 -coastline
printf "\n\n"
go run osmgw.go -m=1 -bpip
printf "\n\n"
go run osmgw.go -m=1 -nbt
printf "\n\n"
go run osmgw.go -m=1 -nbt -bpip
printf "\n\n"

printf "\nBasic grid dataprocessing tests:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=1 -bg
printf "\n\n"
go run osmgw.go -m=1 -bg -coastline
printf "\n\n"
go run osmgw.go -m=1 -bg -bpip
printf "\n\n"
go run osmgw.go -m=1 -bg -nbt
printf "\n\n"
go run osmgw.go -m=1 -bg -nbt -bpip
printf "\n\n"

printf "\nEval dataprocessing tests:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=2
printf "\n\n"
go run osmgw.go -m=2 -lm
printf "\n\n"
go run osmgw.go -m=2 -coastline
printf "\n\n"
go run osmgw.go -m=2 -bpip
printf "\n\n"
go run osmgw.go -m=2 -nbt
printf "\n\n"
go run osmgw.go -m=2 -nbt -bpip
printf "\n\n"

printf "\nEval basic grid dataprocessing tests:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=2 -bg
printf "\n\n"
go run osmgw.go -m=2 -bg -coastline
printf "\n\n"
go run osmgw.go -m=2 -bg -bpip
printf "\n\n"
go run osmgw.go -m=2 -bg -nbt
printf "\n\n"
go run osmgw.go -m=2 -bg -nbt -bpip
printf "\n\n"

printf "\nEval wayfinding tests:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=3 -r=10 -bg
printf "\n\n"
go run osmgw.go -m=3 -r=10 -bg -bpip
printf "\n\n"
go run osmgw.go -m=3 -r=10
printf "\n\n"
go run osmgw.go -m=3 -r=10 -bpip
printf "\n\n"
go run osmgw.go -m=3 -a=1 -r=10
printf "\n\n"
go run osmgw.go -m=3 -a=2 -r=10
printf "\n\n"
go run osmgw.go -m=3 -a=3 -r=10
printf "\n\n"

printf "Eval readpbf test:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=4
printf "\n"

printf "\n===========================================================\n"
printf "Testing done."
