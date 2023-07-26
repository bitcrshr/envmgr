/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        'page-bg': '#f0f5f9',
        'mystic': '#d5e0e7',
        'heather': '#8a9ba8',
        'regent-gray': '#8a9ba8',
        'pale-sky': '#627084'
      },
      
      fontFamily: {
        'barlow': ['Barlow Condensed', 'sans-serif']
      }
    },
  },
  plugins: [],
}

