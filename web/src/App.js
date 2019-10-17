import React from 'react';
import logo from './logo.svg';
import './App.css';
import Search from './components/Search';
import Table from './components/Table';

function App() {
  return (
    <div className="app-body">
      <div id="app-container">
        <Search />
        <Table />
      </div>
    </div>
  );
}

export default App;
