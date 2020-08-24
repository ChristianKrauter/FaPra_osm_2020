## Test all different options
echo "Tests on:" >> test.txt
echo $(date) >> test.txt
echo "--------------------" >> test.txt
echo "--------------------" >> test.txt

printf "Server tests (you have to terminate each server manually):\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "Server tests:" >> test.txt
go run osmgw.go >> test.txt
go run osmgw.go -bpip >> test.txt
go run osmgw.go -bg >> test.txt
go run osmgw.go -bg -bpip >> test.txt

printf "Dataprocessing tests:\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "Dataprocessing tests:" >> test.txt
go run osmgw.go -m=1 >> test.txt
go run osmgw.go -m=1 -lm >> test.txt
go run osmgw.go -m=1 -coastline >> test.txt
go run osmgw.go -m=1 -bpip >> test.txt
go run osmgw.go -m=1 -nbt >> test.txt
go run osmgw.go -m=1 -nbt -bpip >> test.txt

printf "Basic grid dataprocessing tests:\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "Basic grid dataprocessing tests:" >> test.txt
go run osmgw.go -m=1 -bg >> test.txt
go run osmgw.go -m=1 -bg -coastline >> test.txt
go run osmgw.go -m=1 -bg -bpip >> test.txt
go run osmgw.go -m=1 -bg -nbt >> test.txt
go run osmgw.go -m=1 -bg -nbt -bpip >> test.txt

printf "Eval dataprocessing tests:\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "Eval dataprocessing tests:" >> test.txt
go run osmgw.go -m=2 >> test.txt
go run osmgw.go -m=2 -lm >> test.txt
go run osmgw.go -m=2 -coastline >> test.txt
go run osmgw.go -m=2 -bpip >> test.txt
go run osmgw.go -m=2 -nbt >> test.txt
go run osmgw.go -m=2 -nbt -bpip >> test.txt

printf "Eval basic grid dataprocessing tests:\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "Eval basic grid dataprocessing tests:" >> test.txt
go run osmgw.go -m=2 -bg >> test.txt
go run osmgw.go -m=2 -bg -coastline >> test.txt
go run osmgw.go -m=2 -bg -bpip >> test.txt
go run osmgw.go -m=2 -bg -nbt >> test.txt
go run osmgw.go -m=2 -bg -nbt -bpip >> test.txt

printf "Eval wayfinding tests:\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "Eval wayfinding tests:" >> test.txt
go run osmgw.go -m=3 -r=100 -bg >> test.txt
echo "" >> test.txt
go run osmgw.go -m=3 -r=100 -bg -bpip >> test.txt
echo "" >> test.txt
go run osmgw.go -m=3 -r=100 >> test.txt
echo "" >> test.txt
go run osmgw.go -m=3 -r=100 -bpip >> test.txt
echo "" >> test.txt

printf "Eval readpbf test:\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "Eval readpbf test:" >> test.txt
go run osmgw.go -m=4 >> test.txt

printf "Testing done, see test.txt for the results.\n"
echo "" >> test.txt
echo "--------------------" >> test.txt
echo "--------------------" >> test.txt
echo "Testing done, see test.txt for the results." >> test.txt
echo "" >> test.txt
