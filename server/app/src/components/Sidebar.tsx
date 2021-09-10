import clsx from "clsx";
import {
  Alert,
  Button,
  Heading,
  NotificationsIcon,
  Pane,
  Paragraph,
  Position,
  SideSheet,
  Tab,
  Tablist,
  Tooltip
} from "evergreen-ui";
import React from "react";
import { useSelector } from "react-redux";
import { useHistory } from "react-router";
import useLocalStorageState from "use-local-storage-state";
import { RootState } from "../state/store";
import SearchHistory from "./SearchHistory";
import ServerList from "./ServerList";

const Sidebar: React.FC = () => {
  const [selectedIndex, setIndex] = useLocalStorageState("index", 0);
  const connected = useSelector((store: RootState) => store.state.isConnected);
  const history = useHistory();
  const [showDrawer, setShowDrawer] = React.useState(false);

  const handleClick = (index: number) => {
    console.log("Handling click.");
    setIndex(index);
    history.push(index === 0 ? "/" : "/library");
  };

  return (
    <>
      <Pane minWidth={325} background="white">
        <Pane zIndex={1} flexShrink={0} elevation={0}>
          <Pane padding={16} borderBottom="muted">
            <Heading
              size={600}
              display="flex"
              justifyContent="space-between"
              paddingRight="10px">
              OpenBooks
              <Pulse
                enabled={connected}
                onClick={() => setShowDrawer((x) => !x)}
              />
            </Heading>
            <Paragraph size={400} color="muted">
              Download eBooks from IRC Highway
            </Paragraph>
          </Pane>
          <Pane display="flex" padding={8}>
            <Tablist display="flex" flex={1} justifyContent="space-between">
              {["Book Search", "Library"].map((tab, index) => (
                <Tab
                  key={tab}
                  flex={1}
                  isSelected={selectedIndex === index}
                  onSelect={() => handleClick(index)}>
                  {tab}
                </Tab>
              ))}
            </Tablist>
          </Pane>
        </Pane>
        <Pane
          flex="1"
          overflowY={selectedIndex === 0 ? "hidden" : "scroll"}
          padding={8}
          maxHeight="calc(100vh - 78px - 44px)">
          {selectedIndex === 0 ? <SearchHistory /> : <ServerList />}
        </Pane>
      </Pane>

      <SideSheet
        isShown={showDrawer}
        width={325}
        onCloseComplete={() => setShowDrawer(false)}>
        <div className="p-2 font-medium border-b border-gray-300">
          <h2 className="text-2xl">Notifications</h2>
        </div>
        <div className="flex flex-col p-2 gap-2">
          <Alert
            onRemove={() => console.log("removing?")}
            isRemoveable={true}
            intent="none"
            title="There are over 200 integrations available"
          />
          <Alert intent="success" title="Your source is now sending data" />
          <Alert intent="warning" title="Changes will affect all warehouses">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed dod
            eiusmod tempor incididunt ut labore et dolore magna aliqua.
          </Alert>
          <Alert intent="danger" title="We weren’t able to save your changes">
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed dod
            eiusmod tempor incididunt ut labore et dolore magna aliqua.
          </Alert>
        </div>
      </SideSheet>
    </>
  );
};

interface PulseProps {
  enabled: boolean;
  onClick: React.MouseEventHandler<HTMLDivElement>;
}

const Pulse = ({ enabled, onClick }: PulseProps) => {
  return (
    <Tooltip
      position={Position.BOTTOM}
      content={`OpenBooks server ${enabled ? "connected" : "disconnected"}.`}>
      <div onClick={onClick} className="flex justify-evenly items-center">
        <NotificationsIcon
          className={clsx([
            "w-3 h-3 rounded-full cursor-pointer",
            { "text-custom-blue animate-custom-pulse": enabled },
            { "animate-none text-gray-500": !enabled }
          ])}
        />
        {/* <div
          className={clsx([
            "w-2 h-2 rounded-full",
            { "bg-custom-blue animate-custom-pulse": enabled },
            { "animate-none bg-gray-500": !enabled }
          ])}
        /> */}
      </div>
    </Tooltip>
  );
};

export default Sidebar;
