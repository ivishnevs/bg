import Ember from 'ember';
import ENV from './../config/environment';

const {
  Service,
	RSVP,
  set,
  $
} = Ember;

export default Service.extend({
  currentGamerID: null,

	getRooms() {
		return $.getJSON(`${ENV.apiURL}api/v1/rooms/`);
	},

  getRoomDetails(roomId) {
		return $.getJSON(`${ENV.apiURL}api/v1/rooms/${roomId}/`);
  },

  getGameDetails(gameId) {
		return $.getJSON(`${ENV.apiURL}api/v1/games/${gameId}/`);
  },

	saveRoom(id, data) {
    let url = `${ENV.apiURL}api/v1/rooms/${id}/`;
    return new RSVP.Promise((resolve, reject) => {
      $.ajax({
        type: 'PUT',
        url,
        data: JSON.stringify(data),
        success(data) {
          resolve(data);
        },
        error(request, textStatus, error) {
          reject(request, textStatus, error);
        }
      });
    });
  },

  fetchGamerData(gamerId) {
    let url = `${ENV.apiURL}api/v1/gamers/${gamerId}/`;
    let that = this;
    return new RSVP.Promise((resolve, reject) => {
      $.ajax({
        type: 'GET',
        url,
        xhrFields: { withCredentials: true },
        success(data) {
          if (data.redirect_url) {
            set(that, 'currentGamerID', null);
            window.location.replace(data.redirect_url);
          } else {
            set(that, 'currentGamerID', gamerId);
            resolve(data);
          }
        },
        error(request, textStatus, error) {
          reject(request, textStatus, error);
        }
      });
    });
	},

  makeOrderRequest(gamerId, order) {
    let url = `${ENV.apiURL}api/v1/gamers/${gamerId}/`;
    let data = { order: parseInt(order) };
    let that = this;
    return new RSVP.Promise((resolve, reject) => {
      $.ajax({
        type: 'POST',
        xhrFields: { withCredentials: true },
        url,
        data: JSON.stringify(data),
        success(data) {
          if (data.redirect_url) {
            set(that, 'currentGamerID', null);
            window.location.replace(data.redirect_url);
          } else {
            resolve(data);
          }
        },
        error(request, textStatus, error) {
          reject(request, textStatus, error);
        }
      });
    });
  },

	deleteGame(id) {
		let url = `${ENV.apiURL}api/v1/games/${id}/`;
		return new RSVP.Promise((resolve, reject) => {
			$.ajax({
				type: 'DELETE',
        xhrFields: { withCredentials: true },
				url,
				success(data) {
					resolve(data);
				},
				error(request, textStatus, error) {
					reject(request, textStatus, error);
				}
			});
		});
	},

	createGame(gameData) {
    return new RSVP.Promise((resolve, reject) => {
      let url = `${ENV.apiURL}api/v1/games/`;
      let data = JSON.stringify(gameData);
      $.ajax({
        type: 'POST',
        xhrFields: { withCredentials: true },
        url,
        data,
        success(data) {
          resolve(data);
        },
        error(request, textStatus, error) {
          reject(error);
        }
      });
    });
	},

	saveGameSettings(id, data) {
		let url = `${ENV.apiURL}api/v1/games/${id}/`;
		return new RSVP.Promise((resolve, reject) => {
			$.ajax({
				type: 'PUT',
        xhrFields: { withCredentials: true },
				url,
				data: JSON.stringify(data),
				success(data) {
					resolve(data);
				},
				error(request, textStatus, error) {
					reject(request, textStatus, error);
				}
			});
		});
	},

	getGameStatistics(gameId) {
		return $.getJSON(`${ENV.apiURL}api/v1/games/${gameId}/statistics/`);
	},

  signUp(data) {
    let url = `${ENV.apiURL}api/v1/accounts/signup/`;
    return new RSVP.Promise((resolve, reject) => {
      $.ajax({
        type: 'POST',
        url,
        data,
        xhrFields: { withCredentials: true },
        success(data) {
          resolve(data);
        },
        error(request, textStatus, error) {
          reject(request, textStatus, error);
        }
      });
    });
  },

  signIn(data) {
    let url = `${ENV.apiURL}api/v1/accounts/signin/`;
    return new RSVP.Promise((resolve, reject) => {
      $.ajax({
        type: 'POST',
        url,
        data,
        xhrFields: { withCredentials: true },
        success(data) {
          resolve(data);
        },
        error(request, textStatus, error) {
          reject(request, textStatus, error);
        }
      });
    });
  },

  signOut() {
    let url = `${ENV.apiURL}api/v1/accounts/signout/`;
    return new RSVP.Promise((resolve, reject) => {
      $.ajax({
        type: 'GET',
        url,
        xhrFields: { withCredentials: true },
        success(data) {
          resolve(data);
        },
        error(request, textStatus, error) {
          reject(request, textStatus, error);
        }
      });
    });
  },

	getCurrentUser() {
    let url = `${ENV.apiURL}api/v1/accounts/current/`;
    return new RSVP.Promise((resolve, reject) => {
      $.ajax({
        type: 'GET',
        url,
        xhrFields: { withCredentials: true },
        success(data) {
          resolve(data);
        },
        error(request, textStatus, error) {
          reject(request, textStatus, error);
        }
      });
    });
	}
});
