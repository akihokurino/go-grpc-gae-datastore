import {Change, EventContext} from "firebase-functions/lib/cloud-functions";
import {DataSnapshot} from "firebase-functions/lib/providers/database";
const {PubSub} = require('@google-cloud/pubsub');


type Snapshot = { [index: string]: any; };
const pubsubClient = new PubSub();

const syncMessage = async (change: Change<DataSnapshot>, context: EventContext) => {
  const before: Snapshot = change.before.val();
  const after: Snapshot = change.after.val();
  const roomId: string = context.params.roomId;

  let newKeys: string[];
  if (before) {
    newKeys = getArrayDiff(Object.keys(before), Object.keys(after));
  } else {
    newKeys = Object.keys(after);
  }

  const tasks: Promise<void>[] = newKeys.map((key: string): Promise<void> => {
    const data: Snapshot = after[key];
    return sendMessage(
      key,
      roomId,
      data['fromId'],
      data['toId'],
      data['text'],
      data['imageUrl'],
      data['fileUrl'],
      data['createdAt']
    );
  });

  return Promise.all(tasks);
};

const getArrayDiff = (arr1: string[], arr2: string[]): string[] => {
  const arr: string[] = arr1.concat(arr2);
  return arr.filter((v: string): boolean => {
    return !(arr1.indexOf(v) !== -1 && arr2.indexOf(v) !== -1);
  });
};

const sendMessage = (id: string,
                     roomId: string,
                     fromId: string,
                     toId: string,
                     text: string,
                     imageUrl: string,
                     fileUrl: string,
                     createdAt: string): Promise<void> => {

  const data = {
    id,
    roomId,
    fromId,
    toId,
    text,
    imageUrl,
    fileUrl,
    createdAt
  };

  console.log("send message...");
  console.log(JSON.stringify(data));

  return pubsubClient
    .topic("message")
    .publish(Buffer.from(JSON.stringify(data)));
};

export default syncMessage;