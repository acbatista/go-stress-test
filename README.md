# Ferramenta de Teste de Carga em Go

Uma ferramenta de linha de comando escrita em Go para realizar testes de carga em serviços web. Esta ferramenta permite que você especifique o número de requisições e conexões simultâneas para testar o desempenho dos seus serviços web.

## Funcionalidades

- Número configurável de requisições totais
- Nível ajustável de concorrência
- Relatório detalhado com distribuição de códigos de status HTTP
  - Categorização por tipo de resposta (1xx, 2xx, 3xx, 4xx, 5xx)
  - Descrições detalhadas dos códigos de status
  - Contagem de requisições por código
- Medição do tempo total de execução
- Rastreamento de erros

## Compilação e Execução

### Compilação Local

```bash
# Compilar a aplicação
go build -o stress-test

# Executar a aplicação
./stress-test --url=http://seu-servico.com --requests=1000 --concurrency=10
```

### Compilação e Execução com Docker

```bash
# Construir a imagem Docker
docker build -t stress-test .

# Executar o teste de carga usando Docker
docker run stress-test --url=http://seu-servico.com --requests=1000 --concurrency=10
```

## Parâmetros

- `--url`: URL do serviço a ser testado (obrigatório)
- `--requests`: Número total de requisições a serem feitas (obrigatório)
- `--concurrency`: Número de requisições simultâneas (obrigatório)

## Exemplo de Saída

```
Iniciando teste de carga para http://exemplo.com
Total de requisições: 1000
Nível de concorrência: 10

Resultados do Teste:
Tempo total: 5.234s
Total de requisições: 1000

Distribuição de Códigos de Status:

Sucesso (2xx):
  HTTP 200 (OK - Requisição bem sucedida): 985 requisições
  HTTP 201 (Created - Recurso criado com sucesso): 5 requisições

Redirecionamento (3xx):
  HTTP 301 (Moved Permanently - Recurso movido permanentemente): 2 requisições

Erro do Cliente (4xx):
  HTTP 404 (Not Found - Recurso não encontrado): 5 requisições
  HTTP 429 (Too Many Requests - Muitas requisições): 1 requisição

Erro do Servidor (5xx):
  HTTP 500 (Internal Server Error - Erro interno do servidor): 2 requisições

Requisições com falha: 0
```

## Observações

- O nível de concorrência não pode ser maior que o número total de requisições
- Todos os parâmetros são obrigatórios
- A ferramenta exibirá mensagens de erro se alguma requisição falhar
- Os códigos de status HTTP são agrupados por categoria para melhor visualização
- Cada código de status inclui uma descrição detalhada do seu significado

## Estrutura do Projeto

```
.
├── main.go         # Código principal da aplicação
├── go.mod          # Definição do módulo Go
├── Dockerfile      # Arquivo para construção da imagem Docker
└── README.md       # Esta documentação
```

## Como Funciona

A ferramenta utiliza goroutines do Go para realizar requisições HTTP de forma concorrente. O processo funciona da seguinte maneira:

1. A aplicação valida os parâmetros de entrada fornecidos
2. Cria um pool de workers baseado no nível de concorrência especificado
3. Distribui as requisições entre os workers
4. Coleta os resultados de todas as requisições
5. Categoriza e agrupa os códigos de status HTTP
6. Gera um relatório detalhado com as estatísticas

## Requisitos

- Go 1.21 ou superior (para compilação local)
- Docker (para execução em container)

## Tratamento de Erros

A ferramenta lida com diversos tipos de erros, incluindo:
- Parâmetros inválidos ou ausentes
- Falhas de conexão
- Timeouts
- Respostas HTTP inesperadas

Todos os erros são registrados e contabilizados no relatório final, com categorização apropriada dos códigos de status HTTP. 