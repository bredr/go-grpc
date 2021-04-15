const config = {
  proto: "web/web.proto",
  service: "WebService",
  package: "web"
};

const router = {
  CreateFilm: (x) => x,
  UpdateFilm: (x) => x,
  DeleteFilm: () => ({}),
  FindFilms: () => ({ Films: [{ ID: "1", Name: "Star Wars", Genre: "SCI_FI" }] })
}

export {
  config,
  router
}