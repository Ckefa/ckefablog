module.exports = {
  content: [
    './contents/**/*.html',  // Adjust paths as needed
    './handlers/**/*.go',    // Include Go templates if applicable
    './static/**/*.css',     // Scan your CSS files
    './views/**/*.html',     // Include any HTML files in your views directory
  ],
  theme: {
    extend: {
      colors: {
      },
    },
  },
  plugins: [],
};
