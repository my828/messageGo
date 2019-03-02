const mongoose = require('mongoose')

const channelSchema = mongoose.Schema({
    name: {
        type: String,
        unique: true,
        required: [true, 'Name field is required']
    },
    description: {
        type: String,
        default: "",
    },
    private: Boolean,
    members: [],
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
// {
//     id: { type: Number, required: true},
//     userName: { type: String, required: true},
//     firstName: { type: String, required: true},
//     lastName: { type: String, required: true}
// }
module.exports = mongoose.model("Channel", channelSchema)