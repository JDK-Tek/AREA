/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{js,ts,jsx,tsx,html}'],
  theme: {
      extend: {
          fontFamily: {
              anton: ['Anton', 'sans-serif'],
              spartan: ['"League Spartan"', 'sans-serif'],
          },
      },
  },
  plugins: [],
};
