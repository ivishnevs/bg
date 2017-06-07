import Ember from 'ember';
import {roleName} from '../../helpers/role-name';

const {
	Controller,
	computed,
	get,
	set
} = Ember;

export default Controller.extend({
	selectedParam: 'gamerOrder',
	gameStatistics: computed('model', 'selectedParam', function () {
		let gameStatistic = {
			labels: [],
			datasets: []
		};

		let gamerCount = get(this, 'model').length;
		let selectedParam = get(this, 'selectedParam');

		get(this, 'model').sort(function (a, b) {
			return b.gamerRole - b.gamerRole;
		});
		get(this, 'model').forEach(function (obj) {
			obj.stats.sort(function (a, b) {
				return a.step - b.step;
			});
		});

		get(this, 'model')[0].stats.forEach(function (stats) {
			gameStatistic.labels.push(stats.step);
		});
		get(this, 'model').forEach(function(model) {
			let r = Math.pow(model.gamerRole + 2, 5).toString(16);
			let g = Math.pow(model.gamerRole + 2, 6).toString(16);
			let b = Math.pow(model.gamerRole + 2, 7).toString(16);
			let gamerColor = `#${r[0]}${r[1]}${g[0]}${g[1]}${b[0]}${b[1]}`;
			gameStatistic.datasets.push({
				label: roleName([model.gamerRole, gamerCount]),
				data: model.stats.map(function (stat) {
					return stat[selectedParam];
				}),
				fill: false,
				borderColor: gamerColor,
				pointBorderColor: gamerColor,
				pointBackgroundColor: gamerColor,
				backgroundColor: gamerColor,
				tension: 0.1
			});
		});

		return gameStatistic;
	}),
	actions: {
		selectParam(param) {
			set(this, 'selectedParam', param);
		}
	}
});
