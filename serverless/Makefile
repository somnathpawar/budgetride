build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/cabBooking cabBooking/main.go cabBooking/cabs.go cabBooking/lyft.go cabBooking/uber.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
