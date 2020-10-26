import { Pane } from 'evergreen-ui';
import React from 'react';
import Sidebar from './components/Sidebar';
import HistoryProvider from './models/HistoryProvider';
import SearchPage from './pages/SearchPage';

function App() {
  return (
    <HistoryProvider>
      <Pane display="flex" flexFlow="row nowrap" width="100%">
        <Sidebar></Sidebar>
        <SearchPage></SearchPage>
      </Pane>
    </HistoryProvider>
  );
}

export default App;
