import {
  ActionIcon,
  Center,
  Drawer,
  Group,
  Notification,
  Stack,
  Text,
  Tooltip,
  useMantineColorScheme
} from "@mantine/core";
import { AnimatePresence, motion } from "framer-motion";
import { BellSimpleSlash } from "phosphor-react";
import { useSelector } from "react-redux";
import { NotificationType } from "../../state/messages";
import {
  clearNotifications,
  dismissNotification,
  toggleDrawer
} from "../../state/notificationSlice";
import { RootState, useAppDispatch } from "../../state/store";
import { defaultAnimation } from "../../utils/animation";

export default function NotificationDrawer() {
  const { isOpen, notifications } = useSelector(
    (store: RootState) => store.notifications
  );
  const dispatch = useAppDispatch();

  const { colorScheme } = useMantineColorScheme();

  const getIntent = (type: NotificationType): string => {
    switch (type) {
      case NotificationType.NOTIFY:
        return colorScheme === "dark" ? "brand.2" : "brand";
      case NotificationType.DANGER:
        return "red";
      case NotificationType.SUCCESS:
        return "green";
      case NotificationType.WARNING:
        return "yellow";
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
        <Stack
          spacing="xs"
          style={{ overflow: "scroll", height: "calc(100% - 44px)" }}>
          <AnimatePresence mode="popLayout">
            {notifications.map((notif) => (
              <motion.div {...defaultAnimation} key={notif.timestamp}>
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
                <Notification
                  color={getIntent(notif.appearance)}
                  styles={{
                    root: {
                      boxShadow: "none"
                    }
                  }}
                  title={notif.title}
                  onClose={() => dispatch(dismissNotification(notif))}>
                  {notif.detail}
                </Notification>
              </motion.div>
            ))}
          </AnimatePresence>
        </Stack>
      )}
    </Drawer>
  );
}
