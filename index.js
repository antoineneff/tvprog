import fastify from 'fastify'
import { DateTime } from 'luxon'
import { fetchAndParseXML, print } from './utils.js'

const server = fastify()

let programs = null

server.get('/', async () => {
  const TODAY = DateTime.local().setZone('Europe/Paris').toISODate()
  if (!programs[TODAY]) {
    programs = await fetchAndParseXML()
  }
  return print(programs[TODAY])
})

server.get('/fetch', async () => {
  programs = await fetchAndParseXML()
  return { updated: true  }
})

const start = async () => {
  try {
    await server.listen(3000, '0.0.0.0')
    console.log("API started ✓")
    programs = await fetchAndParseXML()
  } catch (error) {
    console.log("Error starting the API", error)
    process.exit(1)
  }
}

start()
