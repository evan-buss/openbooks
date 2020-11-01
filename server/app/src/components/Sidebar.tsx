import useLocalStorage from '@rehooks/local-storage';
import { Heading, Pane, Paragraph, Tab, Tablist } from 'evergreen-ui';
import React from 'react';
import styled from "styled-components";
import Pulse from './Pulse';
import SearchHistory from './SearchHistory';
import ServerList from './ServerList';

const SidebarContent = styled(Pane)`
    max-height: calc(100vh - 78px - 44px);
`;

const Sidebar: React.FC = () => {
    const [selectedIndex, setIndex] = useLocalStorage("index", 0);

    return (
        <>
            <Pane minWidth={325}>
                <Pane zIndex={1} flexShrink={0} elevation={0}>
                    <Pane padding={16} borderBottom="muted">
                        <Heading size={600}
                            display="flex"
                            justifyContent="space-between"
                            paddingRight="10px">
                            OpenBooks
                            <Pulse />
                        </Heading>
                        <Paragraph size={400} color="muted">
                            Download eBooks from IRC Highway
                    </Paragraph>
                    </Pane>
                    <Pane display="flex" padding={8}>
                        <Tablist display="flex" flex={1} justifyContent="space-between">
                            {['Search History', 'Online Servers'].map(
                                (tab, index) => (
                                    <Tab
                                        key={tab}
                                        flex={1}
                                        isSelected={selectedIndex === index}
                                        onSelect={() => setIndex(index)}>
                                        {tab}
                                    </Tab>))
                            }
                        </Tablist>
                    </Pane>
                </Pane>
                <SidebarContent flex="1" overflowY={selectedIndex === 0 ? 'hidden' : 'scroll'} padding={8}>
                    {selectedIndex === 0 ? <SearchHistory /> : <ServerList />}
                </SidebarContent>
            </Pane>
        </>
    )
}

export default Sidebar;
