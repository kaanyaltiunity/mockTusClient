(async () => {
    let bucket
    const client = require("./client")
    try {
        bucket = await client.createBucket() // use if necessary
        console.log("BUCKET ID ",bucket.data.id)
        await client.deleteBucket(bucket.data.id)
        // const entry = await client.createEntry(bucket.data.id) //local
        // const entry = await client.createEntry("c951bc9a-72da-4c84-a62e-0b74b6aa1584") //staging
        // await client.uploadContent(bucket.data.id, entry.data.entryid) //local
        // await client.uploadContent("0844e1d3-a8b1-4488-83ae-93c125816a2c", entry.data.entryid) //staging
    } catch (err) {
        console.error(err.message)
    }
})()