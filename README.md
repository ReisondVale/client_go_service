# API REST de Gerenciamento de Clientes

Este repositório contém uma API REST desenvolvida em Go para o gerenciamento de registros de clientes. A aplicação oferece funcionalidades como inserção, listagem de todos os clientes registrados em ordem alfabética, busca de clientes por nome e a verificacão de existência ou não de cpf_cnpj já cadastrados. Inclui validação de CPF/CNPJ e integração com PostgreSQL para armazenamento de dados, além disso, possui funcionalidade para inserção de clientes através de um arquivo csv.

## Funcionalidades

- **Gerenciamento de Clientes:**
  - Criar registros de clientes.
  - Listar todos os clientes.
  - Busca de clientes por nome.
  - Inserção de clientes através de um arquivo csv.

- **Validação:**
  - Validar números de CPF e CNPJ.

- **Infraestrutura:**
  - Aplicação conteinerizada com Docker.
  - PostgreSQL como banco de dados.

- **Saúde do Servidor:**
  - Endpoint para verificar o status do servidor.

## Requisitos

- Go 1.23.2
- Docker e Docker Compose
- PostgreSQL
- Curl para teste dos endpoints

## Como Começar

### 1. Clone o repositório:
```bash
git clone <repository-url>
cd <repository-directory>
```

### 2. Configure as variáveis de ambiente:

Edite o arquivo `.env` no diretório raiz se preferir:
```env
DATABASE_URL=postgres://user:password@db:5432/client_go_service?sslmode=disable
```

### 3. Compile e execute a aplicação, utilize o Makefile para facilitar execução:

#### Usando Docker:
compile a aplicação:
```bash
make up
```
execute a migration:
```bash
make migrate
```
carregue e insira a base de dados dos clientes do arquivo clients.csv:
```bash
make insert_csv_clients
```

### 4. Acesse a API:

- URL Base: `http://localhost:8080`

## Endpoints da API

### Gerenciamento de Clientes

#### 1. Adicionar um cliente:
- **Endpoint:** `POST /clients`
- **Corpo da Requisição utilizando o curl:**
  ```bash
  curl -X POST http://localhost:8080/clients \
    -H "Content-Type: application/json" \
    -d '{
      "Name": "Jane Doe",
      "CPF_CNPJ": "98765432100"
    }'
  ```
- **Resposta:**
  ```bash
  {"message":"Client created successfully"}%
  ```

#### 2. Listar todos os clientes:
- **Endpoint:** `GET /clients`
- **Corpo da Requisição utilizando o curl:**
  ```bash
  curl -X GET http://localhost:8080/clients
  ```
- **Resposta:**
  A resposta será uma lista com todos os clientes inseridos no banco de dados.
  ```bash
  [{"id":50994,"name":"Jane Doe","cpf_cnpj":"98765432100","blocklist":false}]%
  ```

#### 3. Verificar a existência de um cpf_cnpj cadastrado:
- **Endpoint:** `GET /clients/exists/:cpfCnpj`
- **Corpo da Requisição utilizando o curl:**
  ```bash
  curl -X GET "http://localhost:8080/clients/exists/98765432100"
  ```
- **Resposta:**
  ```bash
  {"exists":true}%
  ```

#### 4. Busca de clientes por nome:
- **Endpoint:** `GET /clients/search`
- **Corpo da Requisição utilizando o curl:**
  ```bash
  curl "http://localhost:8080/clients/search?name=jane"
  ```
- **Resposta:**
  ```bash
  [{"id":50994,"name":"Jane Doe","cpf_cnpj":"98765432100","blocklist":false}]%
  ```

### Saúde do Servidor

#### Verificar o Status do Servidor:
- **Endpoint:** `GET /status`
- **Corpo da Requisição utilizando o curl:**
  ```bash
  curl "http://localhost:8080/status"
  ```
- **Resposta:**
  A resposta retorna a quantidade de requisições feitas, e o tempo em segundos que a aplicção está em funcionamento:
  ```bash
  {"request_count":3,"uptime_seconds":77.157762828}%
  ```

## Regras de Validação

- A validação de CPF/CNPJ garante que o número esteja devidamente formatado e seja válido conforme os padrões brasileiros.

## Em construção:

- Testes unitários em construção.
- Validação através do cpf_cnpj se o cliente está em alguma blocklist.
