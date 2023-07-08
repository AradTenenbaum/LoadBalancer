const express = require("express");
const logs = require("./middlewares/logs");
const app = express();
require("dotenv").config();

const port = process.env.PORT || 5000;

app.use(logs);

app.get("/", (req, res) => {
  res.status(200).send(`Hit ${port}`);
});

app.listen(port, () => {
  console.log(`Server is listening on port ${port}`);
});
