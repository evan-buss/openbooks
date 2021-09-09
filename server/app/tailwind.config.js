module.exports = {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      animation: {
        'custom-pulse': 'pulse-animation 2s infinite'
      },
      keyframes: {
        'pulse-animation': {
          '0%': { 'box-shadow': '0 0 0 0px rgba(0, 0, 0, 0.2)' },
          '100%': { 'box-shadow': '0 0 0 10px rgba(0, 0, 0, 0)' }
        }
      }
    },
  },
  variants: {
    extend: {},
  },
  plugins: [],
}
