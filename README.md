# Kubernetes Secret Manager

Um gerenciador de secrets para Kubernetes que pode ser usado tanto via CLI quanto como servidor HTTP.

## Instalação

```bash
# Clonar o repositório
git clone https://github.com/yourusername/k8s-secret-manager.git
cd k8s-secret-manager

# Construir o projeto
make build
```

## Uso

### CLI

```bash
# Criar um novo secret
./build/k8s-secret-manager create --name mysecret --namespace default --data "key1=value1,key2=value2"

# Atualizar um secret existente
./build/k8s-secret-manager update --name mysecret --namespace default --data "key1=newvalue1"

# Deletar um
```

go mod init github.com/yourusername/k8s-secret-manager
