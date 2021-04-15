const config = {
  proto: "films/films.proto",
  service: "FilmService",
  package: "films"
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