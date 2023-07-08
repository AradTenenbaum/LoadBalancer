const serverLog = require("../utils/log");

const logs = (req, res, next) => {
  res.on("finish", function () {
    serverLog(
      {
        method: req.method,
        path: req.url,
        data: req.body,
        cookies: req.cookies,
        status: res.statusCode,
        message: res.statusMessage,
      },
      "REQUEST"
    );
  });
  next();
};

module.exports = logs;
