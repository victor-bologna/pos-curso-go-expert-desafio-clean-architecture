# Desafio Clean Architecture

Desafio usando Webserver Go, gRPC e GraphQL em uma aplicação com a intenção de criar pedidos (com ID, preço e juros) e retornar o preço final do pedido e listar todos os pedidos gerados nas 3 APIs.

## Tutorial

1. Usando CMD na pasta do projeto com [Docker](https://www.docker.com/) instalado na máquina, executar `docker compose up -d`;
2. Para configurar o [RabbitMQ](https://hub.docker.com/_/rabbitmq):
    1. Acessar http://localhost:15672/ e usar login e senha configurado no docker-compose.yml
    2. Ir na aba `Queues and Streams` e criar nova fila (queue) com nome `orders`
    3. Ir em `Exchanges` > `amq.direct` > `Bind` para a queue `orders`
3. Para rodar o servidor é necessario ir a pasta via CMD `cmd/server` e rodar `go run main.go wire_gen.go`
4. Para rodar os serviços (Necessário o servidor estar rodando):
    1. Webserver REST - Porta 8081:
        1. Na pasta do projeto, acessar a pasta `test`e executar as URLs desejadas (GET e POST orders)
    2. [gRPC](https://grpc.io/) (Necessário instalar [evans](https://github.com/ktr0731/evans?tab=readme-ov-file#installation)) - Porta 50051:
        1. accesar via CMD o comando na pasta principal do projeto: `evans -r repl`
        2. Especificar package: `package pb`
        3. Especificar service: `OrderService`
        4. Chamar serviços `call CreateOrder` ou `call GetAllOrders`
    3. [GraphQL](https://gqlgen.com/) - Porta 8080:
        1. Acessar via Browser a URL: http://localhost:8080/
        2. No editor a esquerda colocar o seguinte code snippet:
        ```
        mutation CreateOrder {
            createOrder(input: {
                ID: "423432rew",
                Price: 3334.4,
                Tax: 45.49
            }) {
                ID
                Price
                Tax
                FinalPrice
            }
        }

        query GetOrders {
            getOrders{
                Orders{
                Price
                Tax
                FinalPrice
                }
            }
        }
        ```
        3. Clicar no play e executar um dos serviços