import * as express from 'express';
import { Router, Request, Response, NextFunction } from 'express';
import { Sleep } from "./sleep";
import JsonLogger from './json.logger';

class App {
  private sleeper: Sleep;
  private app;
  private jsonLoggerManger: JsonLogger; 
  constructor(sleeper: Sleep, jsonLogger: JsonLogger) {
    this.sleeper = sleeper;
    this.jsonLoggerManger = jsonLogger;
    this.init();
  }
  init() {
    this.app = express();
    const router = this.createRouter();
    this.app.use(router);
  }
  createRouter(): Router {
    const router = Router();
    const jsonLogger = this.jsonLoggerManger.getLogger();
    router.use((req: Request, res: Response, next: NextFunction) => {
      const startTime = Date.now();
      jsonLogger.info(null, {path: req.path, startTime: startTime });
      res.on('finish', () => {
        const duration = Date.now() - startTime;
        jsonLogger.info(null, { duration: `${duration}ms`, path: req.path });
      });
      next();
    })
    router.get('/echo', async (req: Request, res: Response) => {
      jsonLogger.info(null, { path: req.path, query: req.query });
      await this.sleeper.sleep(1000);
      const query = req.query;
      res.json({ message: 'echo', query: query });
    });
    return router;
  } 
  setListenPortAndCallBack(port:number, fn: Function) {
    this.app.listen(port, fn);
  } 
  getApp() {
    return this.app;
  }
}

export default App;