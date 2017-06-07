import Ember from 'ember';
import { task, timeout } from 'ember-concurrency';

const {
	Controller,
	computed,
	inject,
	set,
	get
} = Ember;

export default Controller.extend({
	api: inject.service(),
  user: computed({
    get() {
      return null;
    },
    set(key, value) {
      if (!value.Room) {
        return null;
      }
      get(this, 'fetchRoomTask').perform(value.Room.id);
      return value
    }
  }),

  fetchRoomTask: task(function*(roomID) {
    let room = yield get(this, 'api').getRoomDetails(roomID);
    set(this, 'model', room);
  }).drop(),

  createGameTask: task(function*(gameData) {
    yield get(this, 'api').createGame(gameData);
    let user = get(this, 'user');
    get(this, 'fetchRoomTask').perform(user.Room.id);
    set(this, 'modal', false);
  }).drop(),

	setMessageFor(field, msg) {
		let newGameValidationErrors = get(this, 'newGameValidationErrors');
		for (let i = 0; i < newGameValidationErrors.length; i++) {
			if (newGameValidationErrors[i].field === field) {
				set(newGameValidationErrors[i], 'msg', msg);
				return true;
			}
		}
		newGameValidationErrors.pushObject({
			field,
			msg
		});
	},

	fieldValidation(key, value, min, max, errorLabel) {
		if (value >= min && value <= max) {
			this.setMessageFor(`${key}`, '');
		} else if (value < min) {
			this.setMessageFor(`${key}`, `${errorLabel} не может быть меньше ${min}`);
		} else if (value > max) {
			this.setMessageFor(`${key}`, `${errorLabel} не может быть больше ${max}`);
		}
	},

	gamesSortingDesc: ['CreatedAt:desc'],
	games: computed.sort('model.games', 'gamesSortingDesc'),

	newGameValidationErrors: [],
	validationErrors: [],
	currentRoomName: computed('model', {
		get() {
			return get(this, 'model.name');
		},
		set(key, value) {
			return value;
		}
	}),
	currentRoomDescription: computed('model', {
		get() {
			return get(this, 'model.description');
		},
		set(key, value) {
			return value;
		}
	}),
	gamerCount: computed({
		get() {
			return 4;
		},
		set(key, value) {
			this.fieldValidation(key, value, 2, 50, 'Количество игроков');
			return value;
		}
	}),
	stepsNumber: computed({
		get() {
			return 40;
		},
		set(key, value) {
			this.fieldValidation(key, value, 10, 500, 'Количество шагов');
			return value;
		}
	}),
	holdingCost: computed({
		get() {
			return 0.5;
		},
		set(key, value) {
			this.fieldValidation(key, value, 0, 100000, 'Цена');
			return value;
		}
	}),
	backorderCost: computed({
		get() {
			return 1;
		},
		set(key, value) {
			this.fieldValidation(key, value, 0, 100000, 'Штраф');
			return value;
		}
	}),
	demandPattern: computed({
		get() {
			return 1;
		},
		set(key, value) {
			this.fieldValidation(key, value, 1, 1, 'Модель');
			return value;
		}
	}),

	currentRoomNameChanged: computed('currentRoomName', 'model.name', function () {
		return (get(this, 'currentRoomName') !== get(this, 'model.name'));
	}),
	currentRoomDescriptionChanged: computed('currentRoomDescription', 'model.description', function () {
		return (get(this, 'currentRoomDescription') !== get(this, 'model.description'));
	}),

	actions: {
		saveRoomName() {
			get(this, 'api').saveRoom(get(this, 'model.id'), {
				name: get(this, 'currentRoomName'),
				description: get(this, 'model.description')
			});
			set(this, 'model.name', get(this, 'currentRoomName'));
		},
		saveRoomDescription() {
			get(this, 'api').saveRoom(get(this, 'model.id'), {
				name: get(this, 'model.name'),
				description: get(this, 'currentRoomDescription')
			});
			set(this, 'model.description', get(this, 'currentRoomDescription'));
		},
		resetRoomName() {
			set(this, 'currentRoomName', get(this, 'model.name'));
		},
		resetRoomDescription() {
			set(this, 'currentRoomDescription', get(this, 'model.description'));
		},
		createGame() {
			let newGameValidationErrors = get(this, 'newGameValidationErrors');
			let isGameValid = true;
			let data = {
				RoomID: get(this, 'model.id'),	// roomID
				gamerCount: parseInt(get(this, 'gamerCount')),
				stepsNumber: parseInt(get(this, 'stepsNumber')),
				holdingCost: parseFloat(get(this, 'holdingCost')),
				backorderCost: parseFloat(get(this, 'backorderCost')),
				demandPattern: parseFloat(get(this, 'demandPattern'))
			};

			newGameValidationErrors.forEach(function (error) {
				if (isGameValid === true && error.msg.length !== 0) {
					isGameValid = false;
				}
			});
			if (isGameValid) {
				get(this, 'createGameTask').perform(data);
			}
		},
		restartGame(gameID, data) {
			get(this, 'api').saveGameSettings(gameID, data);
		},
		deleteGame(gameID) {
			get(this, 'api').deleteGame(gameID);
		},
		saveGameSettings(gameID, data) {
			get(this, 'api').saveGameSettings(gameID, data);
		},
		transitionToStatistics(gameID) {
			this.transitionToRoute('statistics', gameID);
		}
	}
});
