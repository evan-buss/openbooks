import { Pane, SearchInput } from 'evergreen-ui';
import React, { useEffect, useRef } from 'react';
import SearchResults from '../components/SearchResults';
import { MessageHandler } from '../models/MessageHandler';

type Props = {

}

const SearchPage: React.FC = () => {

    const handler = useRef<MessageHandler | undefined>(undefined);

    useEffect(() => {
        handler.current = new MessageHandler();

        return () => {
            handler.current!.dispose();
        };
    }, []);

    return (
        <Pane display="flex" justifyContent="center" background="tint2" height="100vh">
            <Pane width="65%" display="flex" flexDirection="column" alignItems="center">
                <SearchInput margin={24} placeholder="Search for a book." height={40} width="100%"></SearchInput>
                <SearchResults></SearchResults>
            </Pane>
        </Pane>
    )
}

export default SearchPage;