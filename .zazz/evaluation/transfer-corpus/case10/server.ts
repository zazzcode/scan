import express from 'express';
import {lookup} from './store';
const api = express();
api.post('/session', (req, res) => lookup(req.body.handle));
