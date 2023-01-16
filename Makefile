redis-up:
	docker run --rm --name redis-mm -d -p 6379:6379 redis:7-alpine

redis-down:
	docker stop redis-mm

redis-cli:
	docker exec -it redis-mm redis-cli

test:
	go test ./... -count=1

bench:
	go test ./... -run=NONE -bench=. -count=6 | tee bench-$(shell date +%s)
