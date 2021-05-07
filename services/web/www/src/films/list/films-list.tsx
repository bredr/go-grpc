import React from "react";
import { useFilmsListQuery, Genre } from "../../generated/graphql";

interface Props {
  searchTerm: string;
  genres: Genre[];
}

export const FilmsList = ({ searchTerm, genres = [] }: Props) => {
  const [result] = useFilmsListQuery({
    variables: { searchTerm, genres },
    pause: !searchTerm && genres.length === 0,
    requestPolicy: "cache-and-network",
  });

  const { data, fetching, error } = result;
  if (fetching) return <p>Loading...</p>;
  if (error) return <p>Oh no... {error.message}</p>;

  return (
    <ul>
      {data?.films?.map(({ ID, Name, Genre }) => (
        <li key={ID}>
          <b>{Genre}</b> {Name}{" "}
        </li>
      ))}
    </ul>
  );
};
