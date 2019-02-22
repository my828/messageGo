//put interpreter into strict mode
"use strict";

//require the express and morgan packages
const express = require("express")
const mongoose = require("mongoose")
//const morgan = require("morgan")
const bodyParser = require("body-parser")
//create a new express application
const app = express();

// connect to mongodb
mongoose.connect(
        'mongodb://localhost:27017/app', 
        { userNewUrlParser: true}
    )
    .then(() => console.log("MongoDB connected!"))
    .catch(err => console.log(err)) 

//mongoose.Promise = global.Promise;

//get ADDR environment variable,
//defaulting to ":80"
// const addr = process.env.ADDR || "4000";
// //split host and port using destructuring
// const [host, port] = addr.split(":");

// app.use(morgan("dev"));
// app.use(bodyParser.urlencoded({extended: false}))
app.use(bodyParser.json());

const channelRoute = require('./api/routes/channel')

app.use('/v1/channels/:channelID', channelRoute.channelID)
app.use('/v1/channels', channelRoute.channel)

app.use((req,res,next) => {
    const error = new Error("Not Found");
    error.status = 404;
    next(error);
})

// catch all kinds of error that reach here
app.use((error, req, res, next) => {
    res.status(error.status || 500).json(error.message);
})

app.listen(3000, () => {
    console.log('Server is listening at http://3000...');
})