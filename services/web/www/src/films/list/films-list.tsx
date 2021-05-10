import React from "react";
import { useFilmsListQuery, Genre } from "../../generated/graphql";
import "./list.css";
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

  if (!data?.films?.length) return <p>No films found yet</p>;

  return (
    <>
      <p></p>
      <table>
        <thead>
          <tr>
            <th>Genre</th>
            <th>Name</th>
          </tr>
        </thead>
        <tbody>
          {data?.films?.map(({ ID, Name, Genre }) => (
            <tr key={ID}>
              <td>{Genre}</td>
              <td>{Name}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </>
  );
};
