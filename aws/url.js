var AWS = require('aws-sdk');
var credentials = {
    accessKeyId: process.env.AWS_ACCESS_KEY_ID,
    secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY
};
AWS.config.update({
    credentials: credentials,
    region: process.env.AWS_REGION
});

var s3 = new AWS.S3();

exports.generatePresignedURL = function (req, res) {
    const filename = req.query.filename || 'video.mp4';
    const contentType = req.query.contentType || 'video/mp4';

    const params = {
        Bucket: process.env.S3_BUCKET_NAME,
        Expires: 3600, // URL expires in 1 hour
        Fields: {
            key: filename,
            'Content-Type': contentType,
             acl: 'private'
        },
        Conditions: [
            ['content-length-range', 0, 1000000000], // 0-1GB
            {'acl': 'private'},
            {'success_action_status': '201'},
            ['starts-with', '$key', ''],
            ['starts-with', '$Content-Type', 'video/'],
            {'x-amz-algorithm': 'AWS4-HMAC-SHA256'}
        ]
    };

    s3.createPresignedPost(params, function (err, data) {
        if (err) {
            console.log("Error", err);
            res.status(500).json({
                msg: "Error",
                Error: "Error creating presigned URL"
            });
        } else {
            res.status(200).json(data);
        }
    });
};