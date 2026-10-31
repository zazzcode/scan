import express from 'express';
import * as resources from 'finale-rest';
import {Account} from './model';
const api=express();
api.use(externalGuard());
api.post('/accounts',(req,res,next)=>next());
resources.resource({model:Account,endpoints:['/accounts','/accounts/:id']});
