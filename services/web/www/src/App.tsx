import React from "react";
import "./App.css";
import { createClient, Provider } from "urql";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  NavLink,
} from "react-router-dom";
import { Films } from "./films/films";

const client = createClient({
  url: "/query",
});

function App() {
  return (
    <Provider value={client}>
      <Router>
        <div>
          <nav>
            <NavLink activeClassName="is-active" to="/">
              Films
            </NavLink>
          </nav>
        </div>
        <div className="container">
          <Switch>
            <Route path="/">
              <Films />
            </Route>
          </Switch>
        </div>
      </Router>
    </Provider>
  );
}

export default App;
