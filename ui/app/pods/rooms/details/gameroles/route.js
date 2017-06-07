import Ember from 'ember';

const {
  Route,
  inject,
  get
} = Ember;

export default Route.extend({
  api: inject.service(),

  breadCrumb: {
    title: 'Роли'
  },

  model(params) {
    return get(this, 'api').getGameDetails(params.game_id);
  },

	afterModel(model) {
		if (model.isGameFinished) {
			this.transitionTo('statistics', model.gameId);
		}
	},

  actions: {
    didTransition() {
      this.controllerFor('rooms.details').set('isGameOpen', true);
    },
    transitionToGame(gamerId) {
      this.transitionTo('gamer', gamerId);
    }
  }
});
