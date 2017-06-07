import Ember from 'ember';

const {
  Controller,
  inject,
  get
} = Ember;

export default Controller.extend({
  api: inject.service(),
  signInData: {},

  actions: {
    submit() {
      let data = get(this, 'signInData');
      get(this, 'api').signIn(data).then(() => {
        location.replace('/#/admin');
      });
    }
  }
});
