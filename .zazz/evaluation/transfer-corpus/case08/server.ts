import express from 'express';
import {Sequelize} from 'sequelize';
const database = new Sequelize('sqlite::memory:');
const api = express();
import {accepts} from './guard';
api.post('/session', (req, res) => {
  const value = req.body.handle;
  if (!accepts(value)) return res.status(400).send('bad');
  return database.query(`SELECT id FROM accounts WHERE handle = '${value}'`);
});
