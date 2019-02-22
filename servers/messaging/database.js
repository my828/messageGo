const mongoose = require('mongoose')


var url = 'mongodb://mongo/mydb'
var testSchema = mongoose.Schema({
    name: String
})
Test = exports.Test = mongoose.model('Test', testSchema)

exports.initializeMongo = function() {
    mongoose.connect(url)
    console.log("Trying to connect to " + url)
    var db = mongoose.connection;
    db.on('error', console.error.bind(console, 'connection error'))
    db.once('open', function(){
        console.log("we are connected to the database!")
        addTest();
    })
}
var addTest = function() {
    var test = new Test({
        name: 'Hellooo'
    })

    test.save(function (err, fluffy) {
        if (err) return console.err(err);
        console.log('Successfully added test!')
    })
}