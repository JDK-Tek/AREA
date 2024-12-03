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
            spotify: {
                100: "#19943a",
                200: "#2aa84f",
            },
            weatherunderground: {
                100: "#28292e",
                200: "#3a3b40",
            }
          }
      },
  },
  plugins: [],
};
