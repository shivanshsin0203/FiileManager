
require('dotenv').config();
var express = require('express');
var app = express();
var getURL = require('./url.js');
const { uploadFile, getObjectUrl } = require('./upload.js');
app.set('port', (process.env.PORT || 5000));

app.use(express.static(__dirname + '/public'));


app.get('/geturl', async function(req, res) {
    try {
        const url=await uploadFile("testvideo", 'video/mp4');
        res.status(200).json(url);
    } catch (error) {
        console.log(error);
        res.status(500).json({
            msg: "Error",
            Error: "Error creating presigned URL"
        });
    }
});
app.get('/getpublicurl', async function(req, res) {
    try {
        const url=await getObjectUrl("testvideo");
        res.status(200).json(url);
    } catch (error) {
        console.log(error);
        res.status(500).json({
            msg: "Error",
            Error: "Error creating presigned URL"
        });
    }
});

app.listen(app.get('port'), function() {
    console.log('Node app is running on port', app.get('port'));
});