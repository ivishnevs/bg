import Ember from 'ember';
import { task, timeout } from 'ember-concurrency';

const {
	Controller,
  computed,
	observer,
  inject,
	get,
	set,
  run,
  $
} = Ember;

export default Controller.extend({
  api: inject.service(),
  ws: inject.service(),

	order: computed({
		get() {
			return 0;
		},
		set(key, value) {
			if (value >= 0 && value <= 1000) {
				set(this, 'validationErrorMessage', '');
				set(this, 'validationError', false);
			}
			if (value < 0) {
				set(this, 'validationErrorMessage', 'Заказ не может быть отрицательным');
				set(this, 'validationError', true);
			}
			if (value > 1000) {
				set(this, 'validationErrorMessage', 'Заказ не может быть больше 1000');
				set(this, 'validationError', true);
			}
			return value;
		}
	}),
  _isOrderMade: false,
	isOrderMade: computed('gamerData', '_isOrderMade', function() {
    return get(this, 'gamerData.isStepCompleted') || get(this, '_isOrderMade');
  }),
	resetOrderMadeFlag: observer('model', function () {
		set(this, '_isOrderMade', false);
	}),
	validationError: false,
	validationErrorMessage: '',

	gamersSorting: ['id'],
  gamers: computed.sort('gamerData.gamerMetadata.gamers', 'gamersSorting'),

  gamerData: computed('model', function() {
    return get(this, 'model');
  }),

  channel: computed('gamerData', function() {
    return `game-${get(this, 'gamerData.gamerMetadata.gameID')}`;
  }),

	statsOfLastSteps: computed('gamerData', function () {
		let statsOfLastSteps = [];
		let numberOfLastSteps = 11; // one more than displayed
		let firstStepInStats = get(this, 'gamerData.gamerMetadata.currentStep') - numberOfLastSteps;
		let currentStep = get(this, 'gamerData.gamerMetadata.currentStep');

		get(this, 'gamerData.stats').forEach(function (stats) {
			if (stats.step > firstStepInStats && stats.step !== currentStep) {
				statsOfLastSteps.push(stats);
			}
		});
		statsOfLastSteps.sort(function (a, b) {
			return b.step - a.step;
		});
		return statsOfLastSteps;
	}),

	performWSFlow() {
		get(this, 'ws').subscribeTo(get(this, 'channel'));
		get(this, 'ws').onMessage(this.messageHandler, this);
	},

  messageHandler(e) {
    let event = JSON.parse(e.data);

    if (event.action.type === 'gamer.step-completed') {
      $(`.role-${event.action.data}`).addClass('indicator-ready');
			if (get(this, 'gamerData.gamerMetadata.gamerCount') === $('.indicator-ready').length) {
			  run.later(this, function () {
          this.send('refreshModel');
        }, 500);
			}
    }
  },

  makeOrderTask: task(function * () {
		if (get(this, 'isOrderMade')) {
			return;
		}
		if (get(this, 'validationError')) {
			return;
		}
  	yield timeout(500);
    set(this, '_isOrderMade', true);
    yield get(this, 'api').makeOrderRequest(get(this, 'gamerData.id'), get(this, 'order'));
  }).drop()
});
