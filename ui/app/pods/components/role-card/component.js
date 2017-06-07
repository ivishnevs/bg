import Ember from 'ember';

const {
	Component,
  computed,
	inject,
  get
} = Ember;

export default Component.extend({
	ws: inject.service(),
  api: inject.service(),

	channel: computed(function() {
		return `game-${get(this, 'gameId')}-roles`;
	}),

  cardActiveClass: computed('isActive', 'api.currentGamerID', function () {
    let currentGamerID = get(this, 'api.currentGamerID');
    if (currentGamerID && parseInt(currentGamerID) === get(this, 'gamerId')) {
      return 'enabled';
    }
    if (get(this, 'isActive')) {
      return 'disabled';
    }
    return 'enabled';
  }),

  isCurrentRole: computed(function() {
    console.log(get(this, 'api.currentGamerID'));
    console.log(get(this, 'gamerId'));
    return parseInt(get(this, 'api.currentGamerID')) === get(this, 'gamerId');
  }),

	actions: {
		roleSelected() {
      if ($(`.${get(this, 'title')}`).hasClass('disabled')) {
        return;
      }

			let event = {
				"channel": `${get(this, 'channel')}`,
				"action": {
					"type": "role.selected",
					"data": `${get(this, 'title')}`
				}
			};
			get(this, 'ws').send(event, true);
			// console.log(get(this, 'gamerId'));
			this.sendAction('transitionToGame', get(this, 'gamerId'));
		}
	}
});
