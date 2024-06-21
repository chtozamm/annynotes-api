/** @type {import("prettier").Config} */
export default {
  tailwindConfig: "./tailwind.config.js",
  plugins: ["prettier-plugin-go-template", "prettier-plugin-tailwindcss"],
  overrides: [
    {
      files: ["*.html"],
      options: {
        parser: "go-template",
      },
    },
  ],
};
