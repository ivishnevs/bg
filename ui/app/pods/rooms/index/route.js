import Ember from 'ember';

const {
  Route,
  inject,
  get
} = Ember;

export default Route.extend({
  api: inject.service(),

  breadCrumb: {
    title: 'Комнаты'
  },

  model() {
    return get(this, 'api').getRooms();
  }
});
