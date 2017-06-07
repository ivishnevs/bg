import Ember from 'ember';

const {
	Component,
	computed,
	set,
	get,
	run,
	on
} = Ember;

export default Component.extend({
	initialSetting: on('init', function () {
		run.scheduleOnce('afterRender', this, function() {
			set(this, 'validationErrors', []);
			set(this, 'currentGamerCount', get(this, 'gamerCount'));
			set(this, 'currentStepsNumber', get(this, 'stepsNumber'));
			set(this, 'currentHoldingCost', get(this, 'holdingCost'));
			set(this, 'currentBackorderCost', get(this, 'backorderCost'));
			set(this, 'currentDemandPattern', get(this, 'demandPattern'));
		});
	}),

	setMessageFor(field, msg) {
		let validationErrors = get(this, 'validationErrors');
		for (let i = 0; i < validationErrors.length; i++) {
			if (validationErrors[i].field === field) {
				set(validationErrors[i], 'msg', msg);
				return true;
			}
		}
		validationErrors.pushObject({
			field,
			msg
		});
	},

	fieldValidation(key, value, min, max, errorLabel) {
		if (value >= min && value <= max) {
			this.setMessageFor(`${key}${get(this, 'id')}`, '');
		} else if (value < min) {
			this.setMessageFor(`${key}${get(this, 'id')}`, `Игра ${get(this, 'index') + 1}: ${errorLabel} не может быть меньше ${min}`);
		} else if (value > max) {
			this.setMessageFor(`${key}${get(this, 'id')}`, `Игра ${get(this, 'index') + 1}: ${errorLabel} не может быть больше ${max}`);
		}
	},

	canChangeGamerCount: computed('currentStep', function () {
		if (get(this, 'currentStep') > 1) {
			return true;
		}
			return false;
	}),
	gameExist: true,
	currentGamerCount: computed({
		get() {
			return 0;
		},
		set(key, value) {
			run.scheduleOnce('afterRender', this, function() {
				this.fieldValidation(key, value, 2, 50, 'количество игроков');
			});
			return value;
		}
	}),
	currentStepsNumber: computed({
		get() {
			return 0;
		},
		set(key, value) {
			run.scheduleOnce('afterRender', this, function() {
				this.fieldValidation(key, value, 10, 500, 'количество шагов');
			});
			return value;
		}
	}),
	currentHoldingCost: computed({
		get() {
			return 0;
		},
		set(key, value) {
			run.scheduleOnce('afterRender', this, function() {
				this.fieldValidation(key, value, 0, 100000, 'цена');
			});
			return value;
		}
	}),
	currentBackorderCost: computed({
		get() {
			return 0;
		},
		set(key, value) {
			run.scheduleOnce('afterRender', this, function() {
				this.fieldValidation(key, value, 0, 100000, 'штраф');
			});
			return value;
		}
	}),
	currentDemandPattern: computed({
		get() {
			return 0;
		},
		set(key, value) {
			run.scheduleOnce('afterRender', this, function() {
				this.fieldValidation(key, value, 1, 5, 'модель');
			});
			return value;
		}
	}),

	currentGameSettingsChanged: computed(
		'currentGamerCount',
		'currentStepsNumber',
		'currentHoldingCost',
		'currentBackorderCost',
		'currentDemandPattern',
		'gamerCount',
		'stepsNumber',
		'holdingCost',
		'backorderCost',
		'demandPattern',
		function () {
			return parseInt(get(this, 'gamerCount')) !== parseInt(get(this, 'currentGamerCount')) ||
				parseInt(get(this, 'stepsNumber')) !== parseInt(get(this, 'currentStepsNumber')) ||
				parseFloat(get(this, 'holdingCost')) !== parseFloat(get(this, 'currentHoldingCost')) ||
				parseFloat(get(this, 'backorderCost')) !== parseFloat(get(this, 'currentBackorderCost')) ||
				parseInt(get(this, 'demandPattern')) !== parseInt(get(this, 'currentDemandPattern'));
		}
	),

	actions: {
		restartGame() {
			let gameID = get(this, 'id');
			let data = { status: 'open', occupiedPlaces: 0, currentStep: 1 };
			if (get(this, 'status') === 'in process') {
				if (confirm('В данной игре сейчас играют люди. Вы уверены что хотите сбросить игру?')) {
					this.sendAction('restartGame', gameID, data);
					set(this, 'status', data.status);
					set(this, 'occupiedPlaces', data.occupiedPlaces);
					set(this, 'currentStep', data.currentStep);
				}
				return;
			}
			this.sendAction('restartGame', gameID, data);
			set(this, 'status', data.status);
			set(this, 'occupiedPlaces', data.occupiedPlaces);
			set(this, 'currentStep', data.currentStep);
		},
		deleteGame() {
			if (confirm(`Вы действительно хотите удалить Игру ${get(this, 'gameCount')}?`)) {
				let validationErrors = get(this, 'validationErrors');
				let gameID = get(this, 'id');
				this.sendAction('deleteGame', gameID);
				set(this, 'gameExist', false);
				for (let i = 0; i < validationErrors.length; i++) {
					if (~validationErrors[i].field.indexOf(gameID)) {
						set(validationErrors[i], 'msg', '');
						validationErrors.splice(i, 1);
						i--;
					}
				}
			}
		},
		saveGameSettings() {
			let validationErrors = get(this, 'validationErrors');
			let isGameValid = true;
			let gameID = get(this, 'id');
			let data = {
				gamerCount: parseInt(get(this, 'currentGamerCount')),
				stepsNumber: parseInt(get(this, 'currentStepsNumber')),
				holdingCost: parseFloat(get(this, 'currentHoldingCost')),
				backorderCost: parseFloat(get(this, 'currentBackorderCost')),
				demandPattern: parseFloat(get(this, 'currentDemandPattern'))
			};

			validationErrors.forEach(function (error) {
				if (isGameValid === true && ~error.field.indexOf(gameID) && error.msg.length !== 0) {
					isGameValid = false;
				}
			});
			if (isGameValid) {
				this.sendAction('saveGameSettings', gameID, data);
				set(this, 'gamerCount', get(this, 'currentGamerCount'));
				set(this, 'stepsNumber', get(this, 'currentStepsNumber'));
				set(this, 'holdingCost', get(this, 'currentHoldingCost'));
				set(this, 'backorderCost', get(this, 'currentBackorderCost'));
				set(this, 'demandPattern', get(this, 'currentDemandPattern'));
			}
		},
		resetGameSettings() {
			set(this, 'currentGamerCount', get(this, 'gamerCount'));
			set(this, 'currentStepsNumber', get(this, 'stepsNumber'));
			set(this, 'currentHoldingCost', get(this, 'holdingCost'));
			set(this, 'currentBackorderCost', get(this, 'backorderCost'));
			set(this, 'currentDemandPattern', get(this, 'demandPattern'));
		},
		transitionToStatistics() {
			this.sendAction('transitionToStatistics', get(this, 'id'));
		}
	}
});
