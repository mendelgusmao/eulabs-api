# Eulabs Product API

## Visão Geral

A API de produtos oferece endpoints para gerenciar produtos, incluindo listagem, recuperação, criação, atualização e remoção de produtos. A autenticação é feita através de tokens JWT.

### URL Base para Testes

- `https://ludovic.fawn-beaver.ts.net/api/v1`

### Workspace do Postman

- [Eulabs-ProductAPI.postman](postman/Eulabs-Products-API.postman_collection)
- [Development Environment](postman/Development.postman_environment.json)
- [Production Environment](postman/Production.postman_environment.json)

## Autenticação

Todos os endpoints - exceto a autenticação - requerem um token Bearer. Inclua o token no cabeçalho `Authorization`:

```
Authorization: Bearer <seu_token_jwt>
```

## Endpoints

### Autenticar Usuário

#### Requisição

- **URL**: `/users/authenticate`
- **Método**: `POST`
- **Content-Type**: `application/json`

#### Corpo da requisição

```json
{
  "username": "seu_usuario",
  "password": "sua_senha"
}
```

#### Respostas

- **200 OK**: Retorna um token JWT.

  ```json
  {
    "token": "seu_token_jwt"
  }
  ```

- **401 Unauthorized**: Credenciais inválidas.
- **422 Unprocessable Entity**: Corpo da requisição mal formado.
- **500 Internal Server Error**: Erro no servidor.

### Listar produtos

#### Requisição

- **URL**: `/products`
- **Método**: `GET`
- **Content-Type**: `application/json`
- **Authorization**: `Bearer <seu_token_jwt>`

#### Respostas

- **200 OK**: Retorna uma lista de produtos.

  ```json
  [
    {
      "id": 1,
      "name": "Product 1",
      "description": "Descrição do produto 1",
      "price": 10.0,
      "quantity": 100,
      "category": "Categoria 1",
      "brand": "Marca 1"
    }
  ]
  ```

- **204 No Content**: Não há conteúdo para retornar.
- **500 Internal Server Error**: Erro no servidor.

### Obter produto

#### Requisição

- **URL**: `/products/:id`
- **Método**: `GET`
- **Content-Type**: `application/json`
- **Authorization**: `Bearer <seu_token_jwt>`

#### Respostas

- **200 OK**: Retorna um produto.

  ```json
  {
    "id": 1,
    "name": "Product 1",
    "description": "Descrição do produto 1",
    "price": 10.0,
    "quantity": 100,
    "category": "Categoria 1",
    "brand": "Marca 1"
  }
  ```

- **404 Not Found**: O produto não foi encontrado.
- **500 Internal Server Error**: Erro no servidor.

### Criar produto

#### Requisição

- **URL**: `/products`
- **Método**: `POST`
- **Content-Type**: `application/json`
- **Authorization**: `Bearer <seu_token_jwt>`

#### Corpo da requisição

  ```json
  {
    "id": 1,
    "name": "Product 1",
    "description": "Descrição do produto 1",
    "price": 10.0,
    "quantity": 100,
    "category": "Categoria 1",
    "brand": "Marca 1"
  }
  ```

#### Respostas

- **201 Created**: O produto foi criado.

  ```json
  {
    "id": 1,
    "name": "Product 1",
    "description": "Descrição do produto 1",
    "price": 10.0,
    "quantity": 100,
    "category": "Categoria 1",
    "brand": "Marca 1"
  }
  ```

- **400 Bad Request**: Formato do corpo da requisição inválido.
- **422 Unprocessable Entity**: Formato da requisição incorreto ou falha na validação dos dados.
- **500 Internal Server Error**: Erro no servidor.

### Atualizar produto

#### Requisição

- **URL**: `/products/:id`
- **Método**: `PUT`
- **Content-Type**: `application/json`
- **Authorization**: `Bearer <seu_token_jwt>`

#### Corpo da requisição

  ```json
  {
    "name": "Product 1",
    "description": "Descrição do produto 1",
    "price": 10.0,
    "quantity": 100,
    "category": "Categoria 1",
    "brand": "Marca 1"
  }
  ```

#### Respostas

- **200 OK**: O produto foi atualizado.

  ```json
  {
    "id": 1,
    "name": "Product 1",
    "description": "Descrição do produto 1",
    "price": 10.0,
    "quantity": 100,
    "category": "Categoria 1",
    "brand": "Marca 1"
  }
  ```

- **400 Bad Request**: Formato do id não é numérico ou o corpo da requisição é inválido.
- **422 Unprocessable Entity**: Formato da requisição incorreto ou falha na validação dos dados.
- **404 Not Found**: O produto não foi encontrado.
- **500 Internal Server Error**: Erro no servidor.

### Apagar produto

#### Requisição

- **URL**: `/products/:id`
- **Método**: `DELETE`
- **Authorization**: `Bearer <seu_token_jwt>`

#### Respostas

- **201 No Content**: O produto foi apagado.
- **400 Bad Request**: Formato do id não é numérico ou o corpo da requisição é inválido.
- **404 Not Found**: O produto não foi encontrado.
- **500 Internal Server Error**: Erro no servidor.

## Executar a API com Docker Compose

Para executar a aplicação localmente usando Docker Compose, siga os passos abaixo:

1. Certifique-se de ter o Docker e Docker Compose instalados.

2. Clone o repositório da aplicação:

   ```bash
   git clone https://github.com/mendelgusmao/eulabs-api.git
   cd eulabs-api
   ```

3. Execute o seguinte comando para iniciar a aplicação:

   ```bash
   docker-compose up
   ```

4. Aguarde a aplicação e o banco de dados serem iniciados e o banco de dados ser populado.

5. Acesse a API através de `http://localhost:8080`.