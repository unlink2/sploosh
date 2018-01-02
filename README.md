# Sploosh

A simple Discord Bot that lets you play Sploosh! Kaboom! just like in Wind Waker!

## Setup

### Hosted version:

To add the bot to your Discord Server click this
[link](https://discordapp.com/api/oauth2/authorize?client_id=390599729215700992&permissions=0&scope=bot) .  Discord will then ask you which server to add the bot to.

Note: To disable sounds that the bot plays simply either mute it or add it to a group which does not allow voice activity!

### Running it yourself:

To get the bot yourself simply run

  go get github.com/unlink2/sploosh

### Configuration FIle

The bot expects a config file to be placed next to it's executable.

The file should be called

  config.ini

and look as like this:

  [discord]
  clientid=YOUR CLIENTID
  token=YOUR TOKEN
  [soundwhitelist]
  1=guildid1
  2=guildid2
  ...
  [clientids]
  1=clientid1
  2=clientid2
  ...
  [clientsecrets]
  1=clientsecret1
  2=clientsecret2
  ...

the sound whitelist, clientids and clientsecrets are optional. If left out
sounds will be enabled for every guild and the api will not require client-ids or secrets.

### Optional files

The bot can also play optional soundfiles. The files are called

  sounds/sploosh.dca

and

  sounds/kaboom.dca

These files will be played if the user using the bot is in  a voice channel.

### Optional Emotes

The following optional emotes can be added to your discord to make the bot look nicer.
If the emotes are not present emoji will be used instead.

  skMiddleH
  skMiddleV
  skFrontW
  skFrontN
  skFrontS
  skFrontE
  skBackW
  skBackN
  skBackS
  skBackE
  skBomb
  skBombUsed
  skNotSunk
  skSunk
  skBlank
  skSploosh
  skKaboom
