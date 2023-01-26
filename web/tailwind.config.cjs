/** @type {import('tailwindcss').Config} */
module.exports = {
  darkMode: "class",
  content: ["./src/**/*.{astro,html,js,jsx,md,mdx,svelte,ts,tsx,vue}"],
  theme: {
    extend: {
      fontFamily: {
        cantarell: ["Cantarell", "sans-serif"],
      },
      colors: {
        wireguard: {
          50: "#C02126",
          100: "#AE1E22",
          200: "#9D1B1F",
          300: "#88171A",
          400: "#7A1518",
          500: "#691215",
          600: "#570F11",
          700: "#460C0E",
          800: "#34090A",
          900: "#230607",
        },
      },
    },
  },
  plugins: [],
};
