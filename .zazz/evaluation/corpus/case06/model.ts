import {Model,DataTypes,Sequelize} from 'sequelize';
const database=new Sequelize('sqlite::memory:');
export class Account extends Model {}
Account.init({handle:{type:DataTypes.STRING}},{sequelize:database});
