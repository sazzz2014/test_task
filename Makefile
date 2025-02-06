# Makefile для управления проектом

# Запуск сервера и клиента
up:
	docker-compose up --build

# Запуск клиента
run-client:
	docker-compose run client


# Запуск тестов
test:
	cd server && go test ./internal/server -v