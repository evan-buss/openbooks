const MessageTypes = {
  ERROR: -1,
  CONNECT: 0,
  SEARCH: 1,
  DOWNLOAD: 2
}

function messageRouter(message) {
  switch (message.type) {
    case MessageTypes.ERROR:
      console.log("ERROR")
      console.log(message.details)
      break
    case MessageTypes.CONNECT:
      console.log("CONNECT")
      console.log(message.status)
      break
    case MessageTypes.SEARCH:
      console.log("SEARCH")
      break
    default:
      console.log("INVALID")
  }
}

export {
  MessageTypes,
  messageRouter
}