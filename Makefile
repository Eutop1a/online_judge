# start the environment of online judge
.PHONY: start
start:
	docker-compose up -d

# stop the environment of online judge
.PHONY: stop
stop:
	docker-compose down

# run the submission
.PHONY: user
user:
	go run ./app/judgement/cmd/main.go -o ./judgement

# run the monolithic main function
.PHONY: main
main:
	go run .main.go

.PHONY: all
all: user main