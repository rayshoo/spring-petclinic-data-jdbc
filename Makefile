GIT_VER=$(shell git describe --tags)
.DEFAULT_GOAL := all

all: build down pull run log
.PHONY:all

build:
	./gradlew jib
.PHONY:build

down:
	docker-compose down
.PHONY:run

pull:
	docker-compose pull
.PHONY:build

run:
	docker-compose up -d
.PHONY:run

log: logs
.PHONY:log

logs:
	docker-compose logs -f
.PHONY:logs
