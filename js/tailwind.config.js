/* ================================================
   Rosas & Tortas Confeitaria — Tailwind Config
   Incluir ANTES do tailwind CDN em cada página
   ================================================ */

// Este bloco deve estar num <script> antes do CDN do Tailwind
// Exemplo de uso em cada HTML:
//   <script>/* conteúdo deste arquivo */</script>
//   <script src="https://cdn.tailwindcss.com"></script>

tailwind.config = {
  theme: {
    extend: {
      colors: {
        rose: {
          DEFAULT: '#F4A7B9',
          light:   '#FAD4DE',
          dark:    '#E8789A',
          deep:    '#C9546E',
          50:      '#fff0f4',
          100:     '#FAD4DE',
          200:     '#F4A7B9',
          300:     '#ED80A0',
          400:     '#E8789A',
          500:     '#C9546E',
        },
        mint: {
          DEFAULT: '#A8D5BA',
          light:   '#D4EDE0',
          dark:    '#6FBA96',
          deep:    '#4A9B72',
        },
        cream: {
          DEFAULT: '#FFF8F5',
          card:    '#FFFFFF',
          border:  '#F0E6E0',
        },
      },
      fontFamily: {
        display: ['Playfair Display', 'serif'],
        body:    ['DM Sans', 'sans-serif'],
      },
      boxShadow: {
        soft:  '0 2px 16px rgba(244,167,185,0.15)',
        card:  '0 4px 24px rgba(200,120,140,0.10)',
        hover: '0 8px 32px rgba(200,120,140,0.18)',
      },
    }
  }
}
