redis-up:
	docker run --rm --name redis-mm -d redis:7-alpine

redis-down:
	docker stop redis-mm

redis-cli:
	docker exec -it redis-mm redis-cli

test:
	go test ./... -count=1
