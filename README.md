# Serviço de Posts - Relaunch 

Microserviço responsável por gerenciar todas as funcionalidades relacionadas a posts, comentários e interações na plataforma Relaunch.

## 📋 Sobre o Projeto

Este repositório faz parte da arquitetura de microserviços da plataforma Relaunch, sendo responsável pelo gerenciamento completo do sistema de posts, incluindo criação, edição, comentários, respostas e sistema de likes.

## 🚀 Tecnologias Utilizadas

- **Go** - Linguagem de programação principal
- **MySQL** - Banco de dados relacional para persistência
- **gRPC** - Protocolo de comunicação entre microserviços
- **Protocol Buffers** - Serialização de dados

## 🏗️ Arquitetura

O serviço segue os princípios de Clean Architecture, organizando o código em camadas bem definidas:

```
service-post/
├── config/          # Configurações e variáveis de ambiente
├── handler/         # Handlers gRPC e injeção de dependências
├── repositories/    # Camada de acesso aos dados (MySQL)
├── resource/        # Recursos e transformadores de dados
└── server/          # Configuração do servidor gRPC
```

## 📦 Funcionalidades Principais

### Gerenciamento de Posts
- ✅ Criação de posts com suporte a imagens
- ✅ Atualização parcial de posts
- ✅ Listagem de posts (todos ou por autor)
- ✅ Exclusão com validação de autorização

### Sistema de Comentários
- ✅ Comentários em posts
- ✅ Respostas aninhadas (replies to replies)
- ✅ Exclusão com validação de permissões
- ✅ Listagem hierárquica de comentários

### Sistema de Likes
- ✅ Likes em posts
- ✅ Likes em comentários e respostas
- ✅ Toggle automático (adicionar/remover)
- ✅ Ordenação priorizando interações do usuário

## 🗄️ Estrutura do Banco de Dados

O serviço interage com as seguintes tabelas:

- `posts` - Armazenamento de posts
- `comments` - Comentários dos posts
- `comment_replies` - Respostas aos comentários
- `likes` - Likes em posts
- `comment_likes` - Likes em comentários e respostas

## 🛡️ Segurança

- Validação de autorização em operações sensíveis
- Verificação de propriedade antes de updates/deletes
- Queries parametrizadas para prevenção de SQL Injection
- Tratamento adequado de erros com códigos gRPC

## 🔗 Integração

Este serviço faz parte do ecossistema Relaunch e se comunica com outros microserviços através de gRPC, utilizando a biblioteca compartilhada `lib-relaunch-cot` para modelos e contratos comuns.

## 📝 Padrões de Código

- Uso de interfaces para desacoplamento
- Injeção de dependências
- Tratamento robusto de erros
- Código limpo e manutenível seguindo boas práticas Go
