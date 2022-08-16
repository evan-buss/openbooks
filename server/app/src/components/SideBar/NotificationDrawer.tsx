import {
  ActionIcon,
  Alert,
  Center,
  Drawer,
  Group,
  Stack,
  Text,
  Tooltip,
  useMantineColorScheme
} from "@mantine/core";
import {
  BellSimpleSlash,
  CheckCircle,
  Info,
  Warning,
  WarningCircle
} from "phosphor-react";
import { ReactNode } from "react";
import { useSelector } from "react-redux";
import { NotificationType } from "../../state/messages";
import {
  clearNotifications,
  dismissNotification,
  toggleDrawer
} from "../../state/notificationSlice";
import { RootState, useAppDispatch } from "../../state/store";

const NotificationDrawer = () => {
  const { isOpen, notifications } = useSelector(
    (store: RootState) => store.notifications
  );
  const dispatch = useAppDispatch();

  const { colorScheme } = useMantineColorScheme();

  const getMIntent = (
    type: NotificationType
  ): { icon: ReactNode; color: string } => {
    switch (type) {
      case NotificationType.NOTIFY:
        return {
          icon: <Info weight="fill" size={20} />,
          color: colorScheme === "dark" ? "brand.2" : "brand"
        };
      case NotificationType.DANGER:
        return {
          icon: <WarningCircle weight="fill" size={20} />,
          color: "red"
        };
      case NotificationType.SUCCESS:
        return {
          icon: <CheckCircle weight="fill" size={20} />,
          color: "green"
        };
      case NotificationType.WARNING:
        return {
          icon: <Warning weight="fill" size={20} />,
          color: "yellow"
        };
    }
  };

  return (
    <Drawer
      opened={isOpen}
      onClose={() => dispatch(toggleDrawer())}
      styles={{
        title: {
          width: "100%",
          marginRight: 0
        }
      }}
      title={
        <Group position="apart">
          <Text weight="bold" size="lg">
            Notifications
          </Text>
          <Tooltip label="Clear Notifications" position="left">
            <ActionIcon
              color="brand"
              size="md"
              disabled={notifications.length === 0}
              onClick={() => dispatch(clearNotifications())}>
              <BellSimpleSlash weight="bold" size={18} />
            </ActionIcon>
          </Tooltip>
        </Group>
      }
      withCloseButton={false}
      position="right"
      padding="sm"
      size={350}>
      {notifications.length === 0 ? (
        <Center>
          <Text size="sm" color="dimmed">
            No notifications.
          </Text>
        </Center>
      ) : (
        <Stack spacing="xs">
          {notifications.map((notif) => (
            <div key={notif.timestamp}>
              <Tooltip
                position="left"
                label={new Date(notif.timestamp).toLocaleTimeString("en-US", {
                  timeStyle: "medium"
                })}>
                <Text
                  color="dimmed"
                  size="xs"
                  weight={500}
                  style={{ marginBottom: "0.25rem" }}>
                  {new Date(notif.timestamp).toLocaleTimeString("en-US", {
                    timeStyle: "short"
                  })}
                </Text>
              </Tooltip>
              <Alert
                {...getMIntent(notif.appearance)}
                title={notif.title}
                variant="outline"
                radius="md"
                withCloseButton
                onClose={() => dispatch(dismissNotification(notif))}>
                {notif.detail}
              </Alert>
            </div>
          ))}
        </Stack>
      )}
    </Drawer>
  );
};

export default NotificationDrawer;
