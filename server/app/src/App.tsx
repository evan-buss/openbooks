import { Pane } from 'evergreen-ui';
import React from 'react';
import Sidebar from './components/Sidebar';
import SearchPage from './pages/SearchPage';

function App() {
  return (
    <Pane clearfix display="flex" flexFlow="row nowrap" background="tint1" minHeight="100vh" maxHeight="100vh">
      <Sidebar />
      <SearchPage />
    </Pane>
  );
}

export default App;
