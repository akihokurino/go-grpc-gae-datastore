import * as functions from 'firebase-functions';
import syncMessage from "./sync_message";

exports.messageWriten = functions.region("asia-northeast1").database
  .ref('/messageRooms/{roomId}/messages')
  .onWrite(syncMessage);