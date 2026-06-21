/* =============================================
   Rosas & Tortas Confeitaria — JS Compartilhado
   Funções usadas em todas as páginas do sistema.
   Carregado via <script src="../js/app.js"> no final do <body>.
   ============================================= */

// ---------- NAVEGAÇÃO ----------

// Redireciona para outra página — usado em botões que não são links <a>
function navegarPara(page) {
  window.location.href = page;
}

// ---------- TOAST ----------

// Exibe a notificação flutuante no canto da tela por 3 segundos.
// tipo: 'sucesso' (verde) ou 'erro' (vermelho) — padrão é sucesso.
// O elemento #toast precisa existir no HTML da página.
function showToast(msg, tipo = 'sucesso') {
  const el = document.getElementById('toast');
  if (!el) return; // página sem toast, ignora silenciosamente

  document.getElementById('toast-msg').textContent = msg;

  // Troca o ícone conforme o tipo — ícone vem do Font Awesome
  document.getElementById('toast-icon').className =
    tipo === 'sucesso'
      ? 'fa-solid fa-circle-check text-green-500'
      : 'fa-solid fa-circle-exclamation text-red-400';

  el.classList.remove('hidden');

  // Some depois de 3 segundos automaticamente
  setTimeout(() => el.classList.add('hidden'), 3000);
}

// ---------- MODAL GENÉRICO ----------

// Abre o modal passando o id do elemento — ex: abrirModal('modal-excluir')
// O ?. evita erro se o elemento não existir na página
function abrirModal(id)  { document.getElementById(id)?.classList.remove('hidden'); }
function fecharModal(id) { document.getElementById(id)?.classList.add('hidden'); }

// ---------- FORMATAÇÃO ----------

// Formata número para moeda brasileira — ex: 84.5 → "R$ 84,50"
function formatarMoeda(valor) {
  return 'R$ ' + Number(valor).toFixed(2).replace('.', ',');
}

// ---------- DATA ATUAL ----------

// Retorna a data de hoje por extenso — ex: "sábado, 21 de junho de 2026"
// Usado no cabeçalho do dashboard
function formatarDataHoje() {
  return new Date().toLocaleDateString('pt-BR', {
    weekday: 'long', day: '2-digit', month: 'long', year: 'numeric'
  });
}

// ---------- SIDEBAR: MARCAR ITEM ATIVO + PREENCHER USUÁRIO LOGADO ----------

// Esse bloco roda assim que o HTML termina de carregar em qualquer página.
// Faz duas coisas: destaca o item de menu correto e preenche o usuário logado.
document.addEventListener('DOMContentLoaded', () => {

  // ── 1. Destaca o item ativo no menu lateral ──────────────────────────────
  // Cada página define data-page no <body> — ex: <body data-page="clientes">
  // Cada link do menu tem data-nav — ex: <a data-nav="clientes">
  // Quando os dois batem, a classe 'active' é adicionada, que muda a cor no CSS
  const page = document.body.dataset.page;
  if (page) {
    document.querySelectorAll('.nav-item[data-nav]').forEach(el => {
      if (el.dataset.nav === page) el.classList.add('active');
    });
  }

  // ── 2. Preenche o usuário logado no rodapé do sidebar ───────────────────
  // O login.html salva dois valores no localStorage após autenticação:
  //   - usuario_nome  → nome completo vindo da API
  //   - usuario_email → e-mail vindo da API
  // Se o usuário fez login antes dessa mudança existir, usuario_nome pode não
  // estar salvo. Nesse caso, cai no fallback do login-email (que o sistema
  // antigo salvava quando a opção "lembrar acesso" estava marcada).
  const nome  = localStorage.getItem('usuario_nome')  || '';
  const email = localStorage.getItem('usuario_email') || localStorage.getItem('login-email') || '—';

  // nomeExibido: tenta usar o nome primeiro. Se não tiver, usa o e-mail no lugar.
  // Isso evita mostrar "—" quando o usuário pelo menos tem o e-mail salvo.
  const nomeExibido = nome || email;

  // Busca os três elementos do bloco de usuário no sidebar.
  // O if antes de cada atribuição evita erro em páginas que não têm sidebar
  // (como a tela de login, que não carrega esses ids).
  const avatarEl = document.getElementById('sidebar-avatar');  // círculo com a inicial
  const nomeEl   = document.getElementById('sidebar-usuario'); // linha do nome
  const emailEl  = document.getElementById('sidebar-email');   // linha do e-mail

  // Preenche o nome (ou traço se não tiver nada)
  if (nomeEl)   nomeEl.textContent  = nomeExibido || '—';

  // Preenche o e-mail embaixo do nome
  if (emailEl)  emailEl.textContent = email;

  // Preenche o avatar com a primeira letra maiúscula do nome/e-mail.
  // charAt(0) pega o primeiro caractere — toUpperCase() garante maiúsculo.
  // Se não tiver nada, mostra "?" como placeholder.
  if (avatarEl) avatarEl.textContent = nomeExibido ? nomeExibido.charAt(0).toUpperCase() : '?';
});
