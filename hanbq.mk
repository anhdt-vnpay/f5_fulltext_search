
gen_type:
	./schema/scripts/gen_api.sh

api: 
	CONFIG=$(CURDIR)/config/hanbq.yaml go run main.go api
	
test:
	go run main.go

.PHONY: all test clean
