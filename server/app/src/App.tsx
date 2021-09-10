import {
  Alert,
  Heading,
  IconButton,
  majorScale,
  NotificationsIcon,
  Pane,
  Paragraph,
  SideSheet
} from "evergreen-ui";
import React from "react";
import Sidebar from "./components/Sidebar";
import SearchPage from "./pages/SearchPage";
import { BrowserRouter, Switch, Route } from "react-router-dom";
import LibraryPage from "./pages/LibraryPage";

function App() {
  return (
    <BrowserRouter>
      <div className="flex flex-row flex-nowrap bg-tint1 min-h-screen max-h-screen">
        <Sidebar />
        <Switch>
          <Route path="/library">
            <LibraryPage />
          </Route>
          <Route path="/">
            <SearchPage />
          </Route>
        </Switch>
      </div>
    </BrowserRouter>
  );
}

export default App;
