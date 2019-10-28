import {
  notification,
  message as topMessage
} from 'antd'

export const MessageTypes = {
  ERROR: 0,
  CONNECT: 1,
  SEARCH: 2,
  DOWNLOAD: 3,
  SERVERS: 4,
  WAIT: 5,
  IRCERROR: 6
}

// Message router returns new state objects
// It handles delegation of data to the App component
// depending on the JSON message type
export function messageRouter(message, currentState) {
  if ('error' in message) {
    console.error("ERROR: ABORTING");
    sendNotification("error", "Error", message.details)
    return
  }

  switch (message.type) {
    case MessageTypes.ERROR:
      sendNotification("error", "Error Processing Request", message.details);
      currentState.searchQueries.pop();
      window.localStorage.setItem("queries", JSON.stringify(currentState.searchQueries));
      return {
        loading: false, searchQueries: currentState.searchQueries
      };
    case MessageTypes.IRCERROR:
      sendNotification("error", "Internal Book Server Error", message.status);
      currentState.searchQueries.pop();
      window.localStorage.setItem("queries", JSON.stringify(currentState.searchQueries));
      return {
        loading: false, searchQueries: currentState.searchQueries
      };
    case MessageTypes.WAIT:
      topMessage.loading(message.status, 3)
      break;
    case MessageTypes.SEARCH:
      sendNotification("success",
        "Search Results Received",
        "Select a book to download or search again.");
      window.localStorage.setItem("results", JSON.stringify([...currentState.searchResults, message.books]))
      return {
        items: message.books,
        searchResults: [...currentState.searchResults, message.books],
        loading: false,
      };
    case MessageTypes.DOWNLOAD:

      // Alert the user that the file has downloaded
      sendNotification("success", "Book File Received", message.name);
      // Create download link and click it prompting the save dialog
      saveByteArray(message.name, message.file);
      return {
        loading: false
      };
    case MessageTypes.SERVERS:
      sendNotification("success", "Book Server List Updated", "Only download books from active servers.")
      return {
        servers: message.servers
      };
    default:
      console.error("Unknown Server Message")
  }
}

export function countdownTimer(wait, callback) {
  let downloadTimer = setInterval(() => {
    // Decrement the timeLeft each tick
    wait--;
    callback(wait);
    notification.open({
      key: "timer",
      type: "info",
      message: "Please wait before searching",
      description: wait + " seconds",
      duration: 0
    });

    if (wait <= 0) {
      clearInterval(downloadTimer);
      notification.open({
        key: "timer",
        type: "success",
        message: "Server is ready",
        description: "Enter a search query to get started"
      });
    }
  }, 1000);
}

// saveByteArray creates a link and download popup for the returned file
function saveByteArray(fileName, byte) {
  let link = document.createElement('a');
  link.href = `data:application/octet-stream;base64,${byte}`;
  link.download = fileName;
  link.click();
};

export function sendNotification(type, message, description, duration = 4.5) {
  notification[type]({
    message,
    description,
    duration
  })
}