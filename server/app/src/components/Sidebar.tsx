import { Heading, Pane, Paragraph, Tab, Tablist } from 'evergreen-ui';
import React from 'react';
import { useSelector } from 'react-redux';
import useLocalStorageState from "use-local-storage-state";
import { RootState } from '../state/store';
import Pulse from './Pulse';
import SearchHistory from './SearchHistory';
import ServerList from './ServerList';


const Sidebar: React.FC = () => {
    const [selectedIndex, setIndex] = useLocalStorageState("index", 0);
    const connected = useSelector((store: RootState) => store.state.isConnected);

    return (
        <>
            <Pane minWidth={325} background="white">
                <Pane zIndex={1} flexShrink={0} elevation={0}>
                    <Pane padding={16} borderBottom="muted">
                        <Heading size={600}
                            display="flex"
                            justifyContent="space-between"
                            paddingRight="10px">
                            OpenBooks
                            <Pulse disabled={!connected} />
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
                <Pane flex="1" overflowY={selectedIndex === 0 ? 'hidden' : 'scroll'} padding={8} maxHeight="calc(100vh - 78px - 44px)">
                    {selectedIndex === 0 ? <SearchHistory /> : <ServerList />}
                </Pane>
            </Pane>
        </>
    )
}

export default Sidebar;
