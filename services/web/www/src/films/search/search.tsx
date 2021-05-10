import { Field, Form } from "formik";
import React from "react";
import { Genre } from "../../generated/graphql";
import "./search.css";

function enumKeys<O extends object, K extends keyof O = keyof O>(obj: O): K[] {
  return Object.keys(obj).filter((k) => Number.isNaN(+k)) as K[];
}

export const Search = () => (
  <Form>
    <label htmlFor="searchTerm">Search</label>
    <Field id="searchTerm" name="searchTerm" />
    <label htmlFor="genres">Genres</label>
    <div id="genres" role="group">
      {enumKeys(Genre).map((x) => (
        <label key={x}>
          <Field type="checkbox" name="genres" value={Genre[x]} />
          {x}
        </label>
      ))}
    </div>
    <button type="submit">Submit</button>
  </Form>
);
