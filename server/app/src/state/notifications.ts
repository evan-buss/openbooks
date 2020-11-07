import { toaster } from "evergreen-ui";

export enum NotificationType {
    NOTIFY,
    SUCCESS,
    WARNING,
    DANGER
}

export const displayNotification = (type: NotificationType = NotificationType.NOTIFY, message: string, subText?: string) => {
    switch (type) {
        case NotificationType.NOTIFY:
            toaster.notify(message, {
                description: subText
            });
            break;
        case NotificationType.SUCCESS:
            toaster.success(message, {
                description: subText
            });
            break;
        case NotificationType.WARNING:
            toaster.warning(message, {
                description: subText
            });
            break;
        case NotificationType.DANGER:
            toaster.danger(message, {
                description: subText
            });
            break;
    }
}
