import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { Notification } from "./messages";

interface NotificationState {
  isOpen: boolean;
  notifications: Notification[];
}

const initialState: NotificationState = {
  isOpen: false,
  notifications: []
};

const notificationSlice = createSlice({
  name: "notifications",
  initialState,
  reducers: {
    addNotification(state, action: PayloadAction<Notification>) {
      state.notifications = [action.payload, ...state.notifications];
    },
    dismissNotification(state, action: PayloadAction<Notification>) {
      state.notifications = state.notifications.filter(
        (x) => x.timestamp !== action.payload.timestamp
      );
    },
    toggleDrawer(state) {
      state.isOpen = !state.isOpen;
    },
    clearNotifications(state) {
      state.notifications = [];
    }
  }
});

export const {
  addNotification,
  dismissNotification,
  toggleDrawer,
  clearNotifications
} = notificationSlice.actions;
export default notificationSlice.reducer;
