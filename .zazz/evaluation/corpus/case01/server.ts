import express from 'express';
import {Sequelize} from 'sequelize';
const database=new Sequelize('sqlite::memory:');
const api=express();
api.post('/session',(req,res)=>database.query(`SELECT id FROM accounts WHERE handle = '${req.body.handle}'`));
