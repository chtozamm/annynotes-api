/** @type {import('tailwindcss').Config} */
export default {
  content: ["templates/**/*.html"],
  theme: {
    extend: {
      fontFamily: {
        ringbearer: ["Ringbearer"],
      },
      colors: {
        primary: "#ffb220",
        secondary: "#fffbf7",
      },
    },
  },
  plugins: [],
};
