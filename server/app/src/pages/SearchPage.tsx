import { Button, Pane, SearchInput } from 'evergreen-ui';
import React from 'react';
import { useHistory } from '../models/HistoryProvider';

const SearchPage: React.FC = () => {
    const { addItem } = useHistory()!;

    return (
        <Pane width="65%" background="tint2" display="flex" justifyContent="center">
            <Button onClick={() => addItem({ query: "from home", time: new Date().getTime() })}>Add Item</Button>
            <SearchInput placeholder="Search for a book." height={40} width="90%"></SearchInput>
        </Pane>
    )
}

export default SearchPage;