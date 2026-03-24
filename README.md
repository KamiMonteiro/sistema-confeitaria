# 🌹 Rosas & Tortas Confeitaria — Frontend

Sistema de gestão interno da confeitaria. Interface em HTML + Tailwind CSS, preparada para integração com backend **Go**.

---

## 📁 Estrutura de arquivos

```
rosas-tortas/
│
├── css/
│   └── global.css          ← Estilos compartilhados (botões, inputs, sidebar, tabela, modal...)
│
├── js/
│   ├── app.js              ← Funções globais: sidebar, modal de confirmação, toast, apiFetch()
│   └── tailwind.config.js  ← Referência das cores e tokens (não é importado, é inline em cada HTML)
│
├── login.html              ← Tela de autenticação
├── dashboard.html          ← Tela principal com KPIs e últimos pedidos
├── produtos-lista.html     ← Listagem de produtos com filtros, paginação e exclusão
├── produto-form.html       ← Formulário de cadastro/edição de produto (?modo=novo | ?modo=editar&id=X)
│
└── README.md               ← Este arquivo
```

---

## 🔗 Navegação entre telas

| De                    | Para                              | Como                          |
|-----------------------|-----------------------------------|-------------------------------|
| `login.html`          | `dashboard.html`                  | Botão "Entrar" (após auth)    |
| `dashboard.html`      | `produtos-lista.html`             | Card "Gerenciar Produtos"     |
| `produtos-lista.html` | `produto-form.html?modo=novo`     | Botão "Novo produto"          |
| `produtos-lista.html` | `produto-form.html?modo=editar&id=X` | Ícone de lápis na linha    |
| `produto-form.html`   | `produtos-lista.html`             | Botão "Voltar" ou após salvar |

---

## 🔌 Integração com o backend Go

Todos os pontos de integração estão marcados com comentários `// ── Integração Go ──`.

### Base URL
Definida no topo de `js/app.js`:
```js
const API_BASE = 'http://localhost:8080/api';
```

### Autenticação (JWT)
O token é salvo em `localStorage` com a chave `rt_token` e enviado automaticamente como `Authorization: Bearer <token>` pela função `apiFetch()`.

### Função apiFetch
```js
// Exemplo de uso em qualquer tela:
const produtos = await apiFetch('/produtos');
await apiFetch('/produtos', { method: 'POST', body: JSON.stringify(payload) });
await apiFetch(`/produtos/${id}`, { method: 'DELETE' });
```

### Endpoints esperados (sugestão)
| Método | Endpoint              | Descrição                        |
|--------|-----------------------|----------------------------------|
| POST   | `/auth/login`         | Login → retorna `{ token }`      |
| GET    | `/dashboard/kpis`     | KPIs do dia                      |
| GET    | `/pedidos?limit=5`    | Últimos pedidos                  |
| GET    | `/produtos`           | Lista todos os produtos          |
| GET    | `/produtos/:id`       | Busca produto por ID             |
| POST   | `/produtos`           | Cria novo produto                |
| PUT    | `/produtos/:id`       | Atualiza produto                 |
| DELETE | `/produtos/:id`       | Remove produto                   |

### Modelo PRODUTO (JSON esperado)
```json
{
  "id_produto":     1,
  "nome":           "Torta de Morango",
  "preco_unitario": 89.90,
  "descricao":      "Torta recheada com morango fresco.",
  "categoria":      "Tortas",
  "ativo":          "SIM"
}
```

---

## 🎨 Design System

| Token        | Valor      | Uso                          |
|--------------|------------|------------------------------|
| `rose`       | `#F4A7B9`  | Cor principal (botões, ativo)|
| `rose-deep`  | `#C9546E`  | Textos e ênfase              |
| `rose-light` | `#FAD4DE`  | Fundos suaves, badges        |
| `mint`       | `#A8D5BA`  | Cor secundária (filtrar)     |
| `mint-deep`  | `#4A9B72`  | Badge "Ativo"                |
| Fundo body   | `#FFF8F5`  | Creme suave                  |

**Fontes:**
- Títulos: `Playfair Display` (Google Fonts)
- Corpo: `DM Sans` (Google Fonts)

---

## ▶️ Como rodar localmente

Abrir diretamente no navegador:
```
rosas-tortas/login.html
```

> **Obs:** Para que o `app.js` e o `global.css` sejam carregados corretamente, prefira servir via servidor local:
```bash
# Python 3
python3 -m http.server 3000
# Acessar: http://localhost:3000/login.html
```

**Credenciais de demonstração (mock):**
- E-mail: `admin@loja.com`
- Senha: `123456`
