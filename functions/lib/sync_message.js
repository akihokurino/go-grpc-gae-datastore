"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
const { PubSub } = require('@google-cloud/pubsub');
const pubsubClient = new PubSub();
const syncMessage = (change, context) => __awaiter(void 0, void 0, void 0, function* () {
    const before = change.before.val();
    const after = change.after.val();
    const roomId = context.params.roomId;
    let newKeys;
    if (before) {
        newKeys = getArrayDiff(Object.keys(before), Object.keys(after));
    }
    else {
        newKeys = Object.keys(after);
    }
    const tasks = newKeys.map((key) => {
        const data = after[key];
        return sendMessage(key, roomId, data['fromId'], data['toId'], data['text'], data['imageUrl'], data['fileUrl'], data['createdAt']);
    });
    return Promise.all(tasks);
});
const getArrayDiff = (arr1, arr2) => {
    const arr = arr1.concat(arr2);
    return arr.filter((v) => {
        return !(arr1.indexOf(v) !== -1 && arr2.indexOf(v) !== -1);
    });
};
const sendMessage = (id, roomId, fromId, toId, text, imageUrl, fileUrl, createdAt) => {
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
exports.default = syncMessage;
//# sourceMappingURL=sync_message.js.map