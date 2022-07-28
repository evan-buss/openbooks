import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import Sidebar from "./components/SideBar/Sidebar";
import { LiveSearchPage } from "./pages/LiveSearchPage";
import SearchPage from "./pages/SearchPage";

function App() {
  return (
    <div className="flex flex-row max-h-screen min-h-screen flex-nowrap bg-tint1">
      <Sidebar />

      <BrowserRouter>
        <Routes>
          <Route path="/" element={<SearchPage />} />
          <Route path="/live-search" element={<LiveSearchPage />} />
        </Routes>
      </BrowserRouter>
    </div>
  );
}

export default App;
