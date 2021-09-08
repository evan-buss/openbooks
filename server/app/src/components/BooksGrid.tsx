import {
    Button,
    majorScale,
    Pane,
    SelectMenu,
    SelectMenuItem,
    Spinner,
    Table,
    Text,
    TextDropdownButton
} from 'evergreen-ui';
import React, { useEffect, useMemo, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { MessageType } from '../models/messages';
import { sendMessage } from '../state/stateSlice';
import { RootState } from '../state/store';
import image from '../assets/reading.svg'

const stringContains = (first: string, second: string): boolean => {
    return first.toLowerCase().includes(second.toLowerCase());
}

export const BooksGrid: React.FC = () => {
    const activeItem = useSelector((store: RootState) => store.state.activeItem);
    const servers = useSelector((store: RootState) => new Set<string>(store.state.serverFilters));
    const dispatch = useDispatch();
    const [titleFilter, setTitleFilter] = useState<string>("");
    const [authorFilter, setAuthorFilter] = useState<string>("");
    const [formatFilter, setFormatFilter] = useState<string[]>([]);

    useEffect(() => {
        setFormatFilter(activeItem?.results?.map(x => x.format) ?? [])

        const servers = async () => {
            const data = await fetch('http://localhost:5228/servers');
            console.log(await data.json());
        }

        servers();
    }, [activeItem])

    const filteredBooks = useMemo(() => {
        const books = activeItem?.results ?? []
        return books
            .filter(x => {
                if (servers.size === 0)
                    return true;
                return servers.has(x.server);
            })
            .filter(x => stringContains(x.author, authorFilter))
            .filter(x => stringContains(x.title, titleFilter))
            .filter(x => formatFilter.findIndex(format => format === x.format) !== -1)
    }, [activeItem, servers, authorFilter, titleFilter, formatFilter]);

    const availableExtensions = useMemo(() => {
        const extensionTypes = activeItem?.results?.map(x => x.format) ?? [];
        return extensionTypes
            .filter((x, i, arr) => arr.indexOf(x) === i)
            .map(x => { return { label: x, value: x } as SelectMenuItem; });
    }, [activeItem]);

    if (activeItem === null) {
        return (<>
            <Text fontSize="2em" margin="3em">Search a book to get started.</Text>
            <img width="40%" src={image} alt="person reading" />
        </>);
    }

    if (activeItem !== null && !activeItem.results) {
        return <Pane><Spinner marginX="auto" marginY={120} /></Pane>
    }

    const renderFormatSelectHeader = () => {
        return <Table.TextHeaderCell flexBasis={100} flexGrow={0} flexShrink={0}>
            <SelectMenu
                isMultiSelect
                title="Select Book Formats"
                options={availableExtensions}
                selected={formatFilter}
                onSelect={(item) => setFormatFilter((curr) => [...curr, item.value.toString()])}
                onDeselect={(item) => setFormatFilter((curr) => curr.filter(filter => filter !== item.value))}>
                <TextDropdownButton>FORMATS</TextDropdownButton>
            </SelectMenu>
        </Table.TextHeaderCell>
    }

    const renderBody = () => {
        // Results length is zero. No results for query
        if (activeItem.results?.length === 0) {
            return (<Text textAlign="center" margin={majorScale(3)}>No results for search.</Text>);
        }

        // FilteredRows length is zero. Filter not valid.
        if (filteredBooks.length === 0) {
            return (<Text textAlign="center" margin={majorScale(3)}>No books matching filter.</Text>);
        }

        return filteredBooks.map((book, i) => (
            <Table.Row key={book.full + i} isSelectable>
                <Table.TextCell flexBasis={200} flexGrow={0} flexShrink={0}>{book.server}</Table.TextCell>
                <Table.TextCell flexBasis={250} flexGrow={0} flexShrink={0}>{book.author}</Table.TextCell>
                <Table.TextCell>{book.title}</Table.TextCell>
                <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>{book.format}</Table.TextCell>
                <Table.TextCell flexBasis={100} flexGrow={0} flexShrink={0}>{book.size}</Table.TextCell>
                <Table.Cell flexBasis={150} flexGrow={0} flexShrink={0}>
                    <Button appearance="primary" size="small"
                        onClick={() => dispatch(sendMessage({ type: MessageType.DOWNLOAD, payload: { book: book.full } }))}>
                        Download
                    </Button>
                </Table.Cell>
            </Table.Row>
        ))
    }

    return (
        <Table flex={1} display="flex" flexDirection="column" width="100%">
            <Table.Head background="white">
                <Table.TextHeaderCell flexBasis={200} flexGrow={0} flexShrink={0}>Server</Table.TextHeaderCell>
                <Table.SearchHeaderCell
                    onChange={(e) => setAuthorFilter(e)}
                    placeholder="Author"
                    flexBasis={250}
                    flexGrow={0}
                    flexShrink={0} >Author
                </Table.SearchHeaderCell>
                <Table.SearchHeaderCell onChange={(e) => setTitleFilter(e)} placeholder="Title">Title</Table.SearchHeaderCell>
                {renderFormatSelectHeader()}
                <Table.TextHeaderCell flexBasis={100} flexGrow={0} flexShrink={0}>Size</Table.TextHeaderCell>
                <Table.TextHeaderCell flexBasis={150} flexGrow={0} flexShrink={0} >Download?</Table.TextHeaderCell>
            </Table.Head>
            <Table.VirtualBody>
                {renderBody()}
            </Table.VirtualBody >
        </Table >
    )
}
