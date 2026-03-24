/* =============================================
   Rosas & Tortas Confeitaria — JS Compartilhado
   ============================================= */

// ---------- NAVEGAÇÃO ----------
function navegarPara(page) {
  window.location.href = page;
}

// ---------- TOAST ----------
function showToast(msg, tipo = 'sucesso') {
  const el = document.getElementById('toast');
  if (!el) return;
  document.getElementById('toast-msg').textContent = msg;
  document.getElementById('toast-icon').className =
    tipo === 'sucesso'
      ? 'fa-solid fa-circle-check text-green-500'
      : 'fa-solid fa-circle-exclamation text-red-400';
  el.classList.remove('hidden');
  setTimeout(() => el.classList.add('hidden'), 3000);
}

// ---------- MODAL GENÉRICO ----------
function abrirModal(id)  { document.getElementById(id)?.classList.remove('hidden'); }
function fecharModal(id) { document.getElementById(id)?.classList.add('hidden'); }

// ---------- FORMATAÇÃO ----------
function formatarMoeda(valor) {
  return 'R$ ' + Number(valor).toFixed(2).replace('.', ',');
}

// ---------- DATA ATUAL ----------
function formatarDataHoje() {
  return new Date().toLocaleDateString('pt-BR', {
    weekday: 'long', day: '2-digit', month: 'long', year: 'numeric'
  });
}

// ---------- SIDEBAR: MARCAR ITEM ATIVO ----------
document.addEventListener('DOMContentLoaded', () => {
  const page = document.body.dataset.page;
  if (!page) return;
  document.querySelectorAll('.nav-item[data-nav]').forEach(el => {
    if (el.dataset.nav === page) el.classList.add('active');
  });
});
