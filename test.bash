## Test all different options
printf "Starting tests:"
printf "\n===========================================================\n"

printf "\n\nTest Server (you have to terminate each server manually):"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go
printf "\n"
go run osmgw.go -bpip
printf "\n"
go run osmgw.go -bg
printf "\n"
go run osmgw.go -bg -bpip
printf "\n"

printf "\nTest Dataprocessing ug:"
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

printf "\nTest dataprocessing bg:"
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

printf "\nTest Eval of ug dataprocessing:"
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

printf "\nTest eval of bg dataprocessing:"
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

printf "\nTest eval of wayfinding on ug:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=3 -r=10 -a=0
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=0 -bpip
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=1
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=1 -bpip
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=2
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=2 -bpip
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=3
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=3 -bpip

printf "\nTest eval of wayfinding on bg:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=3 -r=10 -a=0 -bg
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=0 -bg -bpip
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=1 -bg
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=1 -bg -bpip
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=2 -bg
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=2 -bg -bpip
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=3 -bg
printf "\n\n"
go run osmgw.go -m=3 -r=10 -a=3 -bg -bpip

printf "Test eval of readpbf:"
printf "\n-----------------------------------------------------------\n"
go run osmgw.go -m=4
printf "\n"

printf "\n===========================================================\n"
printf "Testing done."
