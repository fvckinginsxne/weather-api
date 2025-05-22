# weather-api
REST API для получения текущей погоды по городу

## Запуск
### 1. Клонировать репозиторий
```
git clone https://github.com/fvckinginsxne/lyrics-library.git
```
### 2. Установка конфига (необходимо получить ключ доступа для [OpenWeatherMap](https://openweathermap.org/api))
```
cp .env.example .env
nano .env 
```
### 3. Запуск приложения
```
go run ./cmd/app/main.go -config=.env
```
### 4. Документация доступна по адресу
```
localhost:8080/swagger/index.html
```
