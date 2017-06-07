import Ember from 'ember';

const {
	Route,
	inject,
	get
} = Ember;

export default Route.extend({
	api: inject.service(),

	model(params) {
		return get(this, 'api').getGameStatistics(params.game_id);
	}
});
