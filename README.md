# Proxy Reverso em Go

Um Load Balancer e Proxy Reverso de alta performance escrito em Go, com suporte a HTTP/2, monitoramento de saúde (Health Check) e recarregamento dinâmico de configurações.

## Configuração

O projeto utiliza um diretório centralizado para armazenar as configurações. O caminho padrão depende do sistema operacional:

- **Linux**: `~/.config/teste-proxy/`
- **Windows**: `%AppData%\teste-proxy\` (geralmente `C:\Users\NomeUsuario\AppData\Roaming\teste-proxy\`)

### Estrutura de Diretórios
```text
teste-proxy/
├── main.json             # Configurações globais do servidor
└── [servicos]/           # Subdiretórios para cada grupo de proxies
    └── sites.json        # Arquivos de configuração dos proxies
```

### 1. Configuração Global (`main.json`)
Este arquivo define os limites do servidor e caminhos de certificados.

```json
{
  "readTimeoutSeconds": 10,
  "writeTimeoutSeconds": 10,
  "idleTimeoutSeconds": 60,
  "readHeaderTimeoutSeconds": 5,
  "maxHeaderMB": 1,
  "httpsOn": true,
  "certFilePath": "/caminho/para/cert.pem",
  "keyFilePath": "/caminho/para/key.pem"
}
```

### 2. Configuração de Proxies
Você pode criar subdiretórios dentro de `teste-proxy/` para organizar seus proxies. Cada arquivo JSON dentro desses subdiretórios define um conjunto de regras.

**Exemplo: `teste-proxy/web-apps/api.json`**
```json
[
  {
    "serverName": "api.exemplo.com",
    "prefix": "/v1",
    "loadBalancer": "round-robin",
    "servers": [
      {
        "url": "http://10.0.0.1:8080",
        "weight": 5,
        "available": true
      },
      {
        "url": "http://localhost:8080",
        "weight": 2,
        "available": true
      }
    ]
  }
]
```

## Load Balancers

O proxy suporta três tipos de algoritmos para distribuir o tráfego:

1.  **`round-robin`**: Distribui as requisições sequencialmente entre os servidores disponíveis.
2.  **`weighted-round-robin`**: Distribui com base no peso (`weight`) definido para cada servidor. Um peso maior recebe mais tráfego.
3.  **`random`**: Escolhe um servidor aleatoriamente de forma distribuída.

## Health Check (Monitoramento de Saúde)

O sistema possui um verificador de integridade ativo que:
- Executa a cada **5 segundos**.
- Realiza uma requisição do tipo `HEAD` para cada servidor configurado.
- Se um servidor não responder ou retornar erro (>= 400), ele é marcado como **OFFLINE**.
- Servidores OFFLINE são removidos automaticamente do rodízio do Load Balancer em tempo real.
- Assim que o servidor volta a responder corretamente, ele retorna ao estado **ONLINE** e volta a receber tráfego.

## Segurança e TLS

### Como habilitar HTTPS
1. Defina `"httpsOn": true` no seu `main.json`.
2. Forneça os caminhos absolutos para seu certificado (`certFilePath`) e chave privada (`keyFilePath`).
3. O servidor iniciará automaticamente na porta **443** (HTTPS) e **80** (HTTP).

### O porquê é seguro?
- **HTTP/2 & h2c**: Suporte nativo para o protocolo moderno, garantindo maior velocidade e multiplexação.
- **Timeouts Controlados**: Permite configurar tempos limites de leitura, escrita e cabeçalhos para prevenir ataques de Slowloris.
- **Isolamento de Configuração**: As configurações são lidas em tempo de execução, permitindo atualizações sem reiniciar o binário principal (Hot Reload).

## Considerações finais
- Este projeto não deve ser usado de forma alguma em produção, ele é apenas um estudo tanto da linguagem quanto do funcionamento de um proxy. portanto ele não é um software completo.
