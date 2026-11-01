import express from 'express';
import * as resources from 'finale-rest';
import {Account} from './model';
const api=express();
api.post('/accounts',(req,res,next)=>{if(req.body.confirm!==undefined){if(!req.body.handle.length){return res.status(400).send('bad')}}next()});
resources.resource({model:Account,endpoints:['/accounts','/accounts/:id']});
