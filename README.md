## A service that receives information about the weather in the cities added to the subscription

### To run an app:

```
make build && make run
```

If you are running the application for the first time, you must apply the migrations to the database:

```
make migrate
```

Create swagger documentation:

```
make swag
```

### API Usage Example:

Add city to the list of subscription (current weather information is also added):

```
curl -X POST http://localhost:8000/weather/ --header "Content-Type: application/json" --data '{"name":"Moscow"}
```
```
{"city":"Moscow"}
```

Get list of all cities in subscription:

```
curl -X DEL http://localhost:8000/weather/ --header "Content-Type: application/json" --data '{"name":"Moscow"}
```
```
{"data":[{"id":1,"city":"Moscow","subscription_date":"2022-10-10T00:00:00Z"}
```

Get average temperature in city (actually accumulated data):

```
curl http://localhost:8000/weather/Moscow
```

```
6
```

Delete city from the list of subscription and all weather's information on this city:

```
curl -X DELETE http://localhost:8000/weather/Moscow 
```
```
{"status":"ok"}
```

### How to config

```cntDayArchive``` - the number of days after which the data is transferred to the archive

```periodicity``` - frequency (in minutes) of adding new data

```unit``` - temperature units (Celsius(C), Fahrenheit(F))