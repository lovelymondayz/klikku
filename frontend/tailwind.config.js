/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        primary: '#ff6b9d',
        secondary: '#c44dff',
        accent: '#ffd93d',
      },
      fontFamily: {
        display: ['var(--font-display)', 'sans-serif'],
        body: ['var(--font-body)', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.5s ease-out',
        'slide-up': 'slideUp 0.5s ease-out',
        'scale-in': 'scaleIn 0.3s ease-out',
        'flash': 'flash 0.3s ease-out',
        'confetti': 'confetti 1s ease-out',
        'pulse-slow': 'pulse 3s ease-in-out infinite',
      },
      keyframes: {
        fadeIn: { '0%': { opacity: '0' }, '100%': { opacity: '1' } },
        slideUp: { '0%': { transform: 'translateY(20px)', opacity: '0' }, '100%': { transform: 'translateY(0)', opacity: '1' } },
        scaleIn: { '0%': { transform: 'scale(0.9)', opacity: '0' }, '100%': { transform: 'scale(1)', opacity: '1' } },
        flash: { '0%': { opacity: '0' }, '50%': { opacity: '1' }, '100%': { opacity: '0' } },
        confetti: { '0%': { transform: 'translateY(0) rotate(0deg)', opacity: '1' }, '100%': { transform: 'translateY(-100px) rotate(720deg)', opacity: '0' } },
      },
    },
  },
  plugins: [],
}
