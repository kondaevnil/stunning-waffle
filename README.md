# Запуск

### Клонирование
mkdir -p vk/ecom
cd vk/ecom

go mod init vk/ecom

go get github.com/lib/pq

### Запуск с Docker Compose
make docker-up

### Или без Makefile
docker-compose up -d --build


# Тестирование

### Запуск всех тестов
make -f Makefile.tests test

### Запуск только юнит-тестов
make -f Makefile.tests test-unit

### Запуск только интеграционных тестов
make -f Makefile.tests test-integration

### Запуск с покрытием кода
make -f Makefile.tests test-coverage

### Запуск с race detection
make -f Makefile.tests test-verbose

## Прямые команды Go

### Все тесты
go test ./tests/... -v

### Юнит-тесты
go test ./tests/unit/... -v

### Интеграционные тесты
go test ./tests/integration/... -v
