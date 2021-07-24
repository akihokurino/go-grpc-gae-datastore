"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const functions = require("firebase-functions");
const sync_message_1 = require("./sync_message");
exports.messageWriten = functions.region("asia-northeast1").database
    .ref('/messageRooms/{roomId}/messages')
    .onWrite(sync_message_1.default);
//# sourceMappingURL=index.js.map