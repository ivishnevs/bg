# WS Event structure

 - We use JSON messages for client/server dialog
 - Each event should consist from two keys: **channel** (not required), **action**
 - Action is a hash with two keys too: **type** and **data**

 E.G.
 '''
 {
     "channel": "channelName",
     "action": {
        "type": "typeName",
        "data": "someData"
     }
 }
 '''

# Valid action types

 - subscribe.to
 - unsubscribe.from
 - role.selected
 - game.step-completed
 - gamer.step-completed
