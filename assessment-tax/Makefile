test:
	 go test -covermode=atomic ./... -coverprofile=coverage.out 
testout: 
	go tool cover -html=coverage.out 
testv:
	 go test -v ./tests/service  
	 go test -v ./tests/repository 
	go test -v ./tests/handler 

mocks-tests/service: ## Mocks
 	mockgen -source=../../repository/admin.go -destination=./mocks/admin_mock.go -package=mocks     
	mockgen -source=../../repository/tax.go -destination=./mocks/tax_mock.go -package=mocks
	mockgen -source=../../service/taxcsv.go -destination=./mocks/taxcsv_mock.go -package=mocks  
migrate-up:
	migrate -database "postgres://postgres:postgres@localhost:5432/ktaxes?sslmode=disable" -path . up