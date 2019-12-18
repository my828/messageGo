//put interpreter into strict mode
"use strict";

//require the express and morgan packages
const express = require("express")
const mongoose = require("mongoose")
const bodyParser = require("body-parser")
const Channel = require("./api/models/channel")
const autoIncrement = require('mongoose-auto-increment')
//create a new express application
const app = express();

let channel = {};

var amqp = require('amqplib/callback_api');
amqp.connect('amqp://' + process.env.QNAME + ":5672/", function(err, conn) {
  conn.createChannel(function(err, ch) {
    let q = process.env.QNAME 

    ch.assertQueue(q, {durable: false});
    // Note: on Node 6 Buffer.from(msg) should be used
    channel.chan = ch
    channel.q = q
  });
});

module.exports = channel

var connection = mongoose.createConnection("mongodb://mongodb:27017/app")
autoIncrement.initialize(connection)

connection.then(() => {
    console.log("MongoDB connected!")
    Channel.findOneAndUpdate(
    {           
        name: "general",
        private: true,
        creator: -1
    }, 
    {
        upsert: true,
    })
    .then( 
        channel => console.log("Saved general channel: " + channel)
    )
    .catch(
        err => console.log(err)
);
}).catch(err => console.log(err))

// connect to mongodb
// mongoose.connect(
//         'mongodb://mongodb:27017/app', 
//         { 
//             userNewUrlParser: true,
//             useFindAndModify: false
//         }
//     )
//     .then(() => {
//         console.log("MongoDB connected!")
//         Channel.findOneAndUpdate(
//         {           
//             name: "general",
//             private: true,
//             creator: -1
//         }, 
//         {
//             upsert: true,
//         })
//         .then( 
//             channel => console.log("Saved general channel: " + channel)
//         )
//         .catch(
//             err => console.log(err)
//     );
//     })
//     .catch(err => console.log(err)) 

let addr = process.env.ADDR || ":3000";

const [host, port] = addr.split(":");

//autoIncrement.initialize()

// app.use(morgan("dev"));
app.use(bodyParser.urlencoded({extended: false}))
app.use(bodyParser.json());

const channelRoute = require('./api/routes/channel')
const messageRoute = require('./api/routes/message')

app.use('/v1/messages/:messageID', messageRoute.message)
app.use('/v1/channels/:channelID/members', channelRoute.channelMembers)
app.use('/v1/channels/:channelID', channelRoute.channelID)
app.use('/v1/channels', channelRoute.channel)


app.use((req, res, next) => {
    const error = new Error("Not Found");
    error.status = 404;
    next(error);
})

// catch all kinds of error that reach here
app.use((error, req, res, next) => {
    res.status(error.status || 500).json(error.message);
})

app.listen(port, host, () => {
    console.log('Server is listening at http://{port}...');
})
