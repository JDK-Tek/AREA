/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,ts,jsx,tsx,html}'],
  theme: {
      extend: {
          fontFamily: {
              anton: ['Anton', 'sans-serif'],
              spartan: ['"League Spartan"', 'sans-serif'],
          },
          colors: {
            chartpurple: {
                100: "#5f18eb",
                200: "#410cab",
                300: "#11032e"
            },
            chartgray: {
                100: "#a6a6a6",
                200: "#222222"
            },
            spotify: {
                100: "#05b348",
                200: "#00ca4e",
            },
            weatherunderground: {
                100: "#28292e",
                200: "#3a3b40",
            },
            instagram: {
                100: "#f1005b",
                200: "#ff1970",
            },
            nasa: {
                100: "#341d4f",
                200: "#3d1f5e",
            },
          }
      },
  },
  plugins: [],
};
