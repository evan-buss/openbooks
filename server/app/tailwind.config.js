module.exports = {
  mode: "jit",
  purge: ["./index.html", "./src/**/*.{vue,js,ts,jsx,tsx}"],
  darkMode: false, // or 'media' or 'class'
  theme: {
    extend: {
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
