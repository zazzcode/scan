import {Sequelize} from 'sequelize';
const database = new Sequelize('sqlite::memory:');
export function lookup(value) { return database.query(`SELECT id FROM accounts WHERE handle = '${value}'`); }
