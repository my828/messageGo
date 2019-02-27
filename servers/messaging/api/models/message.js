const mongoose = require('mongoose')

const messageSchema = mongoose.Schema({
    channelID: {
        type: mongoose.Schema.Types.ObjectId,
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
        type: Number,
        required: true
    },
    editedAt: {
        type: Date,
        dafault: Date.now
    }
})

module.exports = mongoose.model("Message", messageSchema)