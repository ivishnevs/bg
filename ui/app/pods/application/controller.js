import Ember from 'ember';

const {
	Controller,
  computed,
	inject,
	get,
	set,
	on
} = Ember;

export default Controller.extend({
	ws: inject.service(),
	api: inject.service(),

	setupWSConnection: on('init', function () {
		get(this, 'ws').setupWSConn();
	}),

  setCurrentGamerID: on('init', function() {
    if (document.cookie) {
      let gamerID = this.getCookie('gamerID');
      console.log(gamerID);
      set(this, 'api.currentGamerID', gamerID);
    }
  }),

  getCookie(name) {
    let value = "; " + document.cookie;
    let parts = value.split("; " + name + "=");
    if (parts.length === 2) {
      return parts.pop().split(";").shift();
    }
  },

  isReturnToGameShown: computed('currentPath', 'api.currentGamerID', function() {
    return !!get(this, 'api.currentGamerID') && get(this, 'currentPath') !== 'gamer';
  }),

  // isAuthorized: computed('model', function() {
  //   return false;
  // }),
  //
  // user: computed('model', function() {
  //   return get(this, 'model');
  // }),

	// checkingAuthorization: on('init', function() {
    // get(this, 'api').isAuthorized().then((data) => {
    //   let name = data.name;
    //   if (name) {
    //     set(this, 'isAuthorized', true);
    //     set(this, 'user', data);
    //   }
    // });
    // }),

  actions: {
    signOutHandler() {
      get(this, 'api').signOut().then(() => {
        location.replace('/');
      });
    }
  }
});
