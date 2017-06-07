import Ember from 'ember';

const {
  Route,
  inject,
	get,
  set
} = Ember;

export default Route.extend({
	api: inject.service(),

  actions: {
    didTransition() {
      let appController = this.controllerFor('application');
      let adminController = this.controllerFor('admin');
      get(this, 'api').getCurrentUser().then((user) => {
        set(appController, 'user', user);
        set(adminController, 'user', user);
      });
    }
  }
});
