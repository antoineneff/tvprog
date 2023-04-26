import fastify from 'fastify';
import { DateTime } from 'luxon';
import { fetchAndParseXML, print } from './utils.js';

const server = fastify();

let programs = null;

async function updatePrograms() {
  programs = await fetchAndParseXML();
}

async function getProgramsOfTheDay(json = false) {
  const DATE = DateTime.local().setZone('Europe/Paris');
  const TODAY = DATE.toISODate();
  const TOMORROW = DATE.plus({ days: 1 }).toISODate();

  if (!programs[TODAY]) {
    programs = await fetchAndParseXML();
  }
  if (!programs[TOMORROW]) {
    updatePrograms();
  }
  if (json) {
    return programs[TODAY];
  }
  return print(programs[TODAY]);
}

server.get('/', async () => {
  const data = await getProgramsOfTheDay();
  return data;
});

server.get('/json', async () => {
  const data = await getProgramsOfTheDay(true);
  return data;
});

server.get('/fetch', async () => {
  updatePrograms();
  return { updating: true };
});

const start = async () => {
  try {
    await server.listen({ port: 3000, host: '0.0.0.0' });
    console.log('API started ✓');
    updatePrograms();
  } catch (error) {
    console.log('Error starting the API', error);
    process.exit(1);
  }
};

start();
