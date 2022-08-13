import React from "react";
import Sidebar from "./components/SideBar/Sidebar";
import { SidebarNeo } from "./components/SidebarNeo/Sidebar";
import SearchPage from "./pages/SearchPage";

function App() {
  return (
    <div className="flex flex-row max-h-screen min-h-screen flex-nowrap bg-tint1">
      <Sidebar />
      <SidebarNeo />
      <SearchPage />
    </div>
  );
}

export default App;
