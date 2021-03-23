const { MongoClient } = require("mongodb");
const cors = require("cors");
const express = require("express");
const app = express();
app.use(cors());

const port = 3000;
const uri = "mongodb://mongo:27017?retryWrites=true&writeConcern=majority";
const client = new MongoClient(uri);
client.connect().then(conn => {
  app.get("/api/:id", (req, res) => {
    run(req.params.id).then((b) => {
      res.send(b);
    });
  });
  async function run(id) {
    const database = conn.db("dc");
    const collection = database.collection("storage");
    const query = { id: id };
    console.log(id)
    return await collection.findOne(query);
  }
}).catch(err => {
  process.exit(1)
});

app.listen(port, () => {
  console.log(`Example app listening at http://localhost:${port}`);
});
