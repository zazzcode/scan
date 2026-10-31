import express from 'express';
import {Sequelize} from 'sequelize';
const database = new Sequelize('sqlite::memory:');
const api = express();
api.post('/session', (req, res) => {
  let value = req.body.handle;
  if (typeof value !== 'string' || !/^[a-z0-9]+$/.test(value)) return res.status(400).send('bad');
  value = req.body.alternate;
  return database.query(`SELECT id FROM accounts WHERE handle = '${value}'`);
});
