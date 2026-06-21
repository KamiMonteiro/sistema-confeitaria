package main

import (
	"net/http"
	"sistema-confeitaria/database"
	"sistema-confeitaria/handler"
)

// corsHandler envolve um handler específico e libera CORS para o Live Server (127.0.0.1:5500).
// Usado em todas as rotas de API individualmente.
func corsHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		// Preflight do browser — responde OK e para por aqui
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

// corsMiddleware é aplicado no file server dos HTMLs (localhost:5500).
// Separado do corsHandler porque o file server usa http.Handler, não http.HandlerFunc.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Abre o banco SQLite e roda as migrations pra garantir que as tabelas existem
	db := database.ConnectDB()
	defer db.Close()
	database.RunMigrations(db)

	// ── Formas de pagamento ───────────────────────────────────────────
	// Handlers ficam em handler/auth.go
	http.HandleFunc("/api/novo/pagamento",      corsHandler(handler.CriarPagamento(db)))
	http.HandleFunc("/api/atualizar/pagamento", corsHandler(handler.AtualizarPagamento(db)))
	http.HandleFunc("/api/pagamento/listar/",   corsHandler(handler.ConsultarPagamento(db)))       // busca por ID — /api/pagamento/listar/{id}
	http.HandleFunc("/api/pagamento/buscar",    corsHandler(handler.BuscarPagamentoPorDescricao(db))) // filtro por descrição — ?descricao=
	http.HandleFunc("/api/pagamento/excluir/",  corsHandler(handler.ExcluirPagamento(db)))         // /api/pagamento/excluir/{id}
	http.HandleFunc("/api/todos/pagamento",     corsHandler(handler.BuscarTodasFormasPagamento(db)))

	// ── Usuários ─────────────────────────────────────────────────────
	// Handlers ficam em handler/auth.go (junto com Login e formas de pagamento)
	http.HandleFunc("/api/novo/usuario",       corsHandler(handler.CriarUsuario(db)))
	http.HandleFunc("/api/todos/usuario",      corsHandler(handler.BuscarTodosUsuario(db)))
	http.HandleFunc("/api/usuarios/buscar",    corsHandler(handler.BuscarUsuariosComFiltro(db))) // ?nome=
	http.HandleFunc("/api/usuarios/listar/",   corsHandler(handler.UsuarioPorID(db)))            // /api/usuarios/listar/{id}
	http.HandleFunc("/api/usuarios/excluir/",  corsHandler(handler.ExcluirUsuario(db)))          // /api/usuarios/excluir/{id}
	http.HandleFunc("/api/atualizar/usuarios", corsHandler(handler.AtualizarUsuario(db)))

	// ── Clientes ─────────────────────────────────────────────────────
	// Handlers ficam em handler/cliente.go
	http.HandleFunc("/api/novo/cliente",      corsHandler(handler.CriarCliente(db)))
	http.HandleFunc("/api/atualizar/cliente", corsHandler(handler.AtualizarCliente(db)))
	http.HandleFunc("/api/cliente/listar/",   corsHandler(handler.ClientePorID(db)))             // /api/cliente/listar/{id}
	http.HandleFunc("/api/cliente/buscar",    corsHandler(handler.BuscarClientesComFiltro(db)))  // ?nome=
	http.HandleFunc("/api/cliente/excluir/",  corsHandler(handler.ExcluirCliente(db)))           // /api/cliente/excluir/{id}
	http.HandleFunc("/api/todos/cliente",     corsHandler(handler.BuscarTodosClientes(db)))

	// ── Produtos ─────────────────────────────────────────────────────
	// Handlers ficam em handler/produto.go
	http.HandleFunc("/api/novo/produto",      corsHandler(handler.CriarProduto(db)))
	http.HandleFunc("/api/atualizar/produto", corsHandler(handler.AtualizarProduto(db)))
	http.HandleFunc("/api/produto/listar/",   corsHandler(handler.ProdutoPorID(db)))             // /api/produto/listar/{id}
	http.HandleFunc("/api/produto/buscar",    corsHandler(handler.BuscarProdutosComFiltro(db)))  // ?nome=
	http.HandleFunc("/api/produto/excluir/",  corsHandler(handler.ExcluirProduto(db)))           // /api/produto/excluir/{id}
	http.HandleFunc("/api/todos/produto",     corsHandler(handler.BuscarTodosProdutos(db)))

	// ── Pedidos ──────────────────────────────────────────────────────
	// Handlers ficam em handler/pedido.go
	// Criar e atualizar pedido usa transação (pedido + itens juntos) — ver repository/pedido.go
	http.HandleFunc("/api/novo/pedido",      corsHandler(handler.CriarPedido(db)))
	http.HandleFunc("/api/atualizar/pedido", corsHandler(handler.AtualizarPedido(db)))
	http.HandleFunc("/api/pedido/listar/",   corsHandler(handler.PedidoPorID(db)))               // /api/pedido/listar/{id} — retorna pedido + itens
	http.HandleFunc("/api/pedido/buscar",    corsHandler(handler.BuscarPedidosComFiltro(db)))    // ?status=
	http.HandleFunc("/api/pedido/excluir/",  corsHandler(handler.ExcluirPedido(db)))             // /api/pedido/excluir/{id}
	http.HandleFunc("/api/todos/pedido",     corsHandler(handler.BuscarTodosPedidos(db)))

	// ── Auth ─────────────────────────────────────────────────────────
	http.HandleFunc("/api/auth/login", corsHandler(handler.Login(db)))

	// ── Dashboard ─────────────────────────────────────────────────────
	// Retorna KPIs (pedidos hoje, faturamento, produtos ativos, clientes) + últimos 5 pedidos
	http.HandleFunc("/api/dashboard", corsHandler(handler.BuscarDashboard(db)))

	// ── Frontend ─────────────────────────────────────────────────────
	// Serve os arquivos da pasta /html na raiz — acessado via localhost:8080
	// CSS e JS ficam em /css e /js (fora da /html), por isso o Live Server é usado no dev
	fs := http.FileServer(http.Dir("./html"))
	http.Handle("/", corsMiddleware(fs))

	http.ListenAndServe(":8080", nil)
}
