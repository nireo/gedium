import React, { useState, ChangeEvent } from 'react';
import { Post } from '../../../interfaces/post.interfaces';

export const Search: React.FC = () => {
  const [search, setSearch] = useState<string>('');
  const [results, setResults] = useState<Post[]>([]);

  const handleSearch = (event: ChangeEvent<HTMLFormElement>) => {
    event.preventDefault();
  };

  return (
    <div className="container">
      <form onSubmit={handleSearch}>
        <h4 className="text-2xl font-mono text-blue-500 mb-2 mt-4">
          Search for posts
        </h4>
        <div>
          <input
            className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline text-xl"
            value={search}
            placeholder="Search for topic"
            onChange={({ target }) => setSearch(target.value)}
          />
          <button
            type="submit"
            className="bg-blue-500 mb-4 hover:bg-blue-700 text-white font-bold py-2 px-4 mt-4"
          >
            Search
          </button>
        </div>
      </form>
    </div>
  );
};
