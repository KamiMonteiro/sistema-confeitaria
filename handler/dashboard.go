package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sistema-confeitaria/repository"
)

// BuscarDashboard retorna os KPIs e os últimos pedidos em uma única chamada.
// Rota: GET /api/dashboard
func BuscarDashboard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		stats, err := repository.BuscarEstatisticasDashboard(db)
		if err != nil {
			http.Error(w, "Erro ao carregar dashboard", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	}
}
