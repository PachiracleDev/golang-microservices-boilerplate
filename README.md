
# Golang Microservices Boilerplate

Este es el codigo base para construir aplicaciones escalables de alta concurrencia usando RabbitMQ, Fibber y Protobuffers.



## Co levantar el proyecto

#### Crear un archivo .env para inicializar tus variables, usar el .env.example como ejemplo

```bash
  cp .env.example .env
```

#### Levantar las instancias de docker para el RabbitMQ Y Mongodb

```bash
  docker compose up -d
```


#### Instalar las dependencias

```bash
  go get
```
 
#### Levantar los microservicios el API Gateway

#### API Gateway
```bash
   make run-api-service
```

#### Microservicio de auth

```bash
   make run-auth-service
```
