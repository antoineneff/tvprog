import got from 'got';
import camaro from 'camaro';
import { DateTime, Settings } from 'luxon';

// Luxon config
Settings.defaultLocale = 'fr';

// TV programs config
const XML_URL = 'https://xmltv.ch/xmltv/xmltv-tnt.xml';
const MINIMUM_PROGRAM_LENGTH = 35; // List programs that are at least 35 mins long
const ORDERED_CHANNELS = [
  'TF1',
  'France 2',
  'France 3',
  'Canal+',
  'France 5',
  'M6',
  'Arte',
  'C8',
  'W9',
  'TMC',
  'TFX',
  'NRJ 12',
  'France 4',
  'CSTAR',
  "L'Equipe",
  '6ter',
  'RMC Story',
  'RMC Découverte',
  'Chérie 25',
];
const MAX_CHANNEL_LENGTH =
  Math.max(...ORDERED_CHANNELS.map((channel) => channel.length)) + 2;

const channelsTemplate = [
  'tv/channel',
  {
    id: '@id',
    name: 'display-name',
  },
];
const programsTemplate = [
  'tv/programme',
  {
    title: 'title',
    channel: '@channel',
    start: '@start',
    end: '@stop',
  },
];

async function getXML() {
  const xml = await got(XML_URL).text();

  return xml;
}

async function parseXML(xml) {
  const channels = await camaro.transform(xml, channelsTemplate);
  const programs = await camaro.transform(xml, programsTemplate);

  return { channels, programs };
}

async function filterPrograms({ channels, programs }) {
  const filteredPrograms = programs
    .filter((program) => {
      const start = DateTime.fromFormat(program.start, 'yyyyMMddHHmmss ZZZ', {
        zone: 'Europe/Paris',
      });
      const end = DateTime.fromFormat(program.end, 'yyyyMMddHHmmss ZZZ', {
        zone: 'Europe/Paris',
      });

      return (
        start.plus({ minutes: MINIMUM_PROGRAM_LENGTH }) < end &&
        ((start.hour === 20 && start.minute > 45) ||
          (start.hour === 21 && start.minute < 20))
      );
    })
    .map((program) => {
      return {
        ...program,
        start: DateTime.fromFormat(
          program.start,
          'yyyyMMddHHmmss ZZZ',
        ).toMillis(),
        end: DateTime.fromFormat(program.end, 'yyyyMMddHHmmss ZZZ').toMillis(),
      };
    });

  return { channels, programs: filteredPrograms };
}

async function formatPrograms({ channels, programs }) {
  const formatted = programs.reduce((prev, program) => {
    const programDate = DateTime.fromMillis(program.start).toISODate();
    const channel = channels.find((channel) => channel.id === program.channel);
    const { title, start, end } = program;
    const formattedProgram = {
      title,
      start: DateTime.fromMillis(start, {
        zone: 'Europe/Paris',
      }).toLocaleString({ hour: '2-digit', minute: '2-digit' }),
      end: DateTime.fromMillis(end, { zone: 'Europe/Paris' }).toLocaleString({
        hour: '2-digit',
        minute: '2-digit',
      }),
    };
    if (!prev[programDate]) {
      prev[programDate] = { [channel.name]: formattedProgram };
    } else {
      prev[programDate] = {
        ...prev[programDate],
        [channel.name]: formattedProgram,
      };
    }
    return prev;
  }, {});

  return formatted;
}

export async function fetchAndParseXML() {
  const xml = await getXML();
  const parsed = await parseXML(xml);
  const filtered = await filterPrograms(parsed);
  const formatted = await formatPrograms(filtered);

  return formatted;
}

export function print(programs) {
  const LONGEST_TITLE_LENGTH = Math.max(
    ...Object.values(programs).map((program) => program.title.length),
  );
  const MIN_TITLE_LENGTH = 30;
  const MAX_TITLE_LENGTH = 55;
  const TITLE_LENGTH = Math.max(
    Math.min(LONGEST_TITLE_LENGTH, MAX_TITLE_LENGTH) + 2,
    MIN_TITLE_LENGTH,
  );
  const TIMETABLE_LENGTH = '00:00 - 00:00'.length + 2;
  const LINE_LENGTH = MAX_CHANNEL_LENGTH + TITLE_LENGTH + TIMETABLE_LENGTH + 4;
  const PROGRAM_TITLE = `PROGRAMME TV DU ${DateTime.local()
    .setZone('Europe/Paris')
    .toLocaleString()}`;
  let SPACES_BEFORE_TITLE = ' '.repeat(
    Math.ceil((LINE_LENGTH - PROGRAM_TITLE.length - 2) / 2) - 1,
  );
  let SPACES_AFTER_TITLE = SPACES_BEFORE_TITLE;
  if ((LINE_LENGTH - PROGRAM_TITLE.length) % 2 !== 0) {
    SPACES_AFTER_TITLE = ' '.repeat(
      (LINE_LENGTH - PROGRAM_TITLE.length - 2) / 2 - 1,
    );
  }

  let print = `${SPACES_BEFORE_TITLE}┌${'─'.repeat(
    PROGRAM_TITLE.length + 2,
  )}┐${SPACES_AFTER_TITLE}\n`;
  print += `${SPACES_BEFORE_TITLE}│ ${PROGRAM_TITLE} │${SPACES_AFTER_TITLE}\n`;
  print += `┌${'─'.repeat(MAX_CHANNEL_LENGTH)}┬${'─'.repeat(
    SPACES_BEFORE_TITLE.length - MAX_CHANNEL_LENGTH - 2,
  )}┴${'─'.repeat(PROGRAM_TITLE.length + 2)}┴${'─'.repeat(
    TITLE_LENGTH -
      PROGRAM_TITLE.length -
      (SPACES_BEFORE_TITLE.length - MAX_CHANNEL_LENGTH) -
      2,
  )}┬${'─'.repeat(TIMETABLE_LENGTH)}┐\n`;
  print += `│ Chaine${' '.repeat(
    MAX_CHANNEL_LENGTH - 'chaine'.length - 1,
  )}│ Titre${' '.repeat(
    TITLE_LENGTH - 'titre'.length - 1,
  )}│ Horaires${' '.repeat(TIMETABLE_LENGTH - 'horaires'.length - 1)}│\n`;
  print += `├${'─'.repeat(MAX_CHANNEL_LENGTH)}┼${'─'.repeat(
    TITLE_LENGTH,
  )}┼${'─'.repeat(TIMETABLE_LENGTH)}┤\n`;

  for (const channel of ORDERED_CHANNELS) {
    if (programs[channel]) {
      const { title, start, end } = programs[channel];
      const trimmedTitle =
        title.length > 55 ? title.substring(0, 52) + '...' : title;
      print += `│ ${channel}${' '.repeat(
        MAX_CHANNEL_LENGTH - channel.length - 1,
      )}│ ${trimmedTitle}${' '.repeat(
        TITLE_LENGTH - trimmedTitle.length - 1,
      )}│ ${start} - ${end} │\n`;
    }
  }

  print += `└${'─'.repeat(MAX_CHANNEL_LENGTH)}┴${'─'.repeat(
    TITLE_LENGTH,
  )}┴${'─'.repeat(TIMETABLE_LENGTH)}┘\n`;

  return print;
}
