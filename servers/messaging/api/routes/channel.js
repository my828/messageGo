const mongoose = require("mongoose")
const Channel = require("../models/channel")
const Message = require("../models/message")
const ObjectID = require('mongodb').ObjectID;

function channel(req, res, next) {
    //initialize list of channels to store current user's channels
    var user = getUser(req)
    switch(req.method) {
        case "GET":
            Channel.find({
                $and: [
                    {$or: [{members: {$in: user.id}}, {creator: user.id}]},
                    {private: true}
                ]
            })
            .then(
                channels => {
                    res.status(200).json(channels);
                }
            )
            .catch(err => res.status(404).send('No channel found'+err))
            break;
        case "POST": 
            Channel.create(req.body).then(
                channel => res.status(201).json(channel)
            ).catch(err => res.status(404).send('Cannot create channel' + err));
            break;
        default: 
            res.status(405).send("Unsupported Method Type!")
            break;
    }
}

function channelID(req, res, next) {
    var id = req.params.channelID
    var query = req.query.before
    var user = getUser(req)
    var channel = checkAuth(user, id)
    if (user === null) {
        res.status(401).send("Unauthroized User")
        return
    }
    if (channel === null || typeof channel === "string") {
        res.status(403).send("Forbidden User")
        return
    }

    console.log(req.params.before)
    switch(req.method) {
        case "GET":
            if (query !== "") {

            }
            Message.find({channelID: ObjectID(id)})
            .sort({$natural: -1})
            .limit(100)
            .then(
                channels => res.status(200).json(channels)
            )
            .catch(
                err => res.status(404).send('Cannot find user with id' + err)
            );
            break;
        case "POST": 
            var message = new Message({
                channelID: ObjectID(id),
                body: req.body.body,
                creator: user.id
            })
            message.save()
                .then(
                    message => res.status(201).json(message)
                )
                .catch(
                    err => res.status(404).send('Cannot create channel ' + err)
            );
            break;
        case "PATCH": 
            var newChannel = {}
            if (req.body.name) newChannel.name = req.body.name; 
            if (req.body.description) newChannel.description = req.body.description;
            Channel.findOneAndUpdate(
                {
                    _id: id
                }, 
                {
                    $set: newChannel
                }, 
                {
                    returnNewDocument: true
                }, function(err, doc) {
                    if (err) {
                        res.status(404).send("Unable update channel: " + err)
                    } else {
                        res.status(201).json(doc)
                    }
                }
            )
            break;
        case "DELETE": 
            var error = ""
            var isGeneral = false
            if (user.id == "-1") {
                isGeneral = true
            } else {
                Channel.findOneAndDelete({_id: ObjectID(id)})
                .then()
                .catch(err => error = error + " " + err)
                Message.deleteMany({channelID: ObjectID(id)})
                .then()
                .catch(err => error = error + " " + err)
            }
            if (error !== "") {
                res.status(404).send('unable to delete channel: ' + err)
            } else if (isGeneral) {
                res.status(404).send('Please do not delete the general channel!')
            } else {
                res.status(200).send("Delete channel and its messages!")
            }
            break;
        default: 
            res.status(405).send("Unsupported Method Type!")
            break;
    }
}

function channelMembers(req, res, next) {
    var id = req.params.channelID
    var user = getUser(req)
    var channel = checkAuth(user, id)
    if (user === null) {
        res.status(401).send("Unauthroized User")
        return
    }
    if (channel === null || typeof channel === "string") {
        res.status(403).send("Forbidden User")
        return
    }
    switch (req.method) {
        case "POST":
            Channel.updateMany(
                {_id: id}, 
                {$push: {"members": req.body.id}}, 
                {
                    returnNewDocument: true
                }, function(err, doc) {
                    if (err) {
                        res.status(404).send("Unable update member: " + err)
                    } else {
                        res.status(201).json(doc)
                    }
                }
            )
            break;
        case "DELETE":
            Channel.updateMany(
                {
                    _id: id
                },
                {
                    $pull: {
                        members: req.body.id
                    }
                },
                {
                    returnNewDocument: true
                }, function(err, doc) {
                    if (err) {
                        res.status(404).send("Unable delete member: " + err)
                    } else {
                        res.status(201).json(doc)
                    }
                }
            )
            break;
        default: 
            res.status(405).send("Unsupported Method Type!")
            break;
    }
}

function getUser(req) {
    var user = req.get("X-User")
    if (user) {
        user = JSON.parse(user)
    } else {
        return null
    }
    return user
}

function checkAuth(user, id) {
    Channel.findOne({
        $and: [
            {$or: [{members: {$in: user.id}}, {creator: user.id}]},
            {private: true}, 
            {_id: id}
        ]
    }).then(
        // channel => {if (channel === null) {
        //     res.status(403).send("Forbidden User")
        // }
    //}
    channel => {return channel}
    ).catch(err => {return "Error: " + err})
}

module.exports = {channel, channelID, channelMembers};