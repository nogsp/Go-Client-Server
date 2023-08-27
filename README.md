# Go-Client-Server
To start RabbitMQ Container run
```bash
docker run -d --hostname my-rabbit --name some-rabbit rabbitmq:3
```
To start Mosquitto Container run
```bash
docker run --name my-mosquitto -it -p 1883:1883 -p 9001:9001 -v ${HOME}/mosquitto.conf:/mosquitto/config/mosquitto.conf eclipse-mosquitto
```