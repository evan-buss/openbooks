module.exports = {
  mode: "jit",
  purge: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
      animation: {
        "custom-pulse": "pulse-animation 2s infinite"
      },
      keyframes: {
        "pulse-animation": {
          "0%": { "box-shadow": "0 0 0 0px rgba(0, 0, 0, 0.2)" },
          "100%": { "box-shadow": "0 0 0 10px rgba(0, 0, 0, 0)" }
        }
      },
      colors: {
        "custom-blue": "#3366ff",
        "tint1": "#FAFBFF",
        "active-blue": "#ebf0ff",
        "active-text-blue": "#3366FF",
        "hover-blue": "#f4f5f9"
      }
    }
  },
  variants: {
    extend: {}
  },
  plugins: []
};
