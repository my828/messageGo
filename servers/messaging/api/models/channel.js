const mongoose = require('mongoose')

const channelSchema = mongoose.Schema({
    name: {
        type: String,
        required: [true, 'Name field is required']
    },
    description: {
        type: String,
        default: "",
    },
    private: Boolean,
    members: [
        {
            id: { type: Number, required: true},
            userName: { type: String, required: true},
            firstName: { type: String, required: true},
            lastName: { type: String, required: true}
        }
    ],
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
//5c6f3e2b720675272fa0630e
//5c6f79fc3a6c3231a88951f1
module.exports = mongoose.model("Channel", channelSchema)