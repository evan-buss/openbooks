const MessageTypes = {
  ERROR: 0,
  CONNECT: 1,
  SEARCH: 2,
  DOWNLOAD: 3,
  SERVERS: 4,
  WAIT: 5
}

// Message router returns new state objects
// It handles delegation of data to the App component
// depending on the JSON message type
function messageRouter(message) {
  if (message.error) {
    console.error("ERROR: ABORTING")
    return
  }

  switch (message.type) {
    case MessageTypes.ERROR:
      console.log("ERROR")
      console.log(message.details)
      break
    case MessageTypes.WAIT:
      console.log("WAIT")
      return { status: message.status }
    case MessageTypes.CONNECT:
      console.log("CONNECT")
      // Need to show a loading indicator for 30 seconds before the user can search anything
      // setInterval(() => {
        // 
      // }, message.wait)
      return { connectionState: message.status }
    case MessageTypes.SEARCH:
      return { items: message.books, loading: false }
    case MessageTypes.DOWNLOAD:
      saveByteArray(message.name, message.file)
      return { loading: false };
    case MessageTypes.SERVERS:
      break
    default:
      console.error("Unkown Server Message")
  }
}

// saveByteArray creates a link and download popup for the returned file
function saveByteArray(fileName, byte) {
  var link = document.createElement('a');
  link.href = `data:application/octet-stream;base64,${byte}`
  link.download = fileName;
  link.click();
};

export {
  MessageTypes,
  messageRouter
}