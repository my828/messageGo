const mongoose = require('mongoose')

const messageSchema = mongoose.Schema({
    channelID: {
        type: Number,
        required: [true, 'Channel ID is required']
    },
    body: {
        type: String,
        default: "",
    },
    createdAt: {
        type: Date,
        default: Date.now,
    },
    creator: {
        id: { type: Number, required: true},
        userName: { type: String, required: true},
        firstName: { type: String, required: true},
        lastName: { type: String, required: true}
    },
    editedAt: {
        type: Date,
        dafault: Date.now
    }
})

module.exports = mongoose.model("Message", messageSchema)