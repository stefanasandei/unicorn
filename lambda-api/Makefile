ENTRY_DIR=./entry
WORKER_DIR=./worker
BUILD_DIR=./build

.PHONY: worker entry rabbitmq redis build

entry:
	cd ${ENTRY_DIR} && go run .

worker:
	cd ${WORKER_DIR} && go run .

rabbitmq:
	docker run -it --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.12-management

redis:
	docker run -it --rm --name redis -p 6379:6379 redis:latest

build:
	mkdir -p ${BUILD_DIR}
	go build -o ${BUILD_DIR}/entry ${ENTRY_DIR}
	go build -o ${BUILD_DIR}/worker ${WORKER_DIR}

test:
	go test ${ENTRY_DIR}/...
	go test ${WORKER_DIR}/...
