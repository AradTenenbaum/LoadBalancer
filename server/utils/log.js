function serverLog(details, type) {
  console.log({ type: type || "INFO", ...details });
}

module.exports = serverLog;
