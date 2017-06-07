import Ember from 'ember';

const {
  Route,
  inject,
  get,
  set
} = Ember;

export default Route.extend({
  api: inject.service(),

  breadCrumb: {},

  model(params) {
    return get(this, 'api').getRoomDetails(params.room_id);
  },

  afterModel(model) {
    set(this, 'breadCrumb', {
      title: model.name
    });
  },

  actions: {
    didTransition() {
      this.controllerFor('rooms.details').set('isGameOpen', false);
    }
  }
});
