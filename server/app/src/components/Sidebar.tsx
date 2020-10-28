import { Heading, Pane, Paragraph, Tab, Tablist } from 'evergreen-ui';
import React, { useState } from 'react';
import styled from "styled-components";
import SearchHistory from './SearchHistory';
import ServerList from './ServerList';

const Sidebar: React.FC = () => {
    const [selectedIndex, setIndex] = useState(1);

    return (
        <>
            <Pane minWidth={300}>
                <Pane zIndex={1} flexShrink={0} elevation={0}>
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
                                        onSelect={() => setIndex(index)}>
                                        {tab}
                                    </Tab>))
                            }
                        </Tablist>
                    </Pane>
                </Pane>
                <SidebarContent flex="1" scrollable={selectedIndex === 0 ? 1 : 0} overflowY="hidden" padding={8}>
                    {selectedIndex === 0 ? <ServerList /> : <SearchHistory />}
                </SidebarContent>
            </Pane>
        </>
    )
}


const SidebarContent = styled(Pane) <{ scrollable: boolean }>`
    max-height: calc(100vh - 78px - 44px);
    &:hover {
        overflow-y: ${props => props.scrollable ? 'scroll' : 'hidden'};
    }
`

export default Sidebar;