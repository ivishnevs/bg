import Ember from 'ember';

const {
	Component,
  computed,
	inject,
	get
} = Ember;

export default Component.extend({
	ws: inject.service(),

  channel: computed(function() {
    return `game-${get(this, 'model.id')}-roles`;
  }),

	gamersSorting: ['role'],
	gamers: computed.sort('model.gamers', 'gamersSorting'),

	didReceiveAttrs() {
		get(this, 'ws').subscribeTo(get(this, 'channel'));
    get(this, 'ws').onMessage(this.messageHandler);
	},

  messageHandler(e) {
    let event = JSON.parse(e.data);

    if (event.action.type === 'role.selected') {
      $(`.${event.action.data}`).removeClass('enabled').addClass('disabled');
    }
  },

  actions: {
    transitionToGame(gamerId) {
      this.sendAction('transitionToGame', gamerId);
    }
  }
});
