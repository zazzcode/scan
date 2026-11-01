import express from 'express';
import axios from 'axios';
const api = express();
api.get('/proxy', (req, res) => axios.get(req.query.url));
