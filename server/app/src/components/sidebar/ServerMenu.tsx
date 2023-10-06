import {
  ActionIcon,
  Button,
  Group,
  Modal,
  px,
  Stack,
  Switch,
  Text,
  Tooltip
} from "@mantine/core";
import { CheckCircle, Plugs, PlugsConnected } from "@phosphor-icons/react";
import { useId, useState } from "react";
import classes from "./ServerMenu.module.css";
import { useDisclosure } from "@mantine/hooks";
import { useAppDispatch, useAppSelector } from "../../state/store";
import {
  connectToServer,
  IrcServer,
  SelectedServer,
  selectServer
} from "../../state/connectionSlice";
import isEqual from "lodash/isEqual";

export function ServerMenu({ connected }: { connected: boolean }) {
  const [opened, { open, close }] = useDisclosure(false);
  const servers = useAppSelector((state) => state.connection.servers);
  const selectedServer = useAppSelector(
    (state) => state.connection.selectedServer
  );
  const dispatch = useAppDispatch();

  const [selected, setSelected] = useState<SelectedServer>(selectedServer);

  const handleSelect = (server: IrcServer, enableTLS: boolean) => {
    setSelected({ name: server.name, enableTLS });
  };

  const handleConnect = () => {
    dispatch(selectServer(selected));
    dispatch(connectToServer());
    close();
  };

  return (
    <>
      <Tooltip label="Select IRC Server">
        <ActionIcon variant="gradient" color="blue" onClick={open}>
          {connected ? <PlugsConnected size={18} /> : <Plugs size={18} />}
        </ActionIcon>
      </Tooltip>

      <Modal opened={opened} onClose={close} title="IRC Server" centered>
        <Stack gap={px("1rem")}>
          {servers.map((server) => (
            <ServerOption
              key={server.name}
              server={server}
              isSelected={selected.name === server.name}
              onSelect={(ssl) => handleSelect(server, ssl)}
            />
          ))}
          <Button
            style={{ alignSelf: "end" }}
            disabled={isEqual(selectedServer, selected)}
            onClick={handleConnect}>
            Connect
          </Button>
        </Stack>
      </Modal>
    </>
  );
}

interface ServerOptionProps {
  server: IrcServer;
  isSelected: boolean;
  onSelect: (ssl: boolean) => void;
}

function ServerOption({ server, isSelected, onSelect }: ServerOptionProps) {
  const id = useId();
  const [enableTLS, setEnableTLS] = useState(true);
  const serverString = `${server.address}:${
    enableTLS ? server.sslPort : server.port
  }`;

  const handleServerChange = () => {
    onSelect(enableTLS);
  };

  const handleTLSChange = (ssl: boolean) => {
    setEnableTLS(ssl);
    onSelect(ssl);
  };

  return (
    <div className={classes.root}>
      <input
        className={classes.input}
        type="radio"
        name="server"
        value={server.name}
        checked={isSelected}
        onChange={handleServerChange}
        id={id}
      />

      <label htmlFor={id} className={classes.label}>
        <Group gap={12} align="center">
          <CheckCircle size={20} weight="bold" className={classes.icon} />
          <Stack gap={0}>
            <Text>{server.name}</Text>
            <Tooltip position="right" label={serverString}>
              <Text c="gray.5">{server.address}</Text>
            </Tooltip>
          </Stack>
        </Group>

        <div>
          <Switch
            className={classes.switch}
            defaultChecked
            size="xs"
            w={100}
            checked={enableTLS}
            onChange={(e) => handleTLSChange(e.currentTarget.checked)}
            label="TLS?"
          />
        </div>
      </label>
    </div>
  );
}
