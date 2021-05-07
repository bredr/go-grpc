import { Field, Form, Formik } from "formik";
import React, { useState } from "react";
import { Genre } from "../generated/graphql";
import { FilmsList } from "./list/films-list";


export const Films = () => {

  const [{searchTerm, genres}, setState] = useState<{
    searchTerm: string, genres: Genre[]
  }>({searchTerm: "", genres: []}) 

  return (
    <div>
      <Formik initialValues={{searchTerm, genres}} onSubmit={(x) => {
        console.log(x);
        setState(x);
      }}>
        {({values}) => (
          <Form>
            <Field name="searchTerm"/>
            <label htmlFor="genres">Genres</label>
      <div role="group">
        {
          enumKeys(Genre).map(x =>(
          <label key={x}>
            <Field type="checkbox" name="genres" value={Genre[x]}/>
            {x}
          </label>))
        }
      </div>
      <button type="submit">Submit</button>
          </Form>
        )}
      </Formik>
    <FilmsList searchTerm={searchTerm} genres={genres}/>
    </div>
  )
}

function enumKeys<O extends object, K extends keyof O = keyof O>(obj: O): K[] {
  return Object.keys(obj).filter(k => Number.isNaN(+k)) as K[];
}