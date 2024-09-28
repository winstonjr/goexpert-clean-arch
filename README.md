# goexpert-clean-arch

## executar sistema

- rodar o docker para subir as dependências de infraestrutura `docker compose up --detach`
- Conectar no banco de dados e rodar o seguinte script:
```sql
CREATE TABLE orders (
    id varchar(255) NOT NULL,
    price float NOT NULL,
    tax float NOT NULL,
    final_price float NOT NULL,
    PRIMARY KEY (id)
);
```
- Executar `go mod tidy` para baixar as dependências
- Entrar no diretório `cmd/ordersystem`
- Executar usando o comanto `go run main.go wire_gen.go`

## Portas dos Serviços
- API HTTP: `:8000`
- API GRPC: `:50051`
- API GRAPHQL: `:8080`
- MySQL: `:3306`
- RabbitMQ Admin: `:15672`

## Testando GRPC com GRPCurl

### listando os métodos do servidor
```bash
grpcurl -plaintext localhost:50051 list
```

```bash
grpcurl -plaintext localhost:50051 list pb.OrderService
```

### Criando uma compra
```shell
grpcurl -plaintext -d @ localhost:50051 pb.OrderService/CreateOrder <<EOM
{
    "id":"xablau2",
    "price": 69.6,
    "tax": 1.3
}
EOM
```

### Listando compras
```shell
grpcurl -plaintext localhost:50051 pb.OrderService/ListOrder
```


## Testando GraphQL na Interface web

http://localhost:8080/

### Cadastrar orders
```graphql
mutation createOrder {
  createOrder(input: {id:"c", Price: 99, Tax: 1.6}) {
    id,
    Price,
    Tax,
    FinalPrice
  }
}
```

### listar orders
```graphql
query orders {
  orders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

## RabbitMQ
- Criar fila `orders`
- Criar bind `orders` com o `amq.direct`

## Gerando arquivos Proto
```shell
protoc --go_out=. --go-grpc_out=.  internal/infra/grpc/protofiles/order.proto
```

## Atualizando os arquivos graphql

### instalando o gqlgen
```shell
go get -d github.com/99designs/gqlgen
```

### atualizando com o novo schema
- Na mesma pasta que o arquivo `gqlgen.yml` 

```shell
go run github.com/99designs/gqlgen generate
```