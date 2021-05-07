import React from 'react';
import './App.css';
import {createClient, Provider} from "urql";
import {BrowserRouter as Router, Switch, Route, Link } from "react-router-dom";
import { Films } from './films/films';

const client = createClient({
  url: "/query"
})

function App() {
  return (
    <Provider value={client}>
    <Router>
      <div>
        <nav>
          <ul>
            <Link to="/">Films</Link>
          </ul>
        </nav>
      </div>
      <Switch>
        <Route path="/">
          <Films/>
        </Route>
      </Switch>
    </Router>
    </Provider>
  );
}

export default App;
