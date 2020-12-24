module.exports = {
  devServer: {
    proxy: {
      "/img": {
        target: "http://newton.apertoire.org:7623",
      },
    },
  },
};
