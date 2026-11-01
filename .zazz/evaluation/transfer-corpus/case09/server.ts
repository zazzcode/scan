import express from 'express';
import {Sequelize} from 'sequelize';
const database = new Sequelize('sqlite::memory:');
const api = express();
function lookup(value) { return database.query(`SELECT id FROM accounts WHERE handle = '${value}'`); }
api.post('/session', (req, res) => lookup(req.body.handle));
