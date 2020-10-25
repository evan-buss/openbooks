import { Pane } from 'evergreen-ui';
import React from 'react';
import Sidebar from './components/Sidebar';
import HistoryProvider from './models/HistoryProvider';
import SearchPage from './pages/SearchPage';

function App() {
  return (
    <HistoryProvider>
      <Sidebar></Sidebar>

      <Pane display="flex" justifyContent="center">
        <SearchPage></SearchPage>
      </Pane>
    </HistoryProvider>
  );
}

export default App;
