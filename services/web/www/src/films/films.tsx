import { Formik } from "formik";
import React, { useState } from "react";
import { Genre } from "../generated/graphql";
import { FilmsList } from "./list/films-list";
import { Search } from "./search/search";
import "./films.css";

export const Films = () => {
  const [{ searchTerm, genres }, setState] = useState<{
    searchTerm: string;
    genres: Genre[];
  }>({ searchTerm: "", genres: [] });

  return (
    <>
      <div>
        <Formik initialValues={{ searchTerm, genres }} onSubmit={setState}>
          {Search}
        </Formik>
      </div>
      <div>
        <FilmsList searchTerm={searchTerm} genres={genres} />
      </div>
    </>
  );
};
