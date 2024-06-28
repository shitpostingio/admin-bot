
# Shitposting.io `admin-bot`

## Available features

- Ability to blacklist all kinds of media, stickers and sticker packs
- AI powered recognition of NSFW/suggestive content with automated removal
- Anti spam (an user can send a maximum of 11 text/media messages, 6 other messages or a total of 18 messages in a 10 second span)
- Anti userbot (analyzes every user that joins to check for similarities between recent joins)
- Anti flood (reduces API calls when under attack)
- User verification (requires user to press a button to verify they're human)
- Emergency mode (automatically restrict users that join, requiring approval from moderators)
- Automated deletion of non-whitelisted group/channel handles and links
- Automated deletion of messages forwarded from non-whitelisted channels
- Automated deletion of long messages (over 800 characters or with over 15 newlines)
- Automated deletion of commands to prevent spamming
- Automated reports and backups for various actions, including `@admin` mentions
- Logging of ban motivations and automated actions with the possibility to quickly undo them or to confirm them

## Available commands in groups

- `/ban` bans a user.
- `/banh` bans a user given their `@username`.
- `/idban` bans a user given their user id (subject to Telegram's limitations on visibility).
- `/mute` restricts a user from sending messages.
- `/nomedia` restricts a user from sending media, stickers and gifs.
- `/noother` restricts a user from sending stickers and gifs.
- `/blm` adds a media to the blacklist.
- `/bls` adds a sticker to the blacklist.
- `/blsp` adds a sticker pack to the blacklist.
- `/blh` adds one or more handles to the blacklist. **[ONLY FOR DB ADMINS]**

### Commands usage

All commands need to be sent as a reply to the message you want to act upon, otherwise they will be automatically deleted.

#### `/ban`

Bans a user. The syntax to use is `/ban motivation`. The command **will not** work if no motivation is provided. The motivation, along with additional data, will be stored in the database for future use.

#### `/banh`

Bans a user. The syntax to use is `/banh @username motivation`. The command **will not** work if no motivation is provided. The motivation, along with additional data, will be stored in the database for future use.

#### `/idban`

Bans a user. The syntax to use is `/idban userid motivation`. The command **will not** work if no motivation is provided. The motivation, along with additional data, will be stored in the database for future use.

#### `/mute`

Restricts a user to read only for a period of time. The syntax to use is `/mute [duration(e|w|d|h|m)]`.

The duration parameter is optional and, if omitted or the specified duration cannot be parsed, the bot will default the duration to 12 hours. **Restricting an user for under a minute will often lead to the restriction being permanent**.

#### `/nomedia`

Restricts a user from sending media, stickers and gifs for a period of time. The syntax to use is `/(nomedia|nopic) [duration(e|w|d|h|m)]`.

The duration parameter is optional and, if omitted or the specified duration cannot be parsed, the bot will default the duration to 12 hours. **Restricting an user for under a minute will often lead to the restriction being permanent**.

#### `/nosticker`

Restricts a user from sending stickers and gifs for a period of time. The syntax to use is `/nosticker [duration(e|w|d|h|m)]`.

The duration parameter is optional and, if omitted or the specified duration cannot be parsed, the bot will default the duration to 12 hours. **Restricting an user for under a minute will often lead to the restriction being permanent**.

### `/blm`

Blacklists a sticker/photo/video/audio/voice message/video message/animation.

### `/bls`

Blacklists a sticker.

### `/blsp`

Blacklists a sticker pack. In case the sticker does not have one, it'll blacklist the single sticker.

### `/blh`

Blacklists handles. **ONLY DATABASE ADMINS CAN USE THIS COMMAND**.

The bot will look for all handles present in a message: @mentions and t.me / telegram.me links. In case no handles are found, the bot will, in case the message has been forwarded, blacklist the handle of the original poster.

### Available feaures in a private conversation **[DB ADMINS ONLY]**

In a private conversation with the bot, an admin can perform all blacklist actions with a few additions.

- Blacklist multiple things at once by activating blacklist all mode with the command `/blacklistall`
- Blacklist handles by sending the handle as a text message to the bot
- Pardon blacklisted content
- Whitelist photos, videos and animations recognised as unsafe by the AI
- Whitelist a channel with the command `/whitelistchannel`.
- Remove a channel from the whitelist with `/removechannel`.
- Get user informations (ex. ban status, restriction status) by forwarding text messages of an user. In case the user is banned or restricted, a button for a quick unban/unrestriction will be provided as well.
- Mod an user in the chat by using `/mod userid`. Additional information on the user will be provided and a button will need to be clicked in order to complete the action.
- Activate an emergency mode with `/emergencymode [duration]`. For the specified amount of time, users without a profile picture or an username will automatically be limited and an alert will be sent on the report channel. Emergency mode can be toggled off at any time by sending `/emergencymode` again.
