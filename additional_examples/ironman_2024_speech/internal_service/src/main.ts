
import sleep from './sleep';
import App from './app.service';
import JsonLogger from './json.logger';

const jsonLoggerManger = new JsonLogger();
const app = new App({sleep}, jsonLoggerManger);
app.setListenPortAndCallBack(3000, () => {
  jsonLoggerManger.getLogger().info("app listen on 3000");
});
