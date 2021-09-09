import clsx from "clsx";
import { Heading, Pane, Paragraph, Position, Tab, Tablist, Tooltip } from 'evergreen-ui';
import React from 'react';
import { useSelector } from 'react-redux';
import useLocalStorageState from "use-local-storage-state";
import { RootState } from '../state/store';
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
                            <Pulse enabled={connected} />
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

const Pulse = ({ enabled = false }) => {
    return (
      <Tooltip
        position={Position.BOTTOM}
        content={`OpenBooks server ${enabled ? "connected" : "disconnected"}.`}>
        <div className="flex justify-evenly items-center">
          <div
            className={clsx([
              "w-2 h-2 rounded-full",
              { "bg-custom-blue animate-custom-pulse": enabled },
              { "animate-none bg-gray-500": !enabled }
            ])}
          />
        </div>
      </Tooltip>
    );
}

export default Sidebar;
