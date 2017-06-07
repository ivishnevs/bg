import Ember from 'ember';

const {
	inject,
	Route,
	get,
	set
} = Ember;

export default Route.extend({
	api: inject.service(),

	model(params) {
		return get(this, 'api').fetchGamerData(params.gamer_id);
	},

	afterModel(model) {
		if (model.isGameFinished) {
			this.transitionTo('statistics', model.gameId);
		}
	},

  actions: {
    refreshModel() {
      this.refresh();
    },
		didTransition() {
    	set(this.controller, 'order', 0);
			get(this, 'controller').performWSFlow();
		},
  }
});
