import { Button, Pane, SearchInput } from 'evergreen-ui';
import React, { FormEvent, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import styled from 'styled-components';
import { BooksGrid } from '../components/BooksGrid';
import { sendSearch } from '../state/stateSlice';
import { RootState } from '../state/store';

const SearchPage: React.FC = () => {
    const dispatch = useDispatch();
    const [searchQuery, setSearchQuery] = useState("");
    const activeItem = useSelector((store: RootState) => store.state.activeItem);

    const searchHandler = (event: FormEvent) => {
        event.preventDefault();
        if (searchQuery === "") return;
        dispatch(sendSearch(searchQuery));
        setSearchQuery("");
    }

    return (
        <Pane display="flex" justifyContent="center" background="tint1" width="100%">
            <Pane width="85%" display="flex" flexDirection="column" alignItems="center">
                <Form onSubmit={(e) => searchHandler(e)}>
                    <SearchInput
                        disabled={activeItem !== null && !activeItem.results}
                        value={searchQuery}
                        onChange={(e: any) => setSearchQuery(e.target.value)}
                        placeholder="Search for a book."
                        height={40}
                        marginRight={24}
                        width="100%">
                    </SearchInput>
                    <Button type="submit" height={40} appearance="primary" disabled={searchQuery === ""}>Search</Button>
                </Form>
                <BooksGrid />
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
