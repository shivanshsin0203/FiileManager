
require('dotenv').config();
var express = require('express');
var app = express();

const { uploadFile, getObjectUrl } = require('./upload.js');
app.set('port', (process.env.PORT || 5000));

app.use(express.static(__dirname + '/public'));


app.get('/geturl', async function(req, res) {
    try {
        const key=req.query.key;
        const url=await uploadFile(key, 'video/mp4');
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
        const url=await getObjectUrl("unwrapped-shivanshsin0203.mp41728402925303/playlist.m3u8");
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