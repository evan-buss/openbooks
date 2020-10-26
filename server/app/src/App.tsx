import React from 'react';
import Sidebar from './components/Sidebar';
import HistoryProvider from './models/HistoryProvider';
import SearchPage from './pages/SearchPage';

function App() {
  return (
    <HistoryProvider>
      <Sidebar></Sidebar>
      <SearchPage></SearchPage>
    </HistoryProvider>
  );
}

export default App;
