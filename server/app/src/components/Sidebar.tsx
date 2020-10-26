import { Heading, IconButton, MoreIcon, Pane, Paragraph, Position, SideSheet, Tab, Tablist } from 'evergreen-ui';
import React, { useState } from 'react';
import styled from "styled-components";
import SearchHistory from './SearchHistory';
import ServerList from './ServerList';

const Sidebar: React.FC = () => {
    const [selectedIndex, setIndex] = useState(0);
    const [opened, setOpen] = useState(false);

    return (
        <>
            {/* Show / hide side pane */}
            {/* < IconButton height={38}
                icon={MoreIcon}
                margin={16}
                position="absolute"
                onClick={() => setOpen(curr => !curr)} >
            </IconButton > */}
            <Pane width={300}>
                <Pane zIndex={1} flexShrink={0} elevation={0} backgroundColor="white">
                    <Pane padding={16} borderBottom="muted">
                        <Heading size={600}>OpenBooks</Heading>
                        <Paragraph size={400} color="muted">
                            Download eBooks from IRC Highway
                    </Paragraph>
                    </Pane>
                    <Pane display="flex" padding={8}>
                        <Tablist display="flex" flex={1} justifyContent="space-between">
                            {['Online Servers', 'Search History'].map(
                                (tab, index) => (
                                    <Tab
                                        key={tab}
                                        flex={1}
                                        isSelected={selectedIndex === index}
                                        onSelect={() => setIndex(index)}
                                    >
                                        {tab}
                                    </Tab>
                                )
                            )}

                        </Tablist>
                    </Pane>
                </Pane>
                <SidebarContent flex="1" scrollable={selectedIndex === 0} overflowY="hidden" padding={8}>
                    {selectedIndex === 0 ?
                        <ServerList servers={['Hello', 'World']}></ServerList>
                        : <SearchHistory></SearchHistory>}
                </SidebarContent>
            </Pane>
        </>
    )
}


const SidebarContent = styled(Pane) <{ scrollable: boolean }>`
    &:hover {
        overflow-y: ${props => props.scrollable ? 'scroll' : 'hidden'};
    }
`

export default Sidebar;