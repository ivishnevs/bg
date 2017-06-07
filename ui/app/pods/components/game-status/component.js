import Ember from 'ember';

const {
	Component,
	computed,
	get,
	set
} = Ember;

export default Component.extend({
	status: computed('currentStep', 'stepsNumber', function () {
		if (get(this, 'currentStep') === 1) {
			set(this, 'textColor', 'orange');
			return 'Не начата';
		} else if (get(this, 'currentStep') > get(this, 'stepsNumber')) {
			set(this, 'textColor', 'black');
			return 'Завершена';
		} else {
			set(this, 'textColor', 'green');
			return 'Идет';
		}
	}),

	textColor: null
});
