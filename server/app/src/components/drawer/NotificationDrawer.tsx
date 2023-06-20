import {
  ActionIcon,
  Center,
  Drawer,
  Group,
  Notification,
  ScrollArea,
  Stack,
  Text,
  Tooltip,
  useMantineColorScheme
} from "@mantine/core";
import { AnimatePresence, motion } from "framer-motion";
import { BellSimpleSlash } from "@phosphor-icons/react";
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
        },
        body: { width: "100%", height: "calc(100% - 52px)", padding: 0 }
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
                    <Text color="dimmed" size="xs" weight={500} mb={4}>
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
    </Drawer>
  );
}
