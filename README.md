### Docker mysql
```docker run --name wager-mysql -p 3306:3306 -e MYSQL_DATABASE=wager -e MYSQL_ROOT_PASSWORD=12345 -d mysql:latest```

### Mock 
```mockgen -source=domain/wager.go -destination=domain/mocks/wager.go -package=mocks```