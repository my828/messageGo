const mongoose = require("mongoose")
const Channel = require("../models/channel")
const Message = require("../models/message")
const ObjectID = require('mongodb').ObjectID;
const ch = require("../../index")

function message(req, res, next) {
    var id = req.params.messageID
    var user = getUser(req)
    var message = checkAuth(user, id, res)
    if (user === null) {
        res.status(401).send("Unauthroized User")
        return
    }
    if (message === null || typeof message === "string") {
        res.status(403).send("Forbidden User")
        return
    }

    switch (req.method) {
        case "PATCH": 
        Message.findOneAndUpdate(
            {
                _id: id
            }, 
            {
                $set: {body: req.body.body}
            }, 
            {
                returnNewDocument: true
            }, function(err, doc) {
                if (err) {
                    res.status(404).send("Unable update message: " + err)
                } else {
                    res.status(201).json(doc)
                }
                Channel.find({_id: message.channelID}).then(channel => {chan = channel}).catch()

                e = {
                    type: "message-update",
                    message: doc,
                    userIDs: channel.private ? channel.members : []
                }
                ch.chan.sendToQueue(ch.q, new Buffer(JSON.stringify(e)));
            }
        )
        break;
        case "DELETE":
            let mes;
            Channel.find({_id: ObjectID(id)}).then(message => {mes = message}).catch()

            Message.remove({_id: ObjectID(id)})
            .then(
                () => {res.status(200).send("Delete message success!")
                e = {
                    type: "message-update",
                    messageID: id,
                    userIDs: channel.private ? channel.members : []
                }
                ch.chan.sendToQueue(ch.q, new Buffer(JSON.stringify(e)))
            })
            .catch(err => res.status(404).send('unable to delete message: ' + err))
        break;
        default: 
        break;
    }
}

function getUser(req) {
    var user = req.get("X-User")
    if (user) {
        user = JSON.parse(user)
    } else {
        res.status(401).send("Unauthroized User")
        return null
    }
    return user
}

function checkAuth(user, id) {
    Message.findOne({
        $and: [
            {_id: id},
            {creator: user.id}
        ]
    }).then(
        message => {if (message === null) {
            return message
        }
    }).catch(err => {return err})
}
module.exports = {message};