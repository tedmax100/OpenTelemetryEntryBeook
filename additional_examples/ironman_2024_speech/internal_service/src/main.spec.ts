import * as request from 'supertest';
import App from './app.service';
import { Sleep } from './sleep';
import JsonLogger from './json.logger';

describe('test echo service', () => {
  let testApp:App;
  let sleeper:Sleep;
  let jsonLoggerManager: JsonLogger;
  
  beforeAll(()=> {
    jsonLoggerManager = new JsonLogger();
    sleeper = { sleep: async(ts:number)=> {}};
    testApp = new App(sleeper, jsonLoggerManager);
  })
  it('test /echo respoonse with {message: echo}', async () => {
    return request(testApp.getApp()).get('/echo').expect(200).expect({
      message: 'echo'
    });
  })
})