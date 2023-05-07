include .env
export

GIT_VER=$(shell git describe --tags)
.DEFAULT_GOAL := all

all: build down pull run log
.PHONY:all

base:
	docker-compose build base && docker-compose push base
.PHONY:base

build: base
	docker-compose build --no-cache was_builder && docker-compose up was_builder
	# cd was && ./gradlew jib
.PHONY:build

down:
	docker-compose down
.PHONY:run

pull:
	docker-compose pull mysql was
.PHONY:build

run: up
.PHONY:run

up: down pull
	docker-compose up -d mysql was
.PHONY:up

log: logs
.PHONY:log

logs:
	docker-compose logs -f mysql was
.PHONY:logs