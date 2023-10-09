import {
  ActionIcon,
  Button,
  Code,
  Group,
  Modal,
  px,
  Stack,
  Switch,
  Text,
  Tooltip
} from "@mantine/core";
import {
  CheckCircle,
  Plugs,
  PlugsConnected,
  ShieldCheck,
  ShieldSlash
} from "@phosphor-icons/react";
import { useId, useState } from "react";
import classes from "./ServerMenu.module.css";
import { useDisclosure } from "@mantine/hooks";
import { useAppDispatch, useAppSelector } from "../../state/store";
import {
  IrcServer,
  SelectedServer,
  setIrcServer
} from "../../state/connectionSlice";
import { setActiveItem } from "../../state/stateSlice";
import { connectToServer } from "../../state/socketMiddleware";

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
    dispatch(setActiveItem(null));
    dispatch(setIrcServer(selected));
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
              isTLSEnabled={selectedServer.enableTLS}
              onSelect={(ssl) => handleSelect(server, ssl)}
            />
          ))}
          <Button style={{ alignSelf: "end" }} onClick={handleConnect}>
            {selectedServer.name === selected.name ? "Reconnect" : "Connect"}
          </Button>
        </Stack>
      </Modal>
    </>
  );
}

interface ServerOptionProps {
  server: IrcServer;
  isSelected: boolean;
  isTLSEnabled: boolean;
  onSelect: (ssl: boolean) => void;
}

function ServerOption({
  server,
  isSelected,
  isTLSEnabled,
  onSelect
}: ServerOptionProps) {
  const id = useId();
  const [enableTLS, setEnableTLS] = useState(isTLSEnabled);
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
            <Text size="md" fw={500}>
              {server.name} <Code>#{server.channel}</Code>
            </Text>
            <Tooltip position="right" label={serverString}>
              <Text c="gray.5">{server.address}</Text>
            </Tooltip>
          </Stack>
        </Group>

        <Tooltip
          label={enableTLS ? "TLS Enabled" : "TLS Disabled"}
          refProp="rootRef">
          <Switch
            aria-label="Enable TLS?"
            className={classes.switch}
            defaultChecked
            labelPosition="right"
            size="lg"
            onLabel="ON"
            offLabel="OFF"
            thumbIcon={
              enableTLS ? (
                <ShieldCheck weight="bold" color="black" size={14} />
              ) : (
                <ShieldSlash weight="bold" color="black" size={14} />
              )
            }
            checked={enableTLS}
            onChange={(e) => handleTLSChange(e.currentTarget.checked)}
          />
        </Tooltip>
      </label>
    </div>
  );
}
