import Ember from 'ember';

const {
  Controller,
  inject,
  get
} = Ember;

export default Controller.extend({
  api: inject.service(),
  signUpData: {},

  actions: {
    submit() {
      let data = get(this, 'signUpData');
      get(this, 'api').signUp(data).then(() => {
        location.replace('/#/admin');
      });
    }
  }
});
