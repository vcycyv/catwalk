docker volume create --name postgres-data -d local
docker volume create --name mongo-data -d local
docker-compose up -d