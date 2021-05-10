const createProxyMiddleware = require("http-proxy-middleware");
const STUB = process.env.STUB;

if (STUB) {
  const { buildClientSchema } = require('graphql');
  const { ApolloServer } = require('apollo-server');
  const introspectionResult = require('../graphql.schema.json');
  const schema = buildClientSchema(introspectionResult);
  const server = new ApolloServer({
    schema,
    mocks: true,
  });
  module.exports = function (app) {
    server.listen().then(({ url }) => {
      console.log(`🚀 Server ready at ${url}`);
      app.use("/query", createProxyMiddleware({
        target: url,
        changeOrigin: true,
      }))
    });
  }
} else {
  module.exports = function (app) {
    app.use(
      "/query",
      createProxyMiddleware({
        target: "http://localhost:3080",
        changeOrigin: true,
      })
    );
  };
}
