.PHONY:
.SILENT:
.DEFAULT_GOAL := run
run:
	docker-compose up -d --remove-orphans app
stop: 
	docker-compose down -v
