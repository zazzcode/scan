import express from 'express';
import * as resources from 'finale-rest';
import {Account} from './model';
const api=express();
function accepts(value){return typeof value==='string'&&value.trim().length>0}
api.post('/accounts',(req,res,next)=>{if(!accepts(req.body.handle)){return res.status(400).send('bad')}next()});
resources.resource({model:Account,endpoints:['/accounts','/accounts/:id']});
