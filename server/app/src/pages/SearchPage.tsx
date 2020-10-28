import { Button, Pane, SearchInput } from 'evergreen-ui';
import React, { FormEvent, useEffect, useRef, useState } from 'react';
import { useDispatch } from 'react-redux';
import styled from 'styled-components';
import SearchResults from '../components/SearchResults';
import { MessageHandler } from '../models/MessageHandler';
import { addHistoryItem } from '../state/historySlice';

const SearchPage: React.FC = () => {
    const dispatch = useDispatch();
    const [searchQuery, setSearchQuery] = useState("");

    const handler = useRef<MessageHandler | undefined>(undefined);

    useEffect(() => {
        handler.current = new MessageHandler();

        return () => {
            handler.current!.dispose();
        };
    }, []);

    const searchHandler = (event: FormEvent) => {
        event.preventDefault()
        if (searchQuery === "") return;
        dispatch(addHistoryItem({ query: searchQuery, timestamp: new Date().getTime() }));
        setSearchQuery("");
    }

    return (
        <Pane display="flex" justifyContent="center" background="tint1" width="100%">
            <Pane width="65%" display="flex" flexDirection="column" alignItems="center">
                <Form onSubmit={(e) => searchHandler(e)}>
                    <SearchInput
                        value={searchQuery}
                        onChange={(e: any) => setSearchQuery(e.target.value)}
                        placeholder="Search for a book."
                        height={40}
                        marginRight={24}
                        width="100%">
                    </SearchInput>
                    <Button type="submit" height={40} appearance="primary">Search</Button>
                </Form>
                <SearchResults></SearchResults>
            </Pane>
        </Pane>
    )
}

const Form = styled.form`
    width: 100%;
    display: flex;
    flex-flow: row nowrap;
    justify-content: center;
    align-items: center;
    margin: 24px;
`

export default SearchPage;