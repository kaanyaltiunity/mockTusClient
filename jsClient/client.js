
require("dotenv").config()
// const projectId = "cf15bc51-eb49-40c4-be2c-14b766982f2d" //staging
const projectId = "714aa7ae-da40-4b75-8bd2-048c80dd3d93" //local
const xRequestId = "3b98efb2-8ead-495a-a52b-a26421df78b6"
const {v4: uuidv4} = require("uuid")
const xUserVersion = 1
const xUser = {
    genesisId: "3573572419045",
    roles: ["unity-dashboard-developer"],
    resourceRoles: "",
    uuid: "0c85494e-f42d-4157-0000-03400983bde5",
    adsId: "asdfasdfasdf"
}
const xUserType = "USER"
// const gatewayStagingBaseUrl = `https://staging.services.unity.com/api/ccd/management/v1/projects/${projectId}`
// const localGatewayBaseURL = `http://localhost:9000/api/ccd/management/v1/projects/${projectId}`
const cdsBaseURL = "http://localhost:22080/api/v1"
// const cdsBaseURL = "https://content-api-stg.cloud.unity3d.com/api/v1"
const baseUrl = cdsBaseURL
console.log("API_TOKEN", process.env.API_KEY)
const baseHeaders = {
    // "X-User-Type": xUserType,
    // "X-User": JSON.stringify(xUser),
    // "X-User-Version": xUserVersion,
    // "X-Request-Id": xRequestId,
    "Authorization": `Basic ${process.env.API_KEY}` //local test
    // "Authorization": `Bearer ${process.env.BEARER_TOKEN}` //staging bearer
}

const axios = require("axios")
const jws = require("jws")
const fs = require("fs")
const path = require("path")
const crypto = require("crypto")
const tus = require("tus-js-client")


async function createBucket() {
    console.log("Creating bucket")
    console.log(uuidv4().split("-"))
    const createBucketUrl = `${baseUrl}/projects/${projectId}/buckets`
    const data = {
        description: "testDescription",
        name: `testBucketWithName-${uuidv4().split("-")[0]}`,
        projectguid: projectId
    }
    const headers = {
        ...baseHeaders,
        "Content-Type": "application/json"
    }
    headers["X-Jws-Signature"] = makeDetachedJWS(headers, JSON.stringify(data))
    const bucket = await axios.post(createBucketUrl, data, { headers })
    return bucket
}

async function deleteBucket(bucketId) {
    console.log("Deleting bucket")
    const deleteBucketUrl = `${baseUrl}/buckets/${bucketId}`
    const headers = {
        ...baseHeaders,
        "Content-Type": "application/json"
    }
    headers["X-Jws-Signature"] = makeDetachedJWS(headers)
    await axios.delete(deleteBucketUrl, { headers })
}

function createHash() {
    return new Promise((resolve, reject) => {
        const hash = crypto.createHash("MD5")
        const readStream = fs.createReadStream(path.resolve(__dirname, "./hello.txt"))
        readStream.on("data", (data) => {
            hash.update(data)
        }).on("end", () => {
            return resolve(hash.digest("hex"))
        }).on("error", (error) => {
            reject(error)
        })
    })
}

async function createEntry(bucketId) {
    console.log("\n==============================================")
    console.log("Creating Entry")

    const createEntryUrl = `${baseUrl}/buckets/${bucketId}/entries`

    const hashMD5 = await createHash()
    const fileSize = fs.statSync(path.resolve(__dirname, "./hello.txt")).size
    
    const headers = {
        ...baseHeaders,
        "Content-Type": "application/json",
    }
    const data = {
        path: `sample${uuidv4()}.pdf`,
        content_hash: hashMD5,
        content_size: fileSize,
        content_type: "applicaton/pdf"
    }
    // headers["X-Jws-Signature"] = makeDetachedJWS(headers, JSON.stringify(data))

    console.log("CREATING ENTRY")
    console.log("FILE HASH", hashMD5)
    console.log("FILE SIZE", fileSize)
    console.log("HEADERS", headers)
    console.log("==============================================\n")

    return await axios.post(createEntryUrl, data, { headers })
}

async function uploadContent(bucketId, entryId) {
    const chunkSize = 5 * 1024 * 1024
    console.log("UPLOADING")
    console.log("BUCKET ID", bucketId)
    console.log("ENTRY ID", entryId)
    // const contentUploadUrl = `http://localhost:9000/api/ccd/management/v1/test`
    const contentUploadUrl = `${baseUrl}/buckets/${bucketId}/entries/${entryId}/content/`
    const fileBuf = fs.readFileSync(path.resolve(__dirname, "./hello.txt"))
    const headers = {
        ...baseHeaders,
        // "Content-Type": "application/offset+octet-stream"
    }
    // headers["X-Jws-Signature"] = makeDetachedJWS(headers)
    await promiseTus(fileBuf, { chunkSize, headers, endpoint: contentUploadUrl })
}

function promiseTus(fileBuf, options) {
    return new Promise((resolve, reject) => {
        const uploader = new tus.Upload(fileBuf, {
            ...options,
            onError: reject,
            onSuccess: resolve
        })
        uploader.start()
    })
}

function makeDetachedJWS(headers, data) {
    let requestBody;
    if (headers["Content-Type"] === "multipart/form-data" || headers["Content-Type"] === "application/offset+octet-stream" || !data) {
        requestBody = '';
    } else {
        requestBody = data;
    }
    const user = headers['X-User'] || '';
    const payload = requestBody + user;
    const expiresInFiveMinutes = Math.floor(Date.now() / 1000) + 60 * 40;

    const signature = jws.sign({
        header: { alg: 'HS256', exp: expiresInFiveMinutes },
        payload,
        secret: "test-secret-needs-to-be-32-chars"
    });

    return signature
        .split('.')
        .map((part, index) => (index === 1 ? '' : part))
        .join('.');
}


module.exports = {
    createBucket,
    createEntry,
    deleteBucket,
    uploadContent
}