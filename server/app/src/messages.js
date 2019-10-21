import {notification} from 'antd'

const MessageTypes = {
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
function messageRouter(message) {
    if (message.error) {
        console.error("ERROR: ABORTING");
        sendNotification("error", "Error", message.details)
        return
    }

    switch (message.type) {
        case MessageTypes.ERROR:
            console.log("ERROR");
            console.log(message.details);
            sendNotification("error", "Error Processing Request", message.details);
            break;
        case MessageTypes.IRCERROR:
            console.log("IRC ERROR");
            sendNotification("error", "Internal Book Server Error", message.status);
            return {loading: false};
        case MessageTypes.WAIT:
            console.log("WAIT");
            sendNotification("info", message.status, "Please wait. Your request is being processed");
            break;
        // return { status: message.status }
        case MessageTypes.CONNECT:
            console.log("CONNECT");
            sendNotification("success", "Successfully Connected", message.status);
            break;
        case MessageTypes.SEARCH:
            sendNotification("success", "Search Results Received", "Select a book to download or search again.");
            return {
                items: message.books.sort((a, b) => {
                    if (a.server < b.server) {
                        return -1
                    }
                    if (a.server > b.server) {
                        return 1
                    }
                    return 0
                }),
                loading: false,
            };
        case MessageTypes.DOWNLOAD:
            saveByteArray(message.name, message.file);
            return {loading: false, status: "File Downloaded"};
        case MessageTypes.SERVERS:
            break;
        default:
            console.error("Unknown Server Message")
    }
}

// saveByteArray creates a link and download popup for the returned file
function saveByteArray(fileName, byte) {
    let link = document.createElement('a');
    link.href = `data:application/octet-stream;base64,${byte}`;
    link.download = fileName;
    link.click();
};

function sendNotification(type, message, description) {
    notification[type]({
        message: message,
        description: description
    })
}

export {
    MessageTypes,
    messageRouter,
    sendNotification
}