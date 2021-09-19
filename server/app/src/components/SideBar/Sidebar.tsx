import { Heading, Pane, Paragraph, Tab, Tablist, Tooltip } from "evergreen-ui";
import React from "react";
import { useDispatch, useSelector } from "react-redux";
import useLocalStorageState from "use-local-storage-state";
import { RootState } from "../../state/store";
import SearchHistory from "./SearchHistory";
import Pulse from "./Pulse";
import NotificationDrawer from "./NotificationDrawer";
import BookList from "./BookList";
import { toggleDrawer } from "../../state/notificationSlice";
import { IdentificationBadge } from "phosphor-react";

const Sidebar: React.FC = () => {
  const [selectedIndex, setIndex] = useLocalStorageState("index", 0);
  const connected = useSelector((store: RootState) => store.state.isConnected);
  const username = useSelector((store: RootState) => store.state.username);
  const dispatch = useDispatch();

  return (
    <>
      <Pane minWidth={325} background="white">
        <Pane zIndex={1} flexShrink={0} className="resize">
          <Pane padding={16} borderBottom="muted">
            <Heading
              size={600}
              display="flex"
              justifyContent="space-between"
              paddingRight="10px">
              OpenBooks
              <Pulse
                enabled={connected}
                onClick={() => dispatch(toggleDrawer())}
              />
            </Heading>
            <Paragraph size={400} color="muted">
              Download eBooks from IRC Highway
            </Paragraph>
          </Pane>
          <Pane display="flex" padding={8} borderBottom="muted">
            <Tablist display="flex" flex={1} justifyContent="space-between">
              {["Book Search", "Past Downloads"].map((tab, index) => (
                <Tab
                  key={tab}
                  flex={1}
                  isSelected={selectedIndex === index}
                  onSelect={() => setIndex(index)}>
                  {tab}
                </Tab>
              ))}
            </Tablist>
          </Pane>
        </Pane>
        <Pane flex="1" padding={8} flexGrow={1}>
          {selectedIndex === 0 ? <SearchHistory /> : <BookList />}
        </Pane>
        {username && (
          <div className="z-10 absolute bottom-0 left-0 flex w-[325px] border-t bg-white border-gray-200 p-4 text-sm">
            <IdentificationBadge size={24} className="mr-4" />
            {username}
          </div>
        )}
      </Pane>

      <NotificationDrawer />
    </>
  );
};

export default Sidebar;
