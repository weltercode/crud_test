1. Start docker

2. ```docker-compose up --build```
DB init will be automaticly done in Dockerfile by this line "CMD ["sh", "-c", "./initDb && ./main"]"
If not, than run ```go run cmd/init/initDB.go```and than ```docker-compose down``` and again ```docker-compose up```

3. Go to ```http://localhost:8080/``` 