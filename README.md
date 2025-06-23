# T1-SO: Simulador de Sistema Operacional

Este projeto implementa um simulador de sistema operacional com gerenciamento de processos, memória paginada, escalonamento, interrupções, shell interativo e interface web para observabilidade. O objetivo é fornecer uma plataforma didática para estudar conceitos fundamentais de SO.

## Sumário

- [Visão Geral](#visão-geral)
- [Principais Funcionalidades](#principais-funcionalidades)
  - [1. Gerenciamento de Processos](#1-gerenciamento-de-processos)
  - [2. Gerenciamento de Memória (Paginação e Swap)](#2-gerenciamento-de-memória-paginação-e-swap)
  - [3. Escalonador Round-Robin com Quantum](#3-escalonador-round-robin-com-quantum)
  - [4. Simulação de CPU e Interrupções](#4-simulação-de-cpu-e-interrupções)
  - [5. Shell Interativo](#5-shell-interativo)
  - [6. Observabilidade via WebSocket e Interface Web](#6-observabilidade-via-websocket-e-interface-web)
  - [7. Programas de Exemplo](#7-programas-de-exemplo)
- [Como Executar](#como-executar)
- [Comandos do Shell](#comandos-do-shell)
- [Estrutura do Projeto](#estrutura-do-projeto)

---

## Visão Geral

O simulador abstrai os principais componentes de um SO: CPU, Kernel, Gerenciador de Processos, Gerenciador de Memória, Escalonador, Disco (swap), além de prover uma interface de linha de comando (shell) e uma interface web para visualização do estado do sistema.

---

## Principais Funcionalidades

### 1. Gerenciamento de Processos

- **Criação, execução e destruição de processos**.
- Cada processo possui um PCB (Process Control Block) com PID, status, programa, registradores, PC, quantum usado, etc.
- Suporte a múltiplos estados: ADDED, READY, RUNNING, BLOCKED, FINISHED.
- Processos podem ser bloqueados/desbloqueados por operações de I/O ou page-fault.

### 2. Gerenciamento de Memória (Paginação e Swap)

- **Memória paginada**: cada processo possui sua própria tabela de páginas.
- **Alocação sob demanda**: apenas a primeira página do processo é carregada na RAM ao criar o processo.
- **Tratamento de page-fault**: páginas não presentes na RAM são carregadas sob demanda.
- **Swap (Disco)**: páginas podem ser movidas para o disco quando não há quadros livres, utilizando política FIFO para substituição.
- **Gerenciamento de quadros livres** e atualização das tabelas de páginas.

### 3. Escalonador Round-Robin com Quantum

- **Escalonador Round-Robin**: alterna entre processos prontos, garantindo justiça.
- **Quantum configurável**: cada processo executa por um número limitado de ciclos antes de ser preemptado.
- **Gerenciamento automático de quantum** e troca de contexto.

### 4. Simulação de CPU e Interrupções

- **Ciclo de clock configurável**.
- **Execução de instruções** dos programas dos processos.
- **Interrupções simuladas**:
  - Tempo (quantum expirado)
  - I/O (pedido e resposta)
  - Término do processo
- **Tratamento centralizado de interrupções** e atualização do estado dos processos.

### 5. Shell Interativo

- **Interface de linha de comando** para interação com o SO simulado.
- Permite criar, executar, matar processos, listar processos, inspecionar memória, responder I/O, alterar nível de log, etc.
- Comandos disponíveis:
  - `help`
  - `create <programa>`
  - `kill <pid>`
  - `ps -a` / `ps -p <pid>`
  - `dump -p <pid>` / `dump -memory <start> <end>`
  - `exec <pid>` / `exec -a`
  - `log <info|trace|debug>`
  - `toggle-ce`
  - `io`
  - `clear`
  - `exit`

### 6. Observabilidade via WebSocket e Interface Web

- **Servidor WebSocket** integrado.
- **Interface web** (em `static/index.html`) exibe em tempo real:
  - Estado da CPU (PC, registradores, bits de interrupção)
  - Memória (visualização paginada)
  - Tabela de páginas do processo atual
  - Processos prontos e seus estados
- **Atualização automática** a cada ciclo de clock.

### 7. Programas de Exemplo

- Diversos programas de exemplo para simular diferentes padrões de execução:
  - Programas que manipulam registradores, fazem operações matemáticas, string, e I/O.
  - Exemplo de programa que solicita entrada do usuário via I/O.

---

## Como Executar

1. **Requisitos**: Go instalado.
2. **Compilar**: `go build -o main.exe ./src`
3. **Executar**: `./main.exe`
4. **Acessar a interface web**: [http://localhost:8080](http://localhost:8080)

---

## Comandos do Shell

| Comando                        | Descrição                                                        |
|------------------------------- |------------------------------------------------------------------|
| `help`                         | Mostra todos os comandos disponíveis                             |
| `create <programa>`            | Cria um novo processo a partir de um programa                    |
| `kill <pid>`                   | Mata o processo com o PID especificado                           |
| `ps -a`                        | Lista todos os processos                                         |
| `ps -p <pid>`                  | Mostra informações do processo com o PID especificado            |
| `dump -p <pid>`                | Mostra o PCB do processo                                         |
| `dump -memory <start> <end>`   | Mostra o conteúdo da memória entre os endereços indicados        |
| `exec <pid>`                   | Executa o processo com o PID especificado                        |
| `exec -a`                      | Executa todos os processos prontos                               |
| `log <info|trace|debug>`       | Altera o nível de log                                            |
| `toggle-ce`                    | Alterna o modo de execução contínua                              |
| `io`                           | Lista e responde a pedidos de I/O pendentes                      |
| `clear`                        | Limpa a tela                                                    |
| `exit`                         | Sai do shell                                                    |

---

## Estrutura do Projeto

- `src/`
  - **cpu.go**: Simulação da CPU, execução de instruções, tratamento de interrupções.
  - **kernel.go**: Inicialização do kernel, gerenciamento global, pedidos de I/O.
  - **memory-manager.go**: Gerenciamento de memória paginada, swap, page-fault.
  - **memory.go**: Estrutura da memória principal.
  - **process-manager.go**: Gerenciamento de processos e PCBs.
  - **programs.go**: Programas de exemplo.
  - **scheduler.go**: Escalonador round-robin.
  - **shell.go**: Implementação do shell interativo.
  - **ws.go**: Comunicação WebSocket para observabilidade.
  - **static/**: Interface web (HTML/CSS/JS).

---

## Observações

- O projeto é modular e extensível, facilitando a inclusão de novas features.
- Ideal para fins didáticos, trabalhos de SO e experimentação com conceitos de sistemas operacionais.

---

Se precisar de exemplos de uso, detalhes de implementação ou quiser expandir o README com prints ou GIFs, é só pedir!