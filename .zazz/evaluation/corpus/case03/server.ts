import express from 'express';
import {Sequelize} from 'sequelize';
const database=new Sequelize('sqlite::memory:');
const api=express();
function accepts(value){return typeof value==='string'&&/^[a-z0-9]+$/.test(value)}
api.post('/session',(req,res)=>{if(!accepts(req.body.handle)){return res.status(400).send('bad')}database.query(`SELECT id FROM accounts WHERE handle = '${req.body.handle}'`)});
