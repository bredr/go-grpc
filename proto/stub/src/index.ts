import { loadPackageDefinition, Server, ServerCredentials } from "@grpc/grpc-js";
import { loadSync } from "@grpc/proto-loader";

function generateRouter(r: Record<string, (x: any) => any>) {
  return Object.keys(r).
    reduce((acc, method) => ({
      ...acc,
      [method]: (call: any, callback: any) => {
        console.log(`${new Date().toISOString()} - ${method}: ${JSON.stringify(call.request)}`);
        return callback(null, r[method](call.request))
      }
    }), {})
}

async function getServer() {
  const pathToStub = process.env.STUB || "./stub"
  const { config, router } = await import(pathToStub)
  const packageDefinition = loadSync(
    config.proto,
    {
      keepCase: true,
      longs: Number,
      enums: String,
      defaults: true,
      oneofs: true,
      includeDirs: [__dirname + "../../../"]
    }
  );
  const protoDescriptor = loadPackageDefinition(packageDefinition);

  const server = new Server();
  server.addService((protoDescriptor as any)[config.package][config.service].service,
    generateRouter(router)
  );
  return server;
}

(async () => {
  const s = await getServer();
  s.bindAsync("0.0.0.0:50051", ServerCredentials.createInsecure(), () => {
    s.start();
  });

})();
