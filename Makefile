start:
	docker-compose -f docker-compose.yml up -d --build

stop:
	docker-compose -f docker-compose.yml down --volumes

test:
	docker-compose -f docker-compose.yml up -d --build
	-sleep 15
	-go test ./... -v
	docker-compose -f docker-compose.yml down --volumes
