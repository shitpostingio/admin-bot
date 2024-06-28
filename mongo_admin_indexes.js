//
use automod

// bans
// Index on TelegramID, descending
db.bans.createIndex({telegramid: -1})

// hostname
// Index on Host
db.hostnames.createIndex( { host: 1 }, { unique: true } )

// media
// Index on fileid
db.media.createIndex( { fileid: 1 }, { unique: true } )

// Index on fileuniqueid
db.media.createIndex( { fileuniqueid: 1 }, {  sparse: true, unique: true } )

// Multikey index on histogramaverage and histogramsum
db.media.createIndex( { histogramaverage: 1, histogramsum: 1 } )

// moderators
// Index on telegramid
db.moderators.createIndex({telegramid: 1}, { unique: true })

// sources
// Use sparse index to allow for null values while keeping the unique-ness
// Index on telegramid
db.sources.createIndex({telegramid: 1}, { sparse: true, unique: true })

// Index on username
db.sources.createIndex( { username: 1 }, { sparse: true, unique: true })

// stickerpacks
db.stickerpacks.createIndex( { setname: 1 }, { unique: true } )
