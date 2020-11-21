# Among Us Events - Discord Bot

A discord bot for organizing and planning among us games. This bot was born out of the pain of trying to organize weekly discord games with a group of friends. The bot allows the organizer to create an event in a given discord channel. Users of this discord server can rsvp to the event with 3 options (accept,decline,change time request).<br /><br />
<div style="text-align:center"><img src="https://i.imgur.com/xmKToMd.png" width="500" /></div><br /><br />

## Installation

Visit and authorize the bot to your discord server: <br />

    https://discord.com/oauth2/authorize?client_id=779556729343049729&scope=bot&permissions=485440

Once the bot is authorized and added to your discord server you will need to create a discord role called ``amongusbot`` the privileges of this role do not matter so you can make them as you see fit.

This role will need to be added to any user you would like to command the bot.<br />

## Usage

The bot responds to one main command by any user with the ``amongusbot`` role. This command creates a new event in the channel the users runs it in:

    !CreateAmongEvent TITLE OF EVENT HERE

Anything after `"!CreateAmongEvent "` will be the title of the event.

The bot will then create the event with the three options(message reactions) and any memeber of the discord server will be allow to select an option to rsvp.<br />

## Issues and Feature Requests

For any issues or feature requests please use the github issues on this repo.