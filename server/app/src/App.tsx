import React from "react";
import Sidebar from "./components/SideBar/Sidebar";
import SearchPage from "./pages/SearchPage";
import { BrowserRouter, Switch, Route } from "react-router-dom";

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-row flex-nowrap bg-tint1 min-h-screen max-h-screen">
        <Sidebar />
        <Switch>
          <Route path="/">
            <SearchPage />
          </Route>
        </Switch>
      </div>
    </BrowserRouter>
  );
}

export default App;
