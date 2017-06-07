import Ember from 'ember';

const {
	Component,
	computed,
	get
} = Ember;

export default Component.extend({
	step: computed('currentStep', 'stepsNumber', function () {
		if (get(this, 'currentStep') > get(this, 'stepsNumber')) {
			return `${get(this, 'currentStep')-1}/${get(this, 'stepsNumber')}`
		}
		return `${get(this, 'currentStep')}/${get(this, 'stepsNumber')}`
	})
});
