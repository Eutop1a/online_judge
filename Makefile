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
	go build -o judgement ./app/judgement/cmd/main.go

# run the monolithic main function
.PHONY: main
main:
	go build -o online-judge main.go

.PHONY: all
all: user main