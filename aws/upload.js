const { S3Client, PutObjectCommand, GetObjectCommand } = require("@aws-sdk/client-s3");
const { getSignedUrl } = require("@aws-sdk/s3-request-presigner");


const s3Client = new S3Client({
  region: process.env.AWS_REGION, 
  credentials: {
    accessKeyId:  process.env.AWS_ACCESS_KEY_ID,
    secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY,
  },
});

 async function uploadFile(key,contentType) {
  const command = new PutObjectCommand({
    Bucket: process.env.S3_BUCKET_NAME,
    Key: key,
    ContentType: contentType,
  });
  const url=await getSignedUrl(s3Client, command);
  console.log(url);
    return url;
}
async function getObjectUrl(key){
    const command = new GetObjectCommand({
        Bucket: process.env.S3_BUCKET_NAME,
        Key: key,
    });
    const url=await getSignedUrl(s3Client, command);
    console.log(url);
        return url;
}
module.exports = { uploadFile,getObjectUrl };