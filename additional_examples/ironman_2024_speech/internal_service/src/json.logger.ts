
import winston, { createLogger, format, transports } from 'winston';
export default class JsonLogger {
  private jsonLogger: winston.Logger;
  constructor() {
    this.jsonLogger = createLogger({
      transports: [
        new transports.Console({
          format: format.json()
        })
      ],
      format: format.combine(format.timestamp())
    });
  }
  getLogger() {
    return this.jsonLogger;
  }
}