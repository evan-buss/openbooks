import {
  ActionIcon,
  AppShell,
  Center,
  Group,
  Notification,
  ScrollArea,
  Stack,
  Text,
  Tooltip,
  useMantineColorScheme
} from "@mantine/core";
import { AnimatePresence, motion } from "framer-motion";
import { BellSimpleSlash, X } from "@phosphor-icons/react";
import { useSelector } from "react-redux";
import { NotificationType } from "../../state/messages";
import {
  clearNotifications,
  dismissNotification,
  toggleDrawer
} from "../../state/notificationSlice";
import { RootState, useAppDispatch } from "../../state/store";
import { defaultAnimation } from "../../utils/animation";

export function NotificationDrawer() {
  const notifications = useSelector(
    (store: RootState) => store.notifications.notifications
  );
  const dispatch = useAppDispatch();

  const { colorScheme } = useMantineColorScheme();

  const getIntent = (type: NotificationType): string => {
    switch (type) {
      case NotificationType.NOTIFY:
        return colorScheme === "dark" ? "blue.2" : "blue";
      case NotificationType.DANGER:
        return "red";
      case NotificationType.SUCCESS:
        return "green";
      case NotificationType.WARNING:
        return "yellow";
    }
  };

  return (
    <>
      <AppShell.Section p="sm">
        <Group justify="space-between">
          <Text fw="bold" size="lg">
            Notifications
          </Text>
          <Group gap="sm">
            <Tooltip label="Close" position="left">
              <ActionIcon color="dark" onClick={() => dispatch(toggleDrawer())}>
                <X weight="bold" size={18} />
              </ActionIcon>
            </Tooltip>
            <Tooltip label="Clear Notifications" position="left">
              <ActionIcon
                size="md"
                disabled={notifications.length === 0}
                onClick={() => dispatch(clearNotifications())}>
                <BellSimpleSlash weight="bold" size={18} />
              </ActionIcon>
            </Tooltip>
          </Group>
        </Group>
      </AppShell.Section>
      <AppShell.Section grow style={{ overflow: "auto" }}>
        {notifications.length === 0 ? (
          <Center>
            <Text size="sm" c="dimmed">
              No notifications.
            </Text>
          </Center>
        ) : (
          <Stack
            gap="xs"
            style={{
              overflow: "hidden",
              height: "100%"
            }}>
            <ScrollArea h="100%" p={16}>
              <AnimatePresence mode="popLayout">
                {notifications.map((notif) => (
                  <motion.div {...defaultAnimation} key={notif.timestamp}>
                    <Tooltip
                      position="bottom-start"
                      label={new Date(notif.timestamp).toLocaleDateString(
                        "en-US",
                        {
                          dateStyle: "long"
                        }
                      )}>
                      <Text c="dimmed" size="xs" fw={500} mb={4}>
                        {new Date(notif.timestamp).toLocaleTimeString("en-US", {
                          timeStyle: "short"
                        })}
                      </Text>
                    </Tooltip>
                    <Notification
                      color={getIntent(notif.appearance)}
                      mb={6}
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
            </ScrollArea>
          </Stack>
        )}
      </AppShell.Section>
    </>
  );
}
