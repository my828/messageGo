const express = require("express")
const router = express.Router();
const Channel = require("../models/channel")
const Message = require("../models/message")
const ObjectId = require('mongodb').ObjectID;

function channel(req, res, next) {
    //initialize list of channels to store current user's channels
    var user = req.get("X-User")
    console.log("SHOULDN'T BE IN HERE")
    if (user) {
        // parse json here
        user = JSON.parse(user)
    } else {
        // return 401
        res.status(401).send("Unauthroized User")
    }
    switch(req.method) {
        case "GET":
        Channel.find({
            $and: [
                {members: {$in: [user.name]}}, 
                {private: {$eq: true}}, 
                {creator: user.id}
            ]
        })
            .then(
                channels => {
                    res.status(200).json(channels);
                }
            )
            .catch(err => res.status(404).json({msg: 'No channel found'}))
            break;

        case "POST": 
            // var channel = new Channel({

            // }) 
            // channel.save().then().catch();
            Channel.create(req.body).then(function(channel){
                res.status(201).json(channel)
            }).catch(err => res.status(404).json({msg: 'Cannot create channel'}));
            break;
        case "DELETE": 
            Channel.deleteMany({name: "Min4"}).exec((err)=> console.log(err))
            res.status(200).json({remove: "ok"})
            console.log("remove")
            break;
        default: 
            res.status(404).json({msg: "Unsupported Method Type!"})
            break;
    }
}

function channelID(req, res, next) {
    const id = req.params.channelID
    // if private channel and current user not a member
    // 403 
    console.log("In GET channel ID")
    var user = req.get("X-User")
    if (user) {
        user = JSON.parse(user)
    } else {
        res.status(401).send("Unauthroized User")
    }
    Channel.findOne({
        $and: [
            {members: {$in: [user.name]}}, 
            {private: {$eq: true}}, 
            {creator: {$eq: user.id}},
            {_id: [id]}
        ]
    }).then(
        (channel, err) => {if (channel === null) {
            res.status(403).json({msg: "Forbidden User"})
            }
    }).catch(err => res.status(404).json({msg: "Error finding user"}))
    switch(req.method) {
        case "GET":
        Channel.findById(id)
            .then(
                channels => {
                    res.status(200).json(channels)
                    }
            )
            .catch(err => res.status(404).json({msg: 'Cannot find user with id'}))
            break;
        case "POST": 
        var messageBody = new Message({
            body: req.body.body
        })
        Channel.save(req.body)
        res.status(201).json({
            
        })
        break;
        case "PATCH": 
        res.status(201).json({
            
        })
        case "DELETE": 
    }
}

module.exports = {channel, channelID};