const config = {
  proto: "films/films.proto",
  service: "FilmService",
  package: "films"
};

const router = {
  CreateFilm: (x) => x,
  UpdateFilm: (x) => x,
  DeleteFilm: () => ({}),
  FindFilms: () => ({
    Films: [
      { ID: "1", Name: "Star Wars", Genre: "SCI_FI" },
      { ID: "2", Name: "Back to the future", Genre: "SCI_FI" },
      { ID: "3", Name: "Star Trek: The wrath of Khan", Genre: "SCI_FI" },
      { ID: "4", Name: "Star Trek: The search for Spock", Genre: "SCI_FI" },
    ]
  })
}


export {
  config,
  router
}